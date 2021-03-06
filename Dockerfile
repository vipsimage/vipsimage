FROM golang:alpine AS build-base

ENV VIPS_VERSION=8.10.0
# build libvips
#COPY data/vips-${VIPS_VERSION}.tar.gz .
#RUN echo "http://mirrors.aliyun.com/alpine/v3.10/main/" > /etc/apk/repositories

RUN wget https://github.com/libvips/libvips/releases/download/v${VIPS_VERSION}/vips-${VIPS_VERSION}.tar.gz
RUN tar -zxvf vips-${VIPS_VERSION}.tar.gz
RUN apk add g++ make glib-dev expat gtk-doc libjpeg-turbo-dev libpng-dev libwebp-dev giflib-dev librsvg-dev libexif-dev lcms2-dev tiff-dev libheif-dev
RUN cd vips-${VIPS_VERSION} && \
    ./configure --without-OpenEXR --enable-debug=no --disable-static --enable-silent-rules && \
    make install-strip

WORKDIR /go/src/github.com/vipsimage

COPY go.mod .
COPY go.sum .

# cache go package
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download

# build vipsimage
FROM build-base AS pre-build
COPY . .
RUN go build -o vipsimage

# build target image
FROM alpine
MAINTAINER vipsimage@vipsimage.com
WORKDIR /app/

ENV GIN_MODE=release

RUN mkdir -p /data/logs && \
#    echo "http://mirrors.aliyun.com/alpine/v3.10/main/" > /etc/apk/repositories && \
    apk --update --no-cache add fftw glib libltdl expat libjpeg-turbo libpng libwebp giflib librsvg libgsf libexif lcms2 libheif tiff

COPY --from=pre-build /usr/local/lib/* /usr/local/lib/
COPY --from=pre-build /go/src/github.com/vipsimage/vipsimage .
COPY data/images /data/images
COPY data/vipsimage.reference.toml data/vipsimage.toml /data/

EXPOSE 8910

# release memory model
ENV GODEBUG="madvdontneed=1"
CMD ["/app/vipsimage"]