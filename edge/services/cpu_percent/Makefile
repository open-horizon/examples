# Transform the machine arch into some standard values: "arm", "arm64", or "amd64"
SYSTEM_ARCH := $(shell uname -m | sed -e 's/aarch64.*/arm64/' -e 's/x86_64.*/amd64/' -e 's/armv.*/arm/')

# To build for an arch different from the current system, set this env var to 1 of the values in the comment above
ARCH ?= $(SYSTEM_ARCH)

# These variables can be overridden from the environment
VERSION ?= 1.2.2
DOCKER_NAME ?= example_ms_$(ARCH)_cpu
DOCKER_HUB_ID ?= openhorizon

default: build

all: build run

build:
	docker build -t $(DOCKER_NAME):$(VERSION) -f ./Dockerfile.$(ARCH) .

dev:
	-docker network create example
	-docker rm -f $(DOCKER_NAME):$(VERSION)
	docker run -it --name cpu --publish=8347:8347 --net=example --net-alias=cpu --volume `pwd`:/outside $(DOCKER_NAME):$(VERSION) /bin/sh

run:
	-docker network create example
	-docker rm -f $(DOCKER_NAME):$(VERSION)
	docker run -d --name cpu --publish=8347:8347 --net=example --net-alias=cpu --volume `pwd`:/outside $(DOCKER_NAME):$(VERSION)

check:
	curl -s localhost:8347/v1/cpu | jq

# To publish you must have write access to the docker hub openhorizon user
publish:
	docker tag $(DOCKER_NAME):$(VERSION) $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(VERSION)
	docker push $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(VERSION)

clean:
	-docker rm -f $(DOCKER_NAME):$(VERSION)
	-docker rmi cpu
	-docker network rm example

.PHONY: default all build dev run check publish clean
