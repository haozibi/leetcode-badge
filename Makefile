APP?=leetcode-badge

# SHELL := /bin/bash # Use bash syntax
GOOS?=linux
GOARCH?=amd64

VERSION?=$(shell git describe --tags --always)
NOW?=$(shell date -u '+%Y/%m/%d/%I:%M:%S%Z')
PROJECT?=github.com/haozibi/${APP}

LDFLAGS += -X "${PROJECT}/app.BuildTime=${NOW}"
LDFLAGS += -X "${PROJECT}/app.BuildVersion=${VERSION}"
LDFLAGS += -X "${PROJECT}/app.BuildAppName=${APP}"
BUILD_TAGS = ""
BUILD_FLAGS = "-v"
# PROTO_LOCATION = "internal/protocol_pb"

.PHONY: build build-local build-linux clean govet bindata docker-image

default: build

build: clean bindata govet
	CGO_ENABLED=0 GOOS= GOARCH= go build ${BUILD_FLAGS} -ldflags '${LDFLAGS}' -tags '${BUILD_TAGS}' -o ${APP} 


build-linux: clean bindata govet
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build ${BUILD_FLAGS} -ldflags '${LDFLAGS}' -tags '${BUILD_TAGS}' -o ${APP}

bindata: 
	go get github.com/jteeuwen/go-bindata/...
	go-bindata -nomemcopy -pkg=static \
		-debug=false \
		-o=static/static.go \
		static/...

govet: 
	@ go vet . && go fmt ./... && \
	(if [[ "$(gofmt -d $(find . -type f -name '*.go' -not -path "./vendor/*" -not -path "./tests/*" -not -path "./assets/*"))" == "" ]]; then echo "Good format"; else echo "Bad format"; exit 33; fi);

clean: 
	@ rm -fr ${APP} main static/*.go

