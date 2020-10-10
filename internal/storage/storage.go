package storage

import "time"

// UserInfo user info
type UserInfo struct {
	ID          int64  `json:"id" db:"id"`
	UserSlug    string `json:"user_slug" db:"user_slug"`
	IsCN        int    `json:"is_cn" db:"is_cn"`
	RealName    string `json:"real_name" db:"real_name"`
	UserAvatar  string `json:"user_avatar" db:"user_avatar"`
	UpdatedTime int64  `json:"updated_time" db:"updated_time"`
	CreatedTime int64  `json:"created_time" db:"created_time"`
}

// HistoryRecord record info
type HistoryRecord struct {
	ID          int64  `json:"id" db:"id"`
	UserSlug    string `json:"user_slug" db:"user_slug"`
	IsCN        int    `json:"is_cn" db:"is_cn"`
	Ranking     int    `json:"ranking" db:"ranking"`
	SolvedNum   int    `json:"solved_num" db:"solved_num"`
	ZeroTime    int64  `json:"zero_time" db:"zero_time"`
	CreatedTime int64  `json:"created_time" db:"created_time"`
}

type Storage interface {
	GetUserInfo(userSlug string, isCN bool) ([]UserInfo, error)
	SaveUserInfo(info UserInfo) (int64, error)
	ListUserInfo(start, limit int) ([]UserInfo, error)

	ListRecord(userSlug string, isCN bool, start, end time.Time) ([]HistoryRecord, error)
	SaveRecord(info HistoryRecord) error
}
