
all:
	make build
	make run

# Get the 3-character prefix of the local host architecture ("arm" or "x86")
SYSTEM_ARCH := $(shell uname -m | sed 's/\(...\).*/\1/')

build:
	docker build -t cpu -f ./Dockerfile.$(SYSTEM_ARCH) .

dev:
	-docker network create example
	-docker stop cpu
	-docker rm cpu
	docker run -it --name cpu --publish=8347:8347 --net=example --net-alias=cpu --volume `pwd`:/outside cpu /bin/sh

run:
	-docker network create example
	-docker stop cpu
	-docker rm cpu
	docker run -d --name cpu --publish=8347:8347 --net=example --net-alias=cpu --volume `pwd`:/outside cpu

check:
	curl -s localhost:8347/v1/cpu | jq

clean:
	-docker stop cpu
	-docker rm cpu
	-docker rmi cpu
	-docker network rm example

