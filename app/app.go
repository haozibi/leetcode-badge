package app

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/haozibi/leetcode-badge/internal/cache"
	"github.com/haozibi/leetcode-badge/internal/cache/memory"
	"github.com/haozibi/leetcode-badge/internal/cache/redis"
	"github.com/haozibi/leetcode-badge/internal/storage"
	"github.com/haozibi/leetcode-badge/internal/storage/mysql"
	"github.com/haozibi/leetcode-badge/static"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/haozibi/zlog"
	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
)

type APP struct {
	debug bool

	config *Config
	cache  cache.Cache
	store  storage.Storage
	group  *singleflight.Group

	mu        sync.Mutex
	userMap   map[string]time.Time
	recordMap map[string]time.Time

	wg               WaitGroupWrapper
	shutdown         chan struct{}
	shutdownComplete chan struct{}
}

func New(c Config) *APP {

	a := &APP{
		debug:            c.Debug,
		config:           &c,
		group:            new(singleflight.Group),
		userMap:          make(map[string]time.Time),
		recordMap:        make(map[string]time.Time),
		shutdown:         make(chan struct{}),
		shutdownComplete: make(chan struct{}),
	}

	return a
}

func (a *APP) Run() error {

	if err := a.initConfig(); err != nil {
		return err
	}

	err := static.RestoreAssets("./", "static")
	if err != nil {
		zlog.ZError().AnErr("Static", err).Msg("[Init]")
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

	go a.quit(3 * time.Second)
	a.Cron(a.config.CronSpec)

	a.wg.Wrap(func() {
		exitFunc(a.runHTTP())
	})

	a.wg.Wrap(func() {
		exitFunc(a.runMonitor())
	})

	select {
	case <-a.shutdownComplete:
		return nil
	case err1 := <-exit:
		return err1
	}
}

func (a *APP) runHTTP() error {
	r := mux.NewRouter()
	setRouter(r, a, ioutil.Discard)

	srv := &http.Server{
		Addr:         a.config.Address,
		WriteTimeout: 120 * time.Second,
		ReadTimeout:  120 * time.Second,
		Handler:      handlers.RecoveryHandler()(r),
	}

	go func() {
		select {
		case <-a.shutdown:
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			err := srv.Shutdown(ctx)
			if err != nil {
				zlog.ZError().Msgf("[http] Shutdown %+v", err)
			}
			select {
			case <-ctx.Done():
				zlog.ZDebug().Msg("[http] timeout of 3 seconds.")
			}
		}
	}()

	zlog.ZInfo().Str("Addr", a.config.Address).Msg("[http]")
	if err := srv.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (a *APP) initConfig() error {
	var err error

	switch a.config.CacheType {
	case "redis":
		a.cache, err = redis.New(
			a.config.RedisConfig.Address,
			a.config.RedisConfig.Password,
		)
		if err != nil {
			return err
		}
	case "memory":
		a.cache = memory.New()
	default:
		return errors.New("not support cache type: " + a.config.CacheType)
	}

	zlog.ZInfo().Msgf("[cache] type: %s", a.config.CacheType)

	switch a.config.StorageType {
	case "mysql":
		a.store, err = mysql.New(
			a.config.MySQLConfig.DBName,
			a.config.MySQLConfig.User,
			a.config.MySQLConfig.Password,
			a.config.MySQLConfig.Host,
			a.config.MySQLConfig.Port,
		)
		if err != nil {
			return err
		}
	default:
		return errors.New("not support storage type: " + a.config.StorageType)
	}

	zlog.ZInfo().Msgf("[storage] type: %s", a.config.StorageType)

	return nil
}

func (a *APP) quit(out time.Duration) error {

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zlog.ZInfo().Msg("[Server] Shutdown Server...")
	close(a.shutdown)

	a.wg.Wait()
	close(a.shutdownComplete)
	return nil
}

// WaitGroupWrapper wrap sync.WaitGroup
type WaitGroupWrapper struct {
	sync.WaitGroup
}

// Wrap wrap
func (w *WaitGroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}

func recordKey(name string, isCN bool) string {
	return name + strconv.FormatBool(isCN)
}
