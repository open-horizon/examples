# Transform the machine arch into some standard values: "arm", "arm64", or "amd64". Note: ppc64le already returns ppc64le, so we do not need to change that.
SYSTEM_ARCH := $(shell uname -m | sed -e 's/aarch64.*/arm64/' -e 's/x86_64.*/amd64/' -e 's/armv.*/arm/')

# To build for an arch different from the current system, set this env var to 1 of the values in the comment above
ARCH ?= $(SYSTEM_ARCH)

# Container image version and microservice name
PWSMS_VERSION ?= 1.1.0
MS_NAME=pwsms
IMAGE=example_ms_$(MS_NAME)

# Construct our Docker registry path and other architecture-dependent values
REGISTRY ?= openhorizon
REG_PATH := $(REGISTRY)/$(ARCH)
IMAGE_PATH = $(REG_PATH)_$(IMAGE):$(PWSMS_VERSION)


default: build

all: build run

build:
	docker build -t $(IMAGE_PATH) -f ./Dockerfile.$(ARCH) .

build-squash:
	docker build --squash -t $(IMAGE_PATH) -f ./Dockerfile.$(ARCH) .

dev:
	-docker network create example
	-docker stop $(IMAGE) && docker rm $(IMAGE)
	docker run -it --name $(IMAGE) --publish=8357:8357 --privileged --net=example --net-alias=$(MS_NAME) --volume `pwd`:/outside $(IMAGE_PATH) /bin/sh

run:
	-docker network create example
	-docker rm -f $(IMAGE_PATH)
	docker run -d --name $(IMAGE) --publish=8357:8357 --net=example --net-alias=$(MS_NAME) --volume `pwd`:/outside $(IMAGE_PATH)

check:
	curl -s localhost:8357/v1/weather | jq

docker-publish:
	docker push $(IMAGE_PATH)

clean:
	-docker rm -f $(IMAGE_PATH)
	-docker stop $(IMAGE) && docker rm $(IMAGE)
	-docker rmi $(IMAGE)
	-docker network rm example

.PHONY: default all build build-squash dev run check publish exchange-publish clean
