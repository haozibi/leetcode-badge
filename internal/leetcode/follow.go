package leetcode

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/haozibi/leetcode-badge/internal/leetcodecn"
)

type FollowInfo struct {
	Following int // 关注数量
	Followers int // 被关注数量
}

// GetFollow 被关注\关注
func GetFollow(name string) (*FollowInfo, error) {
	if name == "" {
		return nil, errors.Errorf("miss param")
	}

	var (
		uri    = "https://leetcode-cn.com/graphql/noty"
		method = http.MethodPost
		client = http.DefaultClient
		p      = DataFollow{}
	)

	query := func(name string) string {
		return fmt.Sprintf("{\n    \"operationName\": \"followCounts\",\n    \"variables\": {\n        \"userSlug\": \"" + name + "\"\n    },\n    \"query\": \"query followCounts($userSlug: String!) {\\n  followers(userSlug: $userSlug) {\\n    allNum\\n     }\\n  followingEntities(userSlug: $userSlug) {\\n    allNum\\n    }\\n}\\n\"\n}")
	}

	if err := leetcodecn.Send(client, uri, method, query(name), &p); err != nil {
		return nil, err
	}

	return &FollowInfo{
		Followers: p.Followers.AllNum,
		Following: p.FollowingEntities.AllNum,
	}, nil
}
