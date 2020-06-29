# The yolo detector using CUDA
  
SERVICE_NAME:="restyolocuda"
SERVICE_VERSION:="1.0.0"

# Get the Open-Horizon architecture type, and IP address for this host
ARCH:=$(shell ../../helper -a)

#
# Targets for building, developing, running and cleaning this service
#

build: validate-dockerhubid
	@if [ "${ARCH}" != "arm" ]; then \
	  docker build -t $(DOCKERHUB_ID)/$(SERVICE_NAME)_$(ARCH):$(SERVICE_VERSION) -f ./Dockerfile.$(ARCH) .; \
        fi

push: validate-dockerhubid
	@if [ "${ARCH}" != "arm" ]; then \
	  docker push $(DOCKERHUB_ID)/$(SERVICE_NAME)_$(ARCH):$(SERVICE_VERSION); \
        fi

dev: validate-dockerhubid
	-docker rm -f $(SERVICE_NAME) 2>/dev/null || :
	docker run -it \
	  -p 5252:80 \
	  --name $(SERVICE_NAME) \
	  --network mqtt-net --network-alias $(SERVICE_NAME) \
	  -v `pwd`:/outside \
	  $(DOCKERHUB_ID)/$(SERVICE_NAME)_$(ARCH):$(SERVICE_VERSION) /bin/bash

test: validate-dockerhubid
	-docker rm -f $(SERVICE_NAME) 2>/dev/null || :
	docker run --rm -d \
	  -p 5252:80 \
	  --name $(SERVICE_NAME) \
	  --network mqtt-net --network-alias $(SERVICE_NAME) \
	  $(DOCKERHUB_ID)/$(SERVICE_NAME)_$(ARCH):$(SERVICE_VERSION)

clean: validate-dockerhubid
	@docker rm -f ${SERVICE_NAME} 2>/dev/null || :
	@docker rmi $(DOCKERHUB_ID)/$(SERVICE_NAME)_$(ARCH):$(SERVICE_VERSION) 2>/dev/null || :

#
# Targets for publishing this service to an Open-Horizon Exhange
#
# NOTE: You must install the Open-Horizon CLI in order to use these targets!
#

publish-service: validate-dockerhubid
	@if [ "${ARCH}" != "arm" ]; then \
          ARCH=$(ARCH) \
          SERVICE_NAME="$(SERVICE_NAME)" \
          SERVICE_VERSION="$(SERVICE_VERSION)"\
          DOCKER_IMAGE_BASE="$(DOCKERHUB_ID)/$(SERVICE_NAME)"\
          hzn exchange service publish -O -f service.json --pull-image; \
        fi

# This service only builds for arm64 and amd64
publish-all-services: validate-dockerhubid
	docker build -t $(DOCKERHUB_ID)/$(SERVICE_NAME)_arm64:$(SERVICE_VERSION) -f ./Dockerfile.arm64 .
	docker push $(DOCKERHUB_ID)/$(SERVICE_NAME)_arm64:$(SERVICE_VERSION)
	ARCH=arm64 \
          SERVICE_NAME="$(SERVICE_NAME)" \
          SERVICE_VERSION="$(SERVICE_VERSION)"\
          DOCKER_IMAGE_BASE="$(DOCKERHUB_ID)/$(SERVICE_NAME)"\
          hzn exchange service publish -O -f service.json --pull-image
	docker build -t $(DOCKERHUB_ID)/$(SERVICE_NAME)_amd64:$(SERVICE_VERSION) -f ./Dockerfile.amd64 .
	docker push $(DOCKERHUB_ID)/$(SERVICE_NAME)_amd64:$(SERVICE_VERSION)
	ARCH=amd64 \
          SERVICE_NAME="$(SERVICE_NAME)" \
          SERVICE_VERSION="$(SERVICE_VERSION)"\
          DOCKER_IMAGE_BASE="$(DOCKERHUB_ID)/$(SERVICE_NAME)"\
          hzn exchange service publish -O -f service.json --pull-image

#
# Sanity check targets
#

validate-dockerhubid:
	@if [ -z "${DOCKERHUB_ID}" ]; \
          then { echo "***** ERROR: \"DOCKERHUB_ID\" is not set!"; exit 1; }; \
          else echo "  NOTE: Using DockerHubID: \"${DOCKERHUB_ID}\""; \
        fi
	@sleep 1


.PHONY: build push dev test clean publish-service publish-all-services validate-dockerhubid

