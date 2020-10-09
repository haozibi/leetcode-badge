package sqlite

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"

	"github.com/haozibi/leetcode-badge/internal/storage"
)

type User struct {
	ID   int
	Name string
	Age  int
}

type lite struct {
	path string
	db   *sqlx.DB

	tableUserInfo string
	tableRecord   string
}

func New(path string) (storage.Storage, error) {
	if path == "" {
		return nil, errors.Errorf("sqlite miss path")
	}

	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		return nil, errors.Wrapf(err, "sqlite open error, path: %s", path)
	}

	if err := db.Ping(); err != nil {
		return nil, errors.Wrapf(err, "sqlite ping error, path: %s", path)
	}

	return &lite{
		path: path,
		db:   db,
	}, nil
}

func (l *lite) GetUserInfo(userSlug string, isCN bool) ([]storage.UserInfo, error) {

	var list []storage.UserInfo

	err := l.db.Select(&list, fmt.Sprintf("SELECT * FROM %s WHERE user_slug=? AND is_cn=? LIMIT 1", l.tableUserInfo), userSlug, isCN)

	return list, errors.Wrapf(err, "GetUserInfo error, user_slug: %s, is_cn: %v", userSlug, isCN)
}

func (l *lite) SaveUserInfo(info storage.UserInfo) (int64, error) {
	return 0, nil
}
func (l *lite) ListUserInfo(start, limit int) ([]storage.UserInfo, error) {

	sql := fmt.Sprintf("SELECT * FROM %s ORDER BY id LIMIT ? OFFSET ?", l.tableUserInfo)

	var list []storage.UserInfo
	err := l.db.Select(&list, sql, limit, start)
	return list, errors.Wrapf(err, "ListUserInfo error, start: %d, limit: %d", start, limit)
}
func (l *lite) ListRecord(userslug string, iscn bool, start, end time.Time) ([]storage.HistoryRecord, error) {
	return nil, nil
}
func (l *lite) SaveRecord(info storage.HistoryRecord) error {
	return nil
}

// func New(path string) error {
// 	db, err := sqlx.Open("sqlite3", path)
// 	if err != nil {
// 		return errors.Wrapf(err, "sqlite3, path: %s", path)
// 	}
//
// 	var p []User
//
// 	err = db.Select(&p, "SELECT * FROM boy LIMIT 1")
// 	if err != nil {
// 		return err
// 	}
//
// 	spew.Dump(p)
// 	return nil
// }
