package sqlite

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"

	"github.com/haozibi/leetcode-badge/internal/storage"
	"github.com/haozibi/leetcode-badge/internal/tools"
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

	createTable := false

	if !tools.Exists(path) {
		createTable = true
	}

	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		return nil, errors.Wrapf(err, "sqlite open error, path: %s", path)
	}

	if err := db.Ping(); err != nil {
		return nil, errors.Wrapf(err, "sqlite ping error, path: %s", path)
	}

	l := &lite{
		path:          path,
		db:            db,
		tableUserInfo: "user_info",
		tableRecord:   "history_record",
	}

	if createTable {
		if err := l.createTable(); err != nil {
			return nil, err
		}
	}

	return l, nil
}

func (l *lite) GetUserInfo(userSlug string, isCN bool) ([]storage.UserInfo, error) {

	var list []storage.UserInfo

	err := l.db.Select(&list, fmt.Sprintf("SELECT * FROM %s WHERE user_slug=? AND is_cn=? LIMIT 1", l.tableUserInfo), userSlug, isCN)

	return list, errors.Wrapf(err, "GetUserInfo error, user_slug: %s, is_cn: %v", userSlug, isCN)
}

func (l *lite) SaveUserInfo(info storage.UserInfo) (int64, error) {

	sql := fmt.Sprintf("INSERT INTO %s (user_slug, real_name, user_avatar, is_cn, updated_time, created_time) VALUES (?, ?, ?, ?, ?, ?)",
		l.tableUserInfo)

	res, err := l.db.Exec(sql, info.UserSlug, info.RealName, info.UserAvatar, info.IsCN, info.UpdatedTime, info.CreatedTime)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return 0, errors.WithStack(storage.ErrHasExist(info.UserSlug))
		}
		return 0, errors.Wrapf(err, "user: %+v", info)
	}

	return res.LastInsertId()
}
func (l *lite) ListUserInfo(start, limit int) ([]storage.UserInfo, error) {

	sql := fmt.Sprintf("SELECT * FROM %s ORDER BY id LIMIT ? OFFSET ?", l.tableUserInfo)

	var list []storage.UserInfo
	err := l.db.Select(&list, sql, limit, start)
	return list, errors.Wrapf(err, "ListUserInfo error, start: %d, limit: %d", start, limit)
}
func (l *lite) ListRecord(userSlug string, isCN bool, start, end time.Time) ([]storage.HistoryRecord, error) {
	sql := fmt.Sprintf("SELECT * FROM %s WHERE user_slug=? AND is_cn=? AND zero_time<=? AND zero_time>=? ORDER BY id", l.tableRecord)

	var list []storage.HistoryRecord
	err := l.db.Select(&list, sql, userSlug, isCN, end, start)

	return list, errors.Wrapf(err, "ListRecord error, user_slug: %s, is_cn: %v", userSlug, isCN)
}
func (l *lite) SaveRecord(info storage.HistoryRecord) error {

	sql := fmt.Sprintf("INSERT INTO %s (user_slug, is_cn, ranking, solved_num, zero_time, created_time) VALUES (?, ?, ?, ?, ?, ?)",
		l.tableRecord)

	_, err := l.db.Exec(sql, info.UserSlug, info.IsCN, info.Ranking, info.SolvedNum, info.ZeroTime, info.CreatedTime)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return errors.WithStack(storage.ErrHasExist(info.UserSlug))
		}
		return errors.Wrapf(err, "save record error, record: %+v", info)
	}

	return nil
}
