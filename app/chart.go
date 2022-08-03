package app

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/haozibi/leetcode-badge/internal/chart"
	"github.com/haozibi/leetcode-badge/internal/heatmap"
	"github.com/haozibi/leetcode-badge/internal/leetcode"
	"github.com/haozibi/leetcode-badge/internal/statics"
	"github.com/haozibi/leetcode-badge/internal/storage"
)

// HistoryChart history chart
func (a *APP) historyChart(badgeType BadgeType, name string, isCN bool, start, end time.Time) ([]byte, error) {

	var showFn ChartShow

	switch badgeType {
	case BadgeTypeChartRanking:
		showFn = a.showRanking
	case BadgeTypeChartSolved:
		showFn = a.showSolved
	default:
		return nil, errors.Errorf("badge not match")
	}

	userInfo, err := a.getUserProfile(name, isCN)
	if err != nil {
		return nil, err
	}

	if userInfo == nil || userInfo.UserSlug == "" {
		// 用户没找到
		return statics.SVGNotFound(), nil
	}

	list, err := a.getHistoryList(name, isCN, start, end)
	if err != nil {
		return nil, err
	}

	if len(list) < 2 {
		// 数据不足，过几天再来查看
		return statics.GetLackSVG(), nil
	}

	body, err := showFn(list, userInfo.RealName)
	return body, err
}

// ChartShow chart function
type ChartShow func(list []storage.HistoryRecord, name string) ([]byte, error)

func (a *APP) showRanking(list []storage.HistoryRecord, name string) ([]byte, error) {
	cc := make([]chart.RankHistory, len(list))

	for i := 0; i < len(list); i++ {
		cc[i].Date = time.Unix(list[i].ZeroTime, 0)
		cc[i].Rank = list[i].Ranking
	}

	buffer := bytes.NewBuffer([]byte{})
	err := chart.ShowRankHistory(buffer, [][]chart.RankHistory{cc}, name)
	return buffer.Bytes(), errors.Wrap(err, "show ranking")
}

func (a *APP) showSolved(list []storage.HistoryRecord, name string) ([]byte, error) {
	cc := make([]chart.SolvedHistory, len(list))

	for i := 0; i < len(list); i++ {
		cc[i].Date = time.Unix(list[i].ZeroTime, 0)
		cc[i].Num = list[i].SolvedNum
	}

	buffer := bytes.NewBuffer([]byte{})
	err := chart.ShowSolvedHistory(buffer, [][]chart.SolvedHistory{cc}, name)
	return buffer.Bytes(), errors.Wrap(err, "show ranking")
}

func (a *APP) getHistoryList(name string, isCN bool, start, end time.Time) ([]storage.HistoryRecord, error) {

	var info []storage.HistoryRecord
	var err error

	info, err = a.cache.GetHistoryRecord(name, isCN, start, end)
	if err == nil && len(info) != 0 {
		return info, nil
	}

	key := fmt.Sprintf("%s_%v_%d_%d",
		name,
		isCN,
		start.Unix(),
		end.Unix(),
	)

	fn := func() (interface{}, error) {
		list, err := a.store.ListRecord(name, isCN, start, end)
		if err != nil {
			return nil, err
		}

		go func() {
			a.cache.SaveHistoryRecord(name, isCN, start, end, list, 30*time.Minute)
		}()
		return list, nil
	}

	result, err, _ := a.group.Do(key, fn)
	if err != nil {
		return nil, err
	}

	return result.([]storage.HistoryRecord), nil
}

// SubCal SubmissionCalendar
func (a *APP) SubCal(_ BadgeType, name string, isCN bool, w http.ResponseWriter, r *http.Request) {
	if !isCN {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 not found"))
		return
	}

	body, err := a.getSubCal(name, r)
	if err != nil {
		if err == ErrUserNotSupport {
			a.write(w, statics.SVGNotFound())
		} else {
			log.Err(err).
				Str("Name", name).
				Bool("IsCN", isCN).
				Msg("get subcal error")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	a.write(w, body)
}

func (a *APP) getSubCal(name string, r *http.Request) ([]byte, error) {
	var (
		query = r.URL.Query().Encode()
		key   = fmt.Sprintf("subcal_%s_oo_%s", name, query)
		body  []byte
		err   error
	)

	var f func(data map[int64]int, color string) ([]byte, error)

	t := r.URL.Query().Get("type")
	color := r.URL.Query().Get("color")

	switch t {
	case "past-year":
		f = heatmap.PastYear
	default:
		f = heatmap.CurrYear
	}

	body, err = a.cache.GetByteBody(key)
	if err == nil && len(body) != 0 {
		return body, nil
	}

	reqKey := "subcal_" + name
	fn := func() (interface{}, error) {
		// now := time.Now().Year()
		data, err := leetcode.GetSubCal(name)
		if err != nil {
			return nil, err
		}
		if len(data) == 0 {
			return nil, ErrUserNotSupport
		}

		res := make(map[int64]int)
		for k, v := range data {
			i, err := strconv.ParseInt(k, 10, 64)
			if err != nil {
				return nil, errors.Wrapf(err, "value: %s", k)
			}

			res[i] = v
		}

		body, err = f(res, color)
		if err != nil {
			return nil, err
		}

		if len(body) == 0 {
			return nil, nil
		}

		err = a.cache.SaveByteBody(key, body, 5*time.Minute)
		return body, err
	}

	result, err, _ := a.group.Do(reqKey, fn)
	if err != nil {
		return nil, err
	}

	return result.([]byte), nil
}
