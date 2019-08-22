# This imports the variables from horizon/hzn.json. You can ignore these lines, but do not remove them.
-include horizon/.hzn.json.tmp.mk

# Default ARCH to the architecture of this machines (as horizon/golang describes it). Can be overridden.
export ARCH ?= $(shell hzn architecture)

default: build publish-service publish-pattern

# Build the docker image for the current architecture
build:
	docker build -t $(DOCKER_IMAGE_BASE)_$(ARCH):$(SERVICE_VERSION) -f ./Dockerfile.$(ARCH) .

# Publish the service to the Horizon Exchange for the current architecture
publish-service:
	hzn exchange service publish -f horizon/service.definition.json

# Publish the pattern to the Horizon Exchange for the current architecture
publish-pattern:
	hzn exchange pattern publish -f horizon/pattern.json

# new target for new exchange to run on startup to publish only
publish-only:
	hzn exchange service publish -O -P -f horizon/service.definition.json
	hzn exchange pattern publish -f horizon/pattern.json

clean:
	-docker rmi $(DOCKER_IMAGE_BASE)_$(ARCH):$(SERVICE_VERSION) 2> /dev/null || :

# This imports the variables from horizon/hzn.cfg. You can ignore these lines, but do not remove them.
horizon/.hzn.json.tmp.mk: horizon/hzn.json
	@ hzn util configconv -f $< > $@
