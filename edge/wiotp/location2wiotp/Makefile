#
# Makefile for location container
#

export VERSION ?= 2.0.7
IMAGE ?= location2wiotp

# The docker network used for testing (must match the gps container network)
DOCKER_TEST_NETWORK ?= gps-test

# Transform the machine arch into some standard values: "arm", "arm64", or "amd64". Note: ppc64le already returns ppc64le, so we do not need to change that.
SYSTEM_ARCH := $(shell uname -m | sed -e 's/aarch64.*/arm64/' -e 's/x86_64.*/amd64/' -e 's/armv.*/arm/')

# To build for an arch different from the current system, set this env var to 1 of the values in the comment above
export ARCH ?= $(SYSTEM_ARCH)

# Construct our Docker registry path from the architecture
REGISTRY ?= openhorizon
REG_PATH := $(REGISTRY)/$(ARCH)
IMAGE_PATH = $(REG_PATH)_$(IMAGE):$(VERSION)

# Make will search these two relative paths for input and target files
VPATH=bin

# Running make with no target will build the production image
.PHONY: default
default: $(IMAGE)-$(VERSION)

# Omitting the image architecture will default to the local architecture
.PHONY: $(IMAGE)-$(VERSION)
$(IMAGE)-$(VERSION): .$(IMAGE)-$(VERSION)-$(ARCH)

# This is the production image target
.$(IMAGE)-$(VERSION)-$(ARCH): Dockerfile.$(ARCH) *.sh *.pem
	docker build -t $(IMAGE_PATH) -f Dockerfile.$(ARCH) .
	mkdir -p bin
	touch bin/.$(IMAGE)-$(VERSION)-$(ARCH)

# This is the "docker save" tarball, suitable for "docker load"
bin/$(IMAGE)-$(VERSION)-$(ARCH).tar: .$(IMAGE)-$(VERSION)-$(ARCH)
	mkdir -p bin
	docker save -o bin/$(IMAGE)-$(VERSION)-$(ARCH).tar $(IMAGE_PATH)

# Push the production image for this architecture to the resistry, and verify
.PHONY: docker-push
docker-push: .$(IMAGE)-$(VERSION)-$(ARCH)
	docker push $(IMAGE_PATH)

.PHONY: docker-verify
docker-verify:
	docker version --format '{{.Server.Version}}'
	docker pull $(IMAGE_PATH) | grep '^Digest: '

.PHONY: exchange-publish
exchange-publish: .$(IMAGE)-$(VERSION)-$(ARCH)
	: $${HZN_ORG_ID:?} $${HZN_EXCHANGE_USER_AUTH:?} $${PRIVATE_KEY_FILE:?}   # this verifies these env vars are set
	envsubst < horizon/$(IMAGE)-template.json > bin/$(IMAGE).json
	hzn exchange workload publish -k $$PRIVATE_KEY_FILE -f bin/$(IMAGE).json

# Create the docker network
.PHONY: network
network:
	-docker network create $(DOCKER_TEST_NETWORK)

# Run the image in developer mode on the local host
.PHONY: develop
develop: network .$(IMAGE)-$(VERSION)-$(ARCH)
	-docker rm -f $(IMAGE)
	docker run --rm -it --volume=`pwd`:/outside -e HZN_AGREEMENTID=$(HZN_AGREEMENTID) -e HZN_HASH=$(HZN_HASH) -e HZN_EXCHANGE_URL=$(HZN_EXCHANGE_URL) --net=$(DOCKER_TEST_NETWORK) --net-alias=$(IMAGE) --name $(IMAGE) $(IMAGE_PATH) /bin/sh

# Run the container as a daemon
.PHONY: daemon
daemon: network .$(IMAGE)-$(VERSION)-$(ARCH)
	-docker rm -f $(IMAGE)
	docker run -d --volume=`pwd`:/outside -e HZN_AGREEMENTID=$(HZN_AGREEMENTID) -e HZN_HASH=$(HZN_HASH) -e HZN_EXCHANGE_URL=$(HZN_EXCHANGE_URL) --net=$(DOCKER_TEST_NETWORK) --net-alias=$(IMAGE) --name $(IMAGE) $(IMAGE_PATH)

# Listen to the remote MQTT broker to see location messages from this container
# This configuration will echo all messages received from the central MQTT.
.PHONY: listener
listener:
	docker exec -t $(IMAGE) /test.sh

# Long duration test target, emitting detailed stats after every N messages.
# Since the default configuration is a message every 10 seconds, then with:
#     N=180 (i.e., 6 * 30, or 6 messages per minute * 30 minutes)
# the "test" target should emit a status message approximately every half hour.
# In practice messages are received about every 11 or 12 seconds, so N is
# set to 150 below and it tends to give stats about every half hour.
.PHONY: test
test:
	docker exec -t $(IMAGE) /test.sh 150 q

# Clean everything
.PHONY: clean
clean:
	-rm -rf bin bin/location.json
	-docker rm -f $(IMAGE)
	-docker rmi $(IMAGE)-$(ARCH):latest $(IMAGE)-$(ARCH):$(VERSION) $(IMAGE_PATH)
	-docker network rm $(DOCKER_TEST_NETWORK)

