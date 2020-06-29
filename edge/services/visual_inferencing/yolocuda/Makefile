# An app that calls the yolo detector, and pushes data to mqtt and kafka

SERVICE_NAME:="yolocuda"
SERVICE_VERSION:="1.0.0"

# Get the Open-Horizon architecture type, and IP address for this host
ARCH:=$(shell ../../helper -a)
HOST_IP:=$(shell ../../helper -i)
DEFAULT_CAM_URL:="http://${HOST_IP}:8888/"

# If HZN_DEVICE_ID is not set, use the hostname
HZN_DEVICE_ID := $(or $(HZN_DEVICE_ID),$(shell sh -c 'hostname'))

#
# Targets for building, developing, testing and cleaning this service
#

build: validate-dockerhubid
	@if [ "${ARCH}" != "arm" ]; then \
          docker build -t $(DOCKERHUB_ID)/$(SERVICE_NAME)_$(ARCH):$(SERVICE_VERSION) -f ./Dockerfile.$(ARCH) .; \
        fi

push: validate-dockerhubid
	@if [ "${ARCH}" != "arm" ]; then \
	  docker push $(DOCKERHUB_ID)/$(SERVICE_NAME)_$(ARCH):$(SERVICE_VERSION); \
        fi

dev: validate-dockerhubid check-cam-url
	-docker rm -f ${SERVICE_NAME} 2>/dev/null || :
	docker run -it -v `pwd`:/outside \
	   -e HZN_DEVICE_ID="${HZN_DEVICE_ID}" \
	   -e DEFAULT_CAM_URL="${DEFAULT_CAM_URL}" \
	   -e CAM_URL="${CAM_URL}" \
	   -e EVENTSTREAMS_BROKER_URLS="${EVENTSTREAMS_BROKER_URLS}" \
	   -e EVENTSTREAMS_API_KEY="${EVENTSTREAMS_API_KEY}" \
	   -e EVENTSTREAMS_PUB_TOPIC="${EVENTSTREAMS_PUB_TOPIC}" \
	   --name ${SERVICE_NAME} \
	   --network mqtt-net --network-alias ${SERVICE_NAME} \
	   $(DOCKERHUB_ID)/$(SERVICE_NAME)_$(ARCH):$(SERVICE_VERSION) /bin/sh

test: validate-dockerhubid check-cam-url
	-docker rm -f ${SERVICE_NAME} 2>/dev/null || :
	docker run --rm -d \
	   -e HZN_DEVICE_ID="${HZN_DEVICE_ID}" \
	   -e DEFAULT_CAM_URL="${DEFAULT_CAM_URL}" \
	   -e CAM_URL="${CAM_URL}" \
	   -e EVENTSTREAMS_BROKER_URLS="${EVENTSTREAMS_BROKER_URLS}" \
	   -e EVENTSTREAMS_API_KEY="${EVENTSTREAMS_API_KEY}" \
	   -e EVENTSTREAMS_PUB_TOPIC="${EVENTSTREAMS_PUB_TOPIC}" \
	   --name ${SERVICE_NAME} \
	   --network mqtt-net --network-alias ${SERVICE_NAME} \
	   $(DOCKERHUB_ID)/$(SERVICE_NAME)_$(ARCH):$(SERVICE_VERSION)

clean: validate-dockerhubid
	@docker rm -f ${SERVICE_NAME} 2>/dev/null || :
	@docker rmi $(DOCKERHUB_ID)/$(SERVICE_NAME):$(SERVICE_VERSION) 2>/dev/null || :
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

check-cam-url:
	@if [ -z "${CAM_URL}" ]; \
          then echo "  Warning: \"CAM_URL\" is not set! Using default."; \
          else echo "  NOTE: Using camera URL: ${CAM_URL}"; \
        fi
	@sleep 1

check-kafka-creds:
	@if [ -z "${EVENTSTREAMS_BROKER_URLS}" ] || \
          [ -z "${EVENTSTREAMS_API_KEY}" ] || \
          [ -z "${EVENTSTREAMS_PUB_TOPIC}" ]; \
          then echo "  Warning: No EventStreams credentials found! Kafka publication is disabled."; \
          else echo "  NOTE: Publishing to topic: ${EVENTSTREAMS_PUB_TOPIC}"; \
        fi

.PHONY: build dev test clean publish-service publish-all-services validate-dockerhubid check-cam-url check-kafka-creds

