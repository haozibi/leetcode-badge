package leetcode

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/haozibi/leetcode-badge/internal/request"
)

// GetUserProfile get user profile by request leetcode
func GetUserProfile(userName string, isCN bool) (*UserProfile, error) {

	if userName == "" {
		return nil, errors.New("name is nil")
	}

	if isCN {
		return getCNUserProfile(userName)
	}
	return getUserProfile(userName)
}

func getCNUserProfile(name string) (*UserProfile, error) {

	var (
		uri    = "https://leetcode-cn.com/graphql"
		method = http.MethodPost
		client = http.DefaultClient
		p      = LeetCodeUserProfile{}
	)

	var query = func(userName string) string {
		s := fmt.Sprintf("{\"operationName\":\"userPublicProfile\",\"variables\":{\"userSlug\":\"%s\"},\"query\":\"query userPublicProfile($userSlug: String!) {\\nuserProfilePublicProfile(userSlug: $userSlug) {\\nusername\\nhaveFollowed\\nsiteRanking\\nprofile {\\nuserSlug\\nrealName\\nuserAvatar\\nlocation\\ncontestCount\\nasciiCode\\n__typename\\n}\\n submissionProgress {\\ntotalSubmissions\\nwaSubmissions\\nacSubmissions\\nreSubmissions\\notherSubmissions\\nacTotal\\nquestionTotal\\n__typename\\n}\\n__typename\\n}\\n}\\n\"}", userName)
		return s
	}

	if err := Send(client, uri, method, query(name), &p); err != nil {
		return nil, err
	}

	if p.UserProfilePublicProfile.Username == "" {
		return nil, nil
	}

	pp := p.UserProfilePublicProfile
	userProfile := &UserProfile{
		UserSlug:         pp.Profile.UserSlug,
		RealName:         pp.Profile.RealName,
		UserAvatar:       pp.Profile.UserAvatar,
		SiteRanking:      pp.SiteRanking,
		TotalSubmissions: pp.SubmissionProgress.TotalSubmissions,
		AcSubmissions:    pp.SubmissionProgress.AcSubmissions,
		WaSubmissions:    pp.SubmissionProgress.WaSubmissions,
		ReSubmissions:    pp.SubmissionProgress.ReSubmissions,
		OtherSubmissions: pp.SubmissionProgress.OtherSubmissions,
		AcTotal:          pp.SubmissionProgress.AcTotal,
		QuestionTotal:    pp.SubmissionProgress.QuestionTotal,
	}

	return userProfile, nil
}

// 部分数据不全
func getUserProfile(userName string) (*UserProfile, error) {

	url := "https://leetcode.com/graphql"

	var genQueryJSON = func(userName string) io.Reader {

		s := fmt.Sprintf("{\"operationName\":\"getUserProfile\",\"variables\":{\"username\":\"%s\"},\"query\":\"query getUserProfile($username: String!) {\\n  allQuestionsCount {\\n    difficulty\\n    count\\n    __typename\\n  }\\n  matchedUser(username: $username) {\\n    username\\n    socialAccounts\\n    githubUrl\\n    contributions {\\n      points\\n      questionCount\\n      testcaseCount\\n      __typename\\n    }\\n    profile {\\n      realName\\n      websites\\n      countryName\\n      skillTags\\n      company\\n      school\\n      starRating\\n      aboutMe\\n      userAvatar\\n      reputation\\n      ranking\\n      __typename\\n    }\\n    submissionCalendar\\n    submitStats {\\n      acSubmissionNum {\\n        difficulty\\n        count\\n        submissions\\n        __typename\\n      }\\n      totalSubmissionNum {\\n        difficulty\\n        count\\n        submissions\\n        __typename\\n      }\\n      __typename\\n    }\\n    __typename\\n  }\\n}\\n\"}", userName)

		return strings.NewReader(s)
	}

	req, err := http.NewRequest("POST", url, genQueryJSON(userName))
	if err != nil {
		return nil, err
	}

	req.Header.Add("origin", "https://leetcode.com")
	req.Header.Add("user-agent", GetUserAgent())
	req.Header.Add("content-type", "application/json")
	req.Header.Add("referer", "https://leetcode.com")

	body, _, err := request.SendRequest(req)
	if err != nil {
		return nil, err
	}

	var p GetUserProfileResult
	err = json.Unmarshal(body, &p)
	if err != nil {
		return nil, errors.Wrap(err, "json parse")
	}

	pp := p.Data.MatchedUser
	userProfile := &UserProfile{
		RealName:    pp.Profile.RealName,
		UserSlug:    fmt.Sprintf("/%s", userName),
		UserAvatar:  pp.Profile.UserAvatar,
		SiteRanking: pp.Profile.Ranking,
	}
	for _, submission := range pp.SubmitStats.AcSubmissionNum {
		if submission.Difficulty == "All" {
			userProfile.AcSubmissions = submission.Submissions
			userProfile.AcTotal = submission.Count
		}
	}
	for _, submission := range pp.SubmitStats.TotalSubmissionNum {
		if submission.Difficulty == "All" {
			userProfile.TotalSubmissions = submission.Submissions
		}
	}
	for _, submission := range p.Data.AllQuestionsCount {
		if submission.Difficulty == "All" {
			userProfile.QuestionTotal = submission.Count
		}
	}

	return userProfile, nil
}

func cleanText(s string) string {
	s = strings.Replace(s, " ", "", -1)
	s = strings.Replace(s, "\t", "", -1)
	s = strings.Replace(s, "\n", "", -1)
	return s
}
