package leetcode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

func getCNUserProfile(userName string) (*UserProfile, error) {

	url := "https://leetcode-cn.com/graphql"

	var genQueryJSON = func(userName string) io.Reader {

		s := fmt.Sprintf("{\"operationName\":\"userPublicProfile\",\"variables\":{\"userSlug\":\"%s\"},\"query\":\"query userPublicProfile($userSlug: String!) {\\nuserProfilePublicProfile(userSlug: $userSlug) {\\nusername\\nhaveFollowed\\nsiteRanking\\nprofile {\\nuserSlug\\nrealName\\nuserAvatar\\nlocation\\ncontestCount\\nasciiCode\\n__typename\\n}\\n submissionProgress {\\ntotalSubmissions\\nwaSubmissions\\nacSubmissions\\nreSubmissions\\notherSubmissions\\nacTotal\\nquestionTotal\\n__typename\\n}\\n__typename\\n}\\n}\\n\"}", userName)

		return strings.NewReader(s)
	}

	req, err := http.NewRequest("POST", url, genQueryJSON(userName))
	if err != nil {
		return nil, err
	}

	req.Header.Add("origin", "https://leetcode-cn.com")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("referer", "https://leetcode-cn.com")

	body, _, err := request.SendRequest(req)
	if err != nil {
		return nil, err
	}

	var p LeetCodeUserProfile

	err = json.Unmarshal(body, &p)
	if err != nil {
		return nil, errors.Wrap(err, "json parse")
	}

	if p.Data.UserProfilePublicProfile.Username == "" {
		return nil, nil
	}

	pp := p.Data.UserProfilePublicProfile

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

	uri := "https://leetcode.com/" + userName

	req, _ := http.NewRequest("GET", uri, nil)

	req.Header.Add("origin", "https://leetcode.com")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("referer", "https://leetcode.com")

	body, resp, err := request.SendRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	dom, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "goquery new")
	}

	var p UserProfile
	dom.Find(".list-group-item").Each(func(i int, selection *goquery.Selection) {
		if err2 := getData(selection, &p); err2 != nil {
			err = err2
		}
	})

	if err != nil {
		return nil, errors.Wrap(err, "get user profile")
	}

	err = getRanking(string(body), &p)
	if err != nil {
		return nil, errors.Wrap(err, "get user ranking")
	}

	p.RealName = cleanText(dom.Find(".realname").First().Text())
	p.UserSlug = cleanText(dom.Find(".username").First().Text())

	return &p, nil
}

func getData(selection *goquery.Selection, p *UserProfile) (err error) {

	var parse = func(s string) (int, int, error) {

		t := strings.Split(s, "/")
		if len(t) != 2 {
			return 0, 0, errors.New("parse / error")
		}

		n1, err := strconv.Atoi(t[0])
		if err != nil {
			return 0, 0, errors.Wrap(err, "strconv")
		}

		n2, err := strconv.Atoi(t[1])
		if err != nil {
			return 0, 0, errors.Wrap(err, "strconv")
		}
		return n1, n2, nil
	}

	if strings.Contains(selection.Text(), "Solved Question") {
		solved := cleanText(selection.Find("span").Text())
		p.AcTotal, p.QuestionTotal, err = parse(solved)
		if err != nil {
			return err
		}
	}

	if strings.Contains(selection.Text(), "Accepted Submission") {
		accepted := cleanText(selection.Find("span").Text())
		p.AcSubmissions, p.TotalSubmissions, err = parse(accepted)
		if err != nil {
			return err
		}
	}

	return nil
}

func getRanking(body string, p *UserProfile) (err error) {
	x := regexp.MustCompile(`(?m)pc\.init\(([\s\S]*?){`).FindStringSubmatch(body)
	if len(x) != 2 {
		return errors.New("get ranking length invalid")
	}

	l := strings.Replace(x[1], "\n", "", -1)
	l = strings.Replace(l, " ", "", -1)
	l = strings.Replace(l, "'", "", -1)
	ls := strings.Split(l, ",")
	if len(ls) < 7 {
		return errors.New("get ranking length invalid")
	}

	p.SiteRanking, err = strconv.Atoi(ls[6])
	if err != nil {
		return errors.Wrap(err, "strconv")
	}
	return nil
}

func cleanText(s string) string {
	s = strings.Replace(s, " ", "", -1)
	s = strings.Replace(s, "\t", "", -1)
	s = strings.Replace(s, "\n", "", -1)
	return s
}
