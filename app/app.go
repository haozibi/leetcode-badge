package app

import (
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/haozibi/leetcode-badge/internal/cache"
	"github.com/haozibi/leetcode-badge/internal/cache/memory"
	"github.com/haozibi/leetcode-badge/internal/cache/redis"
	"github.com/haozibi/leetcode-badge/static"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/haozibi/zlog"
	"golang.org/x/sync/singleflight"
)

type APP struct {
	debug  bool
	config *Config
	cache  cache.Cache
	group  *singleflight.Group
	err    error
}

func New(c *Config) *APP {

	a := new(APP)
	a.config = c
	a.debug = c.Debug
	a.group = new(singleflight.Group)

	if a.debug {
		zlog.NewBasicLog(os.Stdout, zlog.WithDebug(true))
	}

	switch c.CacheType {
	case cache.CacheRedis:
		a.cache, a.err = redis.New(c.CacheAddr, c.CachePasswd)
	case cache.CacheMemory:
		a.cache = memory.New()
	}

	return a
}

func (a *APP) Run() error {

	if a.err != nil {
		return a.err
	}

	err := static.RestoreAssets("./", "static")
	if err != nil {
		zlog.ZError().AnErr("Static", err).Msg("[Init]")
		return err
	}

	exit := make(chan error)
	var once sync.Once
	var wg WaitGroupWrapper

	exitFunc := func(err error) {
		once.Do(func() {
			if err != nil {
				exit <- err
			}
		})
	}

	wg.Wrap(func() {
		exitFunc(a.runHTTP())
	})

	wg.Wrap(func() {
		exitFunc(a.runMonitor())
	})

	err1 := <-exit
	return err1
}

func (a *APP) runHTTP() error {
	r := mux.NewRouter()
	setRouter(r, a, ioutil.Discard)

	srv := &http.Server{
		Addr:         a.config.ListenAddr,
		WriteTimeout: 120 * time.Second,
		ReadTimeout:  120 * time.Second,
		Handler:      handlers.RecoveryHandler()(r),
	}

	zlog.ZInfo().Str("Addr", a.config.ListenAddr).Msg("[http]")
	if err := srv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

type Config struct {
	CacheType   cache.CacheType
	CacheAddr   string
	CachePasswd string
	ListenAddr  string
	Debug       bool
}

type WaitGroupWrapper struct {
	sync.WaitGroup
}

func (w *WaitGroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}
