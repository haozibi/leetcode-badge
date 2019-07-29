ARG GO_VERSION=1.12.7

FROM golang:${GO_VERSION}-alpine3.9 AS build-env

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

RUN apk --no-cache add build-base git

WORKDIR ${GOPATH}/src/github.com/haozibi/leetcode-badge

COPY . ${GOPATH}/src/github.com/haozibi/leetcode-badge

RUN ls -alh && make build-linux

FROM alpine:3.9

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

RUN apk update && apk add tzdata \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \ 
    && echo "Asia/Shanghai" > /etc/timezone

RUN apk add --update ca-certificates && rm -rf /var/cache/apk/*

COPY --from=build-env go/src/github.com/haozibi/leetcode-badge/leetcode-badge /main

ENV LCHTTPAddr=":5050"

EXPOSE 5050

ENTRYPOINT ["/main","run"]