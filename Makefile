APP?=leetcode-badge

# SHELL := /bin/bash # Use bash syntax
GOOS?=linux
GOARCH?=amd64
BLDDIR=bin

VERSION?=$(shell git describe --tags --always)
COMMIT_HASH?=$(shell git rev-parse --short HEAD 2>/dev/null)
NOW?=$(shell date -u '+%Y-%m-%d %I:%M:%S %Z')
PROJECT?=github.com/haozibi/${APP}

CONTAINER_IMAGE?=registry.cn-beijing.aliyuncs.com/github-public/${APP}
K8S_NAMESPACE=default

LDFLAGS += -X "${PROJECT}/app.BuildTime=${NOW}"
LDFLAGS += -X "${PROJECT}/app.BuildVersion=${VERSION}"
LDFLAGS += -X "${PROJECT}/app.BuildAppName=${APP}"
LDFLAGS += -X "${PROJECT}/app.CommitHash=${COMMIT_HASH}"
BUILD_TAGS = ""
BUILD_FLAGS = "-v"
# PROTO_LOCATION = "internal/protocol_pb"

.PHONY: build build-local build-linux clean govet docker docker-push

default: build

build: clean govet
	CGO_ENABLED=1 GOOS= GOARCH= go build ${BUILD_FLAGS} -ldflags '${LDFLAGS}' -tags '${BUILD_TAGS}' -o ${BLDDIR}/${APP}


build-linux: clean govet
	CGO_ENABLED=1 GOOS=${GOOS} GOARCH=${GOARCH} go build ${BUILD_FLAGS} -ldflags '${LDFLAGS}' -tags '${BUILD_TAGS}' -o ${BLDDIR}/${APP}


.PHONY: docker
docker:
	docker build -t ${CONTAINER_IMAGE}:${VERSION} -f ./Dockerfile .

.PHONY: docker-push
docker-push: docker
	docker push ${CONTAINER_IMAGE}:${VERSION}

govet: 
	@ go vet . && go fmt ./... && \
	(if [[ "$(gofmt -d $(find . -type f -name '*.go' -not -path "./vendor/*" -not -path "./tests/*" -not -path "./assets/*"))" == "" ]]; then echo "Good format"; else echo "Bad format"; exit 33; fi);

clean: 
	@ rm -fr ${BLDDIR} ${APP} main static/*.go

