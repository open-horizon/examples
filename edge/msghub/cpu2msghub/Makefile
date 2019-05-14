# Make targets for building the IBM sample cpu2msghub edge service

# This imports the variables from horizon/hzn.json. You can ignore these lines, but do not remove them.
-include horizon/.hzn.json.tmp.mk

# Default ARCH to the architecture of this machines (as horizon/golang describes it). Can be overridden.
export ARCH ?= $(shell hzn architecture)

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

# Soon we can use 'hzn service log ...' for all platforms
ifeq ($(shell uname -s),Darwin)
GET_SVC_OUTPUT = docker logs `docker ps -q --filter name=cpu2msghub`
else
GET_SVC_OUTPUT = $(SYSLOG_SUDO) grep cpu2msghub[[] /var/log/syslog
endif

# Run and verify the service
test: build
	hzn dev service start -S
	@echo 'Testing service...'
	sleep 20  # may have to increase this time depending on the speed of your machine
	@if $(GET_SVC_OUTPUT) | grep -q 'echo {"nodeID":'; then \
	  hzn dev service stop; \
	  echo 'Service test succeeded.'; \
	else \
	  hzn dev service stop; \
	  echo 'Service test failed!'; \
	  false; \
	fi

# Publish the service to the Horizon Exchange for the current architecture
publish-service:
	hzn exchange service publish -f horizon/service.definition.json

# Build, run and verify, if test succeeds then publish (for the current architecture)
build-test-publish: build test publish-service

# Build/test/publish the service to the Horizon Exchange for 3 architectures and publish a deployment pattern for those architectures
publish-all-arches:
	ARCH=amd64 $(MAKE) build-test-publish
	ARCH=arm $(MAKE) build-test-publish
	ARCH=arm64 $(MAKE) build-test-publish
	hzn exchange pattern publish -f horizon/pattern-all-arches.json

# Create a topic in msg hub for cpu2msghub to send data to
msghub-topic:
	: $${MSGHUB_API_KEY:?} $${MSGHUB_ADMIN_URL:?}
	MSGHUB_TOPIC=cpu2msghub curl -sS -w %{http_code} -H 'Content-Type: application/json' -H "X-Auth-Token: $(MSGHUB_API_KEY)" -d "{ \"name\": \"$$MSGHUB_TOPIC\", \"partitions\": 2 }" $(MSGHUB_ADMIN_URL)/admin/topics

# List the topics currently in msg hub
msghub-topic-list:
	@: $${MSGHUB_API_KEY:?} $${MSGHUB_ADMIN_URL:?}
	@curl -sS -H "X-Auth-Token: $$MSGHUB_API_KEY" $$MSGHUB_ADMIN_URL/admin/topics | jq -r ".[] | .name"

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
