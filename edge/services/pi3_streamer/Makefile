# Transform the machine arch into some standard values: "arm", "arm64", or "amd64"
SYSTEM_ARCH := $(shell uname -m | sed -e 's/aarch64.*/arm64/' -e 's/x86_64.*/amd64/' -e 's/armv.*/arm/')

# To build for an arch different from the current system, set this env var to 1 of the values in the comment above
export ARCH ?= $(SYSTEM_ARCH)

# These variables can be overridden from the environment
export PI3STREAMER_NAME ?= pi3streamer
export PI3STREAMER_VERSION ?= 1.0.0
DOCKER_NAME ?= $(ARCH)_$(PI3STREAMER_NAME)
DOCKER_HUB_ID ?= openhorizon
export MYDOMAIN ?= internetofthings.ibmcloud.com
PORT_NUM ?= 8080
VERT_FLIP ?= 1

default: all

all: build run

build:
	docker build -t $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(PI3STREAMER_VERSION) -f ./Dockerfile.$(ARCH) .

run:
	-docker rm -f $(DOCKER_NAME) 2> /dev/null || :
	docker run -d --name $(DOCKER_NAME) --publish=$(PORT_NUM):$(PORT_NUM) --privileged -v /dev:/dev:ro -e VERT_FLIP=$(VERT_FLIP) -e RESOLUTION="-x 1280 -y 720" --volume `pwd`:/outside $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(PI3STREAMER_VERSION)

dev:
	-docker rm -f $(DOCKER_NAME) 2> /dev/null || :
	docker run -it --name $(DOCKER_NAME) --publish=$(PORT_NUM):$(PORT_NUM) --privileged -v /dev:/dev:ro -e VERT_FLIP=$(VERT_FLIP) -e RESOLUTION="-x 1280 -y 720" --volume `pwd`:/outside $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(PI3STREAMER_VERSION) /bin/sh

check:
	curl -sSL http://localhost:$(PORT_NUM)/static_simple.html | grep MJPG

stop:
	-docker rm -f $(DOCKER_NAME) 2> /dev/null || :

# To publish you must have write access to the docker hub DOCKER_HUB_ID user
publish: build
	docker push $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(PI3STREAMER_VERSION)

# Create/update the metadata in the exchange for this microservice
exchange-publish: build
	hzn exchange microservice publish -k $$PRIVATE_KEY_FILE -f horizon/microservice.definition.json

clean:
	-docker rm -f $(DOCKER_NAME) 2> /dev/null || :
	-docker rmi $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(PI3STREAMER_VERSION) 2> /dev/null || :

.PHONY: default all build run stop publish exchange-publish clean
