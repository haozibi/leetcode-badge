package storage

import "time"

// UserInfo user info
type UserInfo struct {
	ID          int64  `json:"id"`
	UserSlug    string `json:"user_slug"`
	IsCN        int    `json:"is_cn"`
	RealName    string `json:"real_name"`
	UserAvatar  string `json:"user_avatar"`
	UpdatedTime int64  `json:"updated_time"`
	CreatedTime int64  `json:"created_time"`
}

// HistoryRecord record info
type HistoryRecord struct {
	ID          int64  `json:"id"`
	UserSlug    string `json:"user_slug"`
	IsCN        int    `json:"is_cn"`
	Ranking     int    `json:"ranking"`
	SolvedNum   int    `json:"solved_num"`
	ZeroTime    int64  `json:"zero_time"`
	CreatedTime int64  `json:"created_time"`
}

type Storage interface {
	GetUserInfo(userslug string, iscn bool) ([]UserInfo, error)
	SaveUserInfo(info UserInfo) (int64, error)
	ListUserInfo(start, limit int) ([]UserInfo, error)

	ListRecord(userslug string, iscn bool, start, end time.Time) ([]HistoryRecord, error)
	SaveRecord(info HistoryRecord) error
}
