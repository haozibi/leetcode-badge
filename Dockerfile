FROM haozibi/upx AS build-upx

FROM golang:1.15.2-alpine3.12 AS build-env

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories

RUN apk --no-cache add build-base git

# build
ARG BIN_NAME=leetcode-badge
WORKDIR /${BIN_NAME}
ADD go.mod .
ADD go.sum .
RUN export GOPROXY=https://goproxy.cn go mod download
ADD . .
RUN make build-linux

# upx
WORKDIR /data
COPY --from=build-upx /bin/upx /bin/upx
RUN cp /${BIN_NAME}/${BIN_NAME} /data/main
RUN upx -k --best --ultra-brute /data/main

FROM alpine3.12

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories

RUN apk update && apk add tzdata \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \ 
    && echo "Asia/Shanghai" > /etc/timezone

RUN apk add --update ca-certificates && rm -rf /var/cache/apk/*

COPY --from=build-env /data/main /home/main

EXPOSE 5050

ENTRYPOINT ["/home/main","run"]