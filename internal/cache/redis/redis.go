package redis

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"

	"github.com/haozibi/leetcode-badge/internal/cache"
	"github.com/haozibi/leetcode-badge/internal/leetcode"
	"github.com/haozibi/leetcode-badge/internal/storage"
)

type redisCache struct {
	client         *redis.Client
	expirationTime time.Duration
}

// New new redis
func New(addr string, password string) (cache.Cache, error) {

	if addr == "" {
		return nil, errors.New("miss redis address")
	}

	client := redis.NewClient(&redis.Options{
		Addr:       addr,
		Password:   password,
		DB:         0,
		MaxRetries: 2,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, errors.Wrap(err, "link redis")
	}

	return &redisCache{
		client:         client,
		expirationTime: cache.DefaultExpirationTime,
	}, nil
}

func userProfileKey(name string, isCN bool) string {
	return "user_profile:" + name + ":" + strconv.FormatBool(isCN)
}

func (m *redisCache) GetUserProfile(name string, isCN bool) (*leetcode.UserProfile, error) {

	name = userProfileKey(name, isCN)

	v, err := m.client.Get(name).Result()
	if err != nil {
		return nil, errors.Wrap(err, "redis get")
	}

	var p *leetcode.UserProfile

	err = gobDeValue([]byte(v), &p)
	return p, errors.Wrap(err, "redis get user profile")
}

func (m *redisCache) SaveUserProfile(name string, isCN bool, value *leetcode.UserProfile) error {

	name = userProfileKey(name, isCN)

	v, err := gobEnValue(value)
	if err != nil {
		return errors.Wrap(err, "redis set")
	}

	return errors.Wrap(m.client.Set(name, v, m.expirationTime).Err(), "redis save user profile")
}

func (m *redisCache) SaveByteBody(name string, body []byte) error {

	return errors.Wrap(m.client.Set(name, body, 24*7*time.Hour).Err(), "redis save byte body")
}

func (m *redisCache) GetByteBody(name string) ([]byte, error) {

	v, err := m.client.Get(name).Result()
	if err != nil {
		return nil, errors.Wrap(err, "redis get")
	}
	return []byte(v), nil
}

func recordKey(name string, isCN bool, start, end time.Time) string {

	return fmt.Sprintf("record:%s:%v:%d_%d", name, isCN, start.Unix(), end.Unix())
}

func (m *redisCache) GetHistoryRecord(name string, isCN bool, start, end time.Time) ([]storage.HistoryRecord, error) {

	key := recordKey(name, isCN, start, end)
	v, err := m.client.Get(key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "redis get")
	}

	var p []storage.HistoryRecord

	err = gobDeValue([]byte(v), &p)
	return p, errors.Wrap(err, "redis get record history")
}

func (m *redisCache) SaveHistoryRecord(name string, isCN bool, start, end time.Time, info []storage.HistoryRecord) error {

	key := recordKey(name, isCN, start, end)
	v, err := gobEnValue(info)
	if err != nil {
		return errors.Wrap(err, "redis set")
	}
	return errors.Wrap(m.client.Set(key, v, m.expirationTime).Err(), "redis save record history")
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
