# This imports the variables from horizon/hzn.json. You can ignore these lines, but do not remove them.
-include horizon/.hzn.json.tmp.mk

# To build for an arch different from the current system, set this env var to 1 of the values in the comment above
export ARCH ?= $(shell hzn architecture)

default: build publish-service-overwrite

build:
	docker build -t $(DOCKER_IMAGE_BASE)_$(ARCH):$(SERVICE_VERSION) -f ./Dockerfile.$(ARCH) .

build-all-arches:
	ARCH=amd64 $(MAKE) build
	ARCH=arm $(MAKE) build
	ARCH=arm64 $(MAKE) build

publish-service:
	hzn exchange service publish -f horizon/service.definition.json

# target for script - overwrite and pull insitead of push docker image
publish-service-overwrite:
	hzn exchange service publish -O -P -f horizon/service.definition.json

# new target for new exchange to run on startup to publish only
publish-only:
	ARCH=amd64 $(MAKE) publish-service-overwrite
	ARCH=arm $(MAKE) publish-service-overwrite
	ARCH=arm64 $(MAKE) publish-service-overwrite

publish-all-arches:
	ARCH=amd64 $(MAKE) build-test-publish
	ARCH=arm $(MAKE) build-test-publish
	ARCH=arm64 $(MAKE) build-test-publish

clean:
	-docker rmi $(DOCKER_IMAGE_BASE)_$(ARCH):$(SERVICE_VERSION) 2> /dev/null || :

clean-all-arches:
	ARCH=amd64 $(MAKE) clean
	ARCH=arm $(MAKE) clean
	ARCH=arm64 $(MAKE) clean

# This imports the variables from horizon/hzn.cfg. You can ignore these lines, but do not remove them.
horizon/.hzn.json.tmp.mk: horizon/hzn.json
	@ hzn util configconv -f $< > $@
