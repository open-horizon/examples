# Transform the machine arch into some standard values: "arm", "arm64", or "amd64"
SYSTEM_ARCH := $(shell uname -m | sed -e 's/aarch64.*/arm64/' -e 's/x86_64.*/amd64/' -e 's/armv.*/arm/')

# To build for an arch different from the current system, set this env var to 1 of the values in the comment above
export ARCH ?= $(SYSTEM_ARCH)

# Container image version and workload name
export SDR2MSGHUB_VERSION = 0.0.4
export SDR2MSGHUB_NAME = sdr2msghub

# When creating the service definition, we reference these
export SDR_NAME ?= sdr
export SDR_VERSION ?= 0.0.7
export GPS_NAME ?= gps
export GPS_VERSION ?= 2.0.5

# Construct our Docker registry path and other architecture-dependent values
DOCKER_NAME ?= $(ARCH)_$(SDR2MSGHUB_NAME)
export DOCKER_HUB_ID ?= openhorizon
IMAGE_PATH = $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(SDR2MSGHUB_VERSION)

export MYDOMAIN ?= github.com.open-horizon.examples
MSGHUB_TOPIC ?= sdr-audio

build: target/build_$(ARCH)

target/build_$(ARCH): main.go
	docker build -t $(IMAGE_PATH) -f Dockerfile.$(ARCH) ../../
	touch target/build_$(ARCH)

run: target/build_$(ARCH)
	-docker network create sdr
	#-docker rm -f $(IMAGE_PATH)
	docker run -it --rm -e HZN_ORGANIZATION="$(HZN_ORG_ID)" -e USE_GPS="false" -e HZN_DEVICE_ID="$(HZN_DEVICE_ID)" -e MSGHUB_API_KEY="$(MSGHUB_API_KEY)" -e MSGHUB_BROKER_URL="$(MSGHUB_BROKER_URL)" -e MSGHUB_TOPIC="$(MSGHUB_TOPIC)" --name $(SDR2MSGHUB_NAME) --net=sdr --net-alias=$(SDR2MSGHUB_NAME) $(IMAGE_PATH)

# To publish you must have write access to the docker hub openhorizon user
publish: target/build_$(ARCH)
	docker push $(IMAGE_PATH)

# Create/update the metadata in the exchange for this service. This also pushes the image to docker hub.
publish-service: #todo: add this: target/build_$(ARCH)
	: $${HZN_ORG_ID:?} $${HZN_EXCHANGE_USER_AUTH:?} $${PRIVATE_KEY_FILE:?} $${PUBLIC_KEY_FILE:?} $${MSGHUB_BROKER_URL:?}   # this verifies these env vars are set
	hzn exchange service publish -k $$PRIVATE_KEY_FILE -K $$PUBLIC_KEY_FILE -f horizon/service.definition.json

# Create/update the deployment pattern for a set of edge nodes
publish-pattern:
	: $${HZN_ORG_ID:?} $${HZN_EXCHANGE_USER_AUTH:?}   # this verifies these env vars are set
	hzn exchange pattern publish -p $$SDR2MSGHUB_NAME -f horizon/pattern/sdr2msghub.json

fakepub: fake/main.go
	go build -o $@ $<

clean:
	-docker rm -f $(IMAGE_PATH)
	-docker stop $(SDR2MSGHUB_NAME) && docker rm $(SDR2MSGHUB_NAME)
	-docker rmi $(SDR2MSGHUB_NAME)
	-docker network rm example
	-rm target/*
	-rm -f fakepub

.PHONY: default all build dev run check publish publish-service publish-pattern fakepub clean

protobuf: audiomsg.proto
	protoc --go_out=audiolib audiomsg.proto
