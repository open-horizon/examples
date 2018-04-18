# Set this to the organization you created in the Watson IoT Platform
export HZN_ORG_ID=myorg

export ARCH=amd64   # arch of your edge node: amd64, or arm for Raspberry Pi, or arm64 for TX2
export VERSION=1.2.2   # the micorservice version, and also used as the tag for the docker image. Must be in OSGI version format.

export DOCKER_HUB_ID=mydockerhubid   # your docker hub username, , sign up at https://hub.docker.com/sso/start/?next=https://hub.docker.com/

# There is normally no need for you to edit these variables
export WIOTP_DOMAIN=internetofthings.ibmcloud.com
export HZN_EXCHANGE_URL="https://$HZN_ORG_ID.$WIOTP_DOMAIN/api/v0002/edgenode/"
