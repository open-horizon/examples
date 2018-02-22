# Transform the machine arch into some standard values: "arm", "arm64", or "amd64"
SYSTEM_ARCH := $(shell uname -m | sed -e 's/aarch64.*/arm64/' -e 's/x86_64.*/amd64/' -e 's/armv.*/arm/')

# To build for an arch different from the current system, set this env var to 1 of the values in the comment above
ARCH ?= $(SYSTEM_ARCH)

# These variables can be overridden from the environment
VERSION ?= 1.1.8
DOCKER_NAME ?= example_wl_$(ARCH)_cpu2wiotp
DOCKER_HUB_ID ?= openhorizon

default: build

all: build run

build:
	docker build -t $(DOCKER_NAME):$(VERSION) -f ./Dockerfile.$(ARCH) .

dev:
	-docker rm -f $(DOCKER_NAME):$(VERSION)
	docker run -it -e WIOTP_DOMAIN="${WIOTP_DOMAIN}" -e HZN_ORGANIZATION="$(HZN_ORG_ID)" -e WIOTP_API_KEY="$(WIOTP_API_KEY)" -e WIOTP_API_AUTH_TOKEN="$(WIOTP_API_TOKEN)" -e WIOTP_DEVICE_TYPE="$(WIOTP_DEVICE_TYPE)" -e WIOTP_DEVICE_AUTH_TOKEN="$(WIOTP_GW_TOKEN)" -e HZN_DEVICE_ID="$(HZN_DEVICE_ID)" --name cpu2wiotp --net=example --net-alias=cpu2wiotp --volume `pwd`:/outside $(DOCKER_NAME):$(VERSION) /bin/sh

run:
	-docker rm -f $(DOCKER_NAME):$(VERSION)
	docker run -d -e WIOTP_DOMAIN="${WIOTP_DOMAIN}" -e HZN_ORGANIZATION="$(HZN_ORG_ID)" -e WIOTP_API_KEY="$(WIOTP_API_KEY)" -e WIOTP_API_AUTH_TOKEN="$(WIOTP_API_TOKEN)" -e WIOTP_DEVICE_TYPE="$(WIOTP_DEVICE_TYPE)" -e WIOTP_DEVICE_AUTH_TOKEN="$(WIOTP_GW_TOKEN)" -e HZN_DEVICE_ID="$(HZN_DEVICE_ID)" --name cpu2wiotp --net=example --net-alias=cpu2wiotp $(DOCKER_NAME):$(VERSION)

check:
	curl -s localhost:8347/v1/cpu | jq

# To publish you must have write access to the docker hub openhorizon user
publish:
	docker tag $(DOCKER_NAME):$(VERSION) $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(VERSION)
	docker push $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(VERSION)

clean:
	-docker rm -f $(DOCKER_NAME):$(VERSION)
	-docker rmi cpu2wiotp

.PHONY: default all build dev run check publish clean
