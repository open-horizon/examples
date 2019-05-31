#
# Makefile for gps Service container
#
# NOTE: Requires make to be installed, but no other build tools are required.
#
# To test this code:
#   make            # build the "gps" container
#   make gps-daemon # or, to not use GPS hardware, use: "make daemon" instead)
#   make test       # kill it after after functional tests, or run it forever!
#

# This imports the variables from horizon/hzn.json. You can ignore these lines, but do not remove them.
-include horizon/.hzn.json.tmp.mk

# Default ARCH to the architecture of this machines (as horizon/golang describes it)
export ARCH ?= $(shell hzn architecture)

# Docker --device entries enabling access to supported GPS device paths:
#GPS_DEVICE_LIST ?= /dev/ttyACM0:/dev/ttyACM0:rw
#GPSD_DOCKER_DEVICE_ENTRIES := $(addprefix --volume ,$(GPS_DEVICE_LIST))

# Firmware configuration (this is only used for development and testing)
# REST server runs on port 31779 by default (test code needs to know the port)
#HZN_GPS_PORT ?= 80
# An arbitrary manually-provided location southeast of San Jose, CA, USA
#HZN_LAT=37.0
#HZN_LON=-121.0
#HZN_LAT ?= 0.0
#HZN_LON ?= 0.0
#HZN_LOCATION_ACCURACY_KM ?= 10.0

# The docker network used for testing
#DOCKER_TEST_NETWORK ?= gps-test-network

# Cross compile in this context means on amd64 compile for arm or arm64. You must also set ARCH to arm or arm64
#CROSS_COMPILE ?= false

# Construct our Docker registry path and other architecture-dependent values
#REG_PATH := $(DOCKER_HUB_ID)/$(ARCH)
#IMAGE_PATH = $(REG_PATH)_$(IMAGE):$(GPS_VERSION)

# Make will search these two relative paths for input and target files
#VPATH=src:bin

# Running make with no target will build the production image
#.PHONY: default
#default: build

#.PHONY: build
#build: $(IMAGE)-$(GPS_VERSION)

# Omitting the image architecture will default to the local architecture
# NOTE: Must depend on a file in the file system, not an image!
#.PHONY: $(IMAGE)-$(GPS_VERSION)
#$(IMAGE)-$(GPS_VERSION): .$(IMAGE)-$(GPS_VERSION)-$(ARCH)

# Omitting the binary architecture will default to the local architecture
#.PHONY: $(BINARY)
#$(BINARY): $(BINARY).$(ARCH)

# This section has compilation rules for the architecture-specific binaries
# For RPi, most models (ARMv6)
#GOOPTIONS_ARM=CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6
#$(BINARY).arm: $(SOURCE) $(OTHER_SOURCES)
	#docker run --rm -t --volume `pwd`:/outside golang env GOPATH=/outside $(GOOPTIONS_ARM) sh -c "cd /outside && go get github.com/kellydunn/golang-geo && go build -o bin/$@ src/$(SOURCE)"


build:
	docker build -t $(DOCKER_IMAGE_BASE)_$(ARCH):$(SERVICE_VERSION) -f ./Dockerfile.$(ARCH) .

build-all-arches:
	ARCH=amd64 $(MAKE) build
	ARCH=arm $(MAKE) build
	ARCH=arm64 $(MAKE) build

test: build
	hzn dev service start -S
	@echo 'Testing service...'
	sleep 1
	curl -sS -w "%{http_code}" http://ibm.gps:80/v1/gps/location | jq .
	hzn dev service stop

# Create/update the metadata in the exchange for this service
publish-service:
	hzn exchange service publish -f horizon/service.definition.json

build-test-publish: build test publish-service

publish-all-arches:
	ARCH=amd64 $(MAKE) build-test-publish
	ARCH=arm $(MAKE) build-test-publish
	ARCH=arm64 $(MAKE) build-test-publish

# new target for icp exchange to run on startup to publish only
icp-publish:
	ARCH=amd64 $(MAKE) publish-service
	ARCH=arm $(MAKE) publish-service
	ARCH=arm64 $(MAKE) publish-service

clean:
	-docker rmi $(DOCKER_IMAGE_BASE)_$(ARCH):$(SERVICE_VERSION) 2> /dev/null || :

clean-all-arches:
	ARCH=amd64 $(MAKE) clean
	ARCH=arm $(MAKE) clean
	ARCH=arm64 $(MAKE) clean

# This imports the variables from horizon/hzn.cfg. You can ignore these lines, but do not remove them.
horizon/.hzn.json.tmp.mk: horizon/hzn.json
	@ hzn util configconv -f $< > $@

.PHONY: default all build run check stop publish publish-service publish-service-only clean
