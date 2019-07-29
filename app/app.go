package app

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/haozibi/leetcode-badge/internal/cache"
	"github.com/haozibi/leetcode-badge/internal/cache/memory"
	"github.com/haozibi/leetcode-badge/internal/cache/redis"

	"github.com/gorilla/mux"
	"github.com/haozibi/zlog"
)

type APP struct {
	debug  bool
	config *Config
	cache  cache.Cache
	err    error
}

func New(c *Config) *APP {

	a := new(APP)
	a.config = c
	a.debug = c.Debug

	if a.debug {
		zlog.NewBasicLog(os.Stderr, zlog.WithDebug(true))
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

	r := mux.NewRouter()

	setRouter(r, a, ioutil.Discard)

	zlog.ZInfo().Str("Addr", a.config.ListenAddr).Msg("[http]")
	if err := http.ListenAndServe(a.config.ListenAddr, r); err != nil {
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
