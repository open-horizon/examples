VERSION=1.0.4
BUILD_IMAGE=
RUN_IMAGE=
GOARCH=
ARCH=
IMG_NAME=
DOCKER_FILE=
WORKDIR=
DOCKER_ID=vkorn

.PHONY: docker-build build-client-amd build-client-arm build-client

.ONESHELL:
docker-build:
	cd $(WORKDIR)
	docker build --no-cache -t $(IMG_NAME):$(ARCH)-$(VERSION) --build-arg BUILD_IMAGE=$(BUILD_IMAGE) --build-arg RUN_IMAGE=$(RUN_IMAGE) --build-arg GOARCH=$(GOARCH) -f $(DOCKER_FILE) .
	docker push $(IMG_NAME):$(ARCH)-$(VERSION)

create-manifest:
	docker pull $(IMG_NAME):arm32v6-$(VERSION)
	docker pull $(IMG_NAME):amd64-$(VERSION)
	docker manifest create $(IMG_NAME):$(VERSION) $(IMG_NAME):arm32v6-$(VERSION) $(IMG_NAME):amd64-$(VERSION) --amend
	docker manifest annotate $(IMG_NAME):$(VERSION) $(IMG_NAME):arm32v6-$(VERSION) --arch arm
	docker manifest push $(IMG_NAME):$(VERSION)

build-client-amd:
	$(MAKE) IMG_NAME=$(DOCKER_ID)/fft-client DOCKER_FILE=Dockerfile BUILD_IMAGE=golang:1.14.2-alpine3.11 \
		RUN_IMAGE=alpine:3.11 GOARCH=amd64 ARCH=amd64 WORKDIR=./client docker-build

build-client-arm:
	$(MAKE) IMG_NAME=$(DOCKER_ID)/fft-client DOCKER_FILE=Dockerfile BUILD_IMAGE=arm32v6/golang:1.14.2-alpine3.11 \
		RUN_IMAGE=arm32v6/alpine:3.11 GOARCH=arm ARCH=arm32v6 WORKDIR=./client docker-build

build-client: build-client-amd build-client-arm
	$(MAKE) IMG_NAME=$(DOCKER_ID)/fft-client create-manifest

build-server-amd:
	$(MAKE) IMG_NAME=$(DOCKER_ID)/fft-server DOCKER_FILE=./server/Dockerfile BUILD_IMAGE=golang:1.14.2-alpine3.11 \
		RUN_IMAGE=alpine:3.11 GOARCH=amd64 ARCH=amd64 WORKDIR=. docker-build

build-server-arm:
	$(MAKE) IMG_NAME=$(DOCKER_ID)/fft-server DOCKER_FILE=./server/Dockerfile BUILD_IMAGE=arm32v6/golang:1.14.2-alpine3.11 \
		RUN_IMAGE=arm32v6/alpine:3.11 GOARCH=arm ARCH=arm32v6 WORKDIR=. docker-build

build-server: build-server-amd build-server-arm
	$(MAKE) IMG_NAME=$(DOCKER_ID)/fft-server create-manifest