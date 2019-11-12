#!/bin/bash

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
    declare -a nameArray
    declare -a valArray

    # get names of inputs into nameArray
    x=0
    for i in $(jq '.userInput[].inputs[].name' input.json); do
        nameArray[x]="$i"
        x=$((x+1))
    done

    # get values of inputs into valArray
    y=0
    for j in $(jq '.userInput[].inputs[].value' input.json); do
        valArray[y]="$j"
        y=$((y+1))
    done

    # search for new HW_WHO value in updated inputs section
    z=0
    for k in "${nameArray[@]}"; do
        if [ "$k" == "\"HW_WHO\"" ]; then
            HW_WHO=${valArray[$z]}
        fi
        z=$((z+1))
    done

    echo "$HZN_DEVICE_ID says: Hello ${HW_WHO}!!"

    sleep 5
    DATA=$(curl -sL -o /input.json ${AUTH}${CERT}${BASEURL}json/input.json/data)
done
