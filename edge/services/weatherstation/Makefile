# Transform the machine arch into some standard values: "arm", "arm64", or "amd64". Note: ppc64le already returns ppc64le, so we do not need to change that.
SYSTEM_ARCH := $(shell uname -m | sed -e 's/aarch64.*/arm64/' -e 's/x86_64.*/amd64/' -e 's/armv.*/arm/')

# To build for an arch different from the current system, set this env var to 1 of the values in the comment above
ARCH ?= $(SYSTEM_ARCH)

# These variables can be overridden from the environment
export PWSMS_NAME ?= pwsms
export PWSMS_VERSION ?= 1.1.0
DOCKER_NAME ?= $(ARCH)_$(PWSMS_NAME)
DOCKER_HUB_ID ?= openhorizon
export MYDOMAIN ?= internetofthings.ibmcloud.com
PORT_NUM ?= 8357

# PWS Specific variables for dev (replace with your own, or comment, as these override ENV Vars)
PWS_ST_TYPE ?= $(PWS_ST_TYPE)
PWS_MODEL ?= $(PWS_MODEL)
PWS_WU_ID ?= $(PWS_WU_ID)
PWS_WU_KEY ?= $(PWS_WU_KEY)
PWS_WU_RPDF ?= $(PWS_WU_RPDF)


default: all

all: build run

build:
	docker build -t $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(PWSMS_VERSION) -f ./Dockerfile.$(ARCH) .

build-squash:
	docker build --squash -t $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(PWSMS_VERSION) -f ./Dockerfile.$(ARCH) .

run:
	-docker rm -f $(DOCKER_NAME) 2> /dev/null || :
	docker run -d --name $(DOCKER_NAME) --publish=$(PORT_NUM):$(PORT_NUM) --privileged -e PWS_ST_TYPE=$(PWS_ST_TYPE) -e PWS_MODEL=$(PWS_MODEL) -e PWS_WU_ID=$(PWS_WU_ID) -e PWS_WU_KEY=$(PWS_WU_KEY) -e PWS_WU_RPDF=$(PWS_WU_RPDF) -v `pwd`:/outside $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(PWSMS_VERSION)

dev:
	-docker stop $(IMAGE) && docker rm $(IMAGE)
	docker run -it --name $(DOCKER_NAME) --publish=$(PORT_NUM):$(PORT_NUM) --privileged -e PWS_ST_TYPE=$(PWS_ST_TYPE) -e PWS_MODEL=$(PWS_MODEL) -e PWS_WU_ID=$(PWS_WU_ID) -e PWS_WU_KEY=$(PWS_WU_KEY) -e PWS_WU_RPDF=$(PWS_WU_RPDF) -v `pwd`:/outside $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(PWSMS_VERSION) /bin/sh

check:
	curl -s localhost:$(PORT_NUM)/v1/weather | jq

stop:
	-docker rm -f $(DOCKER_NAME) 2> /dev/null || :

# To publish you must have write access to the docker hub DOCKER_HUB_ID user
publish: build
	docker push $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(PWSMS_VERSION)

# Create/update the metadata in the exchange for this microservice
exchange-publish: build
	hzn exchange microservice publish -k $$PRIVATE_KEY_FILE -f horizon/microservice.definition.json

exchange-publish-squash: build-squash
	hzn exchange microservice publish -k $$PRIVATE_KEY_FILE -f horizon/microservice.definition.json

clean:
	-docker rm -f $(DOCKER_NAME) 2> /dev/null || :
	-docker rmi $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(PWSMS_VERSION) 2> /dev/null || :

.PHONY: default all build build-squash dev run check publish exchange-publish exchange-publish-squash clean
