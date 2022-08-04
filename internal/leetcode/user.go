package leetcode

import (
	"fmt"
	"net/http"

	"github.com/haozibi/leetcode-badge/internal/leetcodecn"
	"github.com/haozibi/leetcode-badge/internal/models"
)

// GetUserProfile 部分数据不全
func GetUserProfile(userName string) (*models.UserProfile, error) {

	var (
		uri    = "https://leetcode.com/graphql"
		method = http.MethodPost
		client = http.DefaultClient
		p      = UserProfileData{}
	)

	var genQueryJSON = func(userName string) string {

		s := fmt.Sprintf("{\"operationName\":\"getUserProfile\",\"variables\":{\"username\":\"%s\"},\"query\":\"query getUserProfile($username: String!) {\\n  allQuestionsCount {\\n    difficulty\\n    count\\n    __typename\\n  }\\n  matchedUser(username: $username) {\\n    username\\n    socialAccounts\\n    githubUrl\\n    contributions {\\n      points\\n      questionCount\\n      testcaseCount\\n      __typename\\n    }\\n    profile {\\n      realName\\n      websites\\n      countryName\\n      skillTags\\n      company\\n      school\\n      starRating\\n      aboutMe\\n      userAvatar\\n      reputation\\n      ranking\\n      __typename\\n    }\\n    submissionCalendar\\n    submitStats {\\n      acSubmissionNum {\\n        difficulty\\n        count\\n        submissions\\n        __typename\\n      }\\n      totalSubmissionNum {\\n        difficulty\\n        count\\n        submissions\\n        __typename\\n      }\\n      __typename\\n    }\\n    __typename\\n  }\\n}\\n\"}", userName)

		return s
	}

	if err := leetcodecn.Send(client, uri, method, genQueryJSON(userName), &p); err != nil {
		return nil, err
	}

	if p.MatchedUser.Profile.RealName == "" {
		return nil, nil
	}

	pp := p.MatchedUser
	userProfile := &models.UserProfile{
		RealName:    pp.Profile.RealName,
		UserSlug:    fmt.Sprintf("%s", userName),
		UserAvatar:  pp.Profile.UserAvatar,
		SiteRanking: pp.Profile.Ranking,
	}
	for _, submission := range pp.SubmitStats.AcSubmissionNum {
		if submission.Difficulty == "All" {
			userProfile.AcSubmissions = submission.Submissions
			userProfile.AcTotal = submission.Count
			break
		}
	}
	for _, submission := range pp.SubmitStats.TotalSubmissionNum {
		if submission.Difficulty == "All" {
			userProfile.TotalSubmissions = submission.Submissions
			break
		}
	}
	for _, submission := range p.AllQuestionsCount {
		if submission.Difficulty == "All" {
			userProfile.QuestionTotal = submission.Count
			break
		}
	}

	return userProfile, nil
}
