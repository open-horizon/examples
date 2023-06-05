# This imports the variables from horizon/hzn.json. You can ignore these lines, but do not remove them
-include horizon/.hzn.json.tmp.mk

IMG_NAME=${DOCKER_IMAGE_BASE}

DOCKER_FILE=Dockerfile

DOCKER_CLI=DOCKER_CLI_EXPERIMENTAL=enabled docker

.PHONY: docker-build build-amd build-arm build publish-service test-all-arches publish docker-rmi clean publish-service-overwrite build-test-client

docker-build:
	$(DOCKER_CLI) build --no-cache -t $(IMG_NAME):$(HZN_ARCH)-$(SERVICE_VERSION) --build-arg BUILD_IMAGE=$(BUILD_IMAGE) --build-arg RUN_IMAGE=$(RUN_IMAGE) --build-arg GOARCH=$(GOARCH) -f $(DOCKER_FILE) .
	$(DOCKER_CLI) push $(IMG_NAME):$(HZN_ARCH)-$(SERVICE_VERSION)

build-amd:
	$(MAKE) BUILD_IMAGE=amd64/golang:1.14.2-alpine3.11 \
		RUN_IMAGE=amd64/alpine:3.11 GOARCH=amd64 HZN_ARCH=amd64  docker-build

build-arm:
	$(MAKE) BUILD_IMAGE=arm32v6/golang:1.14.2-alpine3.11 \
		RUN_IMAGE=arm32v6/alpine:3.11 GOARCH=arm HZN_ARCH=arm docker-build

build: build-amd build-arm

publish-service:
	hzn exchange service publish -O -I -v -f horizon/service.definition.json 

# This imports the variables from horizon/hzn.cfg. You can ignore these lines, but do not remove them.
horizon/.hzn.json.tmp.mk: horizon/hzn.json
	@ hzn util configconv -f $< | sed 's/=/?=/' > $@

docker-rmi:
	-$(DOCKER_CLI) rmi $(DOCKER_IMAGE_BASE):$(TAG) 2> /dev/null || :

clean:
	$(MAKE) TAG=amd64-$(SERVICE_VERSION) docker-rmi
	$(MAKE) TAG=arm-$(SERVICE_VERSION) docker-rmi
	$(MAKE) TAG=$(SERVICE_VERSION) docker-rmi

publish-service-overwrite:
	hzn exchange service publish -O -P -f horizon/service.definition.json

test-all-arches:
	exit 0

publish:
	ARCH=amd64 $(MAKE) publish-service
	ARCH=arm $(MAKE) publish-service
	hzn exchange pattern publish -f horizon/pattern-all-arches.json