#!/bin/bash

# Very simple Horizon sample using the MMS feature to update HW_WHO

# ${HZN_ESS_AUTH} is mounted to this container and contains a json file with the credentials for authenticating to the ESS.
USER=$(cat ${HZN_ESS_AUTH} | jq -r ".id")
PW=$(cat ${HZN_ESS_AUTH} | jq -r ".token")

# Passing basic auth creds in base64 encoded form (-u).
AUTH="-u ${USER}:${PW} "

# ${HZN_ESS_CERT} is mounted to this container and contains the client side SSL cert to talk to the ESS API.
CERT="--cacert ${HZN_ESS_CERT} "

BASEURL='--unix-socket '${HZN_ESS_API_ADDRESS}' https://localhost/api/v1/objects/'

declare -a nameArray
declare -a valArray

while true; do
    # get names of inputs into nameArray
    nameArray=($(jq -r '.userInput[].inputs[].name' input.json))

    # get values of inputs into valArray
    valArray=($(jq -r '.userInput[].inputs[].value' input.json))

    # export name/value pairs
    x=0
    for i in "${nameArray[@]}"; do
        eval export ${nameArray[$x]}=${valArray[$x]}
        x=$((x+1))
    done

    echo "$HZN_DEVICE_ID says: Hello ${HW_WHO}!!"
    sleep 5

    # read in new file from the ESS
    DATA=$(curl -sL -o /input.json ${AUTH}${CERT}${BASEURL}json/input.json/data)
done
