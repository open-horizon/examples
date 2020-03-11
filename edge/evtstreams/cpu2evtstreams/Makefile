# Make targets for building the IBM sample cpu2evtstreams edge service

# This imports the variables from horizon/hzn.json. You can ignore these lines, but do not remove them.
-include horizon/.hzn.json.tmp.mk

# Default ARCH to the architecture of this machines (as horizon/golang describes it). Can be overridden.
export ARCH ?= $(shell hzn architecture)

# Configurable parameters passed to serviceTest.sh in "test" target.
export MATCH:='{"nodeID":'
export TIME_OUT:=60

# Build the docker image for the current architecture
build:
ifeq (,$(findstring amd64,$(ARCH)))
	mkdir -p tmp/$(ARCH)
	cp ../../../tools/kafkacat/$(ARCH)/*.rsa.pub ../../../tools/kafkacat/$(ARCH)/kafkacat-*.apk tmp/$(ARCH)
endif
	docker build -t $(DOCKER_IMAGE_BASE)_$(ARCH):$(SERVICE_VERSION) -f ./Dockerfile.$(ARCH) .
ifeq (,$(findstring amd64,$(ARCH)))
	rm -f tmp/$(ARCH)/*.rsa.pub tmp/$(ARCH)/kafkacat-*.apk
endif

# Build the docker image for 3 architectures
build-all-arches:
	ARCH=amd64 $(MAKE) build
	ARCH=arm $(MAKE) build
	ARCH=arm64 $(MAKE) build

# Target for travis to test new PRs
test-all-arches:
	ARCH=amd64 $(MAKE) test
	ARCH=arm $(MAKE) test
	ARCH=arm64 $(MAKE) test

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

# Publish the service to the Horizon Exchange for the current architecture
publish-service:
	hzn exchange service publish -O -f horizon/service.definition.json

# Target for travis to publish service and pattern after PR is merged  
publish:
	ARCH=amd64 $(MAKE) publish-service
	ARCH=arm $(MAKE) publish-service
	ARCH=arm64 $(MAKE) publish-service
	hzn exchange pattern publish -f horizon/pattern-all-arches.json

# Build, run and verify, if test succeeds then publish (for the current architecture)
#build-test-publish: build test publish-service   # <- currently test does not always work on mac
build-test-publish: build test publish-service

# Build/test/publish the service to the Horizon Exchange for 3 architectures and publish a deployment pattern for those architectures
publish-all-arches:
	ARCH=amd64 $(MAKE) build-test-publish
	ARCH=arm $(MAKE) build-test-publish
	ARCH=arm64 $(MAKE) build-test-publish
	hzn exchange pattern publish -f horizon/pattern-all-arches.json

# target for script - overwrite and pull insitead of push docker image
publish-service-overwrite:
	hzn exchange service publish -O -P -f horizon/service.definition.json

# Publish Service Policy target for exchange publish script
publish-service-policy:
	hzn exchange service addpolicy -f horizon/service_policy.json $(HZN_ORG_ID)/$(SERVICE_NAME)_$(SERVICE_VERSION)_$(ARCH)

# Publish Business Policy target for exchange publish script
publish-business-policy:
	hzn exchange business addpolicy -f horizon/business_policy.json --org=$(PUBLIC_ORG) $(SERVICE_NAME)_$(SERVICE_VERSION)

# new target for icp exchange to run on startup to publish only
publish-only:
	ARCH=amd64 $(MAKE) publish-service-overwrite
	ARCH=amd64 $(MAKE) publish-service-policy
	ARCH=arm $(MAKE) publish-service-overwrite
	ARCH=arm $(MAKE) publish-service-policy
	ARCH=arm64 $(MAKE) publish-service-overwrite
	ARCH=arm64 $(MAKE) publish-service-policy
	hzn exchange pattern publish -f horizon/pattern-all-arches.json

# Create a topic in Event Streams for cpu2evtstreams to send data to
EVTSTREAMS_TOPIC ?= cpu2evtstreams
evtstreams-topic:
	: $${EVTSTREAMS_API_KEY:?} $${EVTSTREAMS_ADMIN_URL:?}
	curl -sS -w %{http_code} -H 'Content-Type: application/json' -H "X-Auth-Token: $$EVTSTREAMS_API_KEY" -d "{ \"name\": \"$(EVTSTREAMS_TOPIC)\", \"partitions\": 2 }" $$EVTSTREAMS_ADMIN_URL/admin/topics

# List the topics currently in Event Streams
evtstreams-topic-list:
	@: $${EVTSTREAMS_API_KEY:?} $${EVTSTREAMS_ADMIN_URL:?}
	@curl -sS -H "X-Auth-Token: $$EVTSTREAMS_API_KEY" $$EVTSTREAMS_ADMIN_URL/admin/topics | jq -r ".[] | .name"

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
