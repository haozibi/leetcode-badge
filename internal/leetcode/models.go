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

type LeetCodeErrors struct {
	Errors []Errors `json:"errors"`
}

type Locations struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type Errors struct {
	Message   string      `json:"message"`
	Locations []Locations `json:"locations"`
}

// UserProfile user info
type UserProfile struct {
	UserSlug         string `json:"userSlug"`         // URL path
	RealName         string `json:"realName"`         // 显示的名字
	UserAvatar       string `json:"userAvatar"`       // 头像
	SiteRanking      int    `json:"siteRanking"`      // 排名
	TotalSubmissions int    `json:"totalSubmissions"` // 共提交数
	AcSubmissions    int    `json:"acSubmissions"`    // 提交通过数
	WaSubmissions    int    `json:"waSubmissions"`    // 答案错误数【无用】
	ReSubmissions    int    `json:"reSubmissions"`    // 运行错误【无用】
	OtherSubmissions int    `json:"otherSubmissions"` // 其他错误【无用】
	AcTotal          int    `json:"acTotal"`          // 解决题目数量
	QuestionTotal    int    `json:"questionTotal"`    // 题目总数
}

type GetUserProfileResult struct {
	Data UserProfileData `json:"data"`
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
