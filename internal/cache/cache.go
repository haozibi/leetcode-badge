package cache

import (
	"time"

	gocache "github.com/patrickmn/go-cache"

	"github.com/haozibi/leetcode-badge/internal/leetcode"
	"github.com/haozibi/leetcode-badge/internal/storage"
)

const (
	// DefaultExpirationTime default expiration time
	DefaultExpirationTime = 15 * time.Minute

	NoExpiration = gocache.NoExpiration
)

// Cache cache
type Cache interface {
	GetUserProfile(name string, isCN bool) (*leetcode.UserProfile, error)
	SaveUserProfile(name string, isCN bool, value *leetcode.UserProfile, timeout time.Duration) error

	GetFollow(name string, isCN bool) (*leetcode.FollowInfo, error)
	SaveFollow(name string, isCN bool, value *leetcode.FollowInfo, timeout time.Duration) error

	GetHistoryRecord(name string, isCN bool, start, end time.Time) ([]storage.HistoryRecord, error)
	SaveHistoryRecord(name string, isCN bool, start, end time.Time, info []storage.HistoryRecord, timeout time.Duration) error

	SaveByteBody(name string, body []byte, timeout time.Duration) error
	GetByteBody(name string) ([]byte, error)
}
