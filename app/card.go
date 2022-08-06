package app

import (
	"net/http"
	"time"

	"github.com/haozibi/leetcode-badge/internal/card"
	"github.com/haozibi/leetcode-badge/internal/leetcodecn"
	"github.com/pkg/errors"
)

func (a *APP) getCard(badgeType BadgeType, name string, r *http.Request) ([]byte, error) {

	var f func(string, *http.Request) ([]byte, error)
	switch badgeType {
	case BadgeTypeQuestionProcessCard:
		f = a.getQuestionProcess
	default:
		return nil, errors.Errorf("not found card function")
	}

	query := r.URL.Query().Encode()
	key := badgeType.String() + "_" + name + query

	body, err := a.cache.GetByteBody(key)
	if err == nil && len(body) != 0 {
		return body, nil
	}
	fn := func() (interface{}, error) {
		body, err := f(name, r)
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

func (a *APP) getQuestionProcess(name string, r *http.Request) ([]byte, error) {
	data, err := leetcodecn.GetUserQuestionProgress(name)
	if err != nil {
		return nil, err

	}
	if data == nil {
		return nil, ErrUserNotSupport
	}

	body, err := card.Build(data)
	return body, err
}
