
all:
	make build
	make run

# Transform the machine arch into some standard values: "arm", "arm64", or "x86"
SYSTEM_ARCH := $(shell uname -m | sed -e 's/aarch64.*/arm64/' -e 's/x86_64.*/x86/' -e 's/armv.*/arm/')

# To build for an arch different from the current system, set this env var to 1 of the values in the comment above
ARCH ?= $(SYSTEM_ARCH)

build:
	docker build -t cpu -f ./Dockerfile.$(ARCH) .

dev:
	-docker network create example
	-docker stop cpu
	-docker rm cpu
	docker run -it --name cpu --publish=8347:8347 --net=example --net-alias=cpu --volume `pwd`:/outside cpu /bin/sh

run:
	-docker network create example
	-docker stop cpu
	-docker rm cpu
	docker run -d --name cpu --publish=8347:8347 --net=example --net-alias=cpu --volume `pwd`:/outside cpu

check:
	curl -s localhost:8347/v1/cpu | jq

# These variables can be overridden from the environment
CPU_VERSION ?= 1.2.2
CPU_DOCKER_NAME ?= openhorizon/example_ms_$(ARCH)_cpu

# To publish you must have write access to the docker hub openhorizon user
publish:
	docker build -t $(CPU_DOCKER_NAME):$(CPU_VERSION) -f ./Dockerfile.$(ARCH) .
	docker push $(CPU_DOCKER_NAME):$(CPU_VERSION)

clean:
	-docker stop cpu
	-docker rm cpu
	-docker rmi cpu
	-docker network rm example

.PHONY: all build dev run check publish clean
