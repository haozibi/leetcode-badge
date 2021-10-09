package heatmap

import (
	"bytes"
	_ "embed"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/nikolaydubina/calendarheatmap/charts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/haozibi/leetcode-badge/internal/statics"
)

func Do() {
	counts := make(map[string]int)
	counts["2021-10-09"] = 10
	counts["2021-10-08"] = 0
	counts["2021-10-07"] = 20

	colorscale, err := charts.NewBasicColorscaleFromCSV(bytes.NewBuffer(statics.ColorYellow()))
	if err != nil {
		log.Fatal(err)
	}

	fontFace, err := charts.LoadFontFace(statics.TTF(), opentype.FaceOptions{
		Size:    26,
		DPI:     280,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatal(err)
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
	outputfile, err := os.Create("a.svg")
	if err != nil {
		panic(err)
	}
	if err := charts.WriteHeatmap(conf, outputfile); err != nil {
		panic(err)
	}
	if err := outputfile.Close(); err != nil {
		panic(err)
	}
}
