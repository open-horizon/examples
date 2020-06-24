# A monitor program for the achatina examples

# This imports the variables from horizon/hzn.json. You can ignore these lines, but do not remove them
-include horizon/.hzn.json.tmp.mk

SERVICE_NAME:="monitor"
SERVICE_VERSION:="1.0.0"

# Get the Open-Horizon architecture type, and IP address for this host
ARCH:=$(shell ../../helper -a)

#
# Targets for building, developing, testing and cleaning this service
#

build: validate-dockerhubid
	docker build -t $(DOCKERHUB_ID)/$(SERVICE_NAME)_$(ARCH):$(SERVICE_VERSION) -f ./Dockerfile.$(ARCH) .

push: validate-dockerhubid
	docker push $(DOCKERHUB_ID)/$(SERVICE_NAME)_$(ARCH):$(SERVICE_VERSION)

dev: build validate-dockerhubid
	-docker rm -f $(SERVICE_NAME) 2>/dev/null || :
	docker run -it -v `pwd`:/outside \
           -p 0.0.0.0:5200:5200 \
           --name ${SERVICE_NAME} \
           --network mqtt-net --network-alias ${SERVICE_NAME} \
           $(DOCKERHUB_ID)/$(SERVICE_NAME)_$(ARCH):$(SERVICE_VERSION) /bin/sh

test: validate-dockerhubid
	-docker rm -f $(SERVICE_NAME) 2>/dev/null || :
	docker run -d \
           -p 0.0.0.0:5200:5200 \
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

publish-service: validate-dockerhubid validate-org
	ARCH=$(ARCH) \
          SERVICE_NAME="$(SERVICE_NAME)" \
          SERVICE_VERSION="$(SERVICE_VERSION)"\
          DOCKER_IMAGE_BASE="$(DOCKERHUB_ID)/$(SERVICE_NAME)"\
          hzn exchange service publish -O -f service.json --pull-image

publish-all-services: validate-dockerhubid validate-org
	docker build -t $(DOCKERHUB_ID)/$(SERVICE_NAME)_arm:$(SERVICE_VERSION) -f ./Dockerfile.arm .
	docker push $(DOCKERHUB_ID)/$(SERVICE_NAME)_arm:$(SERVICE_VERSION)
	ARCH=arm \
          SERVICE_NAME="$(SERVICE_NAME)" \
          SERVICE_VERSION="$(SERVICE_VERSION)"\
          DOCKER_IMAGE_BASE="$(DOCKERHUB_ID)/$(SERVICE_NAME)"\
          hzn exchange service publish -O -f service.json --pull-image
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

validate-org:
	@if [ -z "${HZN_ORG_ID}" ]; \
          then { echo "***** ERROR: \"HZN_ORG_ID\" is not set!"; exit 1; }; \
          else echo "  NOTE: Using Exchange Org ID: \"${HZN_ORG_ID}\""; \
        fi
	@sleep 1

# This imports the variables from horizon/hzn.cfg. You can ignore these lines, but do not remove them.
horizon/.hzn.json.tmp.mk: horizon/hzn.json
	@ hzn util configconv -f $< | sed 's/=/?=/' > $@

.PHONY: build dev test clean publish-service publish-all-services validate-dockerhubid validate-org

