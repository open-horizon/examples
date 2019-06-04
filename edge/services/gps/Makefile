#
# Makefile for gps Service container
#

# This imports the variables from horizon/hzn.json. You can ignore these lines, but do not remove them.
-include horizon/.hzn.json.tmp.mk

# Default ARCH to the architecture of this machines (as horizon/golang describes it)
export ARCH ?= $(shell hzn architecture)

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
