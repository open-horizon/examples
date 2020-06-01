VERSION=1.0.4
BUILD_IMAGE=
RUN_IMAGE=
GOARCH=
ARCH=
IMG_NAME=
DOCKER_FILE=
WORKDIR=
DOCKER_ID=vkorn
MQTT_USER_FILE=/tmp/hzn/vmq.users

.PHONY: docker-build build-client-amd build-client-arm build-client build build-test-client test-all-arches

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

build-test-client:
	docker build --no-cache -t fft-test-client -f test-client/Dockerfile .

build: build-server build-client

define gen_mqtt_users =
	set -e
	echo fft-server: $$(echo "server-pass\c" | openssl dgst -sha256) > $(MQTT_USER_FILE)
	echo fft-client: $$(echo "client-pass\c" | openssl dgst -sha256) >> $(MQTT_USER_FILE)
endef

.ONESHELL:
SHELL = /bin/sh
test-all-arches: build build-test-client
	$(gen_mqtt_users)
	cd server/ && hzn dev service start -S -v && \
	docker run -it --rm fft-test-client && \
		{ hzn dev service stop; \
		echo "*** Service test succeeded! ***"; } || \
		{ hzn dev service stop; \
		echo "*** Service test failed! ***"; \
		false ;}

publish:
	exit 0