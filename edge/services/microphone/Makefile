# Transform the machine arch into some standard values: "arm", "arm64", or "x86"
SYSTEM_ARCH := $(shell uname -m | sed -e 's/aarch64.*/arm64/' -e 's/x86_64.*/x86/' -e 's/armv.*/arm/')

# To build for an arch different from the current system, set this env var to 1 of the values in the comment above
ARCH ?= $(SYSTEM_ARCH)

# These variables can be overridden from the environment
VERSION ?= 1.2.6
DOCKER_NAME ?= example_ms_$(ARCH)_microphone
DOCKER_HUB_ID ?= openhorizon

default: build

all: build run

target/build_$(ARCH): Dockerfile.$(ARCH) main.go
	docker build -t $(DOCKER_NAME):$(VERSION) -f Dockerfile.$(ARCH) --build-arg version=$(VERSION) .
	touch target/build_$(ARCH)

run: target/build_$(ARCH)
	-docker network create microphone
	-docker kill microphone
	docker run -it --rm --privileged --name microphone --net=microphone --net-alias=microphone --publish=48926:48926 $(DOCKER_NAME):$(VERSION)

# To publish you must have write access to the docker hub openhorizon user
publish: target/build_$(ARCH)
	docker tag $(DOCKER_NAME):$(VERSION) $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(VERSION)
	docker push $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(VERSION)

clean:
	-docker rm -f $(DOCKER_NAME):$(VERSION)
	-docker rmi cpu
	-docker network rm example

.PHONY: default all build dev run check publish clean
