DOCKER_NAME = cpu2wiotp_workload
DOCKER_TAG = 0.0.1

default: all

all: build run check

build: clean
	docker build -t $(DOCKER_NAME):$(DOCKER_TAG) .

run:
	-docker rm -f $(DOCKER_NAME) 2> /dev/null || :
	docker run -d --name $(DOCKER_NAME) -e SAMPLE_SIZE=5 -e SAMPLE_INTERVAL=1 -e MOCK=true -e VERBOSE=1 -e PUBLISH=false $(DOCKER_NAME):$(DOCKER_TAG)

check:
	docker logs -f $(DOCKER_NAME)

stop:
	-docker rm -f $(DOCKER_NAME) 2> /dev/null || :

clean: stop
	-docker rmi $(DOCKER_NAME):$(DOCKER_TAG)

.PHONY: default all build run check publish clean