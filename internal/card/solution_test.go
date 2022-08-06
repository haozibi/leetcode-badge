package card

import (
	"fmt"
	"testing"

	"github.com/haozibi/leetcode-badge/internal/models"
)

func TestBuild(t *testing.T) {

	data := &models.UserQuestionPrecess{
		Overview: models.UserQuestionProcessStat{
			AcceptedNum: 10,
		},
	}

	_, err := Build(data)
	fmt.Println(err)
}
