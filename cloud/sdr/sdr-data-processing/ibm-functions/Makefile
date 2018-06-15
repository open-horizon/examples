# Targets for creating IBM Functions objects.
# To run these, you must have the 'bx wsk' CLI installed and be logged into it and set the target org and space.

WSK_PACKAGE ?= message-hub-evnts
WSK_ACTION ?= process-message
ACTION_CODE ?= actions/msgreceive.js
WSK_TRIGGER ?= message-hub-events-trgr
WSK_RULE ?= message-hub-events-rule-2

MSGHUB_NAME ?= Message_Hub-rt
MSGHUB_TOPIC ?= sdr-audio

# Upload a new version of the action code
update-action: $(ACTION_CODE)
	bx wsk action update --kind nodejs:8 $(WSK_PACKAGE)/$(WSK_ACTION) $<

# This requires node.js be installed. On Mac: brew install node
test-action: actions/msgreceive_test.js $(ACTION_CODE)
	node $<

# Create a trigger and action. Note: haven't tried this yet!!
create: $(ACTION_CODE)
	bx wsk package refresh
	bx wsk trigger create $(WSK_TRIGGER) --feed Bluemix_$(MSGHUB_NAME)_Bruce_Credentials/messageHubFeed --param isJSONData true --param topic $(MSGHUB_TOPIC)
	bx wsk action create $(WSK_PACKAGE)/$(WSK_ACTION) $< # do not think these flag from igor are needed: --param kafka $(KAFKA_INSTANCE) --param destTopic $(DEST_TOPIC)
	bx wsk rule create $(WSK_RULE) $(WSK_TRIGGER) $(WSK_PACKAGE)/$(WSK_ACTION)

# Remove a trigger and action. Note: haven't tried this yet!!
remove:
	bx wsk package refresh
	bx wsk rule delete --disable $(WSK_RULE)
	bx wsk trigger delete $(WSK_TRIGGER)
	bx wsk action delete $(WSK_PACKAGE)/$(WSK_ACTION)

.PHONY: update-action test-action create remove
