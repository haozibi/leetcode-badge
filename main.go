package main

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/haozibi/leetcode-badge/cmd"
)

const (
	timeFormatUnixNano = "2006-01-02 15:04:05.999999999"
)

func New(w io.Writer) {
	zerolog.TimeFieldFormat = timeFormatUnixNano
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = zerolog.New(&zerolog.ConsoleWriter{
		Out:        w,
		TimeFormat: timeFormatUnixNano,
		NoColor:    false,
	}).With().Stack().CallerWithSkipFrameCount(2).Timestamp().Logger()
}

func init() {
	New(os.Stdout)
}

func main() {

	cmd.Execute()
}
