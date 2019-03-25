# Transform the machine arch into some standard values: "arm", "arm64", or "amd64"
SYSTEM_ARCH := $(shell uname -m | sed -e 's/aarch64.*/arm64/' -e 's/x86_64.*/amd64/' -e 's/armv.*/arm/')

# To build for an arch different from the current system, set this env var to 1 of the values in the comment above
export ARCH ?= $(SYSTEM_ARCH)

# Set these variables as appropriate
export HZN_ORG_ID ?= IBM
export SERVICE_NAME ?= helloworld
export SERVICE_URL ?= github.com.open-horizon.examples.$(SERVICE_NAME)
export SERVICE_VERSION ?= 1.0.0
export DOCKER_HUB_ID ?= openhorizon

# Derived variables
export HZN_PATTERN ?= $(HZN_ORG_ID)/pattern-$(SERVICE_NAME)
export DOCKER_IMAGE ?= $(DOCKER_HUB_ID)/$(ARCH)_$(SERVICE_NAME):$(SERVICE_VERSION)

default: build hznstart

build:
	docker build -t $(DOCKER_IMAGE) -f ./Dockerfile.$(ARCH) .

hznstart:
	hzn dev service start -S

hznstop:
	hzn dev service stop

# Create/update the metadata in the exchange for this service
build-publish-service: build publish-service-only

publish-service:
	: $${HZN_EXCHANGE_USER_AUTH:?} $${PRIVATE_KEY_FILE:?} $${PUBLIC_KEY_FILE:?}   # verify these env vars are set
	hzn exchange service publish -k $$PRIVATE_KEY_FILE -K $$PUBLIC_KEY_FILE -f horizon/service.definition.json

# Create/update the deployment pattern for a set of edge nodes
publish-pattern:
	: $${HZN_EXCHANGE_USER_AUTH:?}   # verify these env vars are set
	hzn exchange pattern publish -p $${HZN_PATTERN#*/} -f horizon/pattern/pattern-helloworld.json

register-node:
	: $${HZN_EXCHANGE_USER_AUTH:?} $${HZN_EXCHANGE_NODE_AUTH:?}   # verify these env vars are set
	hzn register -u $(HZN_EXCHANGE_USER_AUTH) -n $(HZN_EXCHANGE_NODE_AUTH) $(HZN_ORG_ID) $(HZN_PATTERN)   # -f horizon/userinput.json

testmake:
	echo $$DOCKER_IMAGE

clean:
	-docker rm -f $(SERVICE_NAME) 2> /dev/null || :
	-docker rmi $(DOCKER_IMAGE) 2> /dev/null || :

.PHONY: default build hznstart hznstop publish-service publish-service-only publish-pattern register-node testmake clean
