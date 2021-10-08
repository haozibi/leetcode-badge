package app

import (
	"net/url"
	"strconv"

	"github.com/pkg/errors"

	"github.com/haozibi/leetcode-badge/internal/leetcode"
)

func (a *APP) userFollow(name string) (*leetcode.FollowInfo, error) {

	var info *leetcode.FollowInfo
	var err error

	info, err = a.cache.GetFollow(name, true)
	if err == nil && info != nil {
		return info, nil
	}

	key := "follow_" + name
	fn := func() (interface{}, error) {
		info, err = leetcode.GetFollow(name)
		if err != nil {
			return nil, err
		}

		if info == nil {
			info = new(leetcode.FollowInfo)
		}

		err = a.cache.SaveFollow(name, true, info)
		if err != nil {
			return nil, err
		}

		return info, nil
	}

	result, err, _ := a.group.Do(key, fn)
	if err != nil {
		return nil, err
	}

	return result.(*leetcode.FollowInfo), nil
}

func (a *APP) followBadge(value url.Values, typeName BadgeType, info *leetcode.FollowInfo) ([]byte, error) {
	var left, right string

	switch typeName {
	case BadgeTypeFollowers:
		left = "Followers"
		right = strconv.Itoa(info.Followers)
	case BadgeTypeFollowing:
		left = "Following"
		right = strconv.Itoa(info.Following)
	default:
		return nil, errors.Errorf("not match type")
	}

	key := "follow_" + left + "_" + right + "_" + value.Encode()
	return a.getBadge(value, key, left, right, DefaultColor)
}
