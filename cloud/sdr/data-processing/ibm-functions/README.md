# IBM Function actions for SDR Data Processing

## IBM Functions Objects Being Used

- packages:
    - MsgHub feed: /Hovitos_dev/openwhisk-messagehub
    - Our actions: /Hovitos_dev/message-hub-evnts
- action: /Hovitos_dev/message-hub-evnts/process-message
- trigger: /Hovitos_dev/message-received-trigger
- rule: /Hovitos_dev/message-received-rule
- Message Hub instance: msghub-hzndemo (Region: US South, CF Org: Hovitos, CF Space: dev)

You can view the above objects with commands like:
```
ic fn namespace get  # the rest of the cmds assume /Hovitos_dev is the default namespace
ic fn list
ic fn package get --summary openwhisk-messagehub
ic fn package get --summary message-hub-evnts
ic fn trigger get message-received-trigger
ic fn rule get message-received-rule
ic fn action get message-hub-evnts/process-message
```

## Setup

```
export STT_USERNAME="<speech-to-text-user>"
export STT_PASSWORD="<speech-to-text-pw>"
export NLU_USERNAME="<natural-language-understand-user>"
export NLU_PASSWORD="<natural-language-understand-pw>"
export SDR_DB_URL="postgres://user:pw@$host:port/dbname"
```

## Some Commands to Create/Update the IBM Functions Objects

Note: this assumes the /Hovitos_dev/openwhisk-messagehub package has already been created and has your IBM Message Hub credentials in the parameters section.

Note: use the make targets to create the official objects.

```
ic fn trigger update message-received-trigger --param isJSONData true --param isBinaryValue false

ic fn rule create message-received-rule message-received-trigger message-hub-evnts/process-message
```

## Test the Action Locally

**Note: this is not useful as-is, because it doesn't pass audio to the action.**

```
make test-action
```

## Upload the Action to the IBM Functions Service
```
make update-action
```

## Make IBM Functions Invoke the Action
```
# invoke the action via the trigger
ic fn trigger fire message-hub-events-trgr -P actions/test/param.json
# invoke the action directly
ic fn action invoke message-hub-evnts/process-message -b -P actions/test/param.json
```

## See the Actions that Get Invoked
```
ic fn activation poll
```

## Test/Upload the golang action (not currently used)
```
make test-go-action  # Test the Action Locally - this is not useful as-is, because it doesn't pass audio to the action serialized with gob.
make exec.zip  # Upload the Action to the IBM Functions Service
```

## TODOs
- Bind STT, NLU, and Compose instances to our action, so it will have those credentials: https://console.bluemix.net/docs/openwhisk/binding_services.html

## Action Details

- Credentials needed inside the action should be passed as params when creating the action object
- The Functions/MsgHub feed and trigger:
    - https://github.com/IBM/ibm-cloud-functions-message-hub-trigger  and  https://github.com/apache/incubator-openwhisk-package-kafka
- Node.js actions:
    - Node.js packages that are pre-installed in the IBM Funcions Node.js 8 environment: https://console.bluemix.net/docs/openwhisk/openwhisk_reference.html#openwhisk_ref_javascript_environments
    - Packaging your action as a nodejs module: https://console.bluemix.net/docs/openwhisk/openwhisk_actions.html#openwhisk_js_packaged_action
    - Trigger parameters: https://github.com/apache/incubator-openwhisk-package-kafka#creating-a-trigger-that-listens-to-an-ibm-messagehub-instance
    - Watson Node.js SDK package: https://www.npmjs.com/package/watson-developer-cloud
    - 3.7.0 Version of Watson Node.js SDK: https://github.com/watson-developer-cloud/node-sdk/blob/v3.7.0/README.md
    - Watson STT Nodejs SDK: https://www.ibm.com/watson/developercloud/speech-to-text/api/v1/node.html?node
    - Watson NLU Node.js SDK: https://www.ibm.com/watson/developercloud/natural-language-understanding/api/v1/?node#post-analyze
    - JS Protobufs: https://www.npmjs.com/package/protobufjs
    - Using binary data with msg hub/openwhisk triggers: https://medium.com/openwhisk/integrating-openwhisk-with-message-hub-now-with-binary-data-81b5b2dc1d69
- Golang docker actions (not currently used):
    - https://console.bluemix.net/docs/openwhisk/openwhisk_actions.html#openwhisk_actions (sections: "Creating Go actions" and "Creating Docker actions")
    - https://www.ibm.com/blogs/bluemix/2017/01/docker-bluemix-openwhisk/
    - https://console.bluemix.net/docs/openwhisk/openwhisk_reference.html#openwhisk_ref_docker

## Problems I Ran Into
- Difficulty creating package, action, trigger, rule
- Bug in `ic fn package refresh` (now fixed)
- Go action gets msg via command line, so command line length limitation on audio data
- Applicationn error in node.js action when using when using protobufjs to deserialize msg
- Trigger doesn't invoke action when msg > 250K (open support ticket)
- watson-developer-cloud SDK node.js module is older version that i can't find the docs for (this is preventing me from running NLU).
