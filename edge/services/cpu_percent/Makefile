# Transform the machine arch into some standard values: "arm", "arm64", or "amd64"
SYSTEM_ARCH := $(shell uname -m | sed -e 's/aarch64.*/arm64/' -e 's/x86_64.*/amd64/' -e 's/armv.*/arm/')

# To build for an arch different from the current system, set this env var to 1 of the values in the comment above
export ARCH ?= $(SYSTEM_ARCH)

# These variables can be overridden from the environment
export CPU_NAME ?= cpu
export CPU_VERSION ?= 1.2.2
DOCKER_NAME ?= $(ARCH)_$(CPU_NAME)
export DOCKER_HUB_ID ?= openhorizon
export MYDOMAIN ?= github.com.open-horizon.examples
PORT_NUM ?= 8347
IMAGE_PATH ?= $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(CPU_VERSION)

default: all

all: build run check

build:
	docker build -t $(IMAGE_PATH) -f ./Dockerfile.$(ARCH) .

run:
	-docker rm -f $(DOCKER_NAME) 2> /dev/null || :
	docker run -d --name $(DOCKER_NAME) --publish=$(PORT_NUM):$(PORT_NUM) --volume `pwd`:/outside $(IMAGE_PATH)

check:
	curl -sSL http://localhost:$(PORT_NUM)/v1/cpu | jq .

stop:
	-docker rm -f $(DOCKER_NAME) 2> /dev/null || :

# To publish you must have write access to the docker hub DOCKER_HUB_ID user
docker-push: build
	docker push $(IMAGE_PATH)

# Create/update the metadata in the exchange for this service
publish-service: build publish-service-only
publish-service-only:
	: $${HZN_EXCHANGE_USER_AUTH:?} $${PRIVATE_KEY_FILE:?} $${PUBLIC_KEY_FILE:?}   # this verifies these env vars are set
	hzn exchange service publish -k $$PRIVATE_KEY_FILE -K $$PUBLIC_KEY_FILE -f horizon/service.definition.json

clean:
	-docker rm -f $(DOCKER_NAME) 2> /dev/null || :
	-docker rmi $(IMAGE_PATH) 2> /dev/null || :

.PHONY: default all build run check stop publish publish-service publish-service-only clean
