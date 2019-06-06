# Make targets for building the IBM example helloworld edge service

# This imports the variables from horizon/hzn.json. You can ignore these lines, but do not remove them.
-include horizon/.hzn.json.tmp.mk

# Default ARCH to the architecture of this machines (as horizon/golang describes it)
export ARCH ?= $(shell hzn architecture)

# Build the docker image for the current architecture
build:
	docker build -t $(DOCKER_IMAGE_BASE)_$(ARCH):$(SERVICE_VERSION) -f ./Dockerfile.$(ARCH) .

# Build the docker image for 3 architectures
build-all-arches:
	ARCH=amd64 $(MAKE) build
	ARCH=arm $(MAKE) build
	ARCH=arm64 $(MAKE) build

# In travis can only read syslog via sudo
ifeq ($(TRAVIS),true)
SYSLOG_SUDO = sudo
else
SYSLOG_SUDO =
endif
# Soon we can use 'hzn service log ...' for all platforms
ifeq ($(shell uname -s),Darwin)
GET_SVC_OUTPUT = docker logs `docker ps -q --filter name=helloworld`
else
GET_SVC_OUTPUT = $(SYSLOG_SUDO) grep helloworld /var/log/syslog
endif

# Run and verify the service
test: build
	hzn dev service start -S
	@echo 'Testing service...'
	sleep 1
	@if $(GET_SVC_OUTPUT) | grep -q 'says: Hello'; then \
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

# new target for icp exchange to run on startup to publish only
publish-only:
	ARCH=amd64 $(MAKE) publish-service
	ARCH=arm $(MAKE) publish-service
	ARCH=arm64 $(MAKE) publish-service

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
