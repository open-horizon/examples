#
# A Makefile for building and testing the Horizon speedtest2wiotp Service
#
# Usage:
#   $ make
#   $ make test
#

VERSION:='1.0.0'
ARCH := $(shell hzn node list | jq .configuration.architecture)
WIOTP_ORG := $(shell cat wiotp.config | jq -r '.config."wiotp-org"')
WIOTP_API_KEYNAME := $(shell cat wiotp.config | jq -r '.config."wiotp-api-keyname"')
WIOTP_API_TOKEN := $(shell cat wiotp.config | jq -r '.config."wiotp-api-token"')
WIOTP_DEV_TYPE := $(shell cat wiotp.config | jq -r '.config."wiotp-device-type"')
WIOTP_DEV_ID := $(shell cat wiotp.config | jq -r '.config."wiotp-device-id"')
TEMP_SERVICE_DEF_FILE := $(shell mktemp)
TEMP_PATTERN_DEF_FILE := $(shell mktemp)

all: build run

build:
	docker build -f Dockerfile.$(ARCH) -t speedtest2wiotp .

dev: build
	-docker rm -f speedtest2wiotp 2> /dev/null || :
	docker run -it --name speedtest2wiotp -p $(HOST_TEST_PORT):80 --volume `pwd`:/outside speedtest2wiotp /bin/sh

run:
	-docker rm -f speedtest2wiotp 2>/dev/null || :
	docker run -d --name speedtest2wiotp -p $(HOST_TEST_PORT):80 speedtest2wiotp

test:
	@echo 'Subscribing to WIoTP topic. Use "Ctrl-C" to exit...'
	mosquitto_sub -h "$(WIOTP_ORG).messaging.internetofthings.ibmcloud.com" -p 8883 -i "a:$(WIOTP_ORG):$(WIOTP_API_KEYNAME)" -u "$(WIOTP_API_KEYNAME)" -P "$(WIOTP_API_TOKEN)" --capath /etc/ssl/certs -t "iot-2/type/$(WIOTP_DEV_TYPE)/id/$(WIOTP_DEV_ID)/evt/status/fmt/json" -d

exec:
	docker exec -it speedtest2wiotp /bin/sh

stop:
	-docker rm -f speedtest2wiotp 2>/dev/null || :

clean: stop
	-docker rmi speedtest2wiotp 2>/dev/null || :

push:
	docker tag speedtest2wiotp openhorizon/$(ARCH)_speedtest2wiotp:$(VERSION)
	docker push openhorizon/$(ARCH)_speedtest2wiotp:$(VERSION)

service-publish:
	export ARCH=$(ARCH); export VERSION=$(VERSION); cat horizon/service.definition.json | envsubst > $(TEMP_SERVICE_DEF_FILE)
	hzn exchange service publish -k ./*-private.key -K ./*-public.pem -f $(TEMP_SERVICE_DEF_FILE)
	-rm -f $(TEMP_SERVICE_DEF_FILE)

pattern-publish:
	export VERSION=$(VERSION); cat horizon/pattern.definition.json | envsubst > $(TEMP_PATTERN_DEF_FILE)
	hzn exchange pattern publish -k ./*-private.key -K ./*-public.pem -f $(TEMP_PATTERN_DEF_FILE)
	-rm -f $(TEMP_PATTERN_DEF_FILE)

.PHONY: all build dev run test exec stop clean push service-publish pattern-publish

