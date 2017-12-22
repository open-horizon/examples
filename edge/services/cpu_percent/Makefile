
all:
	make build
	make run

# Get the prefix of the local host architecture and set accordingly ("arm", "arm64", or "x86")
SYSTEM_ARCH := $(shell uname -m | sed -e 's/aarch64.*/arm64/' -e 's/x86_64.*/amd64/' -e 's/armv.*/arm/')

# To build for an arch different from the current system, set this env var to the preferred arch before running make
ARCH ?= $(SYSTEM_ARCH)

# Checks for existing containers and networks
existing_ntwk := $(shell docker network ls | grep example | awk '{print $$2}')
existing_cntr := $(shell docker ps -a --no-trunc | awk '{if (NR!=1) print $$2}' | grep -x cpu)
existing_img := $(shell docker images --no-trunc | awk '{if (NR!=1) print $$1}' | grep -x cpu)

build:
	docker build -t cpu -f ./Dockerfile.$(ARCH) .

dev:
        ifeq ($(existing_ntwk), )
	-docker network create example
        endif

        ifeq ($(existing_cntr), cpu)
	-docker stop cpu
	-docker rm cpu
        endif

	docker run -it --name cpu --publish=8347:8347 --net=example --net-alias=cpu --volume `pwd`:/outside cpu /bin/sh

run:
        ifeq ($(existing_ntwk), )
	-docker network create example
        endif

        ifeq ($(existing_cntr), cpu)
	-docker stop cpu
	-docker rm cpu
        endif

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
        ifeq ($(existing_cntr), cpu)
	-docker stop cpu
	-docker rm cpu
        endif

        ifeq ($(existing_img), cpu )
	-docker rmi cpu
        endif

        ifeq ($(existing_ntwk), example )
	-docker network rm example
        endif

.PHONY: all build dev run check publish clean
