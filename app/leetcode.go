package app

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/haozibi/leetcode-badge/internal/leetcode"
	"github.com/haozibi/leetcode-badge/internal/shield"

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

func iscn(r *http.Request) bool {
	if strings.HasPrefix(r.URL.Path, "/v1cn") {
		return true
	}
	return false
}

func (a *APP) getUserProfile(r *http.Request) (*leetcode.UserProfile, error) {

	var info *leetcode.UserProfile
	var err error

	vars := mux.Vars(r)
	name := vars["name"]

	if strings.Count(name, "/") >= 1 {
		return nil, errors.New("name error")
	}

	isCN := iscn(r)

	info, err = a.cache.GetUserProfile(name, isCN)
	if info != nil && err == nil {
		return info, nil
	}

	info, err = leetcode.GetUserProfile(name, isCN)
	if err != nil {
		return nil, err
	}

	if info == nil {
		zlog.ZInfo().Str("User", name).Msg("[profile] user not found")
		return nil, nil
	}

	zlog.ZInfo().Str("User", name).Msg("[profile] success")

	go func() {
		a.cache.SaveUserProfile(name, isCN, info)
	}()

	return info, nil
}

func (a *APP) getBadge(r *http.Request, typeName badgeType, info *leetcode.UserProfile) ([]byte, error) {

	isCN := iscn(r)

	var key, left, right string

	if info == nil {
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

	body, err := a.cache.GetBadge(key)
	if err == nil && len(body) != 0 {
		return body, nil
	}

	badgeBody, err := shield.Badge(r, left, right, DefaultColor)
	if err != nil {
		return nil, err
	}

	go func() {
		a.cache.SaveBadge(key, badgeBody)
	}()

	return badgeBody, nil
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

func (a *APP) Badge(w http.ResponseWriter, r *http.Request) {

	uri := strings.TrimPrefix(r.URL.Path, "/")

	uriList := strings.Split(uri, "/")

	if len(uriList) != 3 {
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

	info, err := a.getUserProfile(r)
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

func (a *APP) Profile(w http.ResponseWriter, r *http.Request) {

	info, err := a.getUserProfile(r)
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
