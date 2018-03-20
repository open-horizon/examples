# Transform the machine arch into some standard values: "arm", "arm64", or "x86"
SYSTEM_ARCH := $(shell uname -m | sed -e 's/aarch64.*/arm64/' -e 's/x86_64.*/amd64/' -e 's/armv.*/arm/')

# To build for an arch different from the current system, set this env var to 1 of the values in the comment above
ARCH ?= $(SYSTEM_ARCH)

# These variables can be overridden from the environment
VERSION ?= 1.0
WL_NAME ?= caffedemo2wiotp
DOCKER_NAME ?= $(ARCH)_example_wl_$(WL_NAME)
DOCKER_HUB_ID ?= openhorizon

default: build

all: build run

build:
	docker build -t $(DOCKER_NAME):$(VERSION) -f ./Dockerfile.$(ARCH) .

dev:
	-docker network create example
	-docker rm -f $(WL_NAME)
	docker run -it -e WIOTP_DOMAIN="${WIOTP_DOMAIN}" -e HZN_ORGANIZATION="$(WIOTP_ORG_ID)" -e WIOTP_API_KEY="$(WIOTP_API_KEY)" -e WIOTP_API_AUTH_TOKEN="$(WIOTP_API_AUTH_TOKEN)" -e WIOTP_DEVICE_TYPE="$(WIOTP_DEVICE_TYPE)" -e WIOTP_DEVICE_AUTH_TOKEN="$(WIOTP_DEVICE_AUTH_TOKEN)" -e HZN_DEVICE_ID="$(HZN_DEVICE_ID)" --name $(WL_NAME) --net=example --net-alias=$(WL_NAME) --volume `pwd`:/outside $(DOCKER_NAME):$(VERSION) /bin/sh

run:
	-docker network create example
	-docker rm -f $(WL_NAME)
	docker run -d -e WIOTP_DOMAIN="${WIOTP_DOMAIN}" -e HZN_ORGANIZATION="$(WIOTP_ORG_ID)" -e WIOTP_API_KEY="$(WIOTP_API_KEY)" -e WIOTP_API_AUTH_TOKEN="$(WIOTP_API_AUTH_TOKEN)" -e WIOTP_DEVICE_TYPE="$(WIOTP_DEVICE_TYPE)" -e WIOTP_DEVICE_AUTH_TOKEN="$(WIOTP_DEVICE_AUTH_TOKEN)" -e HZN_DEVICE_ID="$(HZN_DEVICE_ID)" --name $(WL_NAME) --net=example --net-alias=$(WL_NAME) $(DOCKER_NAME):$(VERSION)

#check:
#	curl -s localhost:8347/v1/caffedemo | jq

# To publish you must have write access to the docker hub openhorizon user
publish:
	docker tag $(DOCKER_NAME):$(VERSION) $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(VERSION)
	docker push $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(VERSION)

clean:
	-docker network rm example
	-docker stop $(WL_NAME) && docker rm $(WL_NAME)
	-docker rmi $(WL_NAME)
	-docker rm -f $(DOCKER_NAME):$(VERSION)

.PHONY: default all build dev run check publish clean
