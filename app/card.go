package app

import (
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/haozibi/leetcode-badge/internal/card"
	"github.com/haozibi/leetcode-badge/internal/i18n"
	"github.com/haozibi/leetcode-badge/internal/leetcodecn"
)

func (a *APP) getCard(badgeType BadgeType, name string, r *http.Request) ([]byte, error) {

	var f func(string, *http.Request, *i18n.Config) ([]byte, error)
	switch badgeType {
	case BadgeTypeQuestionProcessCard:
		f = a.getQuestionProcess
	case BadgeTypeContestRankingCard:
		f = a.getContestRankingInfo
	default:
		return nil, errors.Errorf("not found card function")
	}

	i18Cfg := i18n.Get(r.URL.Query().Get("lang"))
	if i18Cfg == nil {
		i18Cfg = i18n.Get("zh")
	}

	query := r.URL.Query().Encode()
	key := badgeType.String() + "_" + name + query

	body, err := a.cache.GetByteBody(key)
	if err == nil && len(body) != 0 {
		return body, nil
	}
	fn := func() (interface{}, error) {
		body, err := f(name, r, i18Cfg)
		if err != nil {
			return nil, err
		}

		err = a.cache.SaveByteBody(key, body, 5*time.Minute)
		return body, err
	}

	result, err, _ := a.group.Do(key, fn)
	if err != nil {
		return nil, err
	}

	return result.([]byte), nil
}

func (a *APP) getQuestionProcess(name string, r *http.Request, i18Cfg *i18n.Config) ([]byte, error) {
	data, err := leetcodecn.GetUserQuestionProgress(name)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, ErrUserNotSupport
	}

	body, err := card.QuestionProcess(name, data, i18Cfg.QuestionProcess)
	return body, err
}

func (a *APP) getContestRankingInfo(name string, r *http.Request, i18Cfg *i18n.Config) ([]byte, error) {
	data, err := leetcodecn.GetUserContestRankingInfo(name)
	if err != nil {
		return nil, err

	}
	if data == nil {
		return nil, ErrUserNotSupport
	}

	body, err := card.ContestRanking(name, data, i18Cfg.ContestRanking)
	return body, err
}
