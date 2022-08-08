package card

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/pkg/errors"

	"github.com/haozibi/leetcode-badge/internal/models"
	"github.com/haozibi/leetcode-badge/internal/statics"
)

func QuestionProcess(name string, data *models.UserQuestionPrecess) ([]byte, error) {
	var (
		baseLen = 215
	)

	getLen := func(num, total int) int {
		return int(float64(baseLen) * (float64(num) / float64(total)))
	}

	info := QuestionProcessInfo{
		Name: name,

		BaseLen:   baseLen,
		EasyLen:   getLen(data.Easy.AcceptedNum, data.Easy.TotalNum),
		MediumLen: getLen(data.Medium.AcceptedNum, data.Medium.TotalNum),
		HardLen:   getLen(data.Hard.AcceptedNum, data.Hard.TotalNum),

		EasyNum:     fmt.Sprintf("% 4d", data.Easy.AcceptedNum),
		EasyTotal:   data.Easy.TotalNum,
		MediumNum:   fmt.Sprintf("% 4d", data.Medium.AcceptedNum),
		MediumTotal: data.Medium.TotalNum,
		HardNum:     fmt.Sprintf("% 4d", data.Hard.AcceptedNum),
		HardTotal:   data.Hard.TotalNum,

		AcceptNum: data.Overview.AcceptedNum,
	}

	return build(statics.TemplateQuestionProcess(), info)
}

func ContestRanking(name string, data *models.UserContestRankingInfo) ([]byte, error) {
	info := ContestRankingInfo{
		Name:          name,
		Rating:        fmt.Sprintf("%d", int(data.Rating)),
		LocalRanking:  fmt.Sprintf("% 6d", data.LocalRanking),
		GlobalRanking: fmt.Sprintf("% 6d", data.GlobalRanking),

		LocalTotal:  fmt.Sprintf("/%d", data.LocalTotalParticipants),
		GlobalTotal: fmt.Sprintf("/%d", data.GlobalTotalParticipants),

		Top: fmt.Sprintf("%0.2f", 100.0-data.TopPercentage),
	}
	return build(statics.TemplateContestRanking(), info)
}

func build(temp []byte, data interface{}) ([]byte, error) {
	t, err := template.New("foo").Parse(string(temp))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	buf := bytes.NewBuffer(nil)
	err = t.Execute(buf, data)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return buf.Bytes(), nil
}

type ContestRankingInfo struct {
	Name          string
	Rating        string
	LocalRanking  string
	GlobalRanking string

	LocalTotal  string
	GlobalTotal string

	Top string
}

type QuestionProcessInfo struct {
	Name string
	// 长度
	BaseLen   int
	EasyLen   int
	MediumLen int
	HardLen   int

	// 数量
	EasyNum     string
	EasyTotal   int
	MediumNum   string
	MediumTotal int
	HardNum     string
	HardTotal   int

	AcceptNum int
}
