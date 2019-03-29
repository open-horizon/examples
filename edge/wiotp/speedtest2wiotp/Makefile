#
# A Makefile for building and testing the Horizon speedtest2wiotp Service
#
# Usage:
#   $ make
#       ...  then set environment variables so you can subscribe to WIoTP ...
#   $ make test
#

# Set this "semver" VERSION value when you make code changes and re-publish
VERSION:='1.0.0'

# Derived values. No need to edit these.
ARCH := $(shell hzn node list | jq .configuration.architecture)
TEMP_SERVICE_DEF_FILE := $(shell mktemp)
TEMP_PATTERN_DEF_FILE := $(shell mktemp)

all: build run

build:
	docker build -f Dockerfile.$(ARCH) -t speedtest2wiotp .
	docker tag speedtest2wiotp openhorizon/$(ARCH)_speedtest2wiotp:$(VERSION)

dev: build
	-docker rm -f speedtest2wiotp 2> /dev/null || :
	docker run -it --name speedtest2wiotp -p $(HOST_TEST_PORT):80 --volume `pwd`:/outside speedtest2wiotp /bin/sh

run:
	-docker rm -f speedtest2wiotp 2>/dev/null || :
	docker run -d --name speedtest2wiotp -p $(HOST_TEST_PORT):80 speedtest2wiotp

test:
	: $${WIOTP_ORG:?} $${WIOTP_DEVICE_TYPE:?} $${WIOTP_DEVICE_ID:?} $${WIOTP_API_KEYNAME:?} $${WIOTP_API_TOKEN:?}   # this verifies these env vars are set
	@echo 'Subscribing to WIoTP topic. Use "Ctrl-C" to exit...'
	mosquitto_sub -h "$(WIOTP_ORG).messaging.internetofthings.ibmcloud.com" -p 8883 -i "a:$(WIOTP_ORG):$(WIOTP_API_KEYNAME)" -u "$(WIOTP_API_KEYNAME)" -P "$(WIOTP_API_TOKEN)" --capath /etc/ssl/certs -t "iot-2/type/$(WIOTP_DEVICE_TYPE)/id/$(WIOTP_DEVICE_ID)/evt/status/fmt/json" -d | grep '^{.*' | jq

exec:
	docker exec -it speedtest2wiotp /bin/sh

stop:
	-docker rm -f speedtest2wiotp 2>/dev/null || :

clean: stop
	-docker rmi speedtest2wiotp 2>/dev/null || :

push:
	docker push openhorizon/$(ARCH)_speedtest2wiotp:$(VERSION)

service-publish:
	export ARCH=$(ARCH); export VERSION=$(VERSION); cat horizon/service.definition.json | envsubst > $(TEMP_SERVICE_DEF_FILE)
	hzn exchange service publish -k ./*-private.key -K ./*-public.pem -f $(TEMP_SERVICE_DEF_FILE)
	-rm -f $(TEMP_SERVICE_DEF_FILE)

pattern-publish:
	export VERSION=$(VERSION); cat horizon/pattern/speedtest2wiotp.json | envsubst > $(TEMP_PATTERN_DEF_FILE)
	hzn exchange pattern publish -k ./*-private.key -K ./*-public.pem -f $(TEMP_PATTERN_DEF_FILE)
	-rm -f $(TEMP_PATTERN_DEF_FILE)

.PHONY: all build dev run test exec stop clean push service-publish pattern-publish

