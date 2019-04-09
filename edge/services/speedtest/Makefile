#
# A simple Makefile for building and testing the Horizon speedtest Service
#
# Usage:
#   $ make
#   $ make test
#

VERSION:='1.0.0'
ARCH := $(shell hzn node list | jq .configuration.architecture)
HOST_TEST_PORT:=5659
TEMP_SERVICE_DEF_FILE := $(shell mktemp)

all: build run

build:
	docker build -f Dockerfile.$(ARCH) -t speedtest .

dev: build
	-docker rm -f speedtest 2> /dev/null || :
	docker run -it --name speedtest -p $(HOST_TEST_PORT):80 --volume `pwd`:/outside speedtest /bin/sh

run:
	-docker rm -f speedtest 2>/dev/null || :
	docker run -d --name speedtest -p $(HOST_TEST_PORT):80 speedtest

test:
	curl -sS localhost:$(HOST_TEST_PORT)/v1/speedtest | jq

exec:
	docker exec -it speedtest /bin/sh

stop:
	-docker rm -f speedtest 2>/dev/null || :

clean: stop
	-docker rmi speedtest 2>/dev/null || :

push:
	docker tag speedtest openhorizon/$(ARCH)_speedtest:$(VERSION)
	docker push openhorizon/$(ARCH)_speedtest:$(VERSION)

service-publish:
	export ARCH=$(ARCH); export VERSION=$(VERSION); cat horizon/service.definition.json | envsubst > $(TEMP_SERVICE_DEF_FILE)
	hzn exchange service publish -k ./*-private.key -K ./*-public.pem -f $(TEMP_SERVICE_DEF_FILE)
	-rm -f $(TEMP_SERVICE_DEF_FILE)

.PHONY: all build dev run test exec stop clean push service-publish
