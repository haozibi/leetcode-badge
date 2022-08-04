package models

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

type UserQuestionProcessStat struct {
	AcceptedNum  int
	FailedNum    int
	UntouchedNum int
	TotalNum     int
}

type UserQuestionPrecess struct {
	Overview UserQuestionProcessStat
	Easy     UserQuestionProcessStat
	Medium   UserQuestionProcessStat
	Hard     UserQuestionProcessStat
}
