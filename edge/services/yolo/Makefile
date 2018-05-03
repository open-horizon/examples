# Transform the machine arch into some standard values: "arm", "arm64", or "amd64"
SYSTEM_ARCH := $(shell uname -m | sed -e 's/aarch64.*/arm64/' -e 's/x86_64.*/amd64/' -e 's/armv.*/arm/')

# To build for an arch different from the current system, set this env var to 1 of the values in the comment above
export ARCH ?= $(SYSTEM_ARCH)

# These variables can be overridden from the environment
export YOLO_NAME ?= yolov3vnc
export YOLO_VERSION ?= 8.7.1
export VNC_PASSWORD ?= vncpass
DOCKER_NAME ?= $(ARCH)_$(YOLO_NAME)
DOCKER_HUB_ID ?= openhorizon
export MYDOMAIN ?= internetofthings.ibmcloud.com
VNC_PORT ?= 5900
export YOLO_CAMERA ?= 0
export YOLO_HTTP_PORT ?= 8359

default: all

all: build run check

all-squash: build-squash run

build:
	docker build -t $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(YOLO_VERSION) -f ./Dockerfile.$(ARCH) .

build-squash:
	docker build --squash -t $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(YOLO_VERSION) -f ./Dockerfile.$(ARCH) .

run:
	docker rm -f $(DOCKER_NAME) 2> /dev/null || :
	docker run -d --privileged --volume /dev/:/dev:ro -e YOLO_CAMERA=$(YOLO_CAMERA) -e VNC_PASSWORD=$(VNC_PASSWORD) --name $(DOCKER_NAME) --publish=$(YOLO_HTTP_PORT):$(YOLO_HTTP_PORT) --publish=$(VNC_PORT):$(VNC_PORT) --volume `pwd`:/outside $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(YOLO_VERSION)

dev: 
	docker rm -f $(DOCKER_NAME) 2> /dev/null || :
	docker run -it --privileged --volume /dev/:/dev:ro -e YOLO_CAMERA=$(YOLO_CAMERA) -e VNC_PASSWORD=$(VNC_PASSWORD) --name $(DOCKER_NAME) --publish=$(YOLO_HTTP_PORT):$(YOLO_HTTP_PORT) --publish=$(VNC_PORT):$(VNC_PORT) --volume `pwd`:/outside $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(YOLO_VERSION) /bin/bash

check:
	docker exec $(DOCKER_NAME) cat output.json
	curl -sSL http://localhost:$(YOLO_HTTP_PORT)/v1/yolo | jq .

stop:
	-docker rm -f $(DOCKER_NAME) 2> /dev/null || :

# To publish you must have write access to the docker hub DOCKER_HUB_ID user
publish: build
	docker push $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(YOLO_VERSION)

publish-squash: build-squash
	docker push $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(YOLO_VERSION)

# Create/update the metadata in the exchange for this microservice
exchange-publish: build
	hzn exchange microservice publish -k $$PRIVATE_KEY_FILE -f horizon/microservice.definition.json

exchange-publish-squash: build-squash
	hzn exchange microservice publish -k $$PRIVATE_KEY_FILE -f horizon/microservice.definition.json

clean:
	-docker rm -f $(DOCKER_NAME) 2> /dev/null || :
	-docker rmi $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(YOLO_VERSION) 2> /dev/null || :

.PHONY: default all build run check stop publish exchange-publish clean all-squash build-squash publish-squash exchange-publish-squash
