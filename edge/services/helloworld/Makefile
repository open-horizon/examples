# Make targets for building the IBM example helloworld edge service

build:
	source horizon/hzn.env && docker build -t $$DOCKER_IMAGE -f ./Dockerfile.$$ARCH .

clean:
	-source horizon/hzn.env && docker rmi $$DOCKER_IMAGE 2> /dev/null || :

.PHONY: build clean
