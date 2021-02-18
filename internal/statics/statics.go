package statics

import (
	"embed"
)

//go:embed svg/*
var f embed.FS

func GetLackSVG() []byte {
	return readBody("svg/lack.svg")
}

func SVGNotFound() []byte {
	return readBody("svg/notfound.svg")
}

func readBody(name string) []byte {
	body, err := f.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return body
}
