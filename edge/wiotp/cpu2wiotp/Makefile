
all:
	make build
	make run

# Get the 3-character prefix of the local host architecture ("arm" or "x86")
SYSTEM_ARCH := $(shell uname -m | sed 's/\(...\).*/\1/')
# To build for an arch different from the current system, set this env var to arm or x86 before running make
ARCH ?= $(SYSTEM_ARCH)

build:
	docker build -t cpu2wiotp -f ./Dockerfile.$(ARCH) .

dev:
	-docker stop cpu2wiotp
	-docker rm cpu2wiotp
	docker run -it -e WIOTP_DOMAIN="${WIOTP_DOMAIN}" -e WIOTP_ORG_ID="$(WIOTP_ORG_ID)" -e WIOTP_API_KEY="$(WIOTP_API_KEY)" -e WIOTP_API_AUTH_TOKEN="$(WIOTP_API_AUTH_TOKEN)" -e WIOTP_DEVICE_TYPE="$(WIOTP_DEVICE_TYPE)" -e WIOTP_DEVICE_AUTH_TOKEN="$(WIOTP_DEVICE_AUTH_TOKEN)" -e HZN_DEVICE_ID="$(HZN_DEVICE_ID)" --name cpu2wiotp --net=example --net-alias=cpu2wiotp --volume `pwd`:/outside cpu2wiotp /bin/sh

run:
	-docker stop cpu2wiotp
	-docker rm cpu2wiotp
	docker run -d -e WIOTP_DOMAIN="${WIOTP_DOMAIN}" -e WIOTP_ORG_ID="$(WIOTP_ORG_ID)" -e WIOTP_API_KEY="$(WIOTP_API_KEY)" -e WIOTP_API_AUTH_TOKEN="$(WIOTP_API_AUTH_TOKEN)" -e WIOTP_DEVICE_TYPE="$(WIOTP_DEVICE_TYPE)" -e WIOTP_DEVICE_AUTH_TOKEN="$(WIOTP_DEVICE_AUTH_TOKEN)" -e HZN_DEVICE_ID="$(HZN_DEVICE_ID)" --name cpu2wiotp --net=example --net-alias=cpu2wiotp cpu2wiotp

check:
	curl -s localhost:8347/v1/cpu | jq

# These variables can be overridden from the environment
CPU2WIOTP_VERSION ?= 1.1.4
CPU2WIOTP_DOCKER_NAME ?= openhorizon/example_wl_$(ARCH)_cpu2wiotp

# To publish you must have write access to the docker hub openhorizon user
publish:
	docker build -t $(CPU2WIOTP_DOCKER_NAME):$(CPU2WIOTP_VERSION) -f ./Dockerfile.$(ARCH) .
	docker push $(CPU2WIOTP_DOCKER_NAME):$(CPU2WIOTP_VERSION)

clean:
	-docker stop cpu2wiotp
	-docker rm cpu2wiotp
	-docker rmi cpu2wiotp

.PHONY: all build dev run check publish clean
