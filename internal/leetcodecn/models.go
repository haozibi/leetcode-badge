package leetcodecn

import (
	"net/http"
	"reflect"

	"github.com/haozibi/leetcode-badge/internal/models"
)

func init() {
	models.Register(models.RequestTypeUserProfile, models.RequestConfig{
		Desc:     "user_profile",
		URI:      "https://leetcode-cn.com/graphql",
		Method:   http.MethodPost,
		Query:    "{\"operationName\":\"userPublicProfile\",\"variables\":{\"userSlug\":\"%s\"},\"query\":\"query userPublicProfile($userSlug: String!) {\\nuserProfilePublicProfile(userSlug: $userSlug) {\\nusername\\nhaveFollowed\\nsiteRanking\\nprofile {\\nuserSlug\\nrealName\\nuserAvatar\\nlocation\\ncontestCount\\nasciiCode\\n__typename\\n}\\n submissionProgress {\\ntotalSubmissions\\nwaSubmissions\\nacSubmissions\\nreSubmissions\\notherSubmissions\\nacTotal\\nquestionTotal\\n__typename\\n}\\n__typename\\n}\\n}\\n\"}",
		Response: reflect.TypeOf(LeetCodeUserProfile{}),
	})
	models.Register(models.RequestTypeUserQuestionProgress, models.RequestConfig{
		Desc:     "user_request_process",
		URI:      "https://leetcode-cn.com/graphql",
		Method:   http.MethodPost,
		Query:    `{"query": "\n    query userQuestionProgress($userSlug: String!) {\n  userProfileUserQuestionProgress(userSlug: $userSlug) {\n    numAcceptedQuestions {\n      difficulty\n      count\n    }\n    numFailedQuestions {\n      difficulty\n      count\n    }\n    numUntouchedQuestions {\n      difficulty\n      count\n    }\n  }\n}\n    ","variables": {"userSlug": "%s"}}`,
		Response: reflect.TypeOf(LeetCodeUserQuestionProgress{}),
	})
	models.Register(models.RequestTypeUserContestRankingInfo, models.RequestConfig{
		Desc:     "user_contest_rank_info",
		URI:      "https://leetcode.cn/graphql/noj-go/",
		Method:   http.MethodPost,
		Query:    `{"query":"\n    query userContestRankingInfo($userSlug: String!) {\n  userContestRanking(userSlug: $userSlug) {\n    attendedContestsCount\n    rating\n    globalRanking\n    localRanking\n    globalTotalParticipants\n    localTotalParticipants\n    topPercentage\n  }\n  userContestRankingHistory(userSlug: $userSlug) {\n    attended\n    totalProblems\n    trendingDirection\n    finishTimeInSeconds\n    rating\n    score\n    ranking\n    contest {\n      title\n      titleCn\n      startTime\n    }\n  }\n}\n    ","variables":{"userSlug":"%s"}}`,
		Response: reflect.TypeOf(LeetCodeUserContestRankingInfo{}),
	})
}
