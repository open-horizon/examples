.PHONY:	help deploy-stt teardown-stt deploy-nlu teardown-nlu \
		deploy-db teardown-db prereqs deploy-all teardown-all \
		deploy-es teardown-es deploy-func teardown-func deploy-ui teardown-ui

# general
WAIT_RESPONSE:=10

# event streams/message hub
MH_INSTANCE=sdr-poc-es
MH_INSTANCE_CREDS:=$(MH_INSTANCE)-credentials
MH_SDR_TOPIC:=sdr-audio
MH_SDR_TOPIC_PARTIONS:=2
MH_RESPONSE_RETRY:=5

# watson speech-to-text service
STT_INSTANCE:=sdr-poc-speech-to-text
STT_INSTANCE_CREDS:=$(STT_INSTANCE)-credentials
STT_INSTANCE_PLAN:=lite

# watson natural language understanding
NLU_INSTANCE:=sdr-poc-natural-language-understanding
NLU_INSTANCE_CREDS:=$(NLU_INSTANCE)-credentials
NLU_INSTANCE_PLAN:=free

# db
DB_INSTANCE:=sdr-poc-compose-for-postgresql
DB_INSTANCE_CREDS:=$(DB_INSTANCE)-credentials
DB_INSTANCE_PLAN:=Standard
DB_INSTANCE_RETRY:=20
DB_NAME:=sdr
DB_SQL:='CREATE DATABASE $(DB_NAME); \l \dt \c $(DB_NAME) \\CREATE TABLE globalnouns( \
		noun TEXT PRIMARY KEY NOT NULL, \
		sentiment DOUBLE PRECISION NOT NULL, \
		numberofmentions BIGINT NOT NULL, \
		timeupdated timestamp with time zone); \
	CREATE TABLE nodenouns( \
		noun TEXT NOT NULL, \
		edgenode TEXT NOT NULL, \
		sentiment DOUBLE PRECISION NOT NULL, \
		numberofmentions BIGINT NOT NULL, \
		timeupdated timestamp with time zone, \
		PRIMARY KEY(noun, edgenode)); \
	CREATE TABLE stations( \
		edgenode TEXT NOT NULL, \
		frequency REAL NOT NULL, \
		numberofclips BIGINT NOT NULL, \
		dataqualitymetric REAL, \
		timeupdated timestamp with time zone, \
		PRIMARY KEY(edgenode, frequency)); \
	CREATE TABLE edgenodes( \
		edgenode TEXT PRIMARY KEY NOT NULL, \
		latitude REAL NOT NULL, \
		longitude REAL NOT NULL, \
		timeupdated timestamp with time zone);'

# functions
FUNC_PACKAGE:=sdr-poc-message-hub-evnts
FUNC_MH_FEED:=Bluemix_$(MH_INSTANCE)_$(MH_INSTANCE_CREDS)/messageHubFeed
FUNC_TRIGGER:=sdr-poc-message-received-trigger
FUNC_ACTION:=sdr-poc-process-message
FUNC_ACTION_CODE:=../../data-processing/ibm-functions/actions/msgreceive.js
FUNC_RULE:=sdr-poc-message-received-rule

# ui
UI_SRC_PATH:=../../ui/sdr-app
UI_APP_NAME:=sdr-poc-app

SHELL:=bash

PREREQUISITES = jq curl ibmcloud psql sed grep cut npm node

## help				: Display help
help : Makefile
	@sed -n 's/^##//p' $<

## prereqs			: Check for prerequisites
prereqs :
	$(foreach prereq, $(PREREQUISITES),\
		$(if $(shell which $(prereq)),$(info Found $(prereq) prerequisite),$(error No $(prereq) installed found)))
    ifndef MAPBOX_TOKEN
      $(error MAPBOX_TOKEN varialbe is not defined)
    endif
	@echo `date` "All prerequisites are met"

## deploy-es			: Create and configure Event Streams Instance
deploy-es : prereqs
	@echo `date` "Creating Event Streams instance $(MH_INSTANCE)"
	@echo `date` "Current Event Streams instances:"
	@ibmcloud -q service list | grep messagehub
	@if [[ $$(ibmcloud -q service list | grep messagehub | cut -d' ' -f1 | grep -Fx "$(MH_INSTANCE)") ]]; then \
		echo `date` "There is $(MH_INSTANCE) Event Streams instance created already, skipping its creation..."; \
	else \
		echo `date` "Found no Event Streams instance $(MH_INSTANCE), creating..."; \
		ibmcloud -q service create messagehub standard "$(MH_INSTANCE)"; \
		ibmcloud -q service list | grep messagehub; \
	fi
	@echo `date` "Creating credentials $(MH_INSTANCE_CREDS) for $(MH_INSTANCE)"
	@echo `date` "Current credentials for $(MH_INSTANCE):"
	@ibmcloud -q service keys "$(MH_INSTANCE)" | sed -n '3!p' | sed -n '1!p'
	@if [[ $$(ibmcloud service keys "$(MH_INSTANCE)" | cut -d' ' -f1 | grep -Fx "$(MH_INSTANCE_CREDS)") ]]; then \
		echo `date` "There is $(MH_INSTANCE_CREDS) credentials for Event Streams $(MH_INSTANCE) created, skipping its creation..."; \
	else \
		echo `date` "Found no $(MH_INSTANCE_CREDS), creating..."; \
		ibmcloud service key-create "$(MH_INSTANCE)" "$(MH_INSTANCE_CREDS)"; \
	fi
	@while [ true ] ; do \
			response="$$(ibmcloud service key-show "$(MH_INSTANCE)" "$(MH_INSTANCE_CREDS)" | sed -n '3!p' | sed -n '1!p')"; \
			echo `date` "Checking if $(MH_INSTANCE) and $(MH_INSTANCE_CREDS) are ready" ; \
			if [[ "$$response" = *"not found"* ]] ; then \
				echo `date` "$(MH_INSTANCE) and $(MH_INSTANCE_CREDS) seem not ready" ; \
				echo `date` "Retrying in $(MH_RESPONSE_RETRY) seconds" ; \
				sleep $(MH_RESPONSE_RETRY) ; \
				continue; \
			else \
				echo `date` $(MH_INSTANCE) and $(MH_INSTANCE_CREDS) are ready ; \
				break; \
			fi \
	done
	@echo `date` Creating $(MH_SDR_TOPIC) with $(MH_SDR_TOPIC_PARTIONS) partitions on $(MH_INSTANCE)
	@response="$$(ibmcloud service key-show "$(MH_INSTANCE)" "$(MH_INSTANCE_CREDS)" | sed -n '3!p' | sed -n '1!p')" ; \
		admin_url="$$(echo "$$response" | jq -r '.kafka_admin_url')" ; \
		api_key="$$(echo "$$response" | jq -r '.api_key')" ; \
		if [[ $$(curl -s -H 'Accept: application/json' \
					-H 'X-Auth-Token: '"$$api_key" \
   					"$$admin_url"/admin/topics/ | jq -c '.[] | select(.name | . and contains("$(MH_SDR_TOPIC)"))') ]]; then \
			echo `date` "$(MH_SDR_TOPIC) topic is already created on $(MH_INSTANCE)"; \
		else \
			echo `date` "Found no $(MH_SDR_TOPIC), creating..."; \
				curl -s -H 'Content-Type: application/json' -H 'Accept: */*' \
					-H 'X-Auth-Token: '"$$api_key" \
					-d '{ "name": "$(MH_SDR_TOPIC)", "partitions": $(MH_SDR_TOPIC_PARTIONS) }' \
					"$$admin_url"/admin/topics ; \
		fi 
	@echo `date` "Finished creation $(MH_INSTANCE) instance with $(MH_SDR_TOPIC) topic, $(MH_SDR_TOPIC_PARTIONS) partitions"
	@echo `date` "Current Event Streams instances:"
	@ibmcloud -q service list | grep messagehub

## teardown-es			: Delete Event Streams instance
teardown-es : prereqs teardown-func
	@echo `date` "Deleting $(MH_INSTANCE_CREDS) for $(MH_INSTANCE)"
	@echo `date` "Waiting $(WAIT_RESPONSE) seconds..."; \
		sleep $(WAIT_RESPONSE)
	@echo `date` "Current Event Streams instances:"
	@ibmcloud -q service list | grep messagehub
	@if [[ $$(ibmcloud service keys "$(MH_INSTANCE)" | cut -d' ' -f1 | grep -Fx "$(MH_INSTANCE_CREDS)") ]]; then \
		echo `date` "Found $(MH_INSTANCE_CREDS) credentials for $(MH_INSTANCE), deleting..."; \
		ibmcloud service key-delete "$(MH_INSTANCE)" "$(MH_INSTANCE_CREDS)" -f; \
	else \
		echo `date` "There is no $(MH_INSTANCE_CREDS) credentials for $(MH_INSTANCE) to delete, skipping..."; \
	fi
	@echo `date` "Deleting $(MH_INSTANCE)..."
	@if [[ $$(ibmcloud -q service list | cut -d' ' -f1 | grep -Fx "$(MH_INSTANCE)") ]]; then \
		echo `date` "Found $(MH_INSTANCE), deleting..."; \
		ibmcloud service delete "$(MH_INSTANCE)" -f; \
	else \
		echo `date` "Found no $(MH_INSTANCE) Event Streams instance, skipping..."; \
	fi
	@echo `date` "Finished deleting $(MH_INSTANCE) and its credentials $(MH_INSTANCE_CREDS)"
	@echo `date` "Current Event Streams instances:"
	@ibmcloud -q service list | grep messagehub

## deploy-db			: Create and configure Compose for PostgreSQL instance
deploy-db : prereqs
	@echo `date` "Creating $(DB_INSTANCE) Compose PostgreSQL DB with the $(DB_INSTANCE_PLAN) plan"
	@echo `date` "Current Compose for PostgreSQL instances:"
	@ibmcloud -q service list | grep compose-for-postgresql
	@if [[ $$(ibmcloud -q service list | grep compose-for-postgresql | cut -d' ' -f1 | grep -Fx "$(DB_INSTANCE)") ]]; then \
		echo `date` "There is $(DB_INSTANCE) Compose for PostgreSQL instance created already, skipping its creation..."; \
	else \
		echo `date` "Found no Compose for PostgreSQL instance $(DB_INSTANCE) with $(DB_INSTANCE_PLAN) plan, creating..."; \
		ibmcloud -q service create compose-for-postgresql "$(DB_INSTANCE_PLAN)" "$(DB_INSTANCE)"; \
		ibmcloud -q service list | grep compose-for-postgresql; \
	fi
	@echo `date` "Creating credentials $(DB_INSTANCE_CREDS) for $(DB_INSTANCE)"
	@echo `date` "Current credentials for $(DB_INSTANCE):"
	@ibmcloud -q service keys "$(DB_INSTANCE)" | sed -n '3!p' | sed -n '1!p'
	@if [[ $$(ibmcloud service keys "$(DB_INSTANCE)" | cut -d' ' -f1 | grep -Fx "$(DB_INSTANCE_CREDS)") ]]; then \
		echo `date` "There is $(DB_INSTANCE_CREDS) credentials for $(DB_INSTANCE) PostgreSQL instance created, skipping its creation..."; \
	else \
		echo `date` "Found no $(DB_INSTANCE_CREDS), creating..."; \
		ibmcloud service key-create "$(DB_INSTANCE)" "$(DB_INSTANCE_CREDS)"; \
	fi
	@echo `date` "Bootstrapping database $(DB_NAME) on $(DB_INSTANCE)"
	@echo `date` "Checking if $(DB_INSTANCE) PostgreSQL instance is ready..."
	@db_uri="$$(ibmcloud service key-show $(DB_INSTANCE) $(DB_INSTANCE_CREDS) | sed -n '3!p' | sed -n '1!p' | jq -r '.uri')"; \
	while [ true ] ; do \
            response="$$(psql "$$db_uri?sslmode=require" -c '\l' 2>&1 | sed -n 1p )"; \
			if [[ "$$response" = *"Operation timed out"* ]] || [[ "$$response" = *"Connection refused"* ]] ; then \
				echo `date` $(DB_INSTANCE) is not ready, response is "$$response" ; \
				echo `date` "Retrying in $(DB_INSTANCE_RETRY) seconds" ; \
				sleep $(DB_INSTANCE_RETRY) ; \
				continue; \
			else \
				echo `date` "$(DB_INSTANCE) is ready" ; \
				break ; \
			fi; \
    done; \
	if [[ $$(psql "$$db_uri?sslmode=require" -c '\l' | sed '1,3d;$d' | cut -d'|' -f1 | cut -d' ' -f2 | grep -Fx $(DB_NAME)) ]]; then \
		echo `date` "There is $(DB_NAME) database already created on $(DB_INSTANCE) instance, skipping" ; \
	else \
		echo `date` "Found no $(DB_NAME) on $(DB_INSTANCE) instance, creating" ; \
		echo $(DB_SQL) | psql "$$db_uri?sslmode=require"; \
		echo `date` "Created $(DB_NAME) on $(DB_INSTANCE) instance"; \
		echo `date` "Created following tables in $(DB_NAME)"; \
		db_sdr="$$(echo "$$db_uri" | sed -e "s/compose/$(DB_NAME)/g")"; \
		psql "$$db_sdr?sslmode=require" -c '\dt'; \
	fi
	@echo `date` "Creating and configuration $(DB_INSTANCE) Compose for PostgreSQL is finished"
	@echo `date` "Current Compose for PostgreSQL instances:"
	@ibmcloud -q service list | grep compose-for-postgresql

## teardown-db			: Delete Compose for PostgreSQL instance
teardown-db : prereqs teardown-ui teardown-func
	@echo `date` "Deleting $(DB_INSTANCE_CREDS) credentials for $(DB_INSTANCE)"
	@echo `date` "Waiting $(WAIT_RESPONSE) seconds..."; \
		sleep $(WAIT_RESPONSE)
	@echo `date` "Current Compose for PostgreSQL instances:"
	@ibmcloud -q service list | grep compose-for-postgresql
	@if [[ $$(ibmcloud service keys "$(DB_INSTANCE)" | cut -d' ' -f1 | grep -Fx "$(DB_INSTANCE_CREDS)") ]]; then \
		echo `date` "Found $(DB_INSTANCE_CREDS) credentials for $(DB_INSTANCE), deleting..."; \
		ibmcloud service key-delete "$(DB_INSTANCE)" "$(DB_INSTANCE_CREDS)" -f; \
	else \
		echo `date` "There is no $(DB_INSTANCE_CREDS) credentials for $(DB_INSTANCE) to delete, skipping..."; \
	fi
	@echo `date` "Deleting $(DB_INSTANCE)..."
	@if [[ $$(ibmcloud -q service list | grep compose-for-postgresql | cut -d' ' -f1 | grep -Fx "$(DB_INSTANCE)") ]]; then \
		echo `date` "Found $(DB_INSTANCE), deleting..."; \
		ibmcloud service delete "$(DB_INSTANCE)" -f; \
	else \
		echo `date` "Found no $(DB_INSTANCE) Compose for PostgreSQL instance, skipping..."; \
	fi
	@echo `date` "Finished deleting $(DB_INSTANCE) and its credentials $(DB_INSTANCE_CREDS)"
	@echo `date` "Current Compose for PostgreSQL instances:"
	@ibmcloud -q service list | grep compose-for-postgresql

## deploy-stt			: Create Watson Speech-To-Text instance
deploy-stt : prereqs
	@echo `date` "Creating $(STT_INSTANCE) Watson Speech-To-Text service instance with the $(STT_INSTANCE_PLAN) plan"
	@echo `date` "Current Watson Speech-To-Text instances:"
	@ibmcloud -q service list | grep speech_to_text
	@if [[ $$(ibmcloud -q service list | grep speech_to_text | cut -d' ' -f1 | grep -Fx "$(STT_INSTANCE)") ]]; then \
		echo `date` "There is $(STT_INSTANCE) Watson STT instance created already, skipping its creation..."; \
	else \
		echo `date` "Found no Watson STT instance $(STT_INSTANCE) with $(STT_INSTANCE_PLAN) plan, creating..."; \
		ibmcloud -q service create speech_to_text "$(STT_INSTANCE_PLAN)" "$(STT_INSTANCE)"; \
		ibmcloud -q service list | grep speech_to_text; \
	fi
	@echo `date` "Creating credentials $(STT_INSTANCE_CREDS) for $(STT_INSTANCE)"
	@echo `date` "Current credentials for $(STT_INSTANCE):"
	@ibmcloud -q service keys "$(STT_INSTANCE)" | sed -n '3!p' | sed -n '1!p'
	@if [[ $$(ibmcloud service keys "$(STT_INSTANCE)" | cut -d' ' -f1 | grep -Fx "$(STT_INSTANCE_CREDS)") ]]; then \
		echo `date` "There is $(STT_INSTANCE_CREDS) credentials for $(STT_INSTANCE) Watson STT instance created, skipping its creation..."; \
	else \
		echo `date` "Found no $(STT_INSTANCE_CREDS), creating..."; \
		ibmcloud service key-create "$(STT_INSTANCE)" "$(STT_INSTANCE_CREDS)"; \
	fi
	@echo `date` "Creating and configuration $(STT_INSTANCE) Watson STT is finished"
	@echo `date` "Current Watson STT instances:"
	@ibmcloud -q service list | grep speech_to_text

## teardown-stt			: Delete Watson Speech-To-text instance
teardown-stt : prereqs teardown-func
	@echo `date` "Deleting $(STT_INSTANCE_CREDS) credentials for $(STT_INSTANCE)"
	@echo `date` "Waiting $(WAIT_RESPONSE) seconds..."; \
		sleep $(WAIT_RESPONSE)
	@echo `date` "Current Watson STT instances:"
	@ibmcloud -q service list | grep speech_to_text
	@if [[ $$(ibmcloud service keys "$(STT_INSTANCE)" | cut -d' ' -f1 | grep -Fx "$(STT_INSTANCE_CREDS)") ]]; then \
		echo `date` "Found $(STT_INSTANCE_CREDS) credentials for $(STT_INSTANCE), deleting..."; \
		ibmcloud service key-delete "$(STT_INSTANCE)" "$(STT_INSTANCE_CREDS)" -f; \
	else \
		echo `date` "There is no $(STT_INSTANCE_CREDS) credentials for $(STT_INSTANCE) to delete, skipping..."; \
	fi
	@echo `date` "Deleting $(STT_INSTANCE)..."
	@if [[ $$(ibmcloud -q service list | grep speech_to_text | cut -d' ' -f1 | grep -Fx "$(STT_INSTANCE)") ]]; then \
		echo `date` "Found $(STT_INSTANCE), deleting..."; \
		ibmcloud service delete "$(STT_INSTANCE)" -f; \
	else \
		echo `date` "Found no $(STT_INSTANCE) Watson STT instance, skipping..."; \
	fi
	@echo `date` "Finished deleting $(STT_INSTANCE) and its credentials $(STT_INSTANCE_CREDS)"
	@echo `date` "Current Watson STT instances:"
	@ibmcloud -q service list | grep speech_to_text

## deploy-nlu			: Create Watson Natural Language Understanding instance
deploy-nlu : prereqs
	@echo `date` "Creating $(NLU_INSTANCE) Watson Natural Language Understanding service instance with the $(NLU_INSTANCE_PLAN) plan"
	@echo `date` "Current Watson Natural Language Understanding instances:"
	@ibmcloud -q service list | grep natural-language-understanding
	@if [[ $$(ibmcloud -q service list | grep natural-language-understanding | cut -d' ' -f1 | grep -Fx "$(NLU_INSTANCE)") ]]; then \
		echo `date` "There is $(NLU_INSTANCE) Watson NLU instance created already, skipping its creation..."; \
	else \
		echo `date` "Found no Watson NLU instance $(NLU_INSTANCE) with $(NLU_INSTANCE_PLAN) plan, creating..."; \
		ibmcloud -q service create natural-language-understanding "$(NLU_INSTANCE_PLAN)" "$(NLU_INSTANCE)"; \
		ibmcloud -q service list | grep natural-language-understanding; \
	fi
	@echo `date` "Creating credentials $(NLU_INSTANCE_CREDS) for $(NLU_INSTANCE)"
	@echo `date` "Current credentials for $(NLU_INSTANCE):"
	@ibmcloud -q service keys "$(NLU_INSTANCE)" | sed -n '3!p' | sed -n '1!p'
	@if [[ $$(ibmcloud service keys "$(NLU_INSTANCE)" | cut -d' ' -f1 | grep -Fx "$(NLU_INSTANCE_CREDS)") ]]; then \
		echo `date` "There is $(NLU_INSTANCE_CREDS) credentials for $(NLU_INSTANCE) Watson NLU instance created, skipping its creation..."; \
	else \
		echo `date` "Found no $(NLU_INSTANCE_CREDS), creating..."; \
		ibmcloud service key-create "$(NLU_INSTANCE)" "$(NLU_INSTANCE_CREDS)"; \
	fi
	@echo `date` "Creating and configuration $(NLU_INSTANCE) Watson NLU is finished"
	@echo `date` "Current Watson NLU instances:"
	@ibmcloud -q service list | grep natural-language-understanding

## teardown-nlu			: Delete Watson Natural Language Understanding instance
teardown-nlu : prereqs teardown-func
	@echo `date` "Deleting $(NLU_INSTANCE_CREDS) credentials for $(NLU_INSTANCE)"
	@echo `date` "Waiting $(WAIT_RESPONSE) seconds..."; \
		sleep $(WAIT_RESPONSE)
	@echo `date` "Current Watson NLU instances:"
	@ibmcloud -q service list | grep natural-language-understanding
	@if [[ $$(ibmcloud service keys "$(NLU_INSTANCE)" | cut -d' ' -f1 | grep -Fx "$(NLU_INSTANCE_CREDS)") ]]; then \
		echo `date` "Found $(NLU_INSTANCE_CREDS) credentials for $(NLU_INSTANCE), deleting..."; \
		ibmcloud service key-delete "$(NLU_INSTANCE)" "$(NLU_INSTANCE_CREDS)" -f; \
	else \
		echo `date` "There is no $(NLU_INSTANCE_CREDS) credentials for $(NLU_INSTANCE) to delete, skipping..."; \
	fi
	@echo `date` "Deleting $(NLU_INSTANCE)..."
	@if [[ $$(ibmcloud -q service list | grep natural-language-understanding | cut -d' ' -f1 | grep -Fx "$(NLU_INSTANCE)") ]]; then \
		echo `date` "Found $(NLU_INSTANCE), deleting..."; \
		ibmcloud service delete "$(NLU_INSTANCE)" -f; \
	else \
		echo `date` "Found no $(NLU_INSTANCE) Watson NLU instance, skipping..."; \
	fi
	@echo `date` "Finished deleting $(NLU_INSTANCE) and its credentials $(NLU_INSTANCE_CREDS)"
	@echo `date` "Current Watson NLU instances:"
	@ibmcloud -q service list | grep natural-language-understanding

## deploy-func			: Create and configure functions
deploy-func : prereqs deploy-stt deploy-nlu deploy-es deploy-db
	@echo `date` "Creating functions entities"
	@echo `date` "Current functions entities:"
	@ibmcloud -q fn list
	@echo `date` "Creating bindings for $(MH_INSTANCE_CREDS) credentials from $(MH_INSTANCE) and functions"
	@ibmcloud fn package refresh
	@echo `date` "Creating trigger $(FUNC_TRIGGER) for $(FUNC_MH_FEED) feed"
	@echo `date` "Current triggers:"
	@ibmcloud -q fn trigger list
	@org="$$(ibmcloud target | grep Org | cut -d : -f 2 | sed -e 's/^[ \t]*//' | cut -d ' ' -f 1)"; \
	space="$$(ibmcloud target | grep Space | cut -d : -f 2 | sed -e 's/^[ \t]*//' | cut -d ' ' -f 1)"; \
	trigger="/""$$org"_"$$space"/"$(FUNC_TRIGGER)"; \
	if [[ $$(ibmcloud -q fn trigger list | cut -d' ' -f1 | grep -Fx "$$trigger") ]]; then \
		echo `date` "There is $$trigger trigger created already, skipping its creation..."; \
	else \
		echo `date` "Found no $$trigger trigger for $(FUNC_MH_FEED) feed and $(MH_SDR_TOPIC) topic, creating..."; \
		ibmcloud -q fn trigger create "$(FUNC_TRIGGER)" \
			--feed "$(FUNC_MH_FEED)" \
			--param isJSONData true \
			--param isBinaryValue false \
			--param topic "$(MH_SDR_TOPIC)" ; \
		ibmcloud -q fn trigger list; \
	fi
	@echo `date` "Creating $(FUNC_PACKAGE) package for $(FUNC_ACTION) action..."
	@echo `date` "Current packages:"
	@ibmcloud -q fn package list
	@org="$$(ibmcloud target | grep Org | cut -d : -f 2 | sed -e 's/^[ \t]*//' | cut -d ' ' -f 1)"; \
	space="$$(ibmcloud target | grep Space | cut -d : -f 2 | sed -e 's/^[ \t]*//' | cut -d ' ' -f 1)"; \
	package="/""$$org"_"$$space"/"$(FUNC_PACKAGE)"; \
	if [[ $$(ibmcloud -q fn package list | cut -d' ' -f1 | grep -Fx "$$package") ]]; then \
		echo `date` "There is $$package package created already, skipping its creation..."; \
	else \
		echo `date` "Found no $$package package, creating..."; \
		ibmcloud fn package create "$(FUNC_PACKAGE)";\
		ibmcloud -q fn package list; \
	fi
	@echo `date` "Creating $(FUNC_ACTION) action in $(FUNC_PACKAGE) package..."
	@echo `date` "Current actions:"
	@ibmcloud -q fn action list
	@org="$$(ibmcloud target | grep Org | cut -d : -f 2 | sed -e 's/^[ \t]*//' | cut -d ' ' -f 1)"; \
	space="$$(ibmcloud target | grep Space | cut -d : -f 2 | sed -e 's/^[ \t]*//' | cut -d ' ' -f 1)"; \
	action="/""$$org"_"$$space"/"$(FUNC_PACKAGE)"/"$(FUNC_ACTION)"; \
	if [[ $$(ibmcloud -q fn action list | cut -d' ' -f1 | grep -Fx "$$action") ]]; then \
		echo `date` "There is $$action action created already, skipping its creation..."; \
	else \
		echo `date` "Found no $$action action, creating..."; \
		stt_response="$$(ibmcloud service key-show $(STT_INSTANCE) $(STT_INSTANCE_CREDS) | sed -n '3!p' | sed -n '1!p')"; \
		stt_username="$$(echo "$$stt_response" | jq -r '.username')"; \
		stt_password="$$(echo "$$stt_response" | jq -r '.password')"; \
		nlu_response="$$(ibmcloud service key-show $(NLU_INSTANCE) $(NLU_INSTANCE_CREDS) | sed -n '3!p' | sed -n '1!p')"; \
		nlu_username="$$(echo "$$nlu_response" | jq -r '.username')"; \
		nlu_password="$$(echo "$$nlu_response" | jq -r '.password')"; \
		db_response="$$(ibmcloud service key-show $(DB_INSTANCE) $(DB_INSTANCE_CREDS) | sed -n '3!p' | sed -n '1!p' )"; \
		db_uri="$$(echo "$$db_response" | jq -r '.uri')"; \
		db_sdr_uri="$$(echo "$$db_uri" | sed -e "s/compose/$(DB_NAME)/g")"; \
		ibmcloud fn action create $(FUNC_PACKAGE)/$(FUNC_ACTION) $(FUNC_ACTION_CODE) \
			--kind nodejs:8 \
			--memory 512 \
			--timeout 300000 \
			--param watsonSttUsername "$$stt_username" \
			--param watsonSttPassword "$$stt_password" \
			--param watsonNluUsername "$$nlu_username" \
			--param watsonNluPassword "$$nlu_password" \
			--param postgresUrl "$$db_sdr_uri"; \
		ibmcloud -q fn action list; \
	fi
	@echo `date` "Creating $(FUNC_RULE) rule for $(FUNC_TRIGGER) and $(FUNC_ACTION)..."
	@echo `date` "Current rules:"
	@ibmcloud -q fn rule list
	@org="$$(ibmcloud target | grep Org | cut -d : -f 2 | sed -e 's/^[ \t]*//' | cut -d ' ' -f 1)"; \
	space="$$(ibmcloud target | grep Space | cut -d : -f 2 | sed -e 's/^[ \t]*//' | cut -d ' ' -f 1)"; \
	rule="/""$$org"_"$$space"/"$(FUNC_RULE)"; \
	if [[ $$(ibmcloud -q fn rule list | cut -d' ' -f1 | grep -Fx "$$rule") ]]; then \
		echo `date` "There is $$rule rule created already, skipping its creation..."; \
	else \
		echo `date` "Found no $$rule rule, creating..."; \
		ibmcloud fn rule create "$(FUNC_RULE)" "$(FUNC_TRIGGER)" "$(FUNC_PACKAGE)"/"$(FUNC_ACTION)"; \
		ibmcloud -q fn rule list; \
	fi
	@echo `date` "Creating and configuration functions entities is finished"
	@echo `date` "Current functions entities:"
	@ibmcloud -q fn list

## teardown-func			: Delete functions
teardown-func : prereqs
	@echo `date` "Deleting functions"
	@echo `date` "Waiting $(WAIT_RESPONSE) seconds..."; \
		sleep $(WAIT_RESPONSE)
	@echo `date` "Current functions entities:"
	@ibmcloud -q fn list
	@org="$$(ibmcloud target | grep Org | cut -d : -f 2 | sed -e 's/^[ \t]*//' | cut -d ' ' -f 1)"; \
	space="$$(ibmcloud target | grep Space | cut -d : -f 2 | sed -e 's/^[ \t]*//' | cut -d ' ' -f 1)"; \
	rule="/""$$org"_"$$space"/"$(FUNC_RULE)"; \
	if [[ $$(ibmcloud -q fn rule list | cut -d' ' -f1 | grep -Fx "$$rule") ]]; then \
		echo `date` "Found $$rule rule for $(FUNC_TRIGGER) and $(FUNC_ACTION), deleting..."; \
		ibmcloud -q fn rule delete --disable "$(FUNC_RULE)"; \
	else \
		echo `date` "There is no $$rule rule to delete, skipping..."; \
	fi
	@org="$$(ibmcloud target | grep Org | cut -d : -f 2 | sed -e 's/^[ \t]*//' | cut -d ' ' -f 1)"; \
	space="$$(ibmcloud target | grep Space | cut -d : -f 2 | sed -e 's/^[ \t]*//' | cut -d ' ' -f 1)"; \
	trigger="/""$$org"_"$$space"/"$(FUNC_TRIGGER)"; \
	if [[ $$(ibmcloud -q fn trigger list | cut -d' ' -f1 | grep -Fx "$$trigger") ]]; then \
		echo `date` "Found $$trigger trigger for $(FUNC_MH_FEED) feed and $(MH_SDR_TOPIC) topic, deleting..."; \
		ibmcloud -q fn trigger delete "$(FUNC_TRIGGER)"; \
	else \
		echo `date` "There is no $$trigger trigger to delete, skipping..."; \
	fi
	@org="$$(ibmcloud target | grep Org | cut -d : -f 2 | sed -e 's/^[ \t]*//' | cut -d ' ' -f 1)"; \
	space="$$(ibmcloud target | grep Space | cut -d : -f 2 | sed -e 's/^[ \t]*//' | cut -d ' ' -f 1)"; \
	action="/""$$org"_"$$space"/"$(FUNC_PACKAGE)"/"$(FUNC_ACTION)"; \
	if [[ $$(ibmcloud -q fn action list | cut -d' ' -f1 | grep -Fx "$$action") ]]; then \
		echo `date` "Found $$action action, deleting..."; \
		ibmcloud -q fn action delete "$(FUNC_PACKAGE)"/"$(FUNC_ACTION)"; \
	else \
		echo `date` "There is no $$action action to delete, skipping..."; \
	fi
	@org="$$(ibmcloud target | grep Org | cut -d : -f 2 | sed -e 's/^[ \t]*//' | cut -d ' ' -f 1)"; \
	space="$$(ibmcloud target | grep Space | cut -d : -f 2 | sed -e 's/^[ \t]*//' | cut -d ' ' -f 1)"; \
	package="/""$$org"_"$$space"/"$(FUNC_PACKAGE)"; \
	if [[ $$(ibmcloud -q fn package list | cut -d' ' -f1 | grep -Fx "$$package") ]]; then \
		echo `date` "Found $$package package, deleting..."; \
		ibmcloud -q fn package delete "$(FUNC_PACKAGE)"; \
	else \
		echo `date` "There is no $$package package to delete, skipping..."; \
	fi
	@org="$$(ibmcloud target | grep Org | cut -d : -f 2 | sed -e 's/^[ \t]*//' | cut -d ' ' -f 1)"; \
	space="$$(ibmcloud target | grep Space | cut -d : -f 2 | sed -e 's/^[ \t]*//' | cut -d ' ' -f 1)"; \
	package="/""$$org"_"$$space"/Bluemix_"$(MH_INSTANCE)"_"$(MH_INSTANCE_CREDS)"; \
	if [[ $$(ibmcloud -q fn package list | cut -d' ' -f1 | grep -Fx "$$package") ]]; then \
		echo `date` "Found Event Streams instance and function binding, deleting..."; \
		ibmcloud -q fn package delete Bluemix_$(MH_INSTANCE)_$(MH_INSTANCE_CREDS); \
	else \
		echo `date` "There is no Event Streams instance and functions binding to delete, skipping..."; \
	fi
	@echo `date` "Finished deleting functions entities"
	@echo `date` "Current functions entities:"
	@ibmcloud -q fn list

## deploy-ui			: Deploy UI application
deploy-ui : prereqs deploy-db
	@echo `date` "Creating UI application"
	@echo `date` "Current applications:"
	@ibmcloud -q app list
	@echo `date` "Updating $(UI_APP_NAME) client configuration file with the Mapbox token"
	@cp $(UI_SRC_PATH)/client/src/config/settings.template.js $(UI_SRC_PATH)/client/src/config/settings.js
	@sed -ie "s|exports.MAPBOX_TOKEN = .*|exports.MAPBOX_TOKEN = '"$$MAPBOX_TOKEN"'|" $(UI_SRC_PATH)/client/src/config/settings.js
	@echo `date` "Building SDR UI app - $(UI_APP_NAME)..."
	@(cd $(UI_SRC_PATH)/client && npm install && npm run build)
	@echo `date` "Updating $(UI_APP_NAME) configuration file with the $(DB_NAME) DB information"
	@cp $(UI_SRC_PATH)/server/config/settings.template.js $(UI_SRC_PATH)/server/config/settings.js
	@db_response="$$(ibmcloud service key-show $(DB_INSTANCE) $(DB_INSTANCE_CREDS) | sed -n '3!p' | sed -n '1!p' )"; \
	db_uri="$$(echo "$$db_response" | jq -r '.uri')"; \
	db_sdr_uri="$$(echo "$$db_uri" | sed -e "s/compose/$(DB_NAME)/g")"; \
	sed -ie "s|exports.postgresUrl = .*|exports.postgresUrl = '"$$db_sdr_uri"';|" $(UI_SRC_PATH)/server/config/settings.js
	@echo `date` "$(UI_APP_NAME) configuration file with the $(DB_NAME) DB information updated"
	@if [[ $$(ibmcloud -q app list | cut -d' ' -f1 | sed '1,4d' | grep -Fx "$(UI_APP_NAME)") ]]; then \
		echo `date` "There is $(UI_APP_NAME) UI application created already, syncing changes..."; \
		cd $(UI_SRC_PATH) && ibmcloud -q app push $(UI_APP_NAME); \
	else \
		echo `date` "Found no $(UI_APP_NAME) UI application, creating..."; \
		cd $(UI_SRC_PATH) && ibmcloud -q app push $(UI_APP_NAME); \
	fi
	@echo `date` "Creating $(UI_APP_NAME) UI application is finished"
	@ui_url="$$(ibmcloud -q app list | sed '1,4d' | grep "^$(UI_APP_NAME)" | sed 's/^.* \(.*$$\)/\1/' )"; \
	echo `date` "$(UI_APP_NAME) UI application URL is ""$$ui_url"
	@echo `date` "Current applications:"
	@ibmcloud -q app list

## teardown-ui			: Delete UI application
teardown-ui : prereqs
	@echo `date` "Deleting $(UI_APP_NAME) UI application"
	@echo `date` "Waiting $(WAIT_RESPONSE) seconds..."; \
		sleep $(WAIT_RESPONSE)
	@echo `date` "Current applications:"
	@ibmcloud -q app list
	@if [[ $$(ibmcloud -q app list | cut -d' ' -f1 | sed '1,4d' | grep -Fx "$(UI_APP_NAME)") ]]; then \
		echo `date` "Found $(UI_APP_NAME) UI application, deleting..."; \
		ibmcloud -q app delete $(UI_APP_NAME) -f -r; \
	else \
		echo `date` "There is no $(UI_APP_NAME) UI application, skipping..."; \
	fi
	@echo `date` "Finished deleting $(UI_APP_NAME) UI application"
	@echo `date` "Current applications:"
	@ibmcloud -q app list

## deploy-all			: Create all instances
deploy-all : deploy-stt deploy-nlu deploy-es deploy-db deploy-func deploy-ui

## teardown-all			: Delete all instances
teardown-all : teardown-ui teardown-func teardown-db teardown-es teardown-nlu teardown-stt

