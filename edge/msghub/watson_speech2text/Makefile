# This imports the variables from horizon/hzn.json. You can ignore these lines, but do not remove them.
-include horizon/.hzn.json.tmp.mk

# Default ARCH to the architecture of this machines (as horizon/golang describes it). Can be overridden.
export ARCH ?= $(shell hzn architecture)

# Configurable parameters passed to serviceTest.sh in "test" target
export MATCH:='envs:'
export TIME_OUT:=30

# Build the docker image for the current architecture
build:
	docker build -t $(DOCKER_IMAGE_BASE)_$(ARCH):$(SERVICE_VERSION) -f ./Dockerfile.$(ARCH) .

# new testing method
test: build
	hzn dev service start -S
	@echo 'Testing service...'
	./serviceTest.sh $(SERVICE_NAME) $(MATCH) $(TIME_OUT) && \
		{ hzn dev service stop; \
		echo "*** Service test succeeded! ***"; } || \
		{ hzn dev service stop; \
		echo "*** Service test failed! ***"; }

# Publish the service to the Horizon Exchange for the current architecture
publish-service:
	hzn exchange service publish -f horizon/service.definition.json

publish-service-overwrite:
	hzn exchange service publish -O -P -f horizon/service.definition.json

# Publish the pattern to the Horizon Exchange for the current architecture
publish-pattern:
	hzn exchange pattern publish -f horizon/pattern.json

# new target for new exchange to run on startup to publish only
publish-only: publish-service-overwrite publish-pattern

# build, test, and publish
publish: test publish-service publish-pattern

clean:
	-docker rmi $(DOCKER_IMAGE_BASE)_$(ARCH):$(SERVICE_VERSION) 2> /dev/null || :

# This imports the variables from horizon/hzn.cfg. You can ignore these lines, but do not remove them.
horizon/.hzn.json.tmp.mk: horizon/hzn.json
	@ hzn util configconv -f $< > $@
