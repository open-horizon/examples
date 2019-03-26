#
# A simple Makefile for building and testing the Horizon speedtest Service
#
# Usage:
#   $ make
#   $ make test
#

ARCH:=$(shell dpkg --print-architecture)

HOST_TEST_PORT:=5659

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

.PHONY: all build dev run test exec stop clean
