package statics

import (
	"embed"
)

//go:embed *
var f embed.FS

func GetLackSVG() []byte {
	return readBody("svg/lack.svg")
}

func SVGNotFound() []byte {
	return readBody("svg/notfound.svg")
}

func TemplateQuestionProcess() []byte {
	return readBody("template/question_process.svg")
}

func TemplateContestRanking() []byte {
	return readBody("template/contest_ranking.svg")
}

func TTF() []byte {
	return readBody("charts/Sunflower-Medium.ttf")
}

func ColorGreenBlue() []byte {
	return color("green-blue")
}

func ColorYellow() []byte {
	return color("yellow")
}

func ColorPurpleBlue() []byte {
	return color("purple-blue")
}

func color(color string) []byte {
	switch color {
	case "green-blue":
		return readBody("charts/green-blue-9.csv")
	case "yellow":
		return readBody("charts/yellow-green-9.csv")
	case "purple-blue":
		return readBody("charts/purple-blue-9.csv")
	}

	return nil
}

func ListI18n() []I18nInfo {
	fs, err := f.ReadDir("i18n")
	if err != nil {
		panic(err)
	}

	list := make([]I18nInfo, 0, len(fs))
	for _, v := range fs {
		if !v.IsDir() {
			body := readBody("i18n/" + v.Name())
			list = append(list, I18nInfo{
				Name: v.Name(),
				Body: body,
			})
		}
	}

	return list
}

type I18nInfo struct {
	Name string
	Body []byte
}

func readBody(name string) []byte {
	body, err := f.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return body
}
