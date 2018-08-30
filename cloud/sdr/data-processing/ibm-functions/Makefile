# Targets for creating IBM Functions objects.
# To run these, you must have the 'ic fn' CLI installed and be logged into it and set the target org and space.

WSK_PACKAGE ?= message-hub-evnts
WSK_MSGHUB_FEED ?= openwhisk-messagehub/messageHubFeed
WSK_ACTION ?= process-message   # this is the node.js action we are not using
ACTION_CODE ?= actions/msgreceive.js
# PROTO_FILE ?= ../../../../edge/msghub/sdr2msghub/audiomsg.proto
# PROTO_JS ?= actions/audiomsg.js
# PKG_JSON ?= actions/package.json
# ACTION_ZIP ?= msgreceive.zip
WSK_TRIGGER ?= message-received-trigger
WSK_RULE ?= message-received-rule

# not used:
WSK_GO_ACTION ?= process-message-go
ACTION_GO_CODE ?= actions/msgreceive.go

MSGHUB_NAME ?= msghub-hzndemo
MSGHUB_TOPIC ?= sdr-audio

# Build the js version of the protobufs struct definition
# $(PROTO_JS): $(PROTO_FILE)
# 	actions/node_modules/protobufjs/cli/bin/pbjs -t static-module -w commonjs -o $(PROTO_JS) $(PROTO_FILE)

# This requires node.js be installed. On Mac: brew install node
action-test: actions/msgreceive_test.js $(ACTION_CODE)
	: $${STT_USERNAME:?} $${STT_PASSWORD:?} $${NLU_USERNAME:?} $${NLU_PASSWORD:?} $${SDR_DB_URL:?}   # this verifies these env vars are set
	node $<

# Upload a new version of the action code
action-update: $(ACTION_CODE)
	: $${STT_USERNAME:?} $${STT_PASSWORD:?} $${NLU_USERNAME:?} $${NLU_PASSWORD:?} $${SDR_DB_URL:?}   # this verifies these env vars are set
	# rm -f $(ACTION_ZIP)
	# cd actions; zip -qr ../$(ACTION_ZIP) * -x test/\* -x msgreceive_test.js; cd ..
	ic fn action update $(WSK_PACKAGE)/$(WSK_ACTION) $(ACTION_CODE) --kind nodejs:8 --memory 512 --timeout 300000 --param watsonSttUsername $(STT_USERNAME) --param watsonSttPassword $(STT_PASSWORD) --param watsonNluUsername $(NLU_USERNAME) --param watsonNluPassword $(NLU_PASSWORD) --param postgresUrl $(SDR_DB_URL)

# Create a trigger, action, and rule
create: $(ACTION_CODE)
	: $${STT_USERNAME:?} $${STT_PASSWORD:?} $${NLU_USERNAME:?} $${NLU_PASSWORD:?} $${SDR_DB_URL:?}   # this verifies these env vars are set
	#ic fn package refresh  # this gives an auth error for some reason
	ic fn trigger create $(WSK_TRIGGER) --feed $(WSK_MSGHUB_FEED) --param isJSONData true --param isBinaryValue false --param topic $(MSGHUB_TOPIC)
	# rm -f $(ACTION_ZIP)
	# cd actions; zip -qr ../$(ACTION_ZIP) * -x test/\* -x msgreceive_test.js; cd ..
	ic fn action create $(WSK_PACKAGE)/$(WSK_ACTION) $(ACTION_CODE) --kind nodejs:8 --memory 512 --timeout 300000 --param watsonSttUsername $(STT_USERNAME) --param watsonSttPassword $(STT_PASSWORD) --param watsonNluUsername $(NLU_USERNAME) --param watsonNluPassword $(NLU_PASSWORD) --param postgresUrl $(SDR_DB_URL)
	ic fn rule create $(WSK_RULE) $(WSK_TRIGGER) $(WSK_PACKAGE)/$(WSK_ACTION)

# Remove the trigger, action and rule
remove:
	#ic fn package refresh  # this gives an auth error for some reason
	ic fn rule delete --disable $(WSK_RULE)
	ic fn trigger delete $(WSK_TRIGGER)
	ic fn action delete $(WSK_PACKAGE)/$(WSK_ACTION)


# The golang-related targets. Not used anymore!
# Note: this is not useful as-is, because it doesn't pass audio to the action serialized with gob
# This requires go be installed. On Mac: brew install go
test-go-action: $(ACTION_GO_CODE)
	: $${STT_USERNAME:?} $${STT_PASSWORD:?} $${NLU_USERNAME:?} $${NLU_PASSWORD:?} $${SDR_DB_URL:?}   # this verifies these env vars are set
	go run $< '{ "messages": [ { "value": "\"this is my first msg\"", "topic": "$(MSGHUB_TOPIC)", "partition": 1 }, { "value": "\"this is my 2nd msg\"", "topic": "$(MSGHUB_TOPIC)", "partition": 1 } ], "watsonSttUsername": "$(STT_USERNAME)", "watsonSttPassword": "$(STT_PASSWORD)", "watsonNluUsername": "$(NLU_USERNAME)", "watsonNluPassword": "$(NLU_PASSWORD)", "postgresUrl": "$(SDR_DB_URL)" }'

# Build and Upload a new version of the action code
exec.zip: $(ACTION_GO_CODE)
	: $${STT_USERNAME:?} $${STT_PASSWORD:?} $${NLU_USERNAME:?} $${NLU_PASSWORD:?} $${SDR_DB_URL:?}   # this verifies these env vars are set
	GOOS=linux GOARCH=amd64 go build -o exec $<
	zip $@ exec
	ic fn action update $(WSK_PACKAGE)/$(WSK_GO_ACTION) $@ --native --param watsonSttUsername $(STT_USERNAME) --param watsonSttPassword $(STT_PASSWORD) --param watsonNluUsername $(NLU_USERNAME) --param watsonNluPassword $(NLU_PASSWORD) --param postgresUrl $(SDR_DB_URL)

.PHONY: action-update action-test create remove test-go-action
