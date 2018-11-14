# Transform the machine arch into some standard values: "arm", "arm64", or "amd64"
SYSTEM_ARCH := $(shell uname -m | sed -e 's/aarch64.*/arm64/' -e 's/x86_64.*/amd64/' -e 's/armv.*/arm/')

# To build for an arch different from the current system, set this env var to 1 of the values in the comment above
export ARCH ?= $(SYSTEM_ARCH)

# These variables can be overridden from the environment
export CPU2MSGHUB_NAME ?= cpu2msghub
export CPU2MSGHUB_VERSION ?= 1.2.5
# the required services - make sure these versions are up to date!
export CPU_NAME ?= cpu
export CPU_VERSION ?= 1.2.2
export GPS_NAME ?= gps
export GPS_VERSION ?= 2.0.5
DOCKER_NAME ?= $(ARCH)_$(CPU2MSGHUB_NAME)
export DOCKER_HUB_ID ?= openhorizon
export MYDOMAIN ?= github.com.open-horizon.examples

default: all

all: build run check

build:
ifeq (,$(findstring amd64,$(ARCH)))
	mkdir -p tmp/$(ARCH)
	cp ../../../tools/kafkacat/$(ARCH)/*.rsa.pub ../../../tools/kafkacat/$(ARCH)/kafkacat-*.apk tmp/$(ARCH)
endif
	docker build -t $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(CPU2MSGHUB_VERSION) -f ./Dockerfile.$(ARCH) .
ifeq (,$(findstring amd64,$(ARCH)))
	rm -f tmp/$(ARCH)/*.rsa.pub tmp/$(ARCH)/kafkacat-*.apk
endif

run:
	-docker rm -f $(DOCKER_NAME) 2> /dev/null || :
	docker run -d --name $(DOCKER_NAME) -e SAMPLE_SIZE=5 -e SAMPLE_INTERVAL=1 -e MOCK=true -e VERBOSE=1 -e PUBLISH=false --volume `pwd`:/outside $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(CPU2MSGHUB_VERSION)

check:
	docker logs -f $(DOCKER_NAME)

stop:
	-docker rm -f $(DOCKER_NAME) 2> /dev/null || :

# Targets for using the hzn dev environment
hznbuild: build
	hzn dev dependency fetch -p ../../services/cpu_percent/horizon
	hzn dev dependency fetch -p ../../services/gps/horizon
	hzn dev service verify

hznstart: hznbuild
	hzn dev service start

hznstop:
	hzn dev service stop

# Push the docker image to the registry. You must have write access to the docker hub openhorizon user
docker-push: build
	docker push $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(CPU2MSGHUB_VERSION)

# Create/update the metadata in the exchange for this service
publish-service: build publish-service-only
publish-service-only:
	: $${HZN_EXCHANGE_USER_AUTH:?} $${PRIVATE_KEY_FILE:?} $${PUBLIC_KEY_FILE:?}   # this verifies these env vars are set
	hzn exchange service publish -k $$PRIVATE_KEY_FILE -K $$PUBLIC_KEY_FILE -f horizon/service.definition.json

# Create/update the deployment pattern for a set of edge nodes
publish-pattern:
	: $${HZN_ORG_ID:?} $${HZN_EXCHANGE_USER_AUTH:?}   # this verifies these env vars are set
	hzn exchange pattern publish -p $$CPU2MSGHUB_NAME -f horizon/pattern/cpu2msghub.json

clean:
	-docker rm -f $(DOCKER_NAME) 2> /dev/null || :
	-docker rmi $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(CPU2MSGHUB_VERSION) 2> /dev/null || :

.PHONY: default all build run check stop hznbuild hznstart hznstop publish-service publish-service-only publish-pattern clean
