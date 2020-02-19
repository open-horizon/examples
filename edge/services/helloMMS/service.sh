#!/bin/bash

# A very simple Horizon sample edge service that shows how to use an Model Management System (MMS) file with your service.
# In this case we use an MMS file as a config file for this service that can be updated dynamically. The service has a default
# copy of the config file built into the docker image. Once the service starts up it periodically checks for a new version of
# the config file using the local MMS API (aka ESS) that the Horizon agent provides to services. If an updated config file is
# found, it is loaded into the service and the config parameters applied (in this case who to say hello to).

# Of course, MMS can also hold and deliver inference models, which can be used by services in a similar way.

OBJECT_TYPE='bp.hello-mms'
OBJECT_ID=config.json

# ${HZN_ESS_AUTH} is mounted to this container and contains a json file with the credentials for authenticating to the ESS.
USER=$(cat ${HZN_ESS_AUTH} | jq -r ".id")
PW=$(cat ${HZN_ESS_AUTH} | jq -r ".token")

# Passing basic auth creds in base64 encoded form (-u).
AUTH="-u ${USER}:${PW}"

# ${HZN_ESS_CERT} is mounted to this container and contains the client side SSL cert to talk to the ESS API.
CERT="--cacert ${HZN_ESS_CERT}"

SOCKET="--unix-socket ${HZN_ESS_API_ADDRESS}"
BASEURL='https://localhost/api/v1/objects'

# Save original config file from the docker image so we can revert back to it if the MMS file is deleted
cp $OBJECT_ID ${OBJECT_ID}.original

# Repeatedly read config.json (initially from the docker image, but then from MMS) and echo hello
while true; do

    # See if there is a new version of the file
    #echo "DEBUG: Checking for MMS updates"
    #OBJ=$(curl -sSL ${AUTH} ${CERT} $SOCKET $BASEURL/$OBJECT_TYPE/$OBJECT_ID)  # this would result in getting the object metadata every call
    HTTP_CODE=$(curl -sSLw "%{http_code}" -o objects.curl ${AUTH} ${CERT} $SOCKET $BASEURL/$OBJECT_TYPE)  # will only get changes that we haven't acknowledged below
    if [[ "$HTTP_CODE" != '200' && "$HTTP_CODE" != '404' ]]; then echo "Error: HTTP code $HTTP_CODE from: curl -sSLw %{http_code} -o objects.curl ${AUTH} ${CERT} $SOCKET $BASEURL/$OBJECT_TYPE"; fi
    #echo "DEBUG: MMS response=$(cat objects.curl)"
    OBJ_ID=$(jq -r ".[] | select(.objectID == \"$OBJECT_ID\") | .objectID" objects.curl)  # if not found, jq returns 0 exit code, but blank value
    if [[ "$HTTP_CODE" == '200' && "$OBJ_ID" == $OBJECT_ID ]]; then
        #echo "DEBUG: Received new metadata for $OBJ_ID"

        # Handle if the MMS file was deleted
        DELETED=$(jq -r ".[] | select(.objectID == \"$OBJECT_ID\") | .deleted" objects.curl)  # if not found, jq returns 0 exit code, but blank value
        if [[ "$DELETED" == "true" ]]; then
            echo "MMS file $OBJECT_ID was deleted, reverting to original $OBJECT_ID"

            # Acknowledge that we saw that it was deleted, so it won't keep telling us
            HTTP_CODE=$(curl -sSLw "%{http_code}" -X PUT ${AUTH} ${CERT} $SOCKET $BASEURL/$OBJECT_TYPE/$OBJECT_ID/deleted)
            if [[ "$HTTP_CODE" != '200' && "$HTTP_CODE" != '204' ]]; then echo "Error: HTTP code $HTTP_CODE from: curl -sSLw %{http_code} -X PUT ${AUTH} ${CERT} $SOCKET $BASEURL/$OBJECT_TYPE/$OBJECT_ID/deleted"; fi

            # Revert back to the original config file from the docker image
            cp ${OBJECT_ID}.original $OBJECT_ID
        
        else
            echo "Received new $OBJECT_ID from MMS"

            # Read the new file from the MMS
            HTTP_CODE=$(curl -sSLw "%{http_code}" -o $OBJECT_ID ${AUTH} ${CERT} $SOCKET $BASEURL/$OBJECT_TYPE/$OBJECT_ID/data)
            if [[ "$HTTP_CODE" != '200' ]]; then echo "Error: HTTP code $HTTP_CODE from: curl -sSLw %{http_code} -o $OBJECT_ID ${AUTH} ${CERT} $SOCKET $BASEURL/$OBJECT_TYPE/$OBJECT_ID/data"; fi
            #ls -l $OBJECT_ID

            # Acknowledge that we got the new file, so it won't keep telling us
            HTTP_CODE=$(curl -sSLw "%{http_code}" -X PUT ${AUTH} ${CERT} $SOCKET $BASEURL/$OBJECT_TYPE/$OBJECT_ID/received)
            if [[ "$HTTP_CODE" != '200' && "$HTTP_CODE" != '204' ]]; then echo "Error: HTTP code $HTTP_CODE from: curl -sSLw %{http_code} -X PUT ${AUTH} ${CERT} $SOCKET $BASEURL/$OBJECT_TYPE/$OBJECT_ID/received"; fi
        fi
    fi

    # Convert all of the key/value pairs in config.json into bash variable assignments
    #HW_WHO=$(jq -r .HW_WHO $OBJECT_ID)
    eval $(jq -r 'to_entries[] | .key + "=\"" + .value + "\""' $OBJECT_ID)

    echo "$HZN_DEVICE_ID says: Hey there ${HW_WHO}!"
    sleep 5
done
