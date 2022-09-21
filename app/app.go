package app

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/singleflight"

	"github.com/haozibi/leetcode-badge/internal/cache"
	"github.com/haozibi/leetcode-badge/internal/cache/memory"
	"github.com/haozibi/leetcode-badge/internal/i18n"
	"github.com/haozibi/leetcode-badge/internal/storage"
	"github.com/haozibi/leetcode-badge/internal/storage/sqlite"
	"github.com/haozibi/leetcode-badge/internal/tools"
)

type APP struct {
	config Config

	cache cache.Cache
	store storage.Storage

	group     *singleflight.Group
	uMu       sync.Mutex
	rMu       sync.Mutex
	userMap   map[string]time.Time
	recordMap map[string]time.Time

	userInfoMu sync.Mutex
	userInfo   map[string]time.Time

	wg               tools.WaitGroupWrapper
	shutdown         chan struct{}
	shutdownComplete chan struct{}
}

func New(cfg Config) *APP {

	spew.Dump(cfg)

	a := &APP{
		config:           cfg,
		group:            new(singleflight.Group),
		userMap:          make(map[string]time.Time),
		recordMap:        make(map[string]time.Time),
		userInfo:         make(map[string]time.Time),
		shutdown:         make(chan struct{}),
		shutdownComplete: make(chan struct{}),
	}

	return a
}

func (a *APP) Setup() (err error) {
	path := a.config.SqlitePath

	a.cache = memory.New()
	if a.store, err = sqlite.New(path); err != nil {
		return err
	}

	err = i18n.InitI18n()
	return err
}

func (a *APP) Run() error {

	var (
		enable = a.config.EnableCron
	)

	if err := a.Setup(); err != nil {
		return err
	}

	exit := make(chan error)
	var once sync.Once

	exitFunc := func(err error) {
		once.Do(func() {
			if err != nil {
				exit <- err
			}
		})
	}

	go a.quit()
	a.wg.Wrap(func() {
		exitFunc(a.runHTTP())
	})
	a.wg.Wrap(func() {
		exitFunc(a.runInterHTTP())
	})

	if enable {
		a.Cron(CronSpec)
		log.Info().Msg("enable cron")
	}

	select {
	case <-a.shutdownComplete:
		return nil
	case err1 := <-exit:
		return err1
	}
}

func (a *APP) runHTTP() error {

	var (
		address = a.config.Address
	)

	r := mux.NewRouter()
	Router(r, a, ioutil.Discard)
	handle := handlers.RecoveryHandler()(handlers.CompressHandler(r))
	return a.runhttp(address, handle)
}

func (a *APP) runhttp(address string, handle http.Handler) error {

	srv := &http.Server{
		Addr:         address,
		WriteTimeout: 120 * time.Second,
		ReadTimeout:  120 * time.Second,
		Handler:      handle,
	}

	go func() {
		select {
		case <-a.shutdown:
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			err := srv.Shutdown(ctx)
			if err != nil {
				log.Err(err).Msg("shutdown error")
			}
			select {
			case <-ctx.Done():
				log.Debug().Msg("[http] timeout of 3 seconds.")
			}
		}
	}()

	log.Info().Str("Address", address).Msg("http listen")
	if err := srv.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (a *APP) quit() {

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("[Server] Shutdown Server...")
	close(a.shutdown)

	a.wg.Wait()
	close(a.shutdownComplete)
	return
}
