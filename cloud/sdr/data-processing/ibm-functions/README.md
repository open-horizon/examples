# IBM Function actions for SDR Data Processing

## IBM Functions Objects Being Used

- package: /Hovitos_dev/message-hub-evnts
- actions: /Hovitos_dev/message-hub-evnts/process-message-go
- triggers: /Hovitos_dev/message-hub-events-trgr
- rule: /Hovitos_dev/message-hub-evnts-process-message-go_message-hub-events-trgr
- Message Hub instance: Message Hub-rt (Region: US South, CF Org: Hovitos, CF Space: dev)

You can view the above objects with commands like:
```
bx wsk list
bx wsk trigger list
bx wsk trigger get message-hub-events-trgr
bx wsk rule get message-hub-evnts-process-message-go_message-hub-events-trgr
bx wsk action get message-hub-evnts/process-message-go
```

## Setup

```
export STT_USERNAME="<speech-to-text-user>"
export STT_PASSWORD="<speech-to-text-pw>"
```

## Test the Action Locally
```
make test-go-action
```

## Upload the Action to the IBM Functions Service
```
make exec.zip
```

## See the Actions that Get Invoked
```
bx wsk activation poll
```

## Action Details

- Credentials needed inside the action should be passed as params when creating the action object
- Golang docker actions:
    - https://console.bluemix.net/docs/openwhisk/openwhisk_actions.html#openwhisk_actions (sections: "Creating Go actions" and "Creating Docker actions")
    - https://www.ibm.com/blogs/bluemix/2017/01/docker-bluemix-openwhisk/
    - https://console.bluemix.net/docs/openwhisk/openwhisk_reference.html#openwhisk_ref_docker 
- Node.js actions:
    - Node.js packages that are pre-installed in the IBM Funcions Node.js 8 environment: https://console.bluemix.net/docs/openwhisk/openwhisk_reference.html#openwhisk_ref_javascript_environments
    - Watson Node.js package, including speec-to-text: https://www.npmjs.com/package/watson-developer-cloud#speech-to-text