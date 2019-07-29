package redis

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/haozibi/leetcode-badge/internal/cache"
	"github.com/haozibi/leetcode-badge/internal/leetcode"

	"github.com/go-redis/redis"
	"github.com/haozibi/zlog"
	"github.com/pkg/errors"
)

type redisCache struct {
	client         *redis.Client
	expirationTime time.Duration
}

// New new redis
func New(addr string, passwd string) (cache.Cache, error) {

	client := redis.NewClient(&redis.Options{
		Addr:       addr,
		Password:   passwd,
		DB:         0,
		MaxRetries: 2,
	})

	ping, err := client.Ping().Result()
	if err != nil {
		return nil, errors.Wrap(err, "link redis")
	}

	zlog.ZInfo().Str("Ping", ping).Msg("[redis]")

	return &redisCache{
		client:         client,
		expirationTime: cache.DefaultExpirationTime,
	}, nil
}

func userProfileKey(name string, isCN bool) string {
	if isCN {
		return "user_profile_cn_" + name
	}
	return "user_profile_" + name
}

func (m *redisCache) GetUserProfile(name string, isCN bool) (*leetcode.UserProfile, error) {

	name = userProfileKey(name, isCN)

	v, err := m.client.Get(name).Result()
	if err != nil {
		return nil, errors.Wrap(err, "redis get")
	}

	var p *leetcode.UserProfile

	err = gobDeValue([]byte(v), &p)
	return p, errors.Wrap(err, "GetUserProfile")
}

func (m *redisCache) SaveUserProfile(name string, isCN bool, value *leetcode.UserProfile) error {

	name = userProfileKey(name, isCN)

	v, err := gobEnValue(value)
	if err != nil {
		return errors.Wrap(err, "redis set")
	}

	return errors.Wrap(m.client.Set(name, v, m.expirationTime).Err(), "SaveUserProfile")
}

func (m *redisCache) SaveBadge(name string, body []byte) error {

	return errors.Wrap(m.client.Set(name, body, 24*7*time.Hour).Err(), "SaveBadge")
}

func (m *redisCache) GetBadge(name string) ([]byte, error) {

	v, err := m.client.Get(name).Result()
	if err != nil {
		return nil, errors.Wrap(err, "redis get")
	}
	return []byte(v), nil
}

func gobEnValue(value interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(value)
	if err != nil {
		return []byte(""), errors.Wrap(err, "gobEnValue")
	}
	return b.Bytes(), nil
}

func gobDeValue(data []byte, p interface{}) error {
	var b = bytes.NewBuffer(data)
	d := gob.NewDecoder(b)
	return errors.Wrap(d.Decode(p), "gobDeValue")
}
