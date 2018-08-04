# Targets for creating IBM Functions objects.
# To run these, you must have the 'bx wsk' CLI installed and be logged into it and set the target org and space.

WSK_PACKAGE ?= message-hub-evnts
# WSK_ACTION ?= process-message   # this is the node.js action we are not using
# ACTION_CODE ?= actions/msgreceive.js
WSK_GO_ACTION ?= process-message-go
ACTION_GO_CODE ?= actions/msgreceive.go
WSK_TRIGGER ?= message-hub-events-trgr
#WSK_RULE ?= message-hub-events-rule-2
#WSK_RULE ?= message-hub-events-msghub-process-message-go
WSK_RULE ?= message-hub-evnts-process-message-go_message-hub-events-trgr

MSGHUB_NAME ?= Message_Hub-rt
MSGHUB_TOPIC ?= sdr-audio

# This requires node.js be installed. On Mac: brew install node
# test-action: actions/msgreceive_test.js $(ACTION_CODE)
# 	: $${STT_USERNAME:?} $${STT_PASSWORD:?}   # this verifies these env vars are set
# 	node $<

# Upload a new version of the action code
# update-action: $(ACTION_CODE)
# 	: $${STT_USERNAME:?} $${STT_PASSWORD:?}   # this verifies these env vars are set
# 	bx wsk action update $(WSK_PACKAGE)/$(WSK_ACTION) $< --kind nodejs:8 --param watsonSttUsername $(STT_USERNAME) --param watsonSttPassword $(STT_PASSWORD)

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
	bx wsk action update $(WSK_PACKAGE)/$(WSK_GO_ACTION) $@ --native --param watsonSttUsername $(STT_USERNAME) --param watsonSttPassword $(STT_PASSWORD) --param watsonNluUsername $(NLU_USERNAME) --param watsonNluPassword $(NLU_PASSWORD) --param postgresUrl $(SDR_DB_URL)

# Create a trigger and action. Note: haven't tried this yet!! (I did it thru the UI)
create: exec.zip
	: $${STT_USERNAME:?} $${STT_PASSWORD:?}   # this verifies these env vars are set
	bx wsk package refresh
	bx wsk trigger create $(WSK_TRIGGER) --feed Bluemix_$(MSGHUB_NAME)_Bruce_Credentials/messageHubFeed --param isJSONData true --param topic $(MSGHUB_TOPIC)  # not sure yet what to specify for --feed
	#bx wsk action create $(WSK_PACKAGE)/$(WSK_ACTION) $< --kind nodejs:8 --param watsonSttUsername $(STT_USERNAME) --param watsonSttPassword $(STT_PASSWORD)
	bx wsk action create $(WSK_PACKAGE)/$(WSK_GO_ACTION) $< --native --param watsonSttUsername $(STT_USERNAME) --param watsonSttPassword $(STT_PASSWORD)
	bx wsk rule create $(WSK_RULE) $(WSK_TRIGGER) $(WSK_PACKAGE)/$(WSK_GO_ACTION)

# Remove a trigger and action. Note: haven't tried this yet!!
remove:
	bx wsk package refresh
	bx wsk rule delete --disable $(WSK_RULE)
	bx wsk trigger delete $(WSK_TRIGGER)
	#bx wsk action delete $(WSK_PACKAGE)/$(WSK_ACTION)
	bx wsk action delete $(WSK_PACKAGE)/$(WSK_GO_ACTION)

.PHONY: update-action test-action create remove
