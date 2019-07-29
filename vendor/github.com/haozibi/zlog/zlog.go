package zlog

import (
	"io"
	"time"

	"github.com/rs/zerolog"
)

var zlog zerolog.Logger

var (
	// TimeFieldFormat time format
	TimeFieldFormat = time.RFC3339

	// TimeFormatUnixNano time format
	TimeFormatUnixNano = "2006-01-02 15:04:05.999999999"

	// NoColor if set color
	NoColor = false
)

func newWriter(w io.Writer, opts ...LogOpt) *zerolog.ConsoleWriter {

	oo := &option{
		timeFormat: TimeFormatUnixNano,
		color:      false,
	}

	for _, o := range opts {
		o(oo)
	}

	// ConsoleWriter parses the JSON input and writes it in an (optionally) colorized, human-friendly format to Out.
	return &zerolog.ConsoleWriter{
		Out:        w,
		TimeFormat: oo.timeFormat,
		NoColor:    oo.color,
	}
}

// NewBasicLog new basic log
func NewBasicLog(w io.Writer, opts ...LogOpt) {
	output := newWriter(w, opts...)

	zlog = newLog(output, opts...).Logger()
}

// NewJSONLog new log by json format
func NewJSONLog(w io.Writer, opts ...LogOpt) {
	zlog = newLog(w, opts...).Logger()
}

type option struct {
	debug      bool
	deep       int
	timeFormat string
	color      bool
}

// LogOpt log option
type LogOpt func(o *option)

// WithDebug set if debug,debug output line num
func WithDebug(debug bool) LogOpt {
	return func(o *option) {
		o.debug = debug
	}
}

// WithDeep set line deep,default eq 2
func WithDeep(n int) LogOpt {
	return func(o *option) {
		o.deep = n
	}
}

// WithNoColor set if has color
func WithNoColor(nocolor bool) LogOpt {
	return func(o *option) {
		o.color = nocolor
	}
}

// WithTimeFormat set time format when basic format
func WithTimeFormat(format string) LogOpt {
	return func(o *option) {
		o.timeFormat = format
	}
}

func newLog(w io.Writer, opts ...LogOpt) zerolog.Context {
	zerolog.TimeFieldFormat = TimeFormatUnixNano

	oo := &option{
		debug: false,
		deep:  2,
	}

	for _, o := range opts {
		o(oo)
	}

	z := zerolog.New(w).With().Timestamp()

	if oo.debug {
		z = z.CallerWithSkipFrameCount(2)
	}

	return z
}

// ZDebug debug log
func ZDebug() *zerolog.Event {
	return zlog.Debug()
}

// ZInfo info log
func ZInfo() *zerolog.Event {
	return zlog.Info()
}

// ZWarn warn log
func ZWarn() *zerolog.Event {
	return zlog.Warn()
}

// ZError error log
func ZError() *zerolog.Event {
	return zlog.Error()
}

// ZFatal fatal log
func ZFatal() *zerolog.Event {
	return zlog.Fatal()
}

// Debugf debug format
func Debugf(format string, v ...interface{}) {
	ZDebug().Msgf(format, v...)
}

// // Infof info format
// func Infof(format string, v ...interface{}) {
// 	ZInfo().Msgf(format, v...)
// }

// // Warnf warn format
// func Warnf(format string, v ...interface{}) {
// 	ZWarn().Msgf(format, v...)
// }

// // Errorf error format
// func Errorf(format string, v ...interface{}) {
// 	ZError().Msgf(format, v...)
// }

// // Fatalf fatalf
// func Fatalf(format string, v ...interface{}) {
// 	ZFatal().Msgf(format, v...)
// }
