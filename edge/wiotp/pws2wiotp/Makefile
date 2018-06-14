# Transform the machine arch into some standard values: "arm", "arm64", or "x86"
SYSTEM_ARCH := $(shell uname -m | sed -e 's/aarch64.*/arm64/' -e 's/x86_64.*/amd64/' -e 's/armv.*/arm/')

# To build for an arch different from the current system, set this env var to 1 of the values in the comment above
export ARCH ?= $(SYSTEM_ARCH)

# Container image version and workload name
PWS2WIOTP_VERSION ?= 1.1.0
WL_NAME=pws2wiotp
IMAGE=$(WL_NAME)

# Construct our Docker registry path and other architecture-dependent values
DOCKER_HUB_ID ?= openhorizon
REG_PATH := $(DOCKER_HUB_ID)/$(ARCH)
IMAGE_PATH = $(REG_PATH)_$(IMAGE):$(PWS2WIOTP_VERSION)


default: build

all: build run

build:
	docker build -t $(IMAGE_PATH) -f ./Dockerfile.$(ARCH) .

build-squash:
	docker build --squash -t $(IMAGE_PATH) -f ./Dockerfile.$(ARCH) .

dev:
	-docker network create example
	-docker stop $(IMAGE) && docker rm $(IMAGE)
	-docker rm -f $(IMAGE_PATH)
	docker run -it -e WIOTP_DOMAIN="${WIOTP_DOMAIN}" -e HZN_ORGANIZATION="$(HZN_ORG_ID)" -e HZN_DEVICE_ID="$(HZN_DEVICE_ID)" -e WIOTP_DEVICE_AUTH_TOKEN=$(WIOTP_GW_TOKEN) --name $(IMAGE) --net=example --net-alias=$(WL_NAME) --volume `pwd`:/outside $(IMAGE_PATH) /bin/sh

run:
	-docker network create example
	-docker rm -f $(IMAGE_PATH)
	docker run -d -e WIOTP_DOMAIN="${WIOTP_DOMAIN}" -e HZN_ORGANIZATION="$(HZN_ORG_ID)" -e HZN_DEVICE_ID="$(HZN_DEVICE_ID)" -e WIOTP_DEVICE_AUTH_TOKEN=$(WIOTP_GW_TOKEN) --name $(IMAGE) --net=example --net-alias=$(WL_NAME) $(IMAGE_PATH)

check:
	curl -s pwsms:8357/v1/weather | jq

# To publish you must have write access to the docker hub openhorizon user
docker-publish:
	docker push $(IMAGE_PATH)

clean:
	-docker rm -f $(IMAGE_PATH)
	-docker stop $(IMAGE) && docker rm $(IMAGE)
	-docker rmi $(IMAGE)
	-docker network rm example

.PHONY: default all build dev run check docker-publish exchange-publish clean
