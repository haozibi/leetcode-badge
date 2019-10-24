package app

import (
	"time"

	"github.com/haozibi/leetcode-badge/internal/leetcode"
	"github.com/haozibi/leetcode-badge/internal/storage"
	"github.com/haozibi/leetcode-badge/internal/tools"

	"github.com/haozibi/zlog"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
)

// Cron cron
// 30 8 * * * 每天凌晨 8 点 30 分
// "@every 2s"
func (a *APP) Cron(spec string) error {
	c := cron.New(cron.WithChain(
		cron.Recover(cron.DefaultLogger),
	))
	_, err := c.AddFunc(spec, func() {
		num := a.cron()
		zlog.ZInfo().Int("UpdateNum", num).Time("Time", time.Now()).Msg("[cron]")
	})
	if err != nil {
		return errors.Wrap(err, "setup cron")
	}
	c.Start()
	zlog.ZInfo().Msgf("[cron] start cron success: %s", spec)
	return nil
}

func (a *APP) cron() int {

	var (
		i     int
		total int
		limit = 100
	)

	for {
		start := i * limit
		userList, err := a.store.ListUserInfo(start, limit)
		if err != nil {
			zlog.ZError().Msgf("%+v", err)
			continue
		}
		if len(userList) == 0 {
			return total
		}

		zlog.ZInfo().Msgf("[cron] find %d user", len(userList))

		for i := 0; i < len(userList); i++ {

			name := userList[i].UserSlug
			isCN := tools.IntToBool(userList[i].IsCN)

			if v, ok := a.recordMap[recordKey(name, isCN)]; ok &&
				tools.IsToday(v) {
				continue
			}

			err := a.updateHistory(name, isCN)
			if err != nil {
				zlog.ZError().Msgf("[cron] update history %+v", err)
				continue
			}
			a.recordMap[recordKey(name, isCN)] = tools.ZeroTime(time.Now())
			total++
		}
		i++
	}
}

func (a *APP) updateHistory(name string, isCN bool) error {

	info, err := leetcode.GetUserProfile(name, isCN)
	if err != nil {
		return err
	}

	if info == nil {
		return errors.Errorf("%s isCN:%v not found", name, isCN)
	}

	record := storage.HistoryRecord{
		UserSlug:    name,
		IsCN:        tools.BoolToInt(isCN),
		Ranking:     info.SiteRanking,
		SolvedNum:   info.AcTotal,
		ZeroTime:    tools.ZeroTime(time.Now()).Unix(),
		CreatedTime: time.Now().Unix(),
	}

	err = a.store.SaveRecord(record)
	if err != nil {
		if !storage.IsHasExistError(err) {
			return errors.Wrapf(err, "user: %s", info.UserSlug)
		}
	}

	return nil
}
