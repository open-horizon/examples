#!/bin/bash
set -e

# general parameters
WAIT_RESPONSE=10
SERVICE_PREFIX="sdr-poc"        # prefix for all SDR PoC services created in a cloud

# event streams/message hub configuration
MH_INSTANCE="${SERVICE_PREFIX}-es"
MH_INSTANCE_CREDS="${MH_INSTANCE}-credentials"
MH_SDR_TOPIC="sdr-audio"
MH_SDR_TOPIC_PARTIONS=2
MH_RESPONSE_RETRY=5

# watson speech-to-text service
STT_INSTANCE="${SERVICE_PREFIX}-speech-to-text"
STT_INSTANCE_CREDS="${STT_INSTANCE}-credentials"
STT_INSTANCE_PLAN="standard"

# watson natural language understanding
NLU_INSTANCE="${SERVICE_PREFIX}-natural-language-understanding"
NLU_INSTANCE_CREDS="${NLU_INSTANCE}-credentials"
NLU_INSTANCE_PLAN="standard"

# db
DB_INSTANCE="${SERVICE_PREFIX}-compose-for-postgresql"
DB_INSTANCE_CREDS="${DB_INSTANCE}-credentials"
DB_INSTANCE_PLAN="Standard"
DB_INSTANCE_RETRY=20
DB_INSTANCE_TIMEOUT=60*20
DB_NAME="sdr"
DB_SQL='CREATE DATABASE '$DB_NAME'; \l \dt \c '$DB_NAME' \\CREATE TABLE globalnouns( 
		noun TEXT PRIMARY KEY NOT NULL, 
		sentiment DOUBLE PRECISION NOT NULL, 
		numberofmentions BIGINT NOT NULL, 
		timeupdated timestamp with time zone); 
	CREATE TABLE nodenouns( 
		noun TEXT NOT NULL, 
		edgenode TEXT NOT NULL, 
		sentiment DOUBLE PRECISION NOT NULL, 
		numberofmentions BIGINT NOT NULL, 
		timeupdated timestamp with time zone, 
		PRIMARY KEY(noun, edgenode)); 
	CREATE TABLE stations( 
		edgenode TEXT NOT NULL, 
		frequency REAL NOT NULL, 
		numberofclips BIGINT NOT NULL, 
		dataqualitymetric REAL, 
		timeupdated timestamp with time zone, 
		PRIMARY KEY(edgenode, frequency)); 
	CREATE TABLE edgenodes( 
		edgenode TEXT PRIMARY KEY NOT NULL, 
		latitude REAL NOT NULL, 
		longitude REAL NOT NULL, 
		timeupdated timestamp with time zone);'

# functions
FUNC_PACKAGE="${SERVICE_PREFIX}-message-hub-evnts"
FUNC_MH_FEED="Bluemix_${MH_INSTANCE}_${MH_INSTANCE_CREDS}/messageHubFeed"
FUNC_TRIGGER="${SERVICE_PREFIX}-message-received-trigger"
FUNC_ACTION="${SERVICE_PREFIX}-process-message"
FUNC_ACTION_CODE="../../data-processing/ibm-functions/actions/msgreceive.js"
FUNC_RULE="${SERVICE_PREFIX}-message-received-rule"

# ui
UI_SRC_PATH="../../ui/sdr-app"
UI_APP_NAME="${SERVICE_PREFIX}-sdr-poc-app"
UI_APP_ID_INSTANCE="${SERVICE_PREFIX}-app-id"
UI_APP_ID_INSTANCE_PLAN="graduated-tier"
UI_APP_ID_INSTANCE_LOCATION="us-south"
UI_APP_ID_INSTANCE_ALIAS="${UI_APP_ID_INSTANCE}-alias"
UI_APP_ID_INSTANCE_ALIAS_CREDS="${UI_APP_ID_INSTANCE_ALIAS}-credentials"

# variables used for configuring the script and displaying its configuration
config=( MH_INSTANCE MH_INSTANCE_CREDS MH_SDR_TOPIC MH_SDR_TOPIC_PARTIONS MH_RESPONSE_RETRY \
		STT_INSTANCE STT_INSTANCE_CREDS STT_INSTANCE_PLAN \
		NLU_INSTANCE NLU_INSTANCE_CREDS  NLU_INSTANCE_PLAN \
		DB_INSTANCE DB_INSTANCE_CREDS DB_INSTANCE_PLAN DB_NAME \
		FUNC_PACKAGE FUNC_MH_FEED FUNC_TRIGGER FUNC_ACTION FUNC_ACTION_CODE FUNC_RULE \
		UI_SRC_PATH UI_APP_NAME UI_APP_ID_INSTANCE UI_APP_ID_INSTANCE_PLAN \
		UI_APP_ID_INSTANCE_LOCATION UI_APP_ID_INSTANCE_ALIAS UI_APP_ID_INSTANCE_ALIAS_CREDS )

# supported components
components=( stt nlu db es func ui all )

# script prerequisites
prerequisites=( jq curl ibmcloud psql sed grep cut npm node )

# supported actions
DEPLOY="deploy"
TEARDOWN="teardown"

function help() {
  cat << EndOfMessage
$(basename "$0") [ [-i|--install=<component>] || [-u|--uninstall=<component>]] -- deploying cloud part for SDR PoC

where:
	-p | --prereqs					- check for prerequisites
	-c | --config					- show current configuration
	-i= | --install=[component]			- install [component]
	-u= | --uninstall=[component]			- uninstall [component]

Example: ./deploy.sh --install=all
	
EndOfMessage
	echo "Supported components are: "
	for i in "${components[@]}"; do echo -n "${i} "; done
}

function now(){
	echo `date '+%Y-%m-%d %H:%M:%S'`
}

# Show configuration
function show_config() {
	echo `now` "Current configuration:"
	echo "============================================================="
	for i in "${config[@]}"; do echo "${i} = ""${!i}"; done
	if [[ -z "${UI_APP_USER}" ]]; then
		echo "UI_APP_USER environment variable is not defined"
	else
		echo "UI_APP_USER = ${UI_APP_USER}"
	fi
	if [[ -z "${UI_APP_PASSWORD}" ]]; then
		echo "UI_APP_PASSWORD environment variable is not defined"
	else
		echo "UI_APP_PASSWORD is defined"
	fi
	if [[ -z "${MAPBOX_TOKEN}" ]]; then
		echo "MAPBOX_TOKEN environment variable is not defined"
	else
		echo "MAPBOX_TOKEN = ${MAPBOX_TOKEN}"
	fi
	echo "============================================================="
}

# Checking prerequisites
function check_prereqs(){
	echo `now` "Checking prerequisites:"
	echo "============================================================="
	for i in "${prerequisites[@]}"; do
	 	if command -v ${i} >/dev/null 2>&1; then
			echo "Found ${i}"
		else
			echo "Not found ${i}, exiting..."
			exit 1
		fi
	done
	: "${UI_APP_USER:?UI_APP_USER is not set or empty}"
	: "${UI_APP_PASSWORD:?UI_APP_PASSWORD is not set or empty}"
	: "${MAPBOX_TOKEN:?MAPBOX_TOKEN is not set or empty}"
	echo `now` "All prerequisites are met"
	echo "============================================================="
}

# Create Watson Speech-To-Text instance
function deploy_stt_(){
	echo `now` "Creating $STT_INSTANCE Watson Speech-To-Text service instance with the $STT_INSTANCE_PLAN plan"
	echo `now` "Current Watson Speech-To-Text instances:"
	ibmcloud -q service list | grep speech_to_text || :
	if [[ $(ibmcloud -q service list | grep speech_to_text | cut -d' ' -f1 | grep -Fx "$STT_INSTANCE") ]]; then
	 	echo `now` "There is $STT_INSTANCE Watson STT instance created already, skipping its creation..."
	else
	 	echo `now` "Found no Watson STT instance $STT_INSTANCE with the $STT_INSTANCE_PLAN plan, creating..."
	 	ibmcloud -q service create speech_to_text "$STT_INSTANCE_PLAN" "$STT_INSTANCE"
	 	ibmcloud -q service list | grep speech_to_text
	fi
	echo `now` "Creating credentials $STT_INSTANCE_CREDS for $STT_INSTANCE Watson STT instance"
	echo `now` "Current credentials for ${STT_INSTANCE}:"
	ibmcloud -q service keys "$STT_INSTANCE" | sed -n '3!p' | sed -n '1!p'
	if [[ $(ibmcloud service keys "$STT_INSTANCE" | cut -d' ' -f1 | grep -Fx "$STT_INSTANCE_CREDS") ]]; then
	 	echo `now` "There is $STT_INSTANCE_CREDS credentials for $STT_INSTANCE Watson STT instance created, skipping its creation..."
	else
	 	echo `now` "Found no ${STT_INSTANCE_CREDS}, creating..."
	 	ibmcloud service key-create "$STT_INSTANCE" "$STT_INSTANCE_CREDS"
	fi
	echo `now` "Creating and configuring $STT_INSTANCE Watson STT is finished"
	echo `now` "Current Watson STT instances:"
	ibmcloud -q service list | grep speech_to_text || :
}

# Delete Watson Speech-To-text instance
function teardown_stt_(){
	echo `now` "Deleting $STT_INSTANCE Watson STT instance..."
	echo `now` "Waiting $WAIT_RESPONSE seconds..."
	sleep $WAIT_RESPONSE
	echo `now` "Deleting $STT_INSTANCE_CREDS credentials for $STT_INSTANCE"
	echo `now` "Current Watson STT instances:"
	ibmcloud -q service list | grep speech_to_text || :
	if [[ $(ibmcloud service keys "$STT_INSTANCE" | cut -d' ' -f1 | grep -Fx "$STT_INSTANCE_CREDS") ]]; then
	 	echo `now` "Found $STT_INSTANCE_CREDS credentials for $STT_INSTANCE, deleting..."
	 	ibmcloud service key-delete "$STT_INSTANCE" "$STT_INSTANCE_CREDS" -f
	else
	 	echo `now` "There is no $STT_INSTANCE_CREDS credentials for $STT_INSTANCE to delete, skipping..."
	fi
	echo `now` "Deleting $STT_INSTANCE..."
	if [[ $(ibmcloud -q service list | grep speech_to_text | cut -d' ' -f1 | grep -Fx "$STT_INSTANCE") ]]; then
	 	echo `now` "Found $STT_INSTANCE Wastson STT instance, deleting..."
		ibmcloud service delete "$STT_INSTANCE" -f
	else
	 	echo `now` "Found no $STT_INSTANCE Watson STT instance, skipping...";
	fi
	echo `now` "Finished deleting $STT_INSTANCE and its credentials $STT_INSTANCE_CREDS"
	echo `now` "Current Watson STT instances:"
	ibmcloud -q service list | grep speech_to_text || :
}

# Create Watson Natural Language Understanding instance
function deploy_nlu_(){
	echo `now` "Creating $NLU_INSTANCE Watson Natural Language Understanding service instance with the $NLU_INSTANCE_PLAN plan"
	echo `now` "Current Watson Natural Language Understanding instances:"
	ibmcloud -q service list | grep natural-language-understanding || :
	if [[ $(ibmcloud -q service list | grep natural-language-understanding | cut -d' ' -f1 | grep -Fx "$NLU_INSTANCE") ]]; then
	 	echo `now` "There is $NLU_INSTANCE Watson NLU instance created already, skipping its creation..."
	else
	 	echo `now` "Found no Watson NLU instance $NLU_INSTANCE with the $NLU_INSTANCE_PLAN plan, creating..."
	 	ibmcloud -q service create natural-language-understanding "$NLU_INSTANCE_PLAN" "$NLU_INSTANCE"
	 	ibmcloud -q service list | grep natural-language-understanding  || :
	fi
	echo `now` "Creating credentials $NLU_INSTANCE_CREDS for $NLU_INSTANCE"
	echo `now` "Current credentials for $NLU_INSTANCE:"
	ibmcloud -q service keys "$NLU_INSTANCE" | sed -n '3!p' | sed -n '1!p'
	if [[ $(ibmcloud service keys "$NLU_INSTANCE" | cut -d' ' -f1 | grep -Fx "$NLU_INSTANCE_CREDS") ]]; then
	 	echo `now` "There is $NLU_INSTANCE_CREDS credentials for $NLU_INSTANCE Watson NLU instance created, skipping its creation..."
	else
	 	echo `now` "Found no $NLU_INSTANCE_CREDS, creating..."
	 	ibmcloud service key-create "$NLU_INSTANCE" "$NLU_INSTANCE_CREDS"
	fi
	echo `now` "Creating and configuring $NLU_INSTANCE Watson NLU is finished"
	echo `now` "Current Watson NLU instances:"
	ibmcloud -q service list | grep natural-language-understanding || :
}

# Delete Watson Natural Language Understanding instance
function teardown_nlu_(){
	echo `now` "Deleting $NLU_INSTANCE Watson NLU..."
	echo `now` "Waiting $WAIT_RESPONSE seconds..."
	sleep $WAIT_RESPONSE
	echo `now` "Deleting $NLU_INSTANCE_CREDS credentials for $NLU_INSTANCE"
	echo `now` "Current Watson NLU instances:"
	ibmcloud -q service list | grep natural-language-understanding || :
	if [[ $(ibmcloud service keys "$NLU_INSTANCE" | cut -d' ' -f1 | grep -Fx "$NLU_INSTANCE_CREDS") ]]; then
	 	echo `now` "Found $NLU_INSTANCE_CREDS credentials for $NLU_INSTANCE, deleting..."
	 	ibmcloud service key-delete "$NLU_INSTANCE" "$NLU_INSTANCE_CREDS" -f
	else
	 	echo `now` "There is no $NLU_INSTANCE_CREDS credentials for $NLU_INSTANCE to delete, skipping..."
	fi
	echo `now` "Deleting $NLU_INSTANCE..."
	if [[ $(ibmcloud -q service list | grep natural-language-understanding | cut -d' ' -f1 | grep -Fx "$NLU_INSTANCE") ]]; then
	 	echo `now` "Found $NLU_INSTANCE, deleting..."
	 	ibmcloud service delete "$NLU_INSTANCE" -f
	else
	 	echo `now` "Found no $NLU_INSTANCE Watson NLU instance, skipping..."
	fi
	echo `now` "Finished deleting $NLU_INSTANCE and its credentials $NLU_INSTANCE_CREDS"
	echo `now` "Current Watson NLU instances:"
	ibmcloud -q service list | grep natural-language-understanding || :
}

# Create and configure Event Streams Instance
function deploy_es_(){
	echo `now` "Creating Event Streams instance $MH_INSTANCE"
	echo `now` "Current Event Streams instances:"
	ibmcloud -q service list | grep messagehub || :
	if [[ $(ibmcloud -q service list | grep messagehub | cut -d' ' -f1 | grep -Fx "$MH_INSTANCE") ]]; then
	 	echo `now` "There is $MH_INSTANCE Event Streams instance created already, skipping its creation..."
	else
	 	echo `now` "Found no Event Streams instance $MH_INSTANCE, creating..."
	 	ibmcloud -q service create messagehub standard "$MH_INSTANCE"
	 	ibmcloud -q service list | grep messagehub || :
	fi
	echo `now` "Creating credentials $MH_INSTANCE_CREDS for $MH_INSTANCE"
	echo `now` "Current credentials for $MH_INSTANCE:"
	ibmcloud -q service keys "$MH_INSTANCE" | sed -n '3!p' | sed -n '1!p'
	if [[ $(ibmcloud service keys "$MH_INSTANCE" | cut -d' ' -f1 | grep -Fx "$MH_INSTANCE_CREDS") ]]; then
	 	echo `now` "There is $MH_INSTANCE_CREDS credentials for Event Streams $MH_INSTANCE created, skipping its creation..."
	else
	 	echo `now` "Found no $MH_INSTANCE_CREDS, creating..."
	 	ibmcloud service key-create "$MH_INSTANCE" "$MH_INSTANCE_CREDS"
	fi
	while [ true ] ; do
	 	response="$(ibmcloud service key-show "$MH_INSTANCE" "$MH_INSTANCE_CREDS" | sed -n '3!p' | sed -n '1!p')"
	 	echo `now` "Checking if $MH_INSTANCE and $MH_INSTANCE_CREDS are ready"
	 	if [[ "$response" = *"not found"* ]] ; then
	 		echo `now` "$MH_INSTANCE and $MH_INSTANCE_CREDS seem not ready"
	 		echo `now` "Retrying in $MH_RESPONSE_RETRY seconds"
	 		sleep $MH_RESPONSE_RETRY
	 		continue
	 	else
	 		echo `now` "$MH_INSTANCE and $MH_INSTANCE_CREDS are ready"
	 		break
	 	fi
	done
	echo `now` "Creating $MH_SDR_TOPIC with $MH_SDR_TOPIC_PARTIONS partitions on $MH_INSTANCE"
	response="$(ibmcloud service key-show "$MH_INSTANCE" "$MH_INSTANCE_CREDS" | sed -n '3!p' | sed -n '1!p')"
	admin_url="$(echo "$response" | jq -r '.kafka_admin_url')"
	api_key="$(echo "$response" | jq -r '.api_key')"
	if [[ $(curl -s -H 'Accept: application/json' -H 'X-Auth-Token: '"$api_key" "${admin_url}"/admin/topics/ | \
		jq -c '.[] | select(.name | . and contains('"\"$MH_SDR_TOPIC\""'))') ]]; then
	 	echo `now` "$MH_SDR_TOPIC topic is already created on $MH_INSTANCE"
	else
	 	echo `now` "Found no $MH_SDR_TOPIC, creating..."
	 	curl -s -H 'Content-Type: application/json' -H 'Accept: */*' -H 'X-Auth-Token: '"$api_key" \
	 		-d '{ "name": '"\"$MH_SDR_TOPIC\""', "partitions": '"$MH_SDR_TOPIC_PARTIONS"' }' \
				"${admin_url}"/admin/topics
		echo `now` "Current topic(s) on $MH_INSTANCE:"
		curl -s -H 'Accept: application/json' -H 'X-Auth-Token: '"$api_key" \
				"${admin_url}"/admin/topics/
	fi 
	echo `now` "Finished creation $MH_INSTANCE instance with $MH_SDR_TOPIC topic, $MH_SDR_TOPIC_PARTIONS partition(s)"
	echo `now` "Current Event Streams instances:"
	ibmcloud -q service list | grep messagehub || :
}

# Delete Event Streams instance
function teardown_es_(){
	echo `now` "Deleting $MH_INSTANCE Event Streams instance..."
	echo `now` "Waiting $WAIT_RESPONSE seconds..."
	sleep $WAIT_RESPONSE
	echo `now` "Deleting $MH_INSTANCE_CREDS for $MH_INSTANCE"
	echo `now` "Current Event Streams instances:"
	ibmcloud -q service list | grep messagehub || :
	if [[ $(ibmcloud service keys "$MH_INSTANCE" | cut -d' ' -f1 | grep -Fx "$MH_INSTANCE_CREDS") ]]; then
	 	echo `now` "Found $MH_INSTANCE_CREDS credentials for $MH_INSTANCE, deleting..."
	 	ibmcloud service key-delete "$MH_INSTANCE" "$MH_INSTANCE_CREDS" -f
	else
	 	echo `now` "There is no $MH_INSTANCE_CREDS credentials for $MH_INSTANCE to delete, skipping..."
	fi
	echo `now` "Deleting $MH_INSTANCE..."
	if [[ $(ibmcloud -q service list | grep messagehub | cut -d' ' -f1 | grep -Fx "$MH_INSTANCE") ]]; then
	 	echo `now` "Found $MH_INSTANCE, deleting..."
	 	ibmcloud service delete "$MH_INSTANCE" -f
	else
	 	echo `now` "Found no $MH_INSTANCE Event Streams instance, skipping..."
	fi
	echo `now` "Finished deleting $MH_INSTANCE and its credentials $MH_INSTANCE_CREDS"
	echo `now` "Current Event Streams instances:"
	ibmcloud -q service list | grep messagehub || :
}

# Create and configure Compose for PostgreSQL instance
function deploy_db_(){
	echo `now` "Creating $DB_INSTANCE Compose PostgreSQL DB with the $DB_INSTANCE_PLAN plan"
	echo `now` "Current Compose for PostgreSQL instances:"
	ibmcloud -q service list | grep compose-for-postgresql || :
	if [[ $(ibmcloud -q service list | grep compose-for-postgresql | cut -d' ' -f1 | grep -Fx "$DB_INSTANCE") ]]; then
	 	echo `now` "There is $DB_INSTANCE Compose for PostgreSQL instance created already, skipping its creation..."
	else
	 	echo `now` "Found no Compose for PostgreSQL instance $DB_INSTANCE with $DB_INSTANCE_PLAN plan, creating..."
	 	ibmcloud -q service create compose-for-postgresql "$DB_INSTANCE_PLAN" "$DB_INSTANCE"
	 	ibmcloud -q service list | grep compose-for-postgresql || :
	fi
	echo `now` "Creating credentials $DB_INSTANCE_CREDS for $DB_INSTANCE"
	echo `now` "Current credentials for $DB_INSTANCE:"
	ibmcloud -q service keys "$DB_INSTANCE" | sed -n '3!p' | sed -n '1!p'
	if [[ $(ibmcloud service keys "$DB_INSTANCE" | cut -d' ' -f1 | grep -Fx "$DB_INSTANCE_CREDS") ]]; then
	 	echo `now` "There is $DB_INSTANCE_CREDS credentials for $DB_INSTANCE PostgreSQL instance created, skipping its creation..."
	else
	 	echo `now` "Found no $DB_INSTANCE_CREDS, creating..."
	 	ibmcloud service key-create "$DB_INSTANCE" "$DB_INSTANCE_CREDS"
	fi
	echo `now` "Bootstrapping database $DB_NAME on $DB_INSTANCE"
	echo `now` "Checking if $DB_INSTANCE PostgreSQL instance is ready... Please wait"
	db_uri="$(ibmcloud service key-show "$DB_INSTANCE" "$DB_INSTANCE_CREDS" | sed -n '3!p' | sed -n '1!p' | jq -r '.uri')"
	start_time_db_check=`date +%s`
	while [ true ] ; do
		response="$(psql "$db_uri?sslmode=require" -c '\l' 2>&1 | sed -n 1p )"
		current_time_db_check=`date +%s`
	 	if [[ "$response" = *"Operation timed out"* ]] || [[ "$response" = *"Connection refused"* ]] ; then
	 		echo `now` $DB_INSTANCE is not ready, response is "$response"
			echo `now` "current time is ${current_time_db_check}, start time is ${start_time_db_check}"
			if (( current_time_db_check - start_time_db_check > DB_INSTANCE_TIMEOUT )); then
				echo `now` "PostgreSQL instance provisioning timeout of $DB_INSTANCE_TIMEOUT seconds occured"
				exit 1
			fi
	 		echo `now` "Retrying in $DB_INSTANCE_RETRY seconds"
	 		sleep $DB_INSTANCE_RETRY
			continue
	 	else
	 		echo `now` "$DB_INSTANCE is ready"
	 		break
	 	fi
    done
	if [[ $(psql "$db_uri?sslmode=require" -c '\l' | sed '1,3d;$d' | cut -d'|' -f1 | cut -d' ' -f2 | grep -Fx "$DB_NAME") ]]; then
	 	echo `now` "There is $DB_NAME database already created on $DB_INSTANCE instance, skipping"
	else
	 	echo `now` "Found no $DB_NAME on $DB_INSTANCE instance, creating"
		echo $DB_SQL | psql "$db_uri?sslmode=require"
		echo `now` "Created $DB_NAME on $DB_INSTANCE instance"
	 	echo `now` "Created the following tables in $DB_NAME"
	 	db_sdr="$(echo "$db_uri" | sed -e "s/compose/${DB_NAME}/g")"
	 	psql "$db_sdr?sslmode=require" -c '\dt'
	fi
	echo `now` "Creating and configuring $DB_INSTANCE Compose for PostgreSQL is finished"
	echo `now` "Current Compose for PostgreSQL instances:"
	ibmcloud -q service list | grep compose-for-postgresql || :
}

# Delete Compose for PostgreSQL instance
function teardown_db_(){
	echo `now` "Deleting $DB_INSTANCE Compose for PostgreSQL instance..."
	echo `now` "Waiting $WAIT_RESPONSE seconds..."
	sleep $WAIT_RESPONSE
	echo `now` "Deleting $DB_INSTANCE_CREDS credentials for $DB_INSTANCE"
	echo `now` "Current Compose for PostgreSQL instances:"
	ibmcloud -q service list | grep compose-for-postgresql || :
	if [[ $(ibmcloud service keys "$DB_INSTANCE" | cut -d' ' -f1 | grep -Fx "$DB_INSTANCE_CREDS") ]]; then
	 	echo `now` "Found $DB_INSTANCE_CREDS credentials for $DB_INSTANCE, deleting..."
	 	ibmcloud service key-delete "$DB_INSTANCE" "$DB_INSTANCE_CREDS" -f
	else
	 	echo `now` "There is no $DB_INSTANCE_CREDS credentials for $DB_INSTANCE to delete, skipping..."
	fi
	echo `now` "Deleting $DB_INSTANCE..."
	if [[ $(ibmcloud -q service list | grep compose-for-postgresql | cut -d' ' -f1 | grep -Fx "$DB_INSTANCE") ]]; then
	 	echo `now` "Found $DB_INSTANCE, deleting..."
	 	ibmcloud service delete "$DB_INSTANCE" -f
	else
	 	echo `now` "Found no $DB_INSTANCE Compose for PostgreSQL instance, skipping..."
	fi
	echo `now` "Finished deleting $DB_INSTANCE and its credentials $DB_INSTANCE_CREDS"
	echo `now` "Current Compose for PostgreSQL instances:"
	ibmcloud -q service list | grep compose-for-postgresql || :
}

# Create and configure functions
function deploy_func_(){
	echo `now` "Creating functions entities"
	echo `now` "Current functions entities:"
	ibmcloud -q fn list
	echo `now` "Creating bindings for $MH_INSTANCE_CREDS credentials from $MH_INSTANCE Event Streams instance and functions"
	ibmcloud fn package refresh
	if [ $? -eq 0 ]; then
		echo `now` "Successfully created bindings"
	else
		echo `now` "Failed to create bindings, exiting..." >&2
		exit 1
	fi
	echo `now` "Creating trigger $FUNC_TRIGGER for $FUNC_MH_FEED feed"
	echo `now` "Current triggers:"
	ibmcloud -q fn trigger list
	org="$(ibmcloud target | grep Org | cut -d : -f 2 | sed -e 's/^[ \t]*//' | cut -d ' ' -f 1)"
	space="$(ibmcloud target | grep Space | cut -d : -f 2 | sed -e 's/^[ \t]*//' | cut -d ' ' -f 1)"
	trigger="/${org}_${space}/${FUNC_TRIGGER}"
	if [[ $(ibmcloud -q fn trigger list | cut -d' ' -f1 | grep -Fx "$trigger") ]]; then
	 	echo `now` "There is $trigger trigger created already, skipping its creation..."
	else
	 	echo `now` "Found no $trigger trigger for $FUNC_MH_FEED feed and $MH_SDR_TOPIC topic, creating..."
		ibmcloud -q fn trigger create "$FUNC_TRIGGER" \
	 		--feed "$FUNC_MH_FEED" \
	 		--param isJSONData true \
	 		--param isBinaryValue false \
	 		--param topic "$MH_SDR_TOPIC"
		ibmcloud -q fn trigger list
	fi
	echo `now` "Creating $FUNC_PACKAGE package for $FUNC_ACTION action..."
	echo `now` "Current packages:"
	ibmcloud -q fn package list
	package="/${org}_${space}/${FUNC_PACKAGE}"
	if [[ $(ibmcloud -q fn package list | cut -d' ' -f1 | grep -Fx "$package") ]]; then
	 	echo `now` "There is $package package created already, skipping its creation..."
	else
	 	echo `now` "Found no $package package, creating..."
	 	ibmcloud fn package create "$FUNC_PACKAGE"
	 	ibmcloud -q fn package list
	fi
	echo `now` "Creating $FUNC_ACTION action in $FUNC_PACKAGE package..."
	echo `now` "Current actions:"
	ibmcloud -q fn action list
	action="/${org}_${space}/${FUNC_PACKAGE}/${FUNC_ACTION}"
	if [[ $(ibmcloud -q fn action list | cut -d' ' -f1 | grep -Fx "$action") ]]; then
	 	echo `now` "There is $action action created already, skipping its creation..."
	else
	 	echo `now` "Found no $action action, creating..."
	 	stt_response="$(ibmcloud service key-show "$STT_INSTANCE" "$STT_INSTANCE_CREDS" | sed -n '3!p' | sed -n '1!p')"
	 	stt_username="apikey"
	 	stt_password="$(echo "$stt_response" | jq -r '.apikey')"
	 	nlu_response="$(ibmcloud service key-show "$NLU_INSTANCE" "$NLU_INSTANCE_CREDS" | sed -n '3!p' | sed -n '1!p')"
	 	nlu_username="apikey"
	 	nlu_password="$(echo "$nlu_response" | jq -r '.apikey')"
	 	db_response="$(ibmcloud service key-show "$DB_INSTANCE" "$DB_INSTANCE_CREDS" | sed -n '3!p' | sed -n '1!p' )"
	 	db_uri="$(echo "$db_response" | jq -r '.uri')"
	 	db_sdr_uri="$(echo "$db_uri" | sed -e "s/compose/${DB_NAME}/g")"
	 	ibmcloud fn action create "${FUNC_PACKAGE}/${FUNC_ACTION}" "$FUNC_ACTION_CODE" \
	 		--kind nodejs:8 \
	 		--memory 512 \
	 		--timeout 300000 \
	 		--param watsonSttUsername "$stt_username" \
	 		--param watsonSttPassword "$stt_password" \
	 		--param watsonNluUsername "$nlu_username" \
	 		--param watsonNluPassword "$nlu_password" \
	 		--param postgresUrl "$db_sdr_uri"
	 	ibmcloud -q fn action list
	fi
	echo `now` "Creating $FUNC_RULE rule for $FUNC_TRIGGER trigger and $FUNC_ACTION action..."
	echo `now` "Current rules:"
	ibmcloud -q fn rule list
	rule="/${org}_${space}/${FUNC_RULE}"
	if [[ $(ibmcloud -q fn rule list | cut -d' ' -f1 | grep -Fx "$rule") ]]; then \
	 	echo `now` "There is $rule rule created already, skipping its creation..."
	else
	 	echo `now` "Found no $rule rule, creating..."
	 	ibmcloud fn rule create "$FUNC_RULE" "$FUNC_TRIGGER" "${FUNC_PACKAGE}/${FUNC_ACTION}"
	 	ibmcloud -q fn rule list
	fi
	echo `now` "Creating and configuring functions entities is finished"
	echo `now` "Current functions entities:"
	ibmcloud -q fn list
}

# Delete functions
function teardown_func_(){
	echo `now` "Deleting functions"
	echo `now` "Waiting $WAIT_RESPONSE seconds..."
	sleep $WAIT_RESPONSE
	echo `now` "Current functions entities:"
	ibmcloud -q fn list
	org="$(ibmcloud target | grep Org | cut -d : -f 2 | sed -e 's/^[ \t]*//' | cut -d ' ' -f 1)"
	space="$(ibmcloud target | grep Space | cut -d : -f 2 | sed -e 's/^[ \t]*//' | cut -d ' ' -f 1)"
	rule="/${org}_${space}/${FUNC_RULE}"
	if [[ $(ibmcloud -q fn rule list | cut -d' ' -f1 | grep -Fx "$rule") ]]; then
	 	echo `now` "Found $rule rule for $FUNC_TRIGGER trigger and $FUNC_ACTION action, deleting..."
	 	ibmcloud -q fn rule delete --disable "$FUNC_RULE"
	else
	 	echo `now` "There is no $rule rule to delete, skipping..."
	fi
	trigger="/${org}_${space}/${FUNC_TRIGGER}"
	if [[ $(ibmcloud -q fn trigger list | cut -d' ' -f1 | grep -Fx "$trigger") ]]; then
	 	echo `now` "Found $trigger trigger for $FUNC_MH_FEED feed and $MH_SDR_TOPIC topic, deleting..."
	 	ibmcloud -q fn trigger delete "$FUNC_TRIGGER"
	else
	 	echo `now` "There is no $trigger trigger to delete, skipping..."
	fi
	action="/${org}_${space}/${FUNC_PACKAGE}/${FUNC_ACTION}"
	if [[ $(ibmcloud -q fn action list | cut -d' ' -f1 | grep -Fx "$action") ]]; then
	 	echo `now` "Found $action action, deleting..."
	 	ibmcloud -q fn action delete "${FUNC_PACKAGE}/${FUNC_ACTION}"
	else
	 	echo `now` "There is no $action action to delete, skipping..."
	fi
	package="/${org}_${space}/${FUNC_PACKAGE}"
	if [[ $(ibmcloud -q fn package list | cut -d' ' -f1 | grep -Fx "$package") ]]; then
	 	echo `now` "Found $package package, deleting..."
	 	ibmcloud -q fn package delete "$FUNC_PACKAGE"
	else
	 	echo `now` "There is no $package package to delete, skipping..."
	fi
	binding="/${org}_${space}/Bluemix_${MH_INSTANCE}_${MH_INSTANCE_CREDS}"
	if [[ $(ibmcloud -q fn package list | cut -d' ' -f1 | grep -Fx "$binding") ]]; then
	 	echo `now` "Found $binding Event Streams instance and functions binding, deleting..."
	 	ibmcloud -q fn package delete "Bluemix_${MH_INSTANCE}_${MH_INSTANCE_CREDS}"
	else
	 	echo `now` "There is no Event Streams instance and functions binding to delete, skipping..."
	fi
	echo `now` "Finished deleting functions entities"
	echo `now` "Current functions entities:"
	ibmcloud -q fn list
}

# Deploy UI application
function deploy_ui_(){
	echo `now` "Creating UI application"
	
	echo `now` "Creating $UI_APP_ID_INSTANCE App ID instance with the $UI_APP_ID_INSTANCE_PLAN plan for $UI_APP_NAME UI application"
	echo `now` "Current App ID instances:"
	ibmcloud resource service-instances --service-name appid
	## check if an app id instance is created, if not, create one
	if [[ $(ibmcloud resource service-instances --service-name appid --output JSON | \
			jq -r 'select (.!=null) | .[].name' | grep -Fx "$UI_APP_ID_INSTANCE") ]]; then
		echo `now` "There is $UI_APP_ID_INSTANCE App ID instance created already, skipping its creation..."
	else
		echo `now` "Found no App ID instance $UI_APP_ID_INSTANCE with $UI_APP_ID_INSTANCE_PLAN plan, creating..."
		ibmcloud resource service-instance-create "$UI_APP_ID_INSTANCE" appid "$UI_APP_ID_INSTANCE_PLAN" "$UI_APP_ID_INSTANCE_LOCATION"
		ibmcloud resource service-instances --service-name appid
	fi
	
	# alias for the App ID service
	echo `now` "Creating alias $UI_APP_ID_INSTANCE_ALIAS for $UI_APP_ID_INSTANCE..."
	echo `now` "Current aliases for $UI_APP_ID_INSTANCE:"
	ibmcloud resource service-aliases --instance-name "$UI_APP_ID_INSTANCE"
	if [[ $(ibmcloud resource service-aliases --instance-name "$UI_APP_ID_INSTANCE" --output JSON | \
		jq -r 'select (.!=null) | .[].name' | grep -Fx "$UI_APP_ID_INSTANCE_ALIAS") ]]; then
		echo `now` "There is $UI_APP_ID_INSTANCE_ALIAS App ID alias instance for $UI_APP_ID_INSTANCE created already, skipping its creation..."
	else
		echo `now` "Found no $UI_APP_ID_INSTANCE_ALIAS alias for $UI_APP_ID_INSTANCE App ID instance, creating..."
		ibmcloud resource service-alias-create "$UI_APP_ID_INSTANCE_ALIAS" --instance-name "$UI_APP_ID_INSTANCE"
		ibmcloud resource service-aliases --instance-name "$UI_APP_ID_INSTANCE"
	fi
	
	# creating credentials for managing and configuring our App ID instance
	echo `now` "Creating credentials $UI_APP_ID_INSTANCE_ALIAS_CREDS for $UI_APP_ID_INSTANCE_ALIAS alias"
	echo `now` "Current credentials for $UI_APP_ID_INSTANCE_ALIAS:"
	ibmcloud -q service keys "$UI_APP_ID_INSTANCE_ALIAS" | sed -n '3!p' | sed -n '1!p'
	if [[ $(ibmcloud service keys "$UI_APP_ID_INSTANCE_ALIAS" | cut -d' ' -f1 | grep -Fx "$UI_APP_ID_INSTANCE_ALIAS_CREDS") ]]; then
		echo `now` "There is $UI_APP_ID_INSTANCE_ALIAS_CREDS credentials for $UI_APP_ID_INSTANCE_ALIAS App ID alias instance created, skipping its creation..."
	else
		echo `now` "Found no $UI_APP_ID_INSTANCE_ALIAS_CREDS, creating..."
		ibmcloud service key-create "$UI_APP_ID_INSTANCE_ALIAS" "$UI_APP_ID_INSTANCE_ALIAS_CREDS"
		ibmcloud -q service keys "$UI_APP_ID_INSTANCE_ALIAS" | sed -n '3!p' | sed -n '1!p'
	fi

	echo `now` "Configuring $UI_APP_ID_INSTANCE App ID instance..."
	appid_response="$(ibmcloud service key-show "$UI_APP_ID_INSTANCE_ALIAS" "$UI_APP_ID_INSTANCE_ALIAS_CREDS" | sed -n '3!p' | sed -n '1!p')"
	apikey="$(echo $appid_response | jq -r '.apikey')"
	managementUrl="$(echo $appid_response| jq -r '.managementUrl')"
	token="$(curl -s -k -X POST \
		--header "Content-Type: application/x-www-form-urlencoded" \
		--header "Accept: application/json" \
		--data-urlencode "grant_type=urn:ibm:params:oauth:grant-type:apikey" \
		--data-urlencode "apikey=${apikey}" \
		"https://iam.bluemix.net/identity/token" | jq -r '.access_token')"

	echo `now` "Configuring identity providers on $UI_APP_ID_INSTANCE App ID instance..."
	
	echo `now` "Disabling Facebook identity provider..."
	curl -s -X PUT --header 'Content-Type: application/json' --header 'Accept: application/json' \
		--header "Authorization: Bearer $token" \
		-d '{"isActive": false}' ${managementUrl}/config/idps/facebook > /dev/null
	
	echo `now` "Disabling Google identity provider..."
	curl -s -X PUT --header 'Content-Type: application/json' --header 'Accept: application/json' \
		--header "Authorization: Bearer $token" \
		-d '{"isActive": false}' ${managementUrl}/config/idps/google > /dev/null
	
	echo `now` "Disabling anonymous identity provider..."
	curl -s -X PUT --header 'Content-Type: application/json' --header 'Accept: application/json' \
		--header "Authorization: Bearer $token" \
		-d '{"anonymousAccess":{"enabled":false}}' ${managementUrl}/config/tokens > /dev/null
	
	echo `now` "Configuring Cloud Directory Identity provider..."
	curl -s -X PUT --header 'Content-Type: application/json' --header 'Accept: application/json' \
		--header "Authorization: Bearer $token" \
		-d '{"config":{"selfServiceEnabled":true,"interactions":{"identityConfirmation":{"accessMode":"OFF","methods":["email"]},"welcomeEnabled":false,"resetPasswordEnabled":false,"resetPasswordNotificationEnable":false},"identityField":"email","signupEnabled":false},"isActive":true}' \
		${managementUrl}/config/idps/cloud_directory > /dev/null
	
	echo `now` "Registering application $UI_APP_NAME on $UI_APP_ID_INSTANCE App ID instance"
	echo `now` "Current registered applications on $UI_APP_ID_INSTANCE App ID instance:"
	curl -s -X GET --header 'Accept: application/json' --header "Authorization: Bearer $token" ${managementUrl}/applications | \
			jq -c '.applications[].name'

	if [[ $( curl -s -X GET --header 'Accept: application/json' --header "Authorization: Bearer $token" ${managementUrl}/applications | \
			jq -c '.applications[].name | select (. == "'"$UI_APP_NAME"'")') ]]; then
		echo `now` "Application $UI_APP_NAME is already registered on $UI_APP_ID_INSTANCE, skipping"
	else
		echo `now` "Found no $UI_APP_NAME, registering it on $UI_APP_ID_INSTANCE"
		curl -s -X POST --header 'Content-Type: application/json' --header 'Accept: application/json' \
			--header "Authorization: Bearer $token" \
			-d '{"name":"'"$UI_APP_NAME"'"}' ${managementUrl}/applications > /dev/null
	fi

	## get UI app registration params for a config file
	echo `now` "Updating $UI_APP_NAME configuration file with registration parameters from $UI_APP_ID_INSTANCE App ID instance"
	ui_app_id_registration="$(curl -s -X GET --header 'Accept: application/json' --header "Authorization: Bearer $token" ${managementUrl}/applications | \
		jq -c '.applications[] | select (.name == "'"$UI_APP_NAME"'")')"
	## updating config with application registration info from app id
	ui_app_clientId="$(echo $ui_app_id_registration | jq -r '.clientId')"
	ui_app_oauth="$(echo $ui_app_id_registration | jq -r '.oAuthServerUrl')"
	ui_app_profiles="https://appid-profiles.ng.bluemix.net"
	ui_app_secret="$(echo $ui_app_id_registration | jq -r '.secret')"
	ui_app_tenantId="$(echo $ui_app_id_registration | jq -r '.tenantId')"
	cp "${UI_SRC_PATH}/localdev-config.json.sample" "${UI_SRC_PATH}/localdev-config.json"
	sed -i.bak -e "s|\"clientId\": \"{clientId from app id}\",|\"clientId\": \""$ui_app_clientId"\",|" \
	 	-e "s|\"oauthServerUrl\": \"{oauthServerUrl from app id}\",|\"oauthServerUrl\": \""$ui_app_oauth"\",|" \
	 	-e "s|\"profilesUrl\": \"{profilesUrl from app id}\",|\"profilesUrl\": \""$ui_app_profiles"\",|" \
	 	-e "s|\"secret\": \"{secret of service from app id}\",|\"secret\": \""$ui_app_secret"\",|" \
	 	-e "s|\"tenantId\": \"{tenantId from app id}\"|\"tenantId\": \""$ui_app_tenantId"\"|" \
	 	"${UI_SRC_PATH}/localdev-config.json"
	
	## adding a user
	echo `now` "Adding $UI_APP_USER user to $UI_APP_ID_INSTANCE App ID instance"
	if [[ $(curl -s -X GET --header 'Accept: application/json' --header "Authorization: Bearer $token" ${managementUrl}/cloud_directory/Users | \
			jq -c '.Resources[].emails[].value | select(. == "'"$UI_APP_USER"'")') ]]; then
		echo `now` "User $UI_APP_USER is already created, skipping"
	else
		echo `now` "Found no user, creating..."
		curl -s -X POST --header 'Content-Type: application/json' --header 'Accept: application/json' \
		--header "Authorization: Bearer $token" \
		-d '{"emails": [{"value":"'"$UI_APP_USER"'","primary": true}],"userName": "SDRUIUser","password": "'"${UI_APP_PASSWORD}"'"}' \
		${managementUrl}/cloud_directory/Users > /dev/null
	fi

	echo `now` "Current applications:"
	ibmcloud -q app list
	
	echo `now` "Updating $UI_APP_NAME client configuration file with the Mapbox token"
	cp "${UI_SRC_PATH}/client/src/config/settings.template.js" "${UI_SRC_PATH}/client/src/config/settings.js"
	sed -i.bak "s|exports.MAPBOX_TOKEN = .*|exports.MAPBOX_TOKEN = '"$MAPBOX_TOKEN"'|" "${UI_SRC_PATH}/client/src/config/settings.js"
	
	echo `now` "Updating $UI_APP_NAME configuration file with the $DB_NAME DB information and random number"
	cp "${UI_SRC_PATH}/server/config/settings.template.js" "${UI_SRC_PATH}/server/config/settings.js"
	db_response="$(ibmcloud service key-show "$DB_INSTANCE" "$DB_INSTANCE_CREDS" | sed -n '3!p' | sed -n '1!p' )"
	db_uri="$(echo "$db_response" | jq -r '.uri')"
	db_sdr_uri="$(echo "$db_uri" | sed -e "s/compose/${DB_NAME}/g")"
	sed -i.bak -e "s|exports.postgresUrl = .*|exports.postgresUrl = '"$db_sdr_uri"'|" \
		-e "s|exports.appIDSecret = .*|exports.appIDSecret = '"${RANDOM}${RANDOM}"'|" \
		"${UI_SRC_PATH}/server/config/settings.js"
	echo `now` "$UI_APP_NAME configuration file with the $DB_NAME DB information updated"

	echo `now` "Building SDR UI app - $UI_APP_NAME..."
	(cd "${UI_SRC_PATH}/client" && npm install && npm run build)
	
	if [[ $(ibmcloud -q cf apps | cut -d' ' -f1 | sed '1,4d' | grep -Fx "$UI_APP_NAME") ]]; then
	  	echo `now` "There is $UI_APP_NAME UI application created already, syncing changes..."
	  	cd "$UI_SRC_PATH" && ibmcloud -q cf push "$UI_APP_NAME"
	else
	  	echo `now` "Found no $UI_APP_NAME UI application, creating..."
	  	cd "$UI_SRC_PATH" && ibmcloud -q cf push "$UI_APP_NAME" --no-start
	fi
	
	echo `now` "Current applications:"
	ibmcloud -q cf apps

	ui_url="$(ibmcloud -q cf apps | sed '1,4d' | grep "^${UI_APP_NAME}" | sed 's/^.* \(.*$\)/\1/' )"
	echo `now` "$UI_APP_NAME UI application URL is $ui_url"
	echo `now` "Configuring redirect URLs for $UI_APP_NAME on $UI_APP_ID_INSTANCE App ID instance..."
	curl -s -X PUT --header 'Content-Type: application/json' --header 'Accept: application/json' \
		--header "Authorization: Bearer $token" \
		-d '{"redirectUris": [ "'"http://${ui_url}/*"'", "'"https://${ui_url}/*"'"]}' \
		${managementUrl}/config/redirect_uris

	echo `now` "Binding $UI_APP_ID_INSTANCE_ALIAS alias and $UI_APP_NAME UI application"
	echo `now` "Current bindings for $UI_APP_ID_INSTANCE_ALIAS alias:"
	ibmcloud resource service-bindings "$UI_APP_ID_INSTANCE_ALIAS"
	if [[ $(ibmcloud resource service-bindings "$UI_APP_ID_INSTANCE_ALIAS" --output JSON | \
			jq -c '.[].TargetName | select(. == "'"$UI_APP_NAME"'")') ]]; then
		echo `now` "There is binding with $UI_APP_ID_INSTANCE_ALIAS alias for $UI_APP_NAME created already, skipping its creation..."
	else
		echo `now` "Found no binding $UI_APP_NAME with $UI_APP_ID_INSTANCE_ALIAS alias, creating..."
		ibmcloud resource service-binding-create "$UI_APP_ID_INSTANCE_ALIAS" "$UI_APP_NAME" Manager
		ibmcloud resource service-bindings "$UI_APP_ID_INSTANCE_ALIAS"
		echo `now` "Starting $UI_APP_NAME app..."
		ibmcloud -q cf start "$UI_APP_NAME"
	fi

	echo `now` "Creating and configuring $UI_APP_NAME UI application is finished"
	echo `now` "The UI is available at https://${ui_url}"
}

# Delete UI application
teardown_ui_(){
	echo `now` "Deleting $UI_APP_NAME UI application"
	echo `now` "Waiting $WAIT_RESPONSE seconds..."
	sleep $WAIT_RESPONSE
	
	echo `now` "Deleting $UI_APP_ID_INSTANCE App ID instance..."
	echo `now` "Current App ID instances:"
	ibmcloud resource service-instances --service-name appid
	## if there's an app id instance, we can check for alias, binding and its creds
	if [[ $(ibmcloud resource service-instances --service-name appid --output JSON | \
			jq -r 'select (.!=null) | .[].name' | grep -Fx "$UI_APP_ID_INSTANCE") ]]; then
		
		# if there's an alias we can check for binging and its creds
		if [[ $(ibmcloud resource service-aliases --instance-name "$UI_APP_ID_INSTANCE" --output JSON | \
			jq -r 'select (.!=null) | .[].name' | grep -Fx "$UI_APP_ID_INSTANCE_ALIAS") ]]; then
			echo `now` "Found $UI_APP_ID_INSTANCE_ALIAS, checking if there's bindings and its creds to delete..."
			echo `now` "Deleting binding with $UI_APP_ID_INSTANCE_ALIAS alias for $UI_APP_NAME application"
			echo `now` "Current bindings for $UI_APP_ID_INSTANCE_ALIAS alias:"
			ibmcloud resource service-bindings "$UI_APP_ID_INSTANCE_ALIAS"
			if [[ $(ibmcloud resource service-bindings "$UI_APP_ID_INSTANCE_ALIAS" --output JSON | \
					jq -c '.[].TargetName | select(. == "'"$UI_APP_NAME"'")') ]]; then
				echo `now` "Found binding with $UI_APP_ID_INSTANCE_ALIAS alias for $UI_APP_NAME, deleting..."
				ibmcloud resource service-binding-delete "$UI_APP_ID_INSTANCE_ALIAS" "$UI_APP_NAME" -f
			else
				echo `now` "Found no binding $UI_APP_NAME with $UI_APP_ID_INSTANCE_ALIAS alias, skipping..."
			fi

			echo `now` "Deleting credentials $UI_APP_ID_INSTANCE_ALIAS_CREDS for $UI_APP_ID_INSTANCE_ALIAS alias..."
			echo `now` "Current credentials for $UI_APP_ID_INSTANCE_ALIAS:"
			ibmcloud -q service keys "$UI_APP_ID_INSTANCE_ALIAS" | sed -n '3!p' | sed -n '1!p'
			if [[ $(ibmcloud service keys "$UI_APP_ID_INSTANCE_ALIAS" | cut -d' ' -f1 | grep -Fx "$UI_APP_ID_INSTANCE_ALIAS_CREDS") ]]; then
				echo `now` "Found $UI_APP_ID_INSTANCE_ALIAS_CREDS credentials for $UI_APP_ID_INSTANCE_ALIAS App ID alias instance, deleting..."
				ibmcloud service key-delete "$UI_APP_ID_INSTANCE_ALIAS" "$UI_APP_ID_INSTANCE_ALIAS_CREDS" -f
			else
				echo `now` "Found no $UI_APP_ID_INSTANCE_ALIAS_CREDS, skipping..."
			fi
			
			echo `now` "Deleting alias $UI_APP_ID_INSTANCE_ALIAS for $UI_APP_ID_INSTANCE..."
			echo `now` "Current aliases for $UI_APP_ID_INSTANCE:"
			ibmcloud resource service-aliases --instance-name "$UI_APP_ID_INSTANCE"
			echo `now` "Found $UI_APP_ID_INSTANCE_ALIAS App ID alias instance for $UI_APP_ID_INSTANCE, deleting..."
			ibmcloud resource service-alias-delete "$UI_APP_ID_INSTANCE_ALIAS" -f
		else
			echo `now` "Found no $UI_APP_ID_INSTANCE_ALIAS alias for $UI_APP_ID_INSTANCE App ID instance, skipping..."
		fi

		echo `now` "Found $UI_APP_ID_INSTANCE App ID instance, deleting..."
		ibmcloud resource service-instance-delete "$UI_APP_ID_INSTANCE" -f
		echo `now` "Current App ID instances:"
		ibmcloud resource service-instances --service-name appid
	else
		echo `now` "Found no App ID instance $UI_APP_ID_INSTANCE with $UI_APP_ID_INSTANCE_PLAN plan, skipping..."
	fi
	
	# deleting the UI app
	echo `now` "Current applications:"
	ibmcloud -q cf apps
	if [[ $(ibmcloud -q cf apps | cut -d' ' -f1 | sed '1,4d' | grep -Fx "$UI_APP_NAME") ]]; then
	  	echo `now` "Found $UI_APP_NAME UI application, deleting..."
	  	ibmcloud -q cf delete "$UI_APP_NAME" -f -r
		echo `now` "Current applications:"
		ibmcloud -q app list
	else
	  	echo `now` "There is no $UI_APP_NAME UI application, skipping..."
	fi
	echo `now` "Finished deleting $UI_APP_NAME UI application"
		
}

# Create Watson Speech-To-Text instance
function deploy_stt(){
	deploy_stt_
}

# Delete Watson Speech-To-text instance and its dependant services
function teardown_stt(){
	echo `now` "Deleting Watson STT dependant services: functions entities"
	teardown_func_

	teardown_stt_
}

# Create Watson Natural Language Understanding instance
function deploy_nlu(){
	deploy_nlu_
}

# Delete Watson Natural Language Understanding instance and its dependant services
function teardown_nlu(){
	echo `now` "Deleting Watson NLU dependant service: functions entities"
	teardown_func_

	teardown_nlu_
}

# Create and configure Event Streams Instance
function deploy_es(){
	deploy_es_
}

# Delete Event Streams instance and its dependant services
function teardown_es(){
	echo `now` "Deleting Event Streams dependant service: functions entities"
	teardown_func_

	teardown_es_
}

# Create and configure Compose for PostgreSQL instance
function deploy_db(){
	deploy_db_
}

# Delete Compose for PostgreSQL instance and its dependant services
function teardown_db(){
	echo `now` "Deleting Compose for PostgreSQL dependant service: UI application"
	teardown_ui_
	echo `now` "Deleting Compose for PostgreSQL dependant service: functions entities"
	teardown_func_

	teardown_db_
}

# Create and configure functions and services they depend on
function deploy_func(){
	echo `now` "Creating service function entities depend on: Watson STT instance"
	deploy_stt_
	echo `now` "Creating service function entities depend on: Watson NLU instance"
	deploy_nlu_
	echo `now` "Creating service function entities depend on: Event Streams instance"
	deploy_es_
	echo `now` "Creating service function entities depend on: Compose for PostgreSQL instance"
	deploy_db_

	deploy_func_
}

# Delete functions
function teardown_func(){
	teardown_func_
}

# Deploy UI application and services it depends on
function deploy_ui(){
	echo `now` "Creating service UI application depends on: Compose for PostgreSQL instance"
	deploy_db_

	deploy_ui_
}

# Delete UI application
function teardown_ui(){
	teardown_ui_
}

# Create all service instances required for SDR PoC
function deploy_all(){
	echo `now` "Creating all service instances required for running SDR PoC"
	deploy_stt_
	deploy_nlu_
	deploy_es_
	deploy_db_
	deploy_func_
	deploy_ui_
}

# Delete all SDR PoC related service instances
function teardown_all(){
	echo `now` "Deleting all SDR PoC related service instances"
	teardown_ui_
	teardown_func_
	teardown_db_
	teardown_es_
	teardown_nlu_
	teardown_stt_
}

# support one parameter specifying action and a service(s) as its value
if [ $# -eq 0 ]; then
    echo "No arguments was supplied, exiting..."
	help
	exit 1
elif [ $# -ne 1 ]; then
	echo "Too many arguments. Please provide a single argument"
	help
	exit 1
fi

for i in "$@"
do
case $i in
    -i=*|--install=*)
    COMPONENT="${i#*=}"
	ACTION="$DEPLOY"
    break
    ;;
    -u=*|--uninstall=*)
    COMPONENT="${i#*=}"
	ACTION="$TEARDOWN"
    break
    ;;
    -c|--config)
    show_config
    exit 0
    ;;
	-h|--help)
    help
    exit 0
    ;;
    -p|--prereqs)
    check_prereqs
    exit 0
    ;;
    * )
	echo "Unsupported option $i"
    help
    exit 1      # unknown option
    ;;
esac
done

if [[ ! " ${components[@]} " =~ " ${COMPONENT} " ]]; then
	echo "Supported components are: "
	for i in "${components[@]}"; do echo -n "${i} "; done
	echo ""
	echo `now` "The specified component ${COMPONENT} is not supported, exiting..."
	exit 1
fi

start_time=`date +%s`

# show current configuration
show_config
# checking prerequisites
check_prereqs

echo `now` "The current action is $ACTION for $COMPONENT component(s)"
if [ "$ACTION" == "$DEPLOY" ] ||  [ "$ACTION" == "$TEARDOWN" ] ; then
	${ACTION}_${COMPONENT}
else
	echo `now` "The action $ACTION is not supported, exiting..."
	exit 1
fi

end_time=`date +%s`
echo `now` "Execution time is" `expr $end_time - $start_time` "second(s)"