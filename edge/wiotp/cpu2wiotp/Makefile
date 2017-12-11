
all:
	make build
	make run

# Get the 3-character prefix of the local host architecture ("arm" or "x86")
SYSTEM_ARCH := $(shell uname -m | sed 's/\(...\).*/\1/')

build:
	docker build -t cpu2wiotp -f ./Dockerfile.$(SYSTEM_ARCH) .

dev:
	-docker stop cpu2wiotp
	-docker rm cpu2wiotp
	docker run -it -e WIOTP_DOMAIN="${WIOTP_DOMAIN}" -e WIOTP_ORG_ID="$(WIOTP_ORG_ID)" -e WIOTP_API_KEY="$(WIOTP_API_KEY)" -e WIOTP_API_AUTH_TOKEN="$(WIOTP_API_AUTH_TOKEN)" -e WIOTP_DEVICE_TYPE="$(WIOTP_DEVICE_TYPE)" -e WIOTP_DEVICE_AUTH_TOKEN="$(WIOTP_DEVICE_AUTH_TOKEN)" -e HZN_DEVICE_ID="$(HZN_DEVICE_ID)" --name cpu2wiotp --net=example --net-alias=cpu2wiotp --volume `pwd`:/outside cpu2wiotp /bin/sh

run:
	-docker stop cpu2wiotp
	-docker rm cpu2wiotp
	docker run -d -e WIOTP_DOMAIN="${WIOTP_DOMAIN}" -e WIOTP_ORG_ID="$(WIOTP_ORG_ID)" -e WIOTP_API_KEY="$(WIOTP_API_KEY)" -e WIOTP_API_AUTH_TOKEN="$(WIOTP_API_AUTH_TOKEN)" -e WIOTP_DEVICE_TYPE="$(WIOTP_DEVICE_TYPE)" -e WIOTP_DEVICE_AUTH_TOKEN="$(WIOTP_DEVICE_AUTH_TOKEN)" -e HZN_DEVICE_ID="$(HZN_DEVICE_ID)" --name cpu2wiotp --net=example --net-alias=cpu2wiotp cpu2wiotp

check:
	curl -s localhost:8347/v1/cpu | jq

clean:
	-docker stop cpu2wiotp
	-docker rm cpu2wiotp
	-docker rmi cpu2wiotp

