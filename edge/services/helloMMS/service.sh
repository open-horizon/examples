#!/bin/sh

# Very simple Horizon sample edge service.

# ${HZN_ESS_AUTH} is mounted to this container and contains a json file with the credentials for authenticating to the ESS.
USER=$(cat ${HZN_ESS_AUTH} | jq -r ".id")
PW=$(cat ${HZN_ESS_AUTH} | jq -r ".token")

# Passing basic auth creds in base64 encoded form (-u).
AUTH="-u ${USER}:${PW} "

# ${HZN_ESS_CERT} is mounted to this container and contains the client side SSL cert to talk to the ESS API.
CERT="--cacert ${HZN_ESS_CERT} "

BASEURL='--unix-socket '${HZN_ESS_API_ADDRESS}' https://localhost/api/v1/objects/'


while true; do
    userinput=$(cat input.json | jq '.[]' | jq '.[]' | jq '.[]' | jq '.[] .value')
    #echo "$HZN_DEVICE_ID:"

    echo $userinput
    sleep 2

    DATA=$(curl -sL -o /input.json ${AUTH}${CERT}${BASEURL}json/input.json/data)
done
