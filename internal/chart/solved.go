package chart

import (
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"
	chart "github.com/wcharczuk/go-chart"
)

type SolvedHistory struct {
	UserName string
	Date     time.Time
	Num      int
}

// ShowSolvedHistory show solved history
func ShowSolvedHistory(w io.Writer, solved [][]SolvedHistory, names ...string) error {

	if len(solved) == 0 || len(solved) != len(names) {
		return errors.New("miss params")
	}

	series := make([]chart.Series, 0, len(solved))
	ticksTime := make([]time.Time, 0, 7)
	ticksNum := make([]int, 0, 7)

	for i := 0; i < len(solved); i++ {
		xValues := make([]float64, len(solved[i]))
		yValues := make([]float64, len(solved[i]))
		zValues := make([]chart.Value2, len(solved[i]))
		for j := 0; j < len(solved[i]); j++ {
			xValues[j] = float64(solved[i][j].Date.Unix())
			yValues[j] = float64(solved[i][j].Num)
			zValues[j] = chart.Value2{
				XValue: xValues[j],
				YValue: yValues[j],
				Label:  fmt.Sprintf("%0.0f", yValues[j]),
			}
			ticksTime = append(ticksTime, solved[i][j].Date)
			ticksNum = append(ticksNum, solved[i][j].Num)
		}

		// 仅标注最后一个数据
		zValues = zValues[len(zValues)-1:]

		series = append(series, chart.ContinuousSeries{
			Name:    names[i],
			XValues: xValues,
			YValues: yValues,
		})
		series = append(series, chart.AnnotationSeries{
			Annotations: zValues,
		})
	}

	xTicks := getXTicks(ticksTime)
	yTicks := getYTicks(ticksNum)

	return show(w, "Solved History", "Date", "Solved Num", false, series, xTicks, yTicks)
}
