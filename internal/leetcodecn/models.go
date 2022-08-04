package leetcodecn

import (
	"reflect"

	"github.com/haozibi/leetcode-badge/internal/models"
)

func init() {
	models.Register(models.RequestTypeUserProfile, models.RequestConfig{
		Desc:     "user_profile",
		URI:      "https://leetcode-cn.com/graphql",
		Query:    "{\"operationName\":\"userPublicProfile\",\"variables\":{\"userSlug\":\"%s\"},\"query\":\"query userPublicProfile($userSlug: String!) {\\nuserProfilePublicProfile(userSlug: $userSlug) {\\nusername\\nhaveFollowed\\nsiteRanking\\nprofile {\\nuserSlug\\nrealName\\nuserAvatar\\nlocation\\ncontestCount\\nasciiCode\\n__typename\\n}\\n submissionProgress {\\ntotalSubmissions\\nwaSubmissions\\nacSubmissions\\nreSubmissions\\notherSubmissions\\nacTotal\\nquestionTotal\\n__typename\\n}\\n__typename\\n}\\n}\\n\"}",
		Response: reflect.TypeOf(LeetCodeUserProfile{}),
	})

	models.Register(models.RequestTypeUserQuestionProgress, models.RequestConfig{
		Desc:     "user_request_process",
		URI:      "https://leetcode-cn.com/graphql",
		Query:    `{"query": "\n    query userQuestionProgress($userSlug: String!) {\n  userProfileUserQuestionProgress(userSlug: $userSlug) {\n    numAcceptedQuestions {\n      difficulty\n      count\n    }\n    numFailedQuestions {\n      difficulty\n      count\n    }\n    numUntouchedQuestions {\n      difficulty\n      count\n    }\n  }\n}\n    ","variables": {"userSlug": "%s"}}`,
		Response: reflect.TypeOf(LeetCodeUserQuestionProgress{}),
	})
}
