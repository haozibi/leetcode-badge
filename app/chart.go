package app

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/haozibi/leetcode-badge/internal/chart"
	"github.com/haozibi/leetcode-badge/internal/storage"
	"github.com/haozibi/leetcode-badge/static"

	"github.com/gorilla/mux"
	"github.com/haozibi/zlog"
	"github.com/pkg/errors"
)

// HistoryChart history chart
func (a *APP) HistoryChart(w http.ResponseWriter, r *http.Request) {

	uri := strings.TrimPrefix(r.URL.Path, "/")
	uriList := strings.Split(uri, "/")

	if len(uriList) != 4 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	vars := mux.Vars(r)
	name := vars["name"]

	if strings.Count(name, "/") >= 1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	day := 7
	dayStr := r.URL.Query().Get("day")
	if dayStr != "" {
		d, err := strconv.Atoi(dayStr)
		if err != nil {
			fmt.Fprintf(w, "day params must number")
			return
		}
		if d > 30 {
			d = 30
		}
		if d > day {
			day = d
		}
	}

	end := time.Now()
	start := time.Now().AddDate(0, 0, -1*(day-1))

	var showFn ChartShow

	switch strings.ToLower(uriList[2]) {
	case "ranking":
		showFn = a.showRanking
	case "solved":
		showFn = a.showSolved
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}

	userInfo, err := a.getUserProfile(name, a.iscn(r))
	if err != nil {
		zlog.ZError().Msgf("[http] %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if userInfo == nil || userInfo.UserSlug == "" {
		// 用户没找到
		zlog.ZDebug().Msgf("%s %v not found", name, a.iscn(r))
		a.write(w, static.MustAsset("static/svg/notfound.svg"))
		return
	}

	list, err := a.getHistoryList(name, a.iscn(r), start, end)
	if err != nil {
		zlog.ZError().Msgf("[http] %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(list) < 2 {
		// 数据不足，过几天再来查看
		a.write(w, static.MustAsset("static/svg/lackdata.svg"))
		return
	}

	body, err := showFn(list, userInfo.RealName)
	if err != nil {
		zlog.ZError().Msgf("[http] %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	a.write(w, body)
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
			a.cache.SaveHistoryRecord(name, isCN, start, end, list)
		}()
		return list, nil
	}

	result, err, _ := a.group.Do(key, fn)
	if err != nil {
		return nil, err
	}

	return result.([]storage.HistoryRecord), nil
}
