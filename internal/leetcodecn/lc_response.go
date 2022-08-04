package leetcodecn

import (
	"encoding/json"

	"github.com/pkg/errors"
)

var (
	ErrUserNotExist = errors.New("That user does not exist.")
)

// common

type Data struct {
	Data   json.RawMessage `json:"data"`
	Errors []Errors        `json:"errors"`
}

type Locations struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type Errors struct {
	Message   string      `json:"message"`
	Locations []Locations `json:"locations"`
}

// user question process

type LeetCodeUserQuestionProgress struct {
	UserProfileUserQuestionProgress LCUserProfileUserQuestionProgress `json:"userProfileUserQuestionProgress"`
}

type LCNumQuestions struct {
	Difficulty string `json:"difficulty"`
	Count      int    `json:"count"`
}

type LCUserProfileUserQuestionProgress struct {
	NumAcceptedQuestions  []LCNumQuestions `json:"numAcceptedQuestions"`
	NumFailedQuestions    []LCNumQuestions `json:"numFailedQuestions"`
	NumUntouchedQuestions []LCNumQuestions `json:"numUntouchedQuestions"`
}

// user profile

type LeetCodeUserProfile struct {
	UserProfilePublicProfile LCUserProfilePublicProfile `json:"userProfilePublicProfile"`
}

type LCUserProfilePublicProfile struct {
	Username           string               `json:"username"`
	HaveFollowed       interface{}          `json:"haveFollowed"`
	SiteRanking        int                  `json:"siteRanking"`
	Profile            LCProfile            `json:"profile"`
	SubmissionProgress LCSubmissionProgress `json:"submissionProgress"`
	Typename           string               `json:"__typename"`
}

type LCProfile struct {
	UserSlug     string `json:"userSlug"`
	RealName     string `json:"realName"`
	UserAvatar   string `json:"userAvatar"`
	Location     string `json:"location"`
	ContestCount int    `json:"contestCount"`
	ASCIICode    string `json:"asciiCode"`
	Typename     string `json:"__typename"`
}

type LCSubmissionProgress struct {
	TotalSubmissions int    `json:"totalSubmissions"`
	WaSubmissions    int    `json:"waSubmissions"`
	AcSubmissions    int    `json:"acSubmissions"`
	ReSubmissions    int    `json:"reSubmissions"`
	OtherSubmissions int    `json:"otherSubmissions"`
	AcTotal          int    `json:"acTotal"`
	QuestionTotal    int    `json:"questionTotal"`
	Typename         string `json:"__typename"`
}
