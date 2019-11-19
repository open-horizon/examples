#
# Makefile for gps Service container
# 
# When cross-compiling using mac, Dockerfile.arm will fail with golang:1.8, 1.10, 1.11 and 1.12
# So far golang:1.9 has been the only version to successfully 

# This imports the variables from horizon/hzn.json. You can ignore these lines, but do not remove them
-include horizon/.hzn.json.tmp.mk

# Default ARCH to the architecture of this machines (as horizon/golang describes it)
export ARCH ?= $(shell hzn architecture)

# Configurable parameters passed to serviceTest.sh in "test" target
export MATCH:='"latitude":'
export TIME_OUT:=20

build:
	docker build -t $(DOCKER_IMAGE_BASE)_$(ARCH):$(SERVICE_VERSION) -f ./Dockerfile.$(ARCH) .

build-all-arches:
	ARCH=amd64 $(MAKE) build
	ARCH=arm $(MAKE) build
	ARCH=arm64 $(MAKE) build

# new testing method
test: build
	hzn dev service start -S
	@echo 'Testing service...'
	../../../tools/curlServiceTest.sh $(SERVICE_NAME) $(MATCH) $(TIME_OUT) && \
		{ hzn dev service stop; \
		echo "*** Service test succeeded! ***"; } || \
		{ hzn dev service stop; \
		echo "*** Service test failed! ***"; \
		false; }

# Create/update the metadata in the exchange for this service
publish-service:
	hzn exchange service publish -f horizon/service.definition.json

build-test-publish: build test publish-service

publish-all-arches:
	ARCH=amd64 $(MAKE) build-test-publish
	ARCH=arm $(MAKE) build-test-publish
	ARCH=arm64 $(MAKE) build-test-publish

# target for script - overwrite and pull insitead of push docker image
publish-service-overwrite:
	hzn exchange service publish -O -P -f horizon/service.definition.json

# new target for icp exchange to run on startup to publish only
publish-only:
	ARCH=amd64 $(MAKE) publish-service-overwrite
	ARCH=arm $(MAKE) publish-service-overwrite
	ARCH=arm64 $(MAKE) publish-service-overwrite

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
