#!/bin/bash

# Very simple Horizon sample edge service that uses a secret
# This service uses a single secret called "hw_who". In the pattern/policy, this secret will be 
# bound to the organization-level secret "secret-name" in the secrets manager. The secret 
# consists of a key/value pair, but this code only utilizes the value of the secret, and the key 
# is ignored. The code prints out "Hello <hw_who>!"

# filepath to where the secret we are using is stored in the container 
SECRET_NAME="hw_who"
FILEPATH="open-horizon-secrets/$SECRET_NAME"

# ${HZN_ESS_AUTH} is mounted to this container by the Horizon agent and is a json file with the credentials for authenticating to ESS.
USER=$(cat ${HZN_ESS_AUTH} | jq -r ".id")
PW=$(cat ${HZN_ESS_AUTH} | jq -r ".token")

# Some curl parameters for using the ESS REST API
AUTH="-u ${USER}:${PW}"
# ${HZN_ESS_CERT} is mounted to this container by the Horizon agent and the cert clients use to verify the identity of ESS.
CERT="--cacert ${HZN_ESS_CERT}"
SOCKET="--unix-socket ${HZN_ESS_API_ADDRESS}"

# SECRETS API CALLS 
SECRETS_API="https://localhost/api/v1/secrets"

# queries for updated secrets that have not been acknowledged (by POST)
UPDATE_CMD="curl --silent -w %{http_code} -X GET ${SECRETS_API} ${AUTH} ${CERT} ${SOCKET}"

# queries for the details (key/value) of updated secrets
GET_SECRET="curl --silent -w %{http_code} -X GET ${SECRETS_API}/${SECRET_NAME} ${AUTH} ${CERT} ${SOCKET}"

# acknowledges that the updated secret is received (secret will no longer be returned by UPDATE_CMD)
function make_post_cmd {
  POST_SECRET="curl --silent -w %{http_code} -X POST ${SECRETS_API}/$1?received=true ${AUTH} ${CERT} ${SOCKET}"
}

# UPDATING SECRETS

# runs the POST_SECRET command to acknowledge that the secret update
# has been received. the secret will no longer be returned on subsequent 
# calls to the UPDATE_CMD. this function takes in the name of the secret 
# whose update we are acknowledging
function acknowledge_secret_update {
  make_post_cmd "$1"
  HTTP_CODE=$($POST_SECRET)
  if [ "$HTTP_CODE" != "201" ]; then 
    echo "Error: HTTP code $HTTP_CODE from: $POST_CMD"
    exit 1
  fi 
}

# update secret via the API - query for the updated secret details (after the secret 
# is returned by UPDATE_CMD) and store the secret value (not key) in HW_WHO
function update_secret_API {
  # query the API for the updated secret 
  SECRET_DETAILS=$($GET_SECRET)
  HTTP_CODE=${SECRET_DETAILS: -3}
  if [ "$HTTP_CODE" != "200" ]; then 
    echo "Error: HTTP code $HTTP_CODE from: $GET_SECRET"
    exit 1
  fi

  # parse the response for the updated secret value
  HW_WHO=$(echo -e "${SECRET_DETAILS::-3}" | jq '.value')
  HW_WHO=${HW_WHO:1:-1}
}

# update secret via file - read and parse the contents of the open-horizon-secrets/<secret-name>
# file contents are a json with fields "key" and "value", the "value" is stored in HW_WHO
function update_secret_file {
  HW_WHO=$(cat $FILEPATH | jq '.value') 
  HW_WHO=${HW_WHO:1:-1}
}

# read the initial secret value from the file
update_secret_file

# checks for updated secrets and updates the local variable
function check_secret_updates {
  # query the API for any changes 
  UPDATES=$($UPDATE_CMD)
  HTTP_CODE=${UPDATES: -3}

  if [ "$HTTP_CODE" == "200" ]; then 

    # loop through the returned secrets 
    for row in $(echo ${UPDATES::-3} | jq -r '.[]'); do 

      # check if we need the updated secret
      if [ "$row" == "$SECRET_NAME" ]; then 
        update_secret_file # can also use update_secret_API
      fi

      # acknowledge that the secret update has been received 
      acknowledge_secret_update "$row"
    done      
  elif [ "$HTTP_CODE" -ne "404" ]; then 
    # we got a code other than 200 or 404 
    echo "Error: HTTP code $HTTP_CODE from: $UPDATE_CMD"
    exit 1
  fi 
}

while true; do
  # query for any changes in the secret(s)
  check_secret_updates

  # print the greeting
  echo "$HZN_DEVICE_ID says: Hello ${HW_WHO}!"
  sleep 3
done
