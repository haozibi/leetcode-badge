package heatmap

import (
	"bytes"
	_ "embed"
	"image/color"
	"time"

	"github.com/nikolaydubina/calendarheatmap/charts"
	"github.com/pkg/errors"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/haozibi/leetcode-badge/internal/statics"
)

// LastYear 过去一年
func LastYear(data map[int64]int) ([]byte, error) {
	now := time.Now()
	return Build(now.AddDate(-1, 0, 0).Unix(), now.Unix(), data)
}

// CurrYear 当前年
func CurrYear(data map[int64]int) ([]byte, error) {
	now := time.Now()
	s := time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
	e := time.Date(now.Year(), time.December, 31, 23, 59, 59, 0, time.UTC)

	return Build(s.Unix(), e.Unix(), data)
}

func Build(start, end int64, data map[int64]int) ([]byte, error) {

	input := make(map[Date]int, len(data))
	for k, v := range data {
		if k >= start && k <= end {
			t := time.Unix(k, 0)
			input[Date{
				Year:  t.Year(),
				Month: t.Month(),
				Day:   t.Day(),
			}] = v % 5 // only 5 color
		}
	}

	st, et := time.Unix(start, 0), time.Unix(end, 0)
	h := New(nil)
	buf := h.Generate(
		Date{
			Year:  st.Year(),
			Month: st.Month(),
			Day:   st.Day(),
		},
		Date{
			Year:  et.Year(),
			Month: et.Month(),
			Day:   et.Day(),
		},
		input,
	)

	return buf.Bytes(), nil
}

func Do(counts map[string]int) ([]byte, error) {

	colorscale, err := charts.NewBasicColorscaleFromCSV(bytes.NewBuffer(statics.ColorYellow()))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	fontFace, err := charts.LoadFontFace(statics.TTF(), opentype.FaceOptions{
		Size:    26,
		DPI:     280,
		Hinting: font.HintingNone,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	conf := charts.HeatmapConfig{
		Counts:              counts,
		ColorScale:          colorscale,
		DrawMonthSeparator:  true,
		DrawLabels:          true,
		Margin:              30,
		BoxSize:             150,
		MonthSeparatorWidth: 5,
		MonthLabelYOffset:   50,
		TextWidthLeft:       300,
		TextHeightTop:       200,
		TextColor:           color.RGBA{R: 100, G: 100, B: 100, A: 255},
		BorderColor:         color.RGBA{R: 200, G: 200, B: 200, A: 255},
		Locale:              "en_US",
		Format:              "svg",
		FontFace:            fontFace,
		ShowWeekdays: map[time.Weekday]bool{
			time.Monday:    true,
			time.Wednesday: true,
			time.Friday:    true,
		},
	}

	b := bytes.NewBuffer(make([]byte, 0, 100))
	if err = charts.WriteHeatmap(conf, b); err != nil {
		return nil, errors.WithStack(err)
	}

	return b.Bytes(), nil
}
