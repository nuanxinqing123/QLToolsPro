FROM  alpine:3.15 AS builder

ARG TARGETARCH

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /usr/src/QLToolsPro

# 安装项目必要环境
RUN \
  #sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
  apk add --no-cache --update go go-bindata g++ ca-certificates tzdata

COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download

COPY . .

# 打包项目文件
RUN \
  #go-bindata -o=bindata/bindata.go -pkg=bindata ./assets/... && \
  go build -ldflags '-linkmode external -s -w -extldflags "-static"' -o QLToolsPro-linux-$TARGETARCH


# FROM alpine:3.15
FROM ubuntu:22.10

MAINTAINER QLToolsPro "nuanxinqing@gmail.com"

ARG TARGETARCH
ENV TARGET_ARCH=$TARGETARCH

WORKDIR /QLToolsPro

COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/src/QLToolsPro/QLToolsPro-linux-$TARGETARCH /usr/src/QLToolsPro/docker-entrypoint.sh /usr/src/QLToolsPro/sample ./

EXPOSE 6600

ENTRYPOINT ["sh", "docker-entrypoint.sh"]