package chart

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/pkg/errors"
	chart "github.com/wcharczuk/go-chart"
)

const (
	DefaultAD = "Power By https://github.com/haozibi/leetcode-badge"
)

type RankHistory struct {
	UserName string
	Date     time.Time
	Rank     int
}

// ShowRankHistory show ranking history
// 使用 ContinuousSeries 底层还是 float64，如果没有设置 Tick 则程序自动设置
// 自动设置是根据最大值和最小值然后取一定的间隔，导致 x 轴显示并不是所需要的
// https://github.com/wcharczuk/go-chart/blob/45fad0cfb8e64e6314abeeaf381102aad6d9279c/tick.go#L47
// chart 至少需要两个数据
func ShowRankHistory(w io.Writer, ranks [][]RankHistory, names ...string) error {

	if len(ranks) == 0 || len(ranks) != len(names) {
		return errors.New("miss params")
	}

	series := make([]chart.Series, 0, len(ranks))
	ticksTime := make([]time.Time, 0, 7)
	ticksNum := make([]int, 0, 7)

	for i := 0; i < len(ranks); i++ {
		xValues := make([]float64, len(ranks[i]))
		yValues := make([]float64, len(ranks[i]))
		zValues := make([]chart.Value2, len(ranks[i]))
		for j := 0; j < len(ranks[i]); j++ {
			xValues[j] = float64(ranks[i][j].Date.Unix())
			yValues[j] = float64(ranks[i][j].Rank)
			zValues[j] = chart.Value2{
				XValue: xValues[j],
				YValue: yValues[j],
				Label:  fmt.Sprintf("%0.0f", yValues[j]),
			}
			ticksTime = append(ticksTime, ranks[i][j].Date)
			ticksNum = append(ticksNum, ranks[i][j].Rank)
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

	return show(w, "Ranking History", "Date", "Ranking", true, series, xTicks, yTicks)
}

func show(w io.Writer, title, xName, yName string, isYDescending bool, series []chart.Series, xTicks []chart.Tick, yTicks []chart.Tick) error {

	graph := chart.Chart{
		Background: chart.Style{
			Padding: chart.Box{
				Top:  20,
				Left: 20,
			},
		},
		Title: title,
		XAxis: chart.XAxis{
			Name:  DefaultAD,
			Ticks: xTicks,
		},
		YAxis: chart.YAxis{
			Name: yName,
			Range: &chart.ContinuousRange{
				Descending: isYDescending,
			},
			Ticks: yTicks,
		},
		Series: series,
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	return errors.WithStack(graph.Render(chart.SVG, w))
}

func getXTicks(ti []time.Time) []chart.Tick {
	t := make([]chart.Tick, len(ti))

	// 当数据很多时，x 轴间隔输出，保证起始和末尾有输出
	interval := 1
	l := len(ti)
	if l > 7 {
		interval = l / 5

	}

	for i := 0; i < l; i++ {
		t[i] = chart.Tick{
			Value: float64(ti[i].Unix()),
			// Label: ti[i].Format("2006-01-02"),
		}
		if i%interval == 0 || i+1 == len(ti) {
			t[i].Label = ti[i].Format("2006-01-02")
		}
	}

	return t
}

func getYTicks(nums []int) []chart.Tick {

	t := make([]chart.Tick, 0, len(nums))
	m := make(map[int]bool, len(nums))
	for i := 0; i < len(nums); i++ {
		if m[nums[i]] {
			continue
		}
		t = append(t, chart.Tick{
			Value: float64(nums[i]),
			Label: strconv.Itoa(nums[i]),
		})
		m[nums[i]] = true
	}

	if len(m) == 1 {
		n := nums[0]
		t = append(t, chart.Tick{
			Value: 0,
			Label: "0",
		})
		t = append(t, chart.Tick{
			Value: float64(n * 2),
			Label: strconv.Itoa(n * 2),
		})
	}

	return t
}
