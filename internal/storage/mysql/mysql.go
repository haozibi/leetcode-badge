package mysql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/haozibi/gendry/builder"
	"github.com/haozibi/gendry/manager"
	"github.com/haozibi/gendry/scanner"
	"github.com/pkg/errors"

	"github.com/haozibi/leetcode-badge/internal/storage"
	"github.com/haozibi/leetcode-badge/internal/tools"
)

func init() {
	scanner.SetTagName("json")
}

type dbHelp struct {
	db *sql.DB

	tableRecordName string
	tableUserName   string
}

func New(dbName, user, passwd, host string, port int) (storage.Storage, error) {

	if dbName == "" || user == "" ||
		passwd == "" || host == "" {
		return nil, errors.New("miss mysql params")
	}

	d, err := manager.New(dbName, user, passwd, host).Set(
		manager.SetCharset("utf8mb4"),
		manager.SetAllowCleartextPasswords(true),
		manager.SetInterpolateParams(true),
		manager.SetTimeout(1*time.Second),
		manager.SetReadTimeout(1*time.Second)).Port(port).Open(true)
	if err != nil {
		return nil, errors.Wrap(err, "mysql link")
	}

	return &dbHelp{
		db:              d,
		tableRecordName: "history_record",
		tableUserName:   "user_info",
	}, nil
}

func (d *dbHelp) GetUserInfo(userslug string, iscn bool) ([]storage.UserInfo, error) {

	where := map[string]interface{}{
		"user_slug": userslug,
		"is_cn":     iscn,
		"_limit":    []uint{1},
	}

	cond, vals, err := builder.BuildSelect(d.tableUserName, where, nil)
	if err != nil {
		return nil, errors.Wrap(err, "[sql] build")
	}

	rows, err := d.db.Query(cond, vals...)
	if err != nil {
		return nil, errors.Wrap(err, "[sql] query")
	}

	var list []storage.UserInfo

	err = scanner.ScanClose(rows, &list)

	return list, errors.Wrap(err, "[sql] scan")
}

func (d *dbHelp) SaveUserInfo(info storage.UserInfo) (int64, error) {

	var data []map[string]interface{}
	data = append(data, tools.Struct2Map(info))

	cond, vals, err := builder.BuildInsert(d.tableUserName, data)
	if err != nil {
		return 0, errors.Wrap(err, "[sql] build")
	}

	result, err := d.db.Exec(cond, vals...)
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062: Duplicate entry") {
			return 0, errors.WithStack(storage.ErrHasExist(info.UserSlug))
		}
		return 0, errors.Wrap(err, "[sql] exec")
	}

	id, err := result.LastInsertId()
	return id, errors.Wrap(err, "[sql] get id")
}

func (d *dbHelp) ListUserInfo(start, limit int) ([]storage.UserInfo, error) {

	sql := fmt.Sprintf("SELECT * FROM %s ORDER BY id LIMIT {{limit}} OFFSET {{start}}", d.tableUserName)

	cond, vals, err := builder.NamedQuery(sql, map[string]interface{}{
		"start": start,
		"limit": limit,
	})
	if err != nil {
		return nil, errors.Wrap(err, "[sql] build")
	}

	rows, err := d.db.Query(cond, vals...)
	if err != nil {
		return nil, errors.Wrap(err, "[sql] query")
	}

	var list []storage.UserInfo

	err = scanner.ScanClose(rows, &list)

	return list, errors.Wrap(err, "[sql] scan")
}

func (d *dbHelp) ListRecord(userslug string, iscn bool, start, end time.Time) ([]storage.HistoryRecord, error) {

	start = tools.ZeroTime(start)
	end = tools.ZeroTime(end)

	sql := fmt.Sprintf("SELECT * FROM %s WHERE user_slug={{userslug}} AND is_cn={{iscn}} AND zero_time<={{end}} AND zero_time>={{start}} ORDER BY zero_time ASC", d.tableRecordName)

	cond, vals, err := builder.NamedQuery(sql, map[string]interface{}{
		"start":    start.Unix(),
		"end":      end.Unix(),
		"userslug": userslug,
		"iscn":     tools.BoolToInt(iscn),
	})
	if err != nil {
		return nil, errors.Wrap(err, "[sql] build")
	}

	rows, err := d.db.Query(cond, vals...)
	if err != nil {
		return nil, errors.Wrap(err, "[sql] query")
	}

	var list []storage.HistoryRecord

	err = scanner.ScanClose(rows, &list)

	return list, errors.Wrap(err, "[sql] scan")
}

func (d *dbHelp) SaveRecord(info storage.HistoryRecord) error {

	var data []map[string]interface{}
	data = append(data, tools.Struct2Map(info))

	cond, vals, err := builder.BuildInsert(d.tableRecordName, data)
	if err != nil {

		return errors.Wrap(err, "[sql] build")
	}

	_, err = d.db.Exec(cond, vals...)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return errors.WithStack(storage.ErrHasExist(info.UserSlug))
		}
		return errors.Wrap(err, "[sql] exec")
	}

	return nil
}
