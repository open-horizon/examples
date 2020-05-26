#!/usr/bin/env bash

set -e 

img_build=(golang:1.13.6-alpine3.11 arm32v6/golang:1.13.6-alpine3.11)
img_run=(alpine:3.11 arm32v6/alpine:3.11)

VERSION=1.0.1


BUILD_IMAGE=${img_build[0]}
RUN_IMAGE=${img_run[0]}
GOARCH=amd64
ARCH=amd64

IMG_NAME=vkorn/volantmq

docker build --no-cache -t $IMG_NAME:$ARCH-${VERSION} --build-arg BUILD_IMAGE=${BUILD_IMAGE} \
	--build-arg RUN_IMAGE=${RUN_IMAGE} --build-arg GOARCH=${GOARCH} .

docker push $IMG_NAME:$ARCH-${VERSION}


BUILD_IMAGE=${img_build[1]}
RUN_IMAGE=${img_run[1]}
GOARCH=arm
ARCH=arm32v6

docker build --no-cache -t $IMG_NAME:$ARCH-${VERSION} --build-arg BUILD_IMAGE=${BUILD_IMAGE} \
	--build-arg RUN_IMAGE=${RUN_IMAGE} --build-arg GOARCH=${GOARCH} .

docker push $IMG_NAME:$ARCH-${VERSION}

docker manifest create $IMG_NAME:${VERSION} $IMG_NAME:arm32v6-${VERSION} $IMG_NAME:amd64-${VERSION} --amend
docker manifest annotate $IMG_NAME:${VERSION} $IMG_NAME:arm32v6-${VERSION} --arch arm
docker manifest push $IMG_NAME:${VERSION} 