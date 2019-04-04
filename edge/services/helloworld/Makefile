# Make targets for building the IBM example helloworld edge service

# This imports the variables from horizon/hzn.cfg. You can ignore these lines, but do not remove them.
-include horizon/.hzn.cfg.tmp.mk

# Build the docker image for the current architecture
build:
	docker build -t $(DOCKER_IMAGE_BASE)_$(ARCH):$(SERVICE_VERSION) -f ./Dockerfile.$(ARCH) .

# Build the docker image for 3 architectures
build-all-arches:
	ARCH=amd64 $(MAKE) build
	ARCH=arm $(MAKE) build
	ARCH=arm64 $(MAKE) build

# Soon we can use 'hzn service log ...' for all platforms
ifeq ($(shell uname -s),Darwin)
GET_SVC_OUTPUT = docker logs `docker ps -q --filter name=helloworld`
else
GET_SVC_OUTPUT = grep helloworld /var/log/syslog
endif

# Run and verify the service
test: build
	hzn dev service start -S
	@echo 'Testing service...'
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
	hzn exchange service publish -f horizon/service.definition.json -k $(HZN_PRIVATE_KEY_FILE) -K $(HZN_PUBLIC_KEY_FILE)

# Build, run and verify, if test succeeds then publish
build-test-publish: build test publish-service

# Publish the service to the Horizon Exchange for 3 architectures and publish a deployment pattern for those architectures
publish-all-arches:
	ARCH=amd64 $(MAKE) publish-service
	ARCH=arm $(MAKE) publish-service
	ARCH=arm64 $(MAKE) publish-service
	hzn exchange pattern publish -f horizon/pattern/pattern-helloworld-all-arches.json -p pattern-helloworld

clean:
	-docker rmi $(DOCKER_IMAGE_BASE)_$(ARCH):$(SERVICE_VERSION) 2> /dev/null || :

clean-all-archs:
	ARCH=amd64 $(MAKE) clean
	ARCH=arm $(MAKE) clean
	ARCH=arm64 $(MAKE) clean

# Just to experiment with make
testmake:
	@echo "testmake: $(DOCKER_IMAGE_BASE)"

# This imports the variables from horizon/hzn.cfg. You can ignore these lines, but do not remove them.
horizon/.hzn.cfg.tmp.mk: horizon/hzn.cfg
	@sed -e '/^export /d' -e "s/'//g" -e 's/"//g' $< > $@

.PHONY: build clean
