# Transform the machine arch into some standard values: "arm", "arm64", or "x86"
SYSTEM_ARCH := $(shell uname -m | sed -e 's/aarch64.*/arm64/' -e 's/x86_64.*/amd64/' -e 's/armv.*/arm/')

# To build for an arch different from the current system, set this env var to 1 of the values in the comment above
ARCH ?= $(SYSTEM_ARCH)

# These variables can be overridden from the environment
MS_NAME := caffedemo
CAFFEDEMO_VERSION ?= $(CAFFEDEMO_VERSION)
MS_VERSION := $(CAFFEDEMO_VERSION)
DOCKER_NAME ?= $(ARCH)_example_ms_$(MS_NAME)
DOCKER_HUB_ID ?= openhorizon

default: build

all: build run

build:
	docker build -t $(DOCKER_NAME):$(MS_VERSION) -f ./Dockerfile.$(ARCH) .

dev:
	-docker network create example
	-docker rm -f $(MS_NAME)
	docker run -it --name $(MS_NAME) --net=example --privileged --net-alias=$(MS_NAME) --volume /dev:/dev --volume `pwd`:/outside $(DOCKER_NAME):$(MS_VERSION) /bin/bash

run:
	-docker network create example
	-docker rm -f $(MS_NAME)
	docker run -d --name $(MS_NAME) --net=example --privileged --volume /dev:/dev --net-alias=$(MS_NAME) $(DOCKER_NAME):$(MS_VERSION)

#check:
#	curl -s localhost:8347/v1/cpu | jq

# To publish you must have write access to the docker hub openhorizon user
publish:
	docker tag $(DOCKER_NAME):$(MS_VERSION) $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(MS_VERSION)
	docker push $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(MS_VERSION)

clean:
	-docker network rm example
	-docker stop $(MS_NAME) && docker rm $(MS_NAME)
	-docker rmi $(MS_NAME)
	-docker rm -f $(DOCKER_NAME):$(MS_VERSION)

.PHONY: default all build dev run check publish clean
