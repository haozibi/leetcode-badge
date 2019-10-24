package app

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/haozibi/leetcode-badge/internal/leetcode"
	"github.com/haozibi/leetcode-badge/internal/shield"
	"github.com/haozibi/leetcode-badge/internal/storage"
	"github.com/haozibi/leetcode-badge/internal/tools"

	"github.com/gorilla/mux"
	"github.com/haozibi/zlog"
)

const (
	// DefaultColor default color
	DefaultColor = "brightgreen"
)

type badgeType int

const (
	Default badgeType = iota + 1
	Ranking
	Solved
	SolvedRate
	Accepted
	AcceptedRate
	UserNotFound
)

// Profile user profile
func (a *APP) Profile(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	name := vars["name"]

	if strings.Count(name, "/") >= 1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	info, err := a.getUserProfile(name, a.iscn(r))
	if err != nil {
		zlog.ZError().Msgf("%+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := a.getBadge(r, Default, info)
	if err != nil {
		zlog.ZError().Msgf("%+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	a.write(w, body)
}

// Badge badge
func (a *APP) Badge(w http.ResponseWriter, r *http.Request) {

	uri := strings.TrimPrefix(r.URL.Path, "/")
	uriList := strings.Split(uri, "/")

	if len(uriList) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	vars := mux.Vars(r)
	userName := vars["name"]

	if strings.Count(userName, "/") >= 1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	name := Default

	switch strings.ToLower(uriList[1]) {
	case "ranking":
		name = Ranking
	case "solved":
		name = Solved
	case "solved-rate":
		name = SolvedRate
	case "accepted":
		name = Accepted
	case "accepted-rate":
		name = AcceptedRate
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}

	info, err := a.getUserProfile(userName, a.iscn(r))
	if err != nil {
		zlog.ZError().Msgf("%+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := a.getBadge(r, name, info)
	if err != nil {
		zlog.ZError().Msgf("%+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	a.write(w, body)
}

func (a *APP) getUserProfile(name string, isCN bool) (*leetcode.UserProfile, error) {

	var info *leetcode.UserProfile
	var err error

	info, err = a.cache.GetUserProfile(name, isCN)
	if info != nil && err == nil {
		return info, nil
	}

	key := name + "_" + strconv.Itoa(tools.BoolToInt(isCN))

	fn := func() (interface{}, error) {
		info, err = leetcode.GetUserProfile(name, isCN)
		if err != nil {
			return nil, err
		}

		if info == nil {
			info = new(leetcode.UserProfile)
			zlog.ZInfo().Str("User", name).Msg("[profile] user not found from leetcode")
		} else {
			zlog.ZInfo().Str("User", name).Msg("[profile] fetch from leetcode success")
		}

		err = a.cache.SaveUserProfile(name, isCN, info)
		if err != nil {
			return nil, err
		}

		go func() {
			// 必要时改成同步执行
			err := a.saveUser(info, isCN)
			if err != nil {
				zlog.ZError().Msgf("%+v", err)
			}
		}()

		return info, nil
	}

	result, err, _ := a.group.Do(key, fn)
	if err != nil {
		return nil, err
	}

	return result.(*leetcode.UserProfile), nil
}

func (a *APP) getBadge(r *http.Request, typeName badgeType, info *leetcode.UserProfile) ([]byte, error) {

	isCN := a.iscn(r)

	var key, left, right string

	if info == nil || info.UserSlug == "" {
		typeName = UserNotFound
	}

	switch typeName {
	case Default:
		left = "LeetCode CN"
		if !isCN {
			left = "LeetCode"
		}
		right = info.RealName
		key = "Default_" + left + "_" + right
	case Ranking:
		right = strconv.Itoa(info.SiteRanking)
		if info.SiteRanking >= 100000 {
			right = "≥100000"
		}
		left = info.RealName
		key = "Ranking_" + left + "_" + right
	case Solved:
		left, right = "Solved", strconv.Itoa(info.AcTotal)+"/"+strconv.Itoa(info.QuestionTotal)
		key = "Solved_" + left + "_" + right
	case SolvedRate:
		left, right = "Solved", fmt.Sprintf("%.2f％", (float64(info.AcTotal)/float64(info.QuestionTotal))*100)
		key = "SolvedRate_" + left + "_" + right
	case Accepted:
		left, right = "Accepted", strconv.Itoa(info.AcSubmissions)+"/"+strconv.Itoa(info.TotalSubmissions)
		key = "Accepted_" + left + "_" + right
	case AcceptedRate:
		left, right = "Accepted", fmt.Sprintf("%.2f％", (float64(info.AcSubmissions)/float64(info.TotalSubmissions))*100)
		key = "AcceptedRate_" + left + "_" + right
	case UserNotFound:
		vars := mux.Vars(r)
		left, right = vars["name"], "User Not Found"
		key = "UserNotFound_" + left + "_" + right
	}

	key += r.URL.Query().Encode()

	body, err := a.cache.GetByteBody(key)
	if err == nil && len(body) != 0 {
		return body, nil
	}

	badgeBody, err := shield.Badge(r, left, right, DefaultColor)
	if err != nil {
		return nil, err
	}

	go func() {
		a.cache.SaveByteBody(key, badgeBody)
	}()

	return badgeBody, nil
}

func (a *APP) saveUser(info *leetcode.UserProfile, isCN bool) error {
	if info == nil || info.UserSlug == "" {
		return nil
	}

	key := recordKey(info.UserSlug, isCN)

	a.mu.Lock()
	defer a.mu.Unlock()

	if _, ok := a.userMap[key]; ok {
		return nil
	}

	user := storage.UserInfo{
		UserSlug:    info.UserSlug,
		IsCN:        tools.BoolToInt(isCN),
		RealName:    info.RealName,
		UserAvatar:  info.UserAvatar,
		UpdatedTime: 0,
		CreatedTime: time.Now().Unix(),
	}

	_, err := a.store.SaveUserInfo(user)
	if err != nil && !storage.IsHasExistError(err) {
		return err
	}

	zlog.ZInfo().Str("UserName", info.UserSlug).Msg("[http] save user")
	a.userMap[key] = time.Now()

	record := storage.HistoryRecord{
		UserSlug:    info.UserSlug,
		IsCN:        tools.BoolToInt(isCN),
		Ranking:     info.SiteRanking,
		SolvedNum:   info.AcTotal,
		ZeroTime:    tools.ZeroTime(time.Now()).Unix(),
		CreatedTime: time.Now().Unix(),
	}

	err = a.store.SaveRecord(record)
	if err != nil && !storage.IsHasExistError(err) {
		zlog.ZError().AnErr("SaveRecord", err).Msg("[http]")
		return err
	}
	zlog.ZInfo().Str("UserName", info.UserSlug).Msg("[http] save record")
	a.recordMap[key] = time.Now()

	return nil
}

func (a *APP) write(w http.ResponseWriter, body []byte) {

	w.Header().Set("Access-Control-Expose-Headers", "Content-Type, Cache-Control, Expires")
	w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, max-age=0")
	w.Header().Set("Expires", "0")
	w.Header().Set("Pragme", "no-cache")

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (a *APP) iscn(r *http.Request) bool {
	if strings.HasPrefix(r.URL.Path, "/v1cn") {
		return true
	}
	return false
}
