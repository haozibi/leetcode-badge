package memory

import (
	"time"

	"github.com/haozibi/leetcode-badge/internal/cache"
	"github.com/haozibi/leetcode-badge/internal/leetcode"

	gocache "github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
)

var store *gocache.Cache

func init() {
	store = gocache.New(5*time.Minute, 10*time.Minute)
}

type memoryCache struct {
	expirationTime time.Duration
}

// New new
func New() cache.Cache {
	return &memoryCache{
		expirationTime: cache.DefaultExpirationTime,
	}
}

func (m *memoryCache) GetUserProfile(name string, isCN bool) (*leetcode.UserProfile, error) {

	name = userProfileKey(name, isCN)

	if x, found := store.Get(name); found {
		return x.(*leetcode.UserProfile), nil
	}
	return nil, errors.New("not found")
}

func (m *memoryCache) SaveUserProfile(name string, isCN bool, value *leetcode.UserProfile) error {

	name = userProfileKey(name, isCN)

	store.Set(name, value, m.expirationTime)
	return nil
}

func userProfileKey(name string, isCN bool) string {
	if isCN {
		return "user_profile_cn_" + name
	}
	return "user_profile_" + name
}

func (m *memoryCache) SaveBadge(name string, body []byte) error {

	store.Set(name, body, gocache.NoExpiration)
	return nil
}

func (m *memoryCache) GetBadge(name string) ([]byte, error) {

	if x, found := store.Get(name); found {
		return x.([]byte), nil
	}
	return nil, errors.New("not found")
}
