package leetcode

type Profile struct {
	UserSlug     string `json:"userSlug"`
	RealName     string `json:"realName"`
	UserAvatar   string `json:"userAvatar"`
	Location     string `json:"location"`
	ContestCount int    `json:"contestCount"`
	ASCIICode    string `json:"asciiCode"`
	Typename     string `json:"__typename"`
}

type SubmissionProgress struct {
	TotalSubmissions int    `json:"totalSubmissions"`
	WaSubmissions    int    `json:"waSubmissions"`
	AcSubmissions    int    `json:"acSubmissions"`
	ReSubmissions    int    `json:"reSubmissions"`
	OtherSubmissions int    `json:"otherSubmissions"`
	AcTotal          int    `json:"acTotal"`
	QuestionTotal    int    `json:"questionTotal"`
	Typename         string `json:"__typename"`
}

type UserProfilePublicProfile struct {
	Username           string             `json:"username"`
	HaveFollowed       interface{}        `json:"haveFollowed"`
	SiteRanking        int                `json:"siteRanking"`
	Profile            Profile            `json:"profile"`
	SubmissionProgress SubmissionProgress `json:"submissionProgress"`
	Typename           string             `json:"__typename"`
}

type LeetCodeUserProfile struct {
	UserProfilePublicProfile UserProfilePublicProfile `json:"userProfilePublicProfile"`
}

type Locations struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type Errors struct {
	Message   string      `json:"message"`
	Locations []Locations `json:"locations"`
}

type UserProfileData struct {
	MatchedUser       MatchedUser  `json:"matchedUser"`
	AllQuestionsCount []Submission `json:"allQuestionsCount"`
}

type MatchedUser struct {
	Profile     MatchedUserProfile `json:"profile"`
	SubmitStats SubmitStats        `json:"submitStats"`
}

type MatchedUserProfile struct {
	RealName   string `json:"realName"`
	UserAvatar string `json:"userAvatar"`
	Ranking    int    `json:"ranking"`
}

type SubmitStats struct {
	AcSubmissionNum    []Submission `json:"acSubmissionNum"`
	TotalSubmissionNum []Submission `json:"totalSubmissionNum"`
}

type Submission struct {
	Count       int    `json:"count"`
	Difficulty  string `json:"difficulty"`
	Submissions int    `json:"submissions"`
}
