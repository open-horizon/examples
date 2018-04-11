DOCKER_NAME ?= cpu_microservice
DOCKER_TAG ?= 0.0.1

DOCKER_HUB_ID ?= yourDockerId

default: all

all: build run check

build: clean
	docker build -t $(DOCKER_NAME):$(DOCKER_TAG) .
	docker tag $(DOCKER_NAME):$(DOCKER_TAG) $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(DOCKER_TAG)

run:
	-docker rm -f $(DOCKER_NAME) 2> /dev/null || :
	docker run -d --name $(DOCKER_NAME) --publish=8347:8347 --volume `pwd`:/outside $(DOCKER_NAME):$(DOCKER_TAG)

check:
	curl -sSL http://localhost:8347/v1/cpu | jq .

stop:
	-docker rm -f $(DOCKER_NAME) 2> /dev/null || :

publish:
	docker push $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(DOCKER_TAG)

clean: stop
	-docker rmi $(DOCKER_NAME):$(DOCKER_TAG)
	-docker rmi $(DOCKER_HUB_ID)/$(DOCKER_NAME):$(DOCKER_TAG)

.PHONY: default all build run check publish clean
