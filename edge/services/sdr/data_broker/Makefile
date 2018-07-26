# Transform the machine arch into some standard values: "arm", "arm64", or "x86"
SYSTEM_ARCH := $(shell uname -m | sed -e 's/aarch64.*/arm64/' -e 's/x86_64.*/amd64/' -e 's/armv.*/arm/')

# To build for an arch different from the current system, set this env var to 1 of the values in the comment above
export ARCH ?= $(SYSTEM_ARCH)

# Container image version and workload name
VERSION ?= 0.0.1
WL_NAME=sdr_data_broker
IMAGE=$(WL_NAME)


# Construct our Docker registry path and other architecture-dependent values
DOCKER_HUB_ID ?= openhorizon
DOCKER_NAME ?= sdr_databroker_$(ARCH)
REG_PATH := $(DOCKER_HUB_ID)/$(ARCH)
IMAGE_PATH = $(REG_PATH)_$(IMAGE):$(VERSION)

target/build_$(ARCH): main.go
	docker build -t $(DOCKER_NAME):$(VERSION) -f Dockerfile.$(ARCH) ../
	touch target/build_$(ARCH)

run: target/build_$(ARCH)
	-docker network create sdr
	#-docker rm -f $(IMAGE_PATH)
	docker run -it --rm -e WIOTP_DOMAIN="${WIOTP_DOMAIN}" -e HZN_ORGANIZATION="$(HZN_ORG_ID)" -e HZN_DEVICE_ID="$(HZN_DEVICE_ID)" -e MSGHUB_API_KEY="$(MSGHUB_API_KEY)" -e MSGHUB_BROKER_URL="$(MSGHUB_BROKER_URL)" --name $(IMAGE) --net=sdr --net-alias=$(WL_NAME) $(DOCKER_NAME):$(VERSION)

# To publish you must have write access to the docker hub openhorizon user
publish: target/build_$(ARCH)
	docker tag $(DOCKER_NAME):$(VERSION) $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(VERSION)
	docker push $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(VERSION)

fakepub: fake/main.go
	go build -o $@ $<

clean:
	-docker rm -f $(IMAGE_PATH)
	-docker stop $(IMAGE) && docker rm $(IMAGE)
	-docker rmi $(IMAGE)
	-docker network rm example
	-rm target/*
	-rm -f fakepub

.PHONY: default all build dev run check docker-publish exchange-publish fakepub clean
