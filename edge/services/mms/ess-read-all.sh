#!/bin/bash
#
# Example program to retrieve data from the ESS inside a Service:
#  - gets a list of the "unreceived" objects from the ESS,
#  - displays each of them (wrapped in JSON) on stdout, and
#  - then marks each of them as "received".
#

# Local ESS REST API endpoint (i.e., the base URL) -- NOTE: use `https`
ESS_ENDPOINT="https://localhost:80/api/v1"

# Get the ESS credentials from the environment
#  - During development, these variables are set by `hzn dev service start`.
#  - In a registered Edge Node, they are set by the local Horizon Agent.
LOCAL_ESS_USER=$(jq -r ".id" ${HZN_ESS_AUTH:-NONE})
LOCAL_ESS_PSWD=$(jq -r ".token" ${HZN_ESS_AUTH:-NONE})

# Specify the object type to search for
objectType="$1"

# Get a list of all objects of this objectType that have not yet been received
OBJECT_ARRAY=$(mktemp)
HTTP_CODE=`curl \
  -sSL \
  -w '%{http_code}' \
  -o ${OBJECT_ARRAY} \
  -u "${LOCAL_ESS_USER}:${LOCAL_ESS_PSWD}" \
  --cacert ${HZN_ESS_CERT} \
  --unix-socket ${HZN_ESS_API_ADDRESS} \
  "${ESS_ENDPOINT}/objects/${objectType}"`
if [ ${HTTP_CODE:-0} -eq 404 ]; then
  echo "Found no unreceived objects." &> /dev/stderr
  echo '[]'
  exit 0
fi
if [ ${HTTP_CODE:-0} -ne 200 ]; then
  echo "ERROR! Object list failed: ${HTTP_CODE}" &> /dev/stderr
  exit 1
fi

# We have a list of objects to traverse (put the noise on stderr)
COUNT=$(jq '.|length' ${OBJECT_ARRAY})
echo "Found ${COUNT} unreceived object(s)..." &> /dev/stderr

# Loop through the objects in the list
echo '['
n=0
for objectID in $(jq -r '.[].objectID' ${OBJECT_ARRAY}); do

  # Retrieve one object (into a temporary file)
  TMP_FILE=$(mktemp)
  HTTP_CODE=`curl -sSL \
    -o ${TMP_FILE} \
    -w '%{http_code}' \
    -u "${LOCAL_ESS_USER}:${LOCAL_ESS_PSWD}" \
    --cacert ${HZN_ESS_CERT} \
    --unix-socket ${HZN_ESS_API_ADDRESS} \
    "${ESS_ENDPOINT}/objects/${objectType}/${objectID}/data"`
  if [ ${HTTP_CODE:-0} -ne 200 ]; then
    echo "ERROR! Object read failed: ${HTTP_CODE}" &> /dev/stderr
    exit 1
  fi
  
  # Echo the data as JSON
  DATA=`cat ${TMP_FILE}`
  if [ "${n}" -ne "0" ]; then
    echo ','
  fi
  echo -n ' {"objectID":"'
  echo -n "${objectID}"
  echo -n '","data":"'
  echo -n "${DATA}"
  echo -n '"}'
  n=$((n + 1))

  # Send a receipt for this object
  HTTP_CODE=`curl -sSL \
    -X PUT \
    -w '%{http_code}' \
    -u "${LOCAL_ESS_USER}:${LOCAL_ESS_PSWD}" \
    --cacert ${HZN_ESS_CERT} \
    --unix-socket ${HZN_ESS_API_ADDRESS} \
    "${ESS_ENDPOINT}/objects/${objectType}/${objectID}/received"`
  if [ ${HTTP_CODE:-0} -ne 204 ]; then
    echo "ERROR! Object receipt failed: ${HTTP_CODE}" &> /dev/stderr
    exit 1
  fi

  # Cleanup the object file
  rm -f ${TMP_+FILE}

done
echo ' '
echo ']'

# Cleanup the object array
rm -f ${OBJECT_ARRAY}

