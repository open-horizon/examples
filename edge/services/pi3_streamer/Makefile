# Transform the machine arch into some standard values: "arm", "arm64", or "amd64"
SYSTEM_ARCH := $(shell uname -m | sed -e 's/aarch64.*/arm64/' -e 's/x86_64.*/amd64/' -e 's/armv.*/arm/')

# To build for an arch different from the current system, set this env var to 1 of the values in the comment above
ARCH ?= $(SYSTEM_ARCH)

# These variables can be overridden from the environment
VERSION ?= 1.0.0
DOCKER_NAME ?= $(ARCH)_example_ms_pi3streamer
DOCKER_HUB_ID ?= openhorizon

default: build

all: build run

build:
	docker build --squash -t $(DOCKER_NAME):$(VERSION) -f ./Dockerfile.$(ARCH) .
	docker tag $(DOCKER_NAME):$(VERSION) $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(VERSION)

dev:
	-docker network create example
	-docker rm -f $(DOCKER_NAME):$(VERSION)
	docker run -it --name pi3streamer --publish=8080:8080 --privileged -v /dev:/dev:ro -e VERT_FLIP=1 -e RESOLUTION="-x 1280 -y 720" --net=example --net-alias=pi3streamer --volume `pwd`:/outside $(DOCKER_NAME):$(VERSION) /bin/sh

run:
	-docker network create example
	-docker rm -f $(DOCKER_NAME):$(VERSION)
	docker run -d --name pi3streamer --publish=8080:8080 --privileged -v /dev:/dev:ro -e VERT_FLIP=1 -e RESOLUTION="-x 1280 -y 720" --net=example --net-alias=pi3streamer --volume `pwd`:/outside $(DOCKER_NAME):$(VERSION)

#check:
#	curl -s localhost:8347/v1/pi3streamer | jq

hznbuild: build
	mkdir -p horizon_build
	cp -r ./horizon/dependencies ./horizon_build/.
	export ARCH=$(ARCH) && export VERSION=$(VERSION) && export DOCKER_HUB_ID=$(DOCKER_HUB_ID) && envsubst < ./horizon/microservice.definition.json > ./horizon_build/microservice.definition.json
	envsubst < ./horizon/userinput.json > ./horizon_build/userinput.json

hznstart: hznbuild
	hzn dev microservice start -d ./horizon_build

hznstop:
	hzn dev microservice stop -d ./horizon_build

# To publish you must have write access to the docker hub openhorizon user
publish:
	docker push $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(VERSION)

clean:
	-hzn dev microservice stop -d ./horizon_build 2> /dev/null || :
	-rm -fr horizon_build 2> /dev/null || :
	-docker rm -f $(DOCKER_NAME):$(VERSION) 2> /dev/null || :
	-docker rmi pi3streamer 2> /dev/null || :
	-docker network rm example 2> /dev/null || :

.PHONY: default all build dev run check publish clean
