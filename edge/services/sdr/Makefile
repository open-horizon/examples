# Transform the machine arch into some standard values: "arm", "arm64", or "amd64"
SYSTEM_ARCH := $(shell uname -m | sed -e 's/aarch64.*/arm64/' -e 's/x86_64.*/amd64/' -e 's/armv.*/arm/')

# To build for an arch different from the current system, set this env var to 1 of the values in the comment above
export ARCH ?= $(SYSTEM_ARCH)

export SDR_VERSION = 0.0.7
export SDR_NAME = sdr

# These variables can be overridden from the environment
DOCKER_NAME ?= $(ARCH)_$(SDR_NAME)
export DOCKER_HUB_ID ?= openhorizon
export MYDOMAIN ?= github.com.open-horizon.examples

default: build

all: build run


build: target/build_$(ARCH)

target/build_$(ARCH): Dockerfile.$(ARCH) main.go
	docker build -t $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(SDR_VERSION) -f Dockerfile.$(ARCH) --build-arg version=$(SDR_VERSION) .
	touch target/build_$(ARCH)

run: target/build_$(ARCH)
	-docker network create sdr
	-docker kill $(SDR_NAME)
	docker run -it --rm --privileged --name $(SDR_NAME) --net=sdr --net-alias=rtlsdr --publish=5427:5427 $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(SDR_VERSION)

# To publish you must have write access to the docker hub openhorizon user
publish: target/build_$(ARCH)
	docker push $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(SDR_VERSION)

# Create/update the metadata in the exchange for this service. This also pushes the image to docker hub.
publish-service: #todo: add this: target/build_$(ARCH)
	: $${HZN_ORG_ID:?} $${HZN_EXCHANGE_USER_AUTH:?} $${PRIVATE_KEY_FILE:?} $${PUBLIC_KEY_FILE:?}   # this verifies these env vars are set
	hzn exchange service publish -k $$PRIVATE_KEY_FILE -K $$PUBLIC_KEY_FILE -f horizon/service.definition.json

clean:
	rm target/*
	-docker rm -f $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(SDR_VERSION)
	-docker rmi cpu
	-docker network rm example

.PHONY: default all build dev run check publish publish-service clean
