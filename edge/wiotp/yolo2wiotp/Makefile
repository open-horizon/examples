# Transform the machine arch into some standard values: "arm", "arm64", or "amd64"
SYSTEM_ARCH := $(shell uname -m | sed -e 's/aarch64.*/arm64/' -e 's/x86_64.*/amd64/' -e 's/armv.*/arm/')

# To build for an arch different from the current system, set this env var to 1 of the values in the comment above
export ARCH ?= $(SYSTEM_ARCH)

# These variables can be overridden from the environment
export YOLO2WIOTP_NAME ?= yolo2wiotp
export YOLO2WIOTP_VERSION ?= 1.0.0
DOCKER_NAME ?= $(ARCH)_$(YOLO2WIOTP_NAME)
DOCKER_HUB_ID ?= openhorizon
export MYDOMAIN ?= internetofthings.ibmcloud.com
CORE_IOT_VERSION ?= 2.4.0

default: all

all: build run check

build:
	docker build -t $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(YOLO2WIOTP_VERSION) -f ./Dockerfile.$(ARCH) .

run:
	-docker rm -f $(DOCKER_NAME) 2> /dev/null || :
	docker run -d --name $(DOCKER_NAME) -e SAMPLE_INTERVAL=1 -e MOCK=true -e VERBOSE=1 -e PUBLISH=false --volume `pwd`:/outside $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(YOLO2WIOTP_VERSION)

dev:
	-docker rm -f $(DOCKER_NAME) 2> /dev/null || :
	docker run -it --name $(DOCKER_NAME) -e SAMPLE_INTERVAL=1 -e MOCK=true -e VERBOSE=1 -e PUBLISH=false --volume `pwd`:/outside $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(YOLO2WIOTP_VERSION) /bin/sh

check:
	docker logs -f $(DOCKER_NAME)

stop:
	-docker rm -f $(DOCKER_NAME) 2> /dev/null || :

# Targets for using the hzn dev environment
hznbuild: build
	hzn dev dependency fetch -p ../../services/yolo/horizon
	hzn dev dependency fetch -s https://internetofthings.ibmcloud.com/wiotp-edge/microservices/edge-core-iot-microservice --ver $(CORE_IOT_VERSION) -o IBM -a $(ARCH) -k /etc/horizon/trust/publicWIoTPEdgeComponentsKey.pem

hznstart: hznbuild
	hzn dev workload start

hznstop:
	hzn dev workload stop

# Push the docker image to the registry. You must have write access to the docker hub openhorizon user
publish: build
	docker push $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(YOLO2WIOTP_VERSION)

# Create/update the metadata in the exchange for this workload
exchange-publish: build
	hzn exchange workload publish -k $(PRIVATE_KEY_FILE) -f horizon/workload.definition.json

clean:
	-docker rm -f $(DOCKER_NAME) 2> /dev/null || :
	-docker rmi $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(YOLO2WIOTP_VERSION) 2> /dev/null || :

.PHONY: default all build run check stop publish exchange-publish clean
