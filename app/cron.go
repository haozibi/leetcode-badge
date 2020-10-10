package app

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"

	"github.com/haozibi/leetcode-badge/internal/leetcode"
	"github.com/haozibi/leetcode-badge/internal/storage"
	"github.com/haozibi/leetcode-badge/internal/tools"
)

const (
	CronSpec = "30 8 * * *"
)

// Cron cron
// 30 8 * * * 每天凌晨 8 点 30 分
func (a *APP) Cron(spec string) {
	c := cron.New(cron.WithChain(
		cron.Recover(cron.DefaultLogger),
	))
	_, err := c.AddFunc(spec, func() {
		a.cron()

	})
	if err != nil {
		panic(err)
	}
	c.Start()
	log.Info().Str("Spec", spec).Msg("[cron] start success")
}

func (a *APP) cron() {

	var (
		j     int
		total int
		limit = 100
	)

	t1 := time.Now()

	for {
		start := j * limit
		userList, err := a.store.ListUserInfo(start, limit)
		if err != nil {
			log.Err(err).Msg("[cron] list user info error")
			continue
		}
		if len(userList) == 0 {
			break
		}

		log.Debug().Int("UserNum", len(userList)).Msg("[cron] find user")

		for i := 0; i < len(userList); i++ {

			name := userList[i].UserSlug
			isCN := tools.IntToBool(userList[i].IsCN)

			a.rMu.Lock()
			if v, ok := a.recordMap[recordKey(name, isCN)]; ok &&
				tools.IsToday(v) {
				a.rMu.Unlock()
				continue
			}
			a.rMu.Unlock()

			err := a.updateHistory(name, isCN)
			if err != nil {
				log.Err(err).Msg("[cron] update history")
				continue
			}

			zero := tools.ZeroTime(time.Now())

			a.rMu.Lock()
			a.recordMap[recordKey(name, isCN)] = zero
			a.rMu.Unlock()

			total++
			log.Debug().
				Str("Name", name).
				Str("Today", zero.Format("2006-01-02")).
				Bool("IsCN", isCN).
				Msg("[cron] update user success")
		}
		j++
	}

	log.Info().
		Int("UpdateNum", total).
		Str("UseTime", time.Since(t1).String()).
		Msg("[cron] cron success")
}

func (a *APP) updateHistory(name string, isCN bool) error {

	// TODO: 更优化的方法，避免 429 错误
	time.Sleep(100*time.Millisecond + time.Duration(rand.Intn(100))*time.Millisecond)
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
		CreatedTime: time.Now().UnixNano() / 1e6,
	}

	err = a.store.SaveRecord(record)
	if err != nil {
		if !storage.IsHasExistError(err) {
			return errors.Wrapf(err, "user: %s", info.UserSlug)
		}
	}

	return nil
}

func recordKey(name string, isCN bool) string {
	return name + strconv.FormatBool(isCN)
}
