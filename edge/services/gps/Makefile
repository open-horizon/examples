#
# Makefile for gps
#
# NOTE: Requires make to be installed, but no other build tools are required.

export GPS_NAME ?= gps
export GPS_VERSION ?= 2.0.5
export DOCKER_HUB_ID ?= openhorizon
# If you use this domain in the service definition, the resulting url has documentation for this service
export MYDOMAIN ?= github.com.open-horizon.examples
IMAGE ?= $(GPS_NAME)
BINARY=$(IMAGE)
SOURCE=main.go
OTHER_SOURCES:=$(wildcard src/*/*.go)
TEST_SOURCE:=test.go

# Docker --device entries enabling access to supported GPS device paths:
GPS_DEVICE_LIST ?= /dev/ttyACM0:/dev/ttyACM0:rw
GPSD_DOCKER_DEVICE_ENTRIES := $(addprefix --volume ,$(GPS_DEVICE_LIST))

# Firmware configuration (this is only used for development and testing)
# REST server runs on port 3779 by default (test code needs to know the port)
HZN_GPS_PORT ?= 31779
# An arbitrary manually-provided location southeast of San Jose, CA, USA
#HZN_LAT=37.0
#HZN_LON=-121.0
HZN_LAT ?= 0.0
HZN_LON ?= 0.0
HZN_LOCATION_ACCURACY_KM ?= 10.0

# The docker network used for testing
DOCKER_TEST_NETWORK ?= gps-test

# Transform the machine arch into some standard values: "arm", "arm64", or "amd64". Note: ppc64le already returns ppc64le, so we do not need to change that.
SYSTEM_ARCH := $(shell uname -m | sed -e 's/aarch64.*/arm64/' -e 's/x86_64.*/amd64/' -e 's/armv.*/arm/')

# To build for an arch different from the current system, set this env var to 1 of the values in the comment above
export ARCH ?= $(SYSTEM_ARCH)

# Cross compile in this context means on amd64 compile for arm or arm64. You must also set ARCH to arm or arm64
CROSS_COMPILE ?= false

# Construct our Docker registry path and other architecture-dependent values
REG_PATH := $(DOCKER_HUB_ID)/$(ARCH)
IMAGE_PATH = $(REG_PATH)_$(IMAGE):$(GPS_VERSION)

# Make will search these two relative paths for input and target files
VPATH=src:bin

# Running make with no target will build the production image
.PHONY: default
default: build

.PHONY: build
build: $(IMAGE)-$(GPS_VERSION)

# Omitting the image architecture will default to the local architecture
# NOTE: Must depend on a file in the file system, not an image!
.PHONY: $(IMAGE)-$(GPS_VERSION)
$(IMAGE)-$(GPS_VERSION): .$(IMAGE)-$(GPS_VERSION)-$(ARCH)

# Omitting the binary architecture will default to the local architecture
.PHONY: $(BINARY)
$(BINARY): $(BINARY).$(ARCH)

# This section has compilation rules for the architecture-specific binaries
# For RPi (ARMv7)
GOOPTIONS_ARM=CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7
$(BINARY).arm: $(SOURCE) $(OTHER_SOURCES)
	docker run --rm -t --volume `pwd`:/outside golang env GOPATH=/outside $(GOOPTIONS_ARM) sh -c "cd /outside && go get github.com/kellydunn/golang-geo && go build -o bin/$@ src/$(SOURCE)"

# For arm64
GOOPTIONS_ARM64=CGO_ENABLED=0 GOOS=linux GOARCH=arm64
$(BINARY).arm64: $(SOURCE) $(OTHER_SOURCES)
	docker run --rm -t --volume `pwd`:/outside golang env GOPATH=/outside $(GOOPTIONS_ARM64) sh -c "cd /outside && go get github.com/kellydunn/golang-geo && go build -o bin/$@ src/$(SOURCE)"

# For amd64 (this is for 64bit)
GOOPTIONS_AMD64=CGO_ENABLED=0 GOOS=linux GOARCH=amd64
$(BINARY).amd64: $(SOURCE) $(OTHER_SOURCES)
	docker run --rm -t --volume `pwd`:/outside golang env GOPATH=/outside $(GOOPTIONS_AMD64) sh -c "cd /outside && go get github.com/kellydunn/golang-geo && go build -o bin/$@ src/$(SOURCE)"

# For ppc64le
GOOPTIONS_PPC=CGO_ENABLED=0 GOOS=linux GOARCH=ppc64le
$(BINARY).ppc64le: $(SOURCE) $(OTHER_SOURCES)
	docker run --rm -t --volume `pwd`:/outside gloang env GOPATH=/outside $(GOOPTIONS_PPC) sh -c "cd /outside && go get github.com/kellydunn/golang-geo && go build -o bin/$@ src/$(SOURCE)"

# This is the production image target (which depends on the BINARY file)
.$(IMAGE)-$(GPS_VERSION)-$(ARCH): Dockerfile.$(ARCH) $(BINARY).$(ARCH)
	docker build -t $(IMAGE_PATH) -f Dockerfile.$(ARCH) .
	mkdir -p bin
	touch bin/.$(IMAGE)-$(GPS_VERSION)-$(ARCH)

# This is the "docker save" tarball, suitable for "docker load"
bin/$(IMAGE)-$(GPS_VERSION)-$(ARCH).tar: .$(IMAGE)-$(GPS_VERSION)-$(ARCH)
	mkdir -p bin
	docker save -o bin/$(IMAGE)-$(GPS_VERSION)-$(ARCH).tar $(IMAGE_PATH)

# Push the production image for this architecture to the resistry, and verify
.PHONY: docker-push
docker-push: .$(IMAGE)-$(GPS_VERSION)-$(ARCH)
	docker push $(IMAGE_PATH)

.PHONY: docker-verify
docker-verify:
	docker version --format '{{.Server.Version}}'
	docker pull $(IMAGE_PATH) | grep '^Digest: '

.PHONY: publish-service
publish-service:
	: $${HZN_EXCHANGE_USER_AUTH:?} $${PRIVATE_KEY_FILE:?} $${PUBLIC_KEY_FILE:?}   # this verifies these env vars are set
	hzn exchange service publish -k $$PRIVATE_KEY_FILE -K $$PUBLIC_KEY_FILE -f horizon/service.definition.json

# Create the docker network
.PHONY: network
network:
	-docker network create $(DOCKER_TEST_NETWORK)

.PHONY: develop
develop: network .$(IMAGE)-$(GPS_VERSION)-$(ARCH)
	-docker rm -f $(IMAGE)
	docker run --rm -it --privileged --volume=`pwd`:/outside $(GPSD_DOCKER_DEVICE_ENTRIES) -e HZN_GPS_PORT=$(HZN_GPS_PORT) -e HZN_USE_GPS=true -e HZN_LAT=$(HZN_LAT) -e HZN_LON=$(HZN_LON) -e HZN_LOCATION_ACCURACY_KM=$(HZN_LOCATION_ACCURACY_KM) --net=$(DOCKER_TEST_NETWORK) --net-alias=$(IMAGE) --name $(IMAGE) $(IMAGE_PATH) /bin/sh

# Run the container as a daemon using GPS hardware for testing purposes
.PHONY: gps_daemon
gps_daemon: network .$(IMAGE)-$(GPS_VERSION)-$(ARCH)
	-docker rm -f $(IMAGE)
	docker run -d --privileged $(GPSD_DOCKER_DEVICE_ENTRIES) -p $(HZN_GPS_PORT):$(HZN_GPS_PORT) -e HZN_GPS_PORT=$(HZN_GPS_PORT) -e HZN_USE_GPS=true -e HZN_LAT=$(HZN_LAT) -e HZN_LON=$(HZN_LON) -e HZN_LOCATION_ACCURACY_KM=$(HZN_LOCATION_ACCURACY_KM) --net=$(DOCKER_TEST_NETWORK) --net-alias=$(IMAGE) --name $(IMAGE) $(IMAGE_PATH)

# Run the container as a daemon using static coordinates or IP estimate (but not GPS device) testing purposes
.PHONY: daemon
daemon: network .$(IMAGE)-$(GPS_VERSION)-$(ARCH)
	-docker rm -f $(IMAGE)
	docker run -d --privileged -p $(HZN_GPS_PORT):$(HZN_GPS_PORT) -e HZN_GPS_PORT=$(HZN_GPS_PORT) -e HZN_USE_GPS=false -e HZN_LAT=$(HZN_LAT) -e HZN_LON=$(HZN_LON) -e HZN_LOCATION_ACCURACY_KM=$(HZN_LOCATION_ACCURACY_KM) --net=$(DOCKER_TEST_NETWORK) --net-alias=$(IMAGE) --name $(IMAGE) $(IMAGE_PATH)

IMAGE_TEST = $(IMAGE)-test
BINARY_TEST = $(BINARY)-test

$(BINARY_TEST).arm: $(TEST_SOURCE) $(OTHER_SOURCES)
	docker run --rm -t --volume `pwd`:/outside golang env GOPATH=/outside $(GOOPTIONS_ARM) sh -c "cd /outside && go get gopkg.in/go-playground/validator.v8 github.com/eclipse/paho.mqtt.golang && go build -o bin/$@ src/$(TEST_SOURCE)"

$(BINARY_TEST).arm64: $(TEST_SOURCE) $(OTHER_SOURCES)
	docker run --rm -t --volume `pwd`:/outside golang env GOPATH=/outside $(GOOPTIONS_ARM64) sh -c "cd /outside && go get gopkg.in/go-playground/validator.v8 github.com/eclipse/paho.mqtt.golang && go build -o bin/$@ src/$(TEST_SOURCE)"

$(BINARY_TEST).amd64: $(TEST_SOURCE) $(OTHER_SOURCES)
	docker run --rm -t --volume `pwd`:/outside golang env GOPATH=/outside $(GOOPTIONS_AMD64) sh -c "cd /outside && go get gopkg.in/go-playground/validator.v8 github.com/eclipse/paho.mqtt.golang && go build -o bin/$@ src/$(TEST_SOURCE)"

# Build the automated test container
.$(IMAGE_TEST)-$(ARCH): Dockerfile.$(ARCH).test $(BINARY_TEST).$(ARCH)
	make $(IMAGE_TEST)-$(ARCH)
.PHONY: $(IMAGE_TEST)-$(ARCH)
$(IMAGE_TEST)-$(ARCH):
	- docker rm -f $(IMAGE_TEST) 2> /dev/null || :
	docker build -t $(REG_PATH)/$(IMAGE_TEST) -f Dockerfile.$(ARCH).test .
	mkdir -p bin
	touch bin/.$(IMAGE_TEST)-$(ARCH)
	@echo "$(IMAGE_TEST)-$(ARCH) has been built."

# This is the "docker save" tarball, suitable for "docker load"
bin/$(IMAGE_TEST)-$(ARCH).tar: .$(IMAGE_TEST)-$(ARCH)
	mkdir -p bin
	docker save -o bin/$(IMAGE_TEST)-$(ARCH).tar $(REG_PATH)/$(IMAGE_TEST)

.PHONY: docker-push-test
docker-push-test: .$(IMAGE_TEST)-$(ARCH)
	docker push $(REG_PATH)/$(IMAGE_TEST)
	docker version --format '{{.Server.Version}}'
	docker pull $(REG_PATH)/$(IMAGE_TEST) | grep '^Digest: '

#
# Run the test image in developer mode on the local host
#   - opens a shell into the test container, with this host directory mounted
#   - inside that shell you will typically work in that directory.  E.g.:
#        # cd /outside
#        # vi src/test.go ...
#     then compile and run the tests with go:
#        # go run src/test.go
#
.PHONY: develop_test
develop_test:  network .$(IMAGE_TEST)
	-docker rm -f $(IMAGE_TEST)
	docker run --rm -it --volume=`pwd`:/outside -e HZN_GPS_PORT=$(HZN_GPS_PORT) --net=$(DOCKER_TEST_NETWORK) --net-alias=$(IMAGE_TEST) --name $(IMAGE_TEST) $(IMAGE_TEST):$(GPS_VERSION) /bin/bash

# Run the automated tests
.PHONY: test
test:  network .$(IMAGE_TEST)
	-docker rm -f $(IMAGE_TEST)
	docker run --rm -t --volume `pwd`:/outside -e HZN_GPS_PORT=$(HZN_GPS_PORT) --net=$(DOCKER_TEST_NETWORK) --net-alias=$(IMAGE_TEST) --name $(IMAGE_TEST) $(IMAGE_TEST):$(GPS_VERSION)

# Clean everything
.PHONY: clean
clean:
	-rm -rf bin pkg bin/gps.json
	-rm -rf src/github.com src/gopkg.in
	-docker rm -f $(IMAGE)
	-docker rmi $(IMAGE)-$(ARCH):latest $(IMAGE)-$(ARCH):$(GPS_VERSION) $(REG_PATH)/$(IMAGE):$(GPS_VERSION)
	-docker rm -f $(IMAGE_TEST)
	-docker rmi $(IMAGE_TEST)
	-docker network rm $(DOCKER_TEST_NETWORK)

# Listen to the data stream from the hardware (assumes ublox7 at /dev/ttyACM0)
.PHONY: check_hw
check_hw:
	cat /dev/ttyACM0 | egrep 'GPRMC|GPGGA'

# Listen to the data stream from the gpsd (assumes device at /dev/ttyACM0)
# Note: must have previously built and installed gpsd wherever you run this.
.PHONY: check_gpsd
check_gpsd:
	gpsd -D5 -N -n /dev/ttyACM0

# Check the output from the service:
.PHONY: check
check:
	 curl -s `docker inspect gps | grep IPAddress | tail -1 | sed 's/.*": "//;s/".*//'`:31779/v1/gps | jq

