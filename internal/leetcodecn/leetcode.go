package leetcodecn

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/pkg/errors"

	"github.com/haozibi/leetcode-badge/internal/models"
)

func GetUserQuestionProgress(name string) (*models.UserQuestionPrecess, error) {
	ptr, err := Get(models.RequestTypeUserQuestionProgress, name)
	if err != nil {
		return nil, err
	}
	p := ptr.(*LeetCodeUserQuestionProgress)
	if len(p.UserProfileUserQuestionProgress.NumUntouchedQuestions) == 0 {
		return nil, nil
	}

	res := models.UserQuestionPrecess{}

	for _, v := range p.UserProfileUserQuestionProgress.NumAcceptedQuestions {
		if v.Difficulty == "EASY" {
			res.Easy.AcceptedNum += v.Count
		}
		if v.Difficulty == "MEDIUM" {
			res.Medium.AcceptedNum += v.Count
		}
		if v.Difficulty == "HARD" {
			res.Hard.AcceptedNum += v.Count
		}
		res.Overview.AcceptedNum += v.Count
	}

	for _, v := range p.UserProfileUserQuestionProgress.NumFailedQuestions {
		if v.Difficulty == "EASY" {
			res.Easy.FailedNum += v.Count
		}
		if v.Difficulty == "MEDIUM" {
			res.Medium.FailedNum += v.Count
		}
		if v.Difficulty == "HARD" {
			res.Hard.FailedNum += v.Count
		}
		res.Overview.FailedNum += v.Count
	}

	for _, v := range p.UserProfileUserQuestionProgress.NumUntouchedQuestions {
		if v.Difficulty == "EASY" {
			res.Easy.UntouchedNum += v.Count
		}
		if v.Difficulty == "MEDIUM" {
			res.Medium.UntouchedNum += v.Count
		}
		if v.Difficulty == "HARD" {
			res.Hard.UntouchedNum += v.Count
		}
		res.Overview.UntouchedNum += v.Count
	}

	res.Overview.TotalNum = res.Overview.AcceptedNum + res.Overview.FailedNum + res.Overview.UntouchedNum
	res.Easy.TotalNum = res.Easy.AcceptedNum + res.Easy.FailedNum + res.Easy.UntouchedNum
	res.Medium.TotalNum = res.Medium.AcceptedNum + res.Medium.FailedNum + res.Medium.UntouchedNum
	res.Hard.TotalNum = res.Hard.AcceptedNum + res.Hard.FailedNum + res.Hard.UntouchedNum

	return &res, nil
}

func GetUserProfile(name string) (*models.UserProfile, error) {
	ptr, err := Get(models.RequestTypeUserProfile, name)
	if err != nil {
		return nil, err
	}
	p := ptr.(*LeetCodeUserProfile)
	if p.UserProfilePublicProfile.Username == "" {
		return nil, nil
	}

	pp := p.UserProfilePublicProfile
	userProfile := &models.UserProfile{
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

func Get(rt models.RequestType, name string) (any, error) {
	if name == "" {
		return nil, errors.Errorf("miss name")
	}

	cfg, ok := models.LeetCodeRequestMap[rt]
	if !ok {
		return nil, errors.Errorf("request type not found, request type: %d", rt)
	}

	ptr := reflect.New(cfg.Response).Interface()

	var (
		uri    = cfg.URI
		method = http.MethodPost
		client = http.DefaultClient
		query  = fmt.Sprintf(cfg.Query, name)
	)

	err := Send(client, uri, method, query, ptr)
	return ptr, errors.Wrapf(err, "rt: %s, method: %s, name: %s, uri: %s, query: %s", cfg.Desc, method, name, uri, query)
}
