# Transform the machine arch into some standard values: "arm", "arm64", or "amd64"
SYSTEM_ARCH := $(shell uname -m | sed -e 's/aarch64.*/arm64/' -e 's/x86_64.*/amd64/' -e 's/armv.*/arm/')

# To build for an arch different from the current system, set this env var to 1 of the values in the comment above
ARCH ?= $(SYSTEM_ARCH)

# These variables can be overridden from the environment
VERSION ?= 1.0.0
DOCKER_NAME ?= example_wl_$(ARCH)_pi3streamer2wiotp
DOCKER_HUB_ID ?= openhorizon

default: build

all: build run

build:
	docker build -t $(DOCKER_NAME):$(VERSION) -f ./Dockerfile.$(ARCH) .
	docker tag $(DOCKER_NAME):$(VERSION) $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(VERSION)

dev:
	-docker rm -f $(DOCKER_NAME):$(VERSION)
	docker run -it -e WIOTP_DOMAIN="${WIOTP_DOMAIN}" -e HZN_ORGANIZATION="$(HZN_ORG_ID)" -e HZN_DEVICE_ID="$(HZN_DEVICE_ID)" --name pi3streamer2wiotp --net=example --net-alias=pi3streamer2wiotp --volume `pwd`:/outside $(DOCKER_NAME):$(VERSION) /bin/sh

run:
	-docker rm -f $(DOCKER_NAME):$(VERSION)
	docker run -d -e WIOTP_DOMAIN="${WIOTP_DOMAIN}" -e HZN_ORGANIZATION="$(HZN_ORG_ID)" -e HZN_DEVICE_ID="$(HZN_DEVICE_ID)" --name pi3streamer2wiotp --net=example --net-alias=pi3streamer2wiotp $(DOCKER_NAME):$(VERSION)

check:
	curl -Is http://localhost:8080/static_simple.html | grep MJPG

hznbuild: build
	mkdir -p horizon_build/dependencies
	set -x && export ARCH=$(ARCH) && cd ./horizon/dependencies && for filename in $$(find . -name "*.microservice.definition.json"); do envsubst < ./$$filename > ../../horizon_build/dependencies/$$filename; done && cd ../..
	export ARCH=$(ARCH) && export VERSION=$(VERSION) && export DOCKER_HUB_ID=$(DOCKER_HUB_ID) && export WIOTP_DOMAIN=$(WIOTP_DOMAIN) && envsubst < ./horizon/workload.definition.json > ./horizon_build/workload.definition.json
	export WIOTP_DOMAIN=$(WIOTP_DOMAIN) && envsubst < ./horizon/userinput.json > ./horizon_build/userinput.json

hznstart: hznbuild
	hzn dev workload start -d ./horizon_build

hznstop:
	hzn dev workload stop -d ./horizon_build

# To publish you must have write access to the docker hub openhorizon user
publish:
	docker push $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(VERSION)

clean:
	-hzn dev workload stop -d ./horizon_build 2> /dev/null || :
	-rm -fr horizon_build 2> /dev/null || :
	-docker rm -f $(DOCKER_NAME):$(VERSION) 2> /dev/null || :
	-docker rmi pi3streamer2wiotp 2> /dev/null || :

.PHONY: default all build dev run check publish clean
