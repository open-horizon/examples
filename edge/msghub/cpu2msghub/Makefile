# Transform the machine arch into some standard values: "arm", "arm64", or "amd64"
SYSTEM_ARCH := $(shell uname -m | sed -e 's/aarch64.*/arm64/' -e 's/x86_64.*/amd64/' -e 's/armv.*/arm/')

# To build for an arch different from the current system, set this env var to 1 of the values in the comment above
export ARCH ?= $(SYSTEM_ARCH)

# These variables can be overridden from the environment
export HZN_ORG_ID=IBM
export SERVICE_NAME ?= cpu2msghub
export SERVICE_URL ?= github.com.open-horizon.examples.$(SERVICE_NAME)
export SERVICE_VERSION ?= 1.2.6
export DOCKER_HUB_ID ?= openhorizon

# Derived variables
export HZN_PATTERN ?= pattern-$(SERVICE_NAME)
export DOCKER_IMAGE ?= $(DOCKER_HUB_ID)/$(ARCH)_$(SERVICE_NAME):$(SERVICE_VERSION)

default: hznstart

build:
ifeq (,$(findstring amd64,$(ARCH)))
	mkdir -p tmp/$(ARCH)
	cp ../../../tools/kafkacat/$(ARCH)/*.rsa.pub ../../../tools/kafkacat/$(ARCH)/kafkacat-*.apk tmp/$(ARCH)
endif
	docker build -t $(DOCKER_IMAGE) -f ./Dockerfile.$(ARCH) .
ifeq (,$(findstring amd64,$(ARCH)))
	rm -f tmp/$(ARCH)/*.rsa.pub tmp/$(ARCH)/kafkacat-*.apk
endif

run:
	-docker rm -f $(SERVICE_NAME) 2> /dev/null || :
	docker run -d --name $(SERVICE_NAME) -e SAMPLE_SIZE=5 -e SAMPLE_INTERVAL=1 -e MOCK=true -e VERBOSE=1 -e PUBLISH=false -e HZN_ORGANIZATION=$(HZN_ORG_ID) -e HZN_PATTERN=$(HZN_PATTERN) -e HZN_DEVICE_ID=$(shell hostname -s) --volume `pwd`:/outside $(DOCKER_IMAGE)

check:
	docker logs -f $(SERVICE_NAME)

stop:
	-docker rm -f $(SERVICE_NAME) 2> /dev/null || :

msghub-topic:
	: $${MSGHUB_API_KEY:?} $${MSGHUB_ADMIN_URL:?}
	export MSGHUB_TOPIC="$${HZN_ORG_ID//[@#%()+=:,<>\/]/_}.$${HZN_ORG_ID//[@#%()+=:,<>\/]/_}_$${HZN_PATTERN//[@#%()+=:,<>\/]/_}" && curl -sS -w %{http_code} -H 'Content-Type: application/json' -H "X-Auth-Token: $(MSGHUB_API_KEY)" -d "{ \"name\": \"$$MSGHUB_TOPIC\", \"partitions\": 2 }" $(MSGHUB_ADMIN_URL)/admin/topics

msghub-topic-list:
	: $${MSGHUB_API_KEY:?} $${MSGHUB_ADMIN_URL:?}
	curl -sS -H "X-Auth-Token: $$MSGHUB_API_KEY" $$MSGHUB_ADMIN_URL/admin/topics | jq -r ".[] | .name"

# Targets for using the hzn dev environment
hznbuild: build
	#todo: these 2 fetch cmds will not be needed once anax issue 785 is completed
	hzn dev dependency fetch --arch=$(ARCH) --org IBM --url github.com.open-horizon.examples.cpu
	hzn dev dependency fetch --arch=$(ARCH) --org IBM --url github.com.open-horizon.examples.gps
	hzn dev service verify

hznstart: hznbuild
	: $${MSGHUB_API_KEY:?} $${MSGHUB_BROKER_URL:?}   # verify these env vars are set
	hzn dev service start

hznstop:
	hzn dev service stop

# Push the docker image to the registry. You must have write access to the docker hub openhorizon user
docker-push: build
	docker push $(DOCKER_IMAGE)

# Create/update the metadata in the exchange for this service
publish-service: build publish-service-only
publish-service-only:
	: $${HZN_EXCHANGE_USER_AUTH:?} $${PRIVATE_KEY_FILE:?} $${PUBLIC_KEY_FILE:?}   # verify these env vars are set
	hzn exchange service publish -k $$PRIVATE_KEY_FILE -K $$PUBLIC_KEY_FILE -f horizon/service.definition.json

# Create/update the deployment pattern for a set of edge nodes
publish-pattern:
	: $${HZN_EXCHANGE_USER_AUTH:?}   # verify these env vars are set
	hzn exchange pattern publish -p $(HZN_PATTERN) -f horizon/pattern/pattern-cpu2msghub.json

register-node:
	: $${HZN_EXCHANGE_USER_AUTH:?} $${EXCHANGE_NODEAUTH:?} $${MSGHUB_API_KEY:?} $${MSGHUB_BROKER_URL:?}   # verify these env vars are set
	hzn register -u $(HZN_EXCHANGE_USER_AUTH) -n $(EXCHANGE_NODEAUTH) $(HZN_ORG_ID) $(HZN_ORG_ID)/$(HZN_PATTERN) -f horizon/userinput.json

testmake:
	echo $$DOCKER_IMAGE

clean:
	-docker rm -f $(SERVICE_NAME) 2> /dev/null || :
	-docker rmi $(DOCKER_IMAGE) 2> /dev/null || :

.PHONY: default all build run check stop hznbuild hznstart hznstop publish-service publish-service-only publish-pattern register-node testmake clean
