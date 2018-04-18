# Set these values to the objects and credentials you created in the Watson IoT Platform
export HZN_ORG_ID=myorg
export WIOTP_GW_TYPE=mygwtype
export WIOTP_GW_ID=mygwinstance
export WIOTP_GW_TOKEN='mygwinstancetoken'
export WIOTP_API_KEY='a-myapikeyrandomchars'
export WIOTP_API_TOKEN='myapikeytoken'

# This variable must be set appropriately for your specific Edge Node
export ARCH=amd64   # or arm for Raspberry Pi, or arm64 for TX2
export VERSION=1.2.2   # used as the tag for the docker image

# There is normally no need for you to edit these variables
export WIOTP_DOMAIN=internetofthings.ibmcloud.com
export HZN_EXCHANGE_URL="https://$HZN_ORG_ID.$WIOTP_DOMAIN/api/v0002/edgenode/"
export HZN_DEVICE_ID="g@${WIOTP_GW_TYPE}@$WIOTP_GW_ID"
export WIOTP_CLIENT_ID_APP="a:$HZN_ORG_ID:$WIOTP_GW_TYPE$WIOTP_GW_ID"
export WIOTP_CLIENT_ID_GW="g:$HZN_ORG_ID:$WIOTP_GW_TYPE:$WIOTP_GW_ID"
export HZN_EXCHANGE_USER_AUTH="$WIOTP_API_KEY:$WIOTP_API_TOKEN"
export HZN_EXCHANGE_API_AUTH="$WIOTP_API_KEY:$WIOTP_API_TOKEN"