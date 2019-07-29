package cache

import (
	"time"

	"github.com/haozibi/leetcode-badge/internal/leetcode"
)

type CacheType string

const (
	CacheRedis  CacheType = "redis"
	CacheMemory CacheType = "memory"
)

func (c CacheType) String() string {
	return string(c)
}

const (
	DefaultExpirationTime = 15 * time.Minute
)

type Cache interface {
	GetUserProfile(name string, isCN bool) (*leetcode.UserProfile, error)
	SaveUserProfile(name string, isCN bool, value *leetcode.UserProfile) error
	SaveBadge(name string, body []byte) error
	GetBadge(name string) ([]byte, error)
}
