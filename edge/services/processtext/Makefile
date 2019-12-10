# This imports the variables from horizon/hzn.json. You can ignore these lines, but do not remove them.
-include horizon/.hzn.json.tmp.mk

# Default ARCH to the architecture of this machines (as horizon/golang describes it). Can be overridden.
export ARCH ?= $(shell hzn architecture)

# Configurable parameters passed to serviceTest.sh in "test" target
export MATCH:='[OVA]'
export TIME_OUT:=60

build:
	docker build -t $(DOCKER_IMAGE_BASE)_$(ARCH):$(SERVICE_VERSION) -f ./Dockerfile.$(ARCH) .

# Run and verify the service
test: build
	hzn dev service start -S
	@echo 'Testing service...'
	../../../tools/serviceTest.sh $(SERVICE_NAME) $(MATCH) $(TIME_OUT) && \
		{ hzn dev service stop; \
		echo "*** Service test succeeded! ***"; } || \
		{ hzn dev service stop; \
		echo "*** Service test failed! ***"; \
		false; }

publish-service:
	hzn exchange service publish -O -f horizon/service.definition.json

# target for script - overwrite and pull insitead of push docker image
publish-service-overwrite:
	 hzn exchange service publish -O -P -f horizon/service.definition.json

publish-pattern:
	hzn exchange pattern publish -f horizon/pattern.json

publish-only: publish-service-overwrite publish-pattern

clean:
	-docker rmi $(DOCKER_IMAGE_BASE)_$(ARCH):$(SERVICE_VERSION) 2> /dev/null || :

# This imports the variables from horizon/hzn.cfg. You can ignore these lines, but do not remove them.
horizon/.hzn.json.tmp.mk: horizon/hzn.json
	@ hzn util configconv -f $< > $@

