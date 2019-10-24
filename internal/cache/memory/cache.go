package memory

import (
	"fmt"
	"time"

	"github.com/haozibi/leetcode-badge/internal/cache"
	"github.com/haozibi/leetcode-badge/internal/leetcode"
	"github.com/haozibi/leetcode-badge/internal/storage"

	gocache "github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
)

type memoryCache struct {
	store          *gocache.Cache
	expirationTime time.Duration
}

// New new
func New() cache.Cache {
	return &memoryCache{
		store:          gocache.New(5*time.Minute, 10*time.Minute),
		expirationTime: cache.DefaultExpirationTime,
	}

}

func (m *memoryCache) GetUserProfile(name string, isCN bool) (*leetcode.UserProfile, error) {

	name = userProfileKey(name, isCN)

	if x, found := m.store.Get(name); found {
		return x.(*leetcode.UserProfile), nil
	}
	return nil, errors.New("not found")
}

func (m *memoryCache) SaveUserProfile(name string, isCN bool, value *leetcode.UserProfile) error {

	name = userProfileKey(name, isCN)

	m.store.Set(name, value, m.expirationTime)
	return nil
}

func userProfileKey(name string, isCN bool) string {
	if isCN {
		return "user_profile_cn_" + name
	}
	return "user_profile_" + name
}

func historyKey(name string, isCN bool, start, end time.Time) string {
	return fmt.Sprintf("%s_%v_%d_%d", name, isCN, start.Unix(), end.Unix())
}

func (m *memoryCache) GetHistoryRecord(name string, isCN bool, start, end time.Time) ([]storage.HistoryRecord, error) {

	key := historyKey(name, isCN, start, end)
	if x, found := m.store.Get(key); found {
		return x.([]storage.HistoryRecord), nil
	}
	return nil, errors.New("not found")
}

func (m *memoryCache) SaveHistoryRecord(name string, isCN bool, start, end time.Time, info []storage.HistoryRecord) error {
	key := historyKey(name, isCN, start, end)
	m.store.Set(key, info, m.expirationTime)
	return nil
}

func (m *memoryCache) SaveByteBody(name string, body []byte) error {

	m.store.Set(name, body, gocache.NoExpiration)
	return nil
}

func (m *memoryCache) GetByteBody(name string) ([]byte, error) {

	if x, found := m.store.Get(name); found {
		return x.([]byte), nil
	}
	return nil, errors.New("not found")
}
