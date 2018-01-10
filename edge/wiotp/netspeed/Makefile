all:
	make build
	make run

# Get the 3-character prefix of the local host architecture ("arm" or "x86")
SYSTEM_ARCH := $(shell uname -m | sed -e 's/aarch64.*/arm64/' -e 's/x86_64.*/amd64/' -e 's/armv.*/arm/')

# To build for an arch different from the current system, set this env var to arm or x86 before running make
ARCH ?= $(SYSTEM_ARCH)
PORT := 8349

build:
	docker build -t netspeed2wiotp -f ./Dockerfile.$(ARCH) .

dev:
	-docker stop netspeed2wiotp
	-docker rm netspeed2wiotp
	docker run -it -e WIOTP_DOMAIN="${WIOTP_DOMAIN}" -e WIOTP_ORG_ID="$(WIOTP_ORG_ID)" -e WIOTP_API_KEY="$(WIOTP_API_KEY)" -e WIOTP_API_AUTH_TOKEN="$(WIOTP_API_AUTH_TOKEN)" -e WIOTP_DEVICE_TYPE="$(WIOTP_DEVICE_TYPE)" -e WIOTP_DEVICE_AUTH_TOKEN="$(WIOTP_DEVICE_AUTH_TOKEN)" -e HZN_DEVICE_ID="$(HZN_DEVICE_ID)" --name netspeed2wiotp --net=example-netspeed --net-alias=netspeed2wiotp --volume `pwd`:/outside netspeed2wiotp /bin/sh

run:
	-docker stop netspeed2wiotp
	-docker rm netspeed2wiotp
	docker run -d -e WIOTP_DOMAIN="${WIOTP_DOMAIN}" -e WIOTP_ORG_ID="$(WIOTP_ORG_ID)" -e WIOTP_API_KEY="$(WIOTP_API_KEY)" -e WIOTP_API_AUTH_TOKEN="$(WIOTP_API_AUTH_TOKEN)" -e WIOTP_DEVICE_TYPE="$(WIOTP_DEVICE_TYPE)" -e WIOTP_DEVICE_AUTH_TOKEN="$(WIOTP_DEVICE_AUTH_TOKEN)" -e HZN_DEVICE_ID="$(HZN_DEVICE_ID)" --name netspeed2wiotp --net=example-netspeed --net-alias=netspeed2wiotp netspeed2wiotp

check:
	curl -s localhost:$(PORT)/v1/netspeed | jq

# These variables can be overridden from the environment
NETSPEED2WIOTP_VERSION ?= 2.5
NETSPEED2WIOTP_DOCKER_NAME ?= openhorizon/example_wl_$(ARCH)_netspeed2wiotp

# To publish you must have write access to the docker hub openhorizon user
publish:
	docker build -t $(NETSPEED2WIOTP_DOCKER_NAME):$(NETSPEED2WIOTP_VERSION) -f ./Dockerfile.$(ARCH) .
	docker push $(NETSPEED2WIOTP_DOCKER_NAME):$(NETSPEED2WIOTP_VERSION)

clean:
	-docker stop netspeed2wiotp
	-docker rm netspeed2wiotp
	-docker rmi netspeed2wiotp

.PHONY: all build dev run check publish clean
