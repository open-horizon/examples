# Make targets for building the IBM example helloworld edge service

# This imports the variables from horizon/hzn.json. You can ignore these lines, but do not remove them.
-include horizon/.hzn.json.tmp.mk

# Default ARCH to the architecture of this machines (as horizon/golang describes it)
export ARCH ?= $(shell hzn architecture)

# Configurable parameters passed to serviceTest.sh in "test" target
export MATCH:='Hello'
export TIME_OUT:=10

# Build the docker image for the current architecture
build:
	docker build -t $(DOCKER_IMAGE_BASE)_$(ARCH):$(SERVICE_VERSION) -f ./Dockerfile.$(ARCH) .

# Build the docker image for 3 architectures
build-all-arches:
	ARCH=amd64 $(MAKE) build
	ARCH=arm $(MAKE) build
	ARCH=arm64 $(MAKE) build

# Run and verify the service
test: build
	hzn dev service start -S
	@echo 'Testing service...'
	../../../tools/serviceTest.sh $(SERVICE_NAME) $(MATCH) $(TIME_OUT) && \
		{ hzn dev service stop; \
		echo "*** Service test succeeded! ***"; } || \
		{ hzn dev service stop; \
		echo "*** Service test failed! ***"; }

# Publish the service to the Horizon Exchange for the current architecture
publish-service:
	hzn exchange service publish -f horizon/service.definition.json

publish-pattern:
	hzn exchange pattern publish -f horizon/pattern.json

# Build, run and verify, if test succeeds then publish (for the current architecture)
build-test-publish: build publish-service

# Build/test/publish the service to the Horizon Exchange for 3 architectures and publish a deployment pattern for those architectures
publish-all-arches:
	ARCH=amd64 $(MAKE) build-test-publish
	ARCH=arm $(MAKE) build-test-publish
	ARCH=arm64 $(MAKE) build-test-publish
	hzn exchange pattern publish -f horizon/pattern-all-arches.json

# target for script - overwrite and pull insitead of push docker image
publish-service-overwrite:
	hzn exchange service publish -O -P -f horizon/service.definition.json

# target to publish new input.json file to mms
publish-mms-object:
	hzn mms object publish -m object.json -f input.json

# target to list mms object
list-mms-object:
	hzn mms object list -t json -i input.json -d

# target to delete input.json file in mms
delete-mms-object:
	hzn mms object delete -t json --id input.json

# new target for exchange publish script to run on startup to publish only
publish-only:
	ARCH=amd64 $(MAKE) publish-service-overwrite
	ARCH=arm $(MAKE) publish-service-overwrite
	ARCH=arm64 $(MAKE) publish-service-overwrite
	hzn exchange pattern publish -f horizon/pattern-all-arches.json

clean:
	-docker rmi $(DOCKER_IMAGE_BASE)_$(ARCH):$(SERVICE_VERSION) 2> /dev/null || :

clean-all-archs:
	ARCH=amd64 $(MAKE) clean
	ARCH=arm $(MAKE) clean
	ARCH=arm64 $(MAKE) clean

# This imports the variables from horizon/hzn.cfg. You can ignore these lines, but do not remove them.
horizon/.hzn.json.tmp.mk: horizon/hzn.json
	@ hzn util configconv -f $< > $@

.PHONY: build build-all-arches test publish-service build-test-publish publish-all-arches clean clean-all-archs
