FROM golang:alpine AS build-base
WORKDIR /go/src/github.com/vipsimage

COPY go.mod .
COPY go.sum .

RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod download

FROM build-base AS pre-build
COPY . .
RUN go build -o vipsimage

FROM alpine
MAINTAINER vipsimage@vipsimage.com
WORKDIR /app/

ENV GIN_MODE=release

RUN mkdir -p /data/{logs,images} && \
    echo "http://mirrors.aliyun.com/alpine/v3.4/main/" > /etc/apk/repositories && \
    apk --update --no-cache add libvips

COPY --from=pre-build /go/src/github.com/vipsimage/vipsimage .
COPY data/vipsimage.reference.toml data/vipsimage.toml /data/

EXPOSE 8910

CMD ["/app/vipsimage"]