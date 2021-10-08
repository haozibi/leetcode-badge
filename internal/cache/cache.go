package cache

import (
	"time"

	"github.com/haozibi/leetcode-badge/internal/leetcode"
	"github.com/haozibi/leetcode-badge/internal/storage"
)

const (
	// DefaultExpirationTime default expiration time
	DefaultExpirationTime = 15 * time.Minute
)

// Cache cache
type Cache interface {
	GetUserProfile(name string, isCN bool) (*leetcode.UserProfile, error)
	SaveUserProfile(name string, isCN bool, value *leetcode.UserProfile) error

	GetFollow(name string, isCN bool) (*leetcode.FollowInfo, error)
	SaveFollow(name string, isCN bool, value *leetcode.FollowInfo) error

	GetHistoryRecord(name string, isCN bool, start, end time.Time) ([]storage.HistoryRecord, error)
	SaveHistoryRecord(name string, isCN bool, start, end time.Time, info []storage.HistoryRecord) error

	SaveByteBody(name string, body []byte) error
	GetByteBody(name string) ([]byte, error)
}
