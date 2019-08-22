# This imports the variables from horizon/hzn.json. You can ignore these lines, but do not remove them.
-include horizon/.hzn.json.tmp.mk

# To build for an arch different from the current system, set this env var to 1 of the values in the comment above
export ARCH ?= $(shell hzn architecture)

default: build publish-service

# Build the docker image for the current architecture
build:
	docker build -t $(DOCKER_IMAGE_BASE)_$(ARCH):$(SERVICE_VERSION) -f ./Dockerfile.$(ARCH) .

# Create/update the metadata in the exchange for this service
publish-service:
	hzn exchange service publish -f horizon/service.definition.json

# new target for new exchange to run on startup to publish only
publish-only:
	hzn exchange service publish -O -P -f horizon/service.definition.json

clean:
	-docker rmi $(DOCKER_IMAGE_BASE)_$(ARCH):$(SERVICE_VERSION) 2> /dev/null || :

# This imports the variables from horizon/hzn.cfg. You can ignore these lines, but do not remove them.
horizon/.hzn.json.tmp.mk: horizon/hzn.json
	@ hzn util configconv -f $< > $@
