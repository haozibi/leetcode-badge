package app

import (
	"bytes"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/haozibi/leetcode-badge/internal/chart"
	"github.com/haozibi/leetcode-badge/internal/storage"
	"github.com/haozibi/leetcode-badge/internal/tools"
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
		return tools.SVCNotFound, nil
	}

	list, err := a.getHistoryList(name, isCN, start, end)
	if err != nil {
		return nil, err
	}

	if len(list) < 2 {
		// 数据不足，过几天再来查看
		return tools.SVGLackData, nil
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
