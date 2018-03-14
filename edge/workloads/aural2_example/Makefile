# Transform the machine arch into some standard values: "arm", "arm64", or "amd64"
SYSTEM_ARCH := $(shell uname -m | sed -e 's/aarch64.*/arm64/' -e 's/x86_64.*/amd64/')

# To build for an arch different from the current system, set this env var to 1 of the values in the comment above
ARCH ?= $(SYSTEM_ARCH)

# These variables can be overridden from the environment
VERSION ?= 1.0.3
DOCKER_NAME ?= example_wl_$(ARCH)_aural2_example
DOCKER_HUB_ID ?= openhorizon

default: build

all: build run

aural2_example: main.go
	CC=$(which musl-gcc) go build --ldflags '-w -linkmode external -extldflags "-static"' -o aural2_example main.go

build: aural2_example
	docker build -t $(DOCKER_NAME):$(VERSION) -f ./Dockerfile.$(ARCH) .

dev:
	-docker rm -f $(DOCKER_NAME):$(VERSION)
	docker run -it --name aural2_example --net=example --net-alias=aural2_example --volume `pwd`:/outside $(DOCKER_NAME):$(VERSION) /bin/sh

run:
	-docker rm -f $(DOCKER_NAME):$(VERSION)
	docker run -d --name aural2_example --net=example --net-alias=aural2_example $(DOCKER_NAME):$(VERSION)

# To publish you must have write access to the docker hub openhorizon user
publish: build
	docker tag $(DOCKER_NAME):$(VERSION) $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(VERSION)
	docker push $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(VERSION)

clean:
	-docker rm -f $(DOCKER_NAME):$(VERSION)
	-docker rmi aural2_example

.PHONY: default all build dev run check publish clean
