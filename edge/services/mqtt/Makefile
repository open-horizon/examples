# A simple local MQTT broker for inter-service communication

# This imports the variables from horizon/hzn.json. You can ignore these lines, but do not remove them
-include horizon/.hzn.json.tmp.mk

SERVICE_NAME:="mqtt"
SERVICE_VERSION:="1.0.0"

# Get the Open-Horizon architecture type, and IP address for this host
ARCH:=$(shell ../../helper -a)

#
# Targets for building, testing and cleaning this service
#

build: validate-dockerhubid
	docker build -t $(DOCKERHUB_ID)/$(SERVICE_NAME)_$(ARCH):$(SERVICE_VERSION) -f ./Dockerfile.$(ARCH) .

push: validate-dockerhubid
	docker push $(DOCKERHUB_ID)/$(SERVICE_NAME)_$(ARCH):$(SERVICE_VERSION)

# How to perform a quick self-test of the MQTT broker
# - Verified on ubuntu and Rapsbian
# - NOTE: requires `mosquito-clients` package to be installed on the host
# - do this to test:
#    1. start a broker instance: `make test`
#    2. in one terminal on the host, subscribe with: `make test-sub`
#    3. then in another termninal on the host, publish with: `make test-pub`
#    4. You should see `Hello, world.` appear in the subscribed terminal
#    5. Terminate the subscriber with Ctrl-C, then: `make test-end`
test: validate-dockerhubid
	@docker network create mqtt-net 2>/dev/null || :
	-docker rm -f $(SERVICE_NAME) 2>/dev/null || :
	docker run --rm -d -p 127.0.0.1:1883:1883 --name $(SERVICE_NAME) --network mqtt-net --network-alias $(SERVICE_NAME) $(DOCKERHUB_ID)/$(SERVICE_NAME)_$(ARCH):$(SERVICE_VERSION)
test-sub:
	mosquitto_sub -h 127.0.0.1 -p 1883 -t /foo
test-pub:
	mosquitto_pub -h 127.0.0.1 -p 1883 -t /foo -m 'Hello, world.'
test-end:
	-docker rm -f $(SERVICE_NAME)

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

.PHONY: build test test-sub test-pub test-end clean publish-service publish-all-services validate-dockerhubid validate-org

