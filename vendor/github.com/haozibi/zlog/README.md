![logo](logo.png)

[![Build Status](https://travis-ci.org/haozibi/zlog.svg?branch=master)](https://travis-ci.org/haozibi/zlog) [![Coverage Status](https://coveralls.io/repos/github/haozibi/zlog/badge.svg?branch=master)](https://coveralls.io/github/haozibi/zlog?branch=master) [![GoDoc](https://godoc.org/github.com/haozibi/zlog?status.svg)](https://godoc.org/github.com/haozibi/zlog) [![Go Report Card](https://goreportcard.com/badge/github.com/haozibi/zlog)](https://goreportcard.com/report/github.com/haozibi/zlog) [![license](https://img.shields.io/github/license/haozibi/zlog.svg)](https://github.com/haozibi/zlog)

# zlog

Just Log Basic On [zerolog](https://github.com/rs/zerolog)

## Install

```shell
$ go get -u github.com/haozibi/zlog
```

## Demo

```go
package main

import (
	"os"

	"github.com/haozibi/zlog"
)

func init() {

	zlog.NewBasicLog(os.Stdout,WithNoColor(true), WithDebug(true))
	// zlog.NewJSONLog(os.Stdout)
}

func main() {
	zlog.ZInfo().
		Int("z", 100-1).
		Msg("just do it")

	zlog.ZDebug().
		Float64("f", 3.1415926).
		Msgf("hello %s", "zlog")
}
```
