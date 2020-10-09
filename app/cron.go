package app

import (
	"time"

	"github.com/rs/zerolog/log"

	"github.com/haozibi/leetcode-badge/internal/leetcode"
	"github.com/haozibi/leetcode-badge/internal/storage"
	"github.com/haozibi/leetcode-badge/internal/tools"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
)

const (
	CronSpec = "30 8 * * *"
)

// Cron cron
// 30 8 * * * 每天凌晨 8 点 30 分
// "@every 2s"
func (a *APP) Cron(spec string) {
	c := cron.New(cron.WithChain(
		cron.Recover(cron.DefaultLogger),
	))
	_, err := c.AddFunc(spec, func() {
		num := a.cron()
		log.Info().Int("UpdateNum", num).Time("Time", time.Now()).Msg("[cron]")
	})
	if err != nil {
		panic(err)
	}
	c.Start()
	log.Info().Str("Spec", spec).Msg("[cron] start success")
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
			log.Err(err).Msg("[cron] list user info error")
			continue
		}
		if len(userList) == 0 {
			return total
		}

		log.Debug().Int("UserNum", len(userList)).Msg("[cron] find user")

		for i := 0; i < len(userList); i++ {

			name := userList[i].UserSlug
			isCN := tools.IntToBool(userList[i].IsCN)

			if v, ok := a.recordMap[recordKey(name, isCN)]; ok &&
				tools.IsToday(v) {
				continue
			}

			err := a.updateHistory(name, isCN)
			if err != nil {
				log.Err(err).Msg("[cron] update history")
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
