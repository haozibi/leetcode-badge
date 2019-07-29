package app

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/haozibi/leetcode-badge/internal/cache"
	"github.com/haozibi/leetcode-badge/internal/cache/memory"
	"github.com/haozibi/leetcode-badge/internal/cache/redis"
	"github.com/haozibi/leetcode-badge/static"

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
