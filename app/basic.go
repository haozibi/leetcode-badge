package app

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/haozibi/leetcode-badge/internal/leetcode"
	"github.com/haozibi/leetcode-badge/internal/leetcodecn"
	"github.com/haozibi/leetcode-badge/internal/models"
	"github.com/haozibi/leetcode-badge/internal/shield"
	"github.com/haozibi/leetcode-badge/internal/storage"
	"github.com/haozibi/leetcode-badge/internal/tools"
)

const (
	// DefaultColor default color
	DefaultColor = "brightgreen"
)

func buildLeetCodeKey(name string, isCN bool) string {
	return name + "_" + strconv.Itoa(tools.BoolToInt(isCN))
}

func (a *APP) getUserProfile(name string, isCN bool) (*models.UserProfile, error) {

	var info *models.UserProfile
	var err error

	info, err = a.cache.GetUserProfile(name, isCN)
	if err == nil && info != nil {
		return info, nil
	}

	key := buildLeetCodeKey(name, isCN)

	fn := func() (interface{}, error) {
		if isCN {
			info, err = leetcodecn.GetUserProfile(name)
		} else {
			info, err = leetcode.GetUserProfile(name)
		}
		if err != nil {
			return nil, err
		}

		if info == nil {
			info = new(models.UserProfile)
		}

		err = a.cache.SaveUserProfile(name, isCN, info, 5*time.Minute)
		if err != nil {
			return nil, err
		}

		go func() {
			// 必要时改成同步执行
			err := a.saveUser(info, isCN)
			if err != nil {
				log.Err(err).Msg("save user profile")
			}
		}()

		return info, nil
	}

	result, err, _ := a.group.Do(key, fn)
	if err != nil {
		return nil, err
	}

	return result.(*models.UserProfile), nil
}

// basicBadge 根据信息获取 badge
func (a *APP) basicBadge(value url.Values, isCN bool, typeName BadgeType, info *models.UserProfile) ([]byte, error) {

	var key, left, right string

	if info == nil || info.UserSlug == "" {
		typeName = BadgeTypeUserNotFound
	}

	switch typeName {
	case BadgeTypeProfile:
		left = "LeetCode CN"
		if !isCN {
			left = "LeetCode"
		}
		right = info.RealName
		key = "Default_" + left + "_" + right
	case BadgeTypeRanking:
		right = strconv.Itoa(info.SiteRanking)
		if info.SiteRanking >= 100000 {
			right = "≥100000"
		}
		left = info.RealName
		key = "Ranking_" + left + "_" + right
	case BadgeTypeSolved:
		left, right = "Solved", strconv.Itoa(info.AcTotal)+"/"+strconv.Itoa(info.QuestionTotal)
		key = "Solved_" + left + "_" + right
	case BadgeTypeSolvedRate:
		left, right = "Solved", fmt.Sprintf("%.2f％", (float64(info.AcTotal)/float64(info.QuestionTotal))*100)
		key = "SolvedRate_" + left + "_" + right
	case BadgeTypeAccepted:
		left, right = "Accepted", strconv.Itoa(info.AcSubmissions)+"/"+strconv.Itoa(info.TotalSubmissions)
		key = "Accepted_" + left + "_" + right
	case BadgeTypeAcceptedRate:
		left, right = "Accepted", fmt.Sprintf("%.2f％", (float64(info.AcSubmissions)/float64(info.TotalSubmissions))*100)
		key = "AcceptedRate_" + left + "_" + right
	case BadgeTypeUserNotFound:
		left, right = "LeetCodeBadge", "User Not Found"
		key = "UserNotFound_" + left + "_" + right
	}

	key += "_" + strconv.Itoa(tools.BoolToInt(isCN)) + "_" + value.Encode()

	return a.getBadge(value, key, left, right, DefaultColor)
}

func (a *APP) getBadge(value url.Values, key, left, right string, color string) ([]byte, error) {
	body, err := a.cache.GetByteBody(key)
	if err == nil && len(body) != 0 {
		return body, nil
	}

	fn := func() (interface{}, error) {
		badgeBody, err := shield.Badge(value, left, right, color)
		if err != nil {
			return nil, err
		}
		if err = a.cache.SaveByteBody(key, badgeBody, 7*24*time.Hour); err != nil {
			log.Err(err).Str("Key", key).Str("Left", left).Str("Right", right).Str("Color", color).Msg("save badge data error")
		}
		return badgeBody, err
	}

	result, err, _ := a.group.Do(key, fn)
	if err != nil {
		return nil, err
	}

	return result.([]byte), nil
}

func (a *APP) saveUser(info *models.UserProfile, isCN bool) error {
	if info == nil || info.UserSlug == "" {
		return nil
	}

	key := recordKey(info.UserSlug, isCN)

	a.uMu.Lock()
	if _, ok := a.userMap[key]; ok {
		a.uMu.Unlock()
		return nil
	}
	a.uMu.Unlock()

	user := storage.UserInfo{
		UserSlug:    info.UserSlug,
		IsCN:        tools.BoolToInt(isCN),
		RealName:    info.RealName,
		UserAvatar:  info.UserAvatar,
		UpdatedTime: 0,
		CreatedTime: time.Now().UnixNano() / 1e6,
	}

	_, err := a.store.SaveUserInfo(user)
	if err != nil && !storage.IsHasExistError(err) {
		return err
	}

	log.Info().Str("UserName", info.UserSlug).Bool("IsCN", isCN).Msg("[http] save user")

	a.uMu.Lock()
	a.userMap[key] = time.Now()
	a.uMu.Unlock()

	record := storage.HistoryRecord{
		UserSlug:    info.UserSlug,
		IsCN:        tools.BoolToInt(isCN),
		Ranking:     info.SiteRanking,
		SolvedNum:   info.AcTotal,
		ZeroTime:    tools.ZeroTime(time.Now()).Unix(),
		CreatedTime: time.Now().UnixNano() / 1e6,
	}

	err = a.store.SaveRecord(record)
	if err != nil && !storage.IsHasExistError(err) {
		return err
	}
	log.Info().Str("UserName", info.UserSlug).Bool("IsCN", isCN).Msg("[http] save record")

	a.rMu.Lock()
	a.recordMap[key] = time.Now()
	a.rMu.Unlock()

	return nil
}

func (a *APP) write(w http.ResponseWriter, body []byte) {

	w.Header().Set("Access-Control-Expose-Headers", "Content-Type, Cache-Control, Expires")
	w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, max-age=0")
	w.Header().Set("Expires", "0")
	w.Header().Set("Pragme", "no-cache")

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}
