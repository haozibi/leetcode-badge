package card

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"

	"github.com/haozibi/leetcode-badge/internal/models"
	"github.com/haozibi/leetcode-badge/internal/statics"
)

func Build(data *models.UserQuestionPrecess) ([]byte, error) {
	var (
		baseLen = 215
	)

	getLen := func(num, total int) int {
		return int(float64(baseLen) * (float64(num) / float64(total)))
	}

	info := Info{
		BaseLen:   baseLen,
		EasyLen:   getLen(data.Easy.AcceptedNum, data.Easy.TotalNum),
		MediumLen: getLen(data.Medium.AcceptedNum, data.Medium.TotalNum),
		HardLen:   getLen(data.Hard.AcceptedNum, data.Hard.TotalNum),

		EasyNum:     data.Easy.AcceptedNum,
		EasyTotal:   data.Easy.TotalNum,
		MediumNum:   data.Medium.AcceptedNum,
		MediumTotal: data.Medium.TotalNum,
		HardNum:     data.Hard.AcceptedNum,
		HardTotal:   data.Hard.TotalNum,

		AcceptNum: data.Overview.AcceptedNum,
	}

	t, err := template.New("foo").Parse(string(statics.TemplateQuestionProcess()))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	buf := bytes.NewBuffer(nil)
	err = t.Execute(buf, info)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return buf.Bytes(), nil
}

type Info struct {
	// 长度
	BaseLen   int
	EasyLen   int
	MediumLen int
	HardLen   int

	// 数量
	EasyNum     int
	EasyTotal   int
	MediumNum   int
	MediumTotal int
	HardNum     int
	HardTotal   int

	AcceptNum int
}
