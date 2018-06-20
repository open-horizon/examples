# IBM Function actions for SDR Data Processing

## IBM Functions Objects Being Used

- package: /Hovitos_dev/message-hub-evnts
- actions: /Hovitos_dev/message-hub-evnts/process-message
- triggers: /Hovitos_dev/message-hub-events-trgr
- rule: /Hovitos_dev/message-hub-events-rule-2
- Message Hub instance: Message Hub-rt (Region: US South, CF Org: Hovitos, CF Space: dev)

## Setup

```
export WATSON_STT_USERNAME="<speech-to-text-user>"
export WATSON_STT_PASSWORD="<speech-to-text-pw>"
```

## Test the Action Locally
```
make test-action
```

## Upload the Action to the IBM Functions Service
```
make update-action
```

## See the Actions that Get Invoked
```
bx wsk activation poll
```

## Action Details

- Node.js packages that are pre-installed in the IBM Funcions Node.js 8 environment: https://console.bluemix.net/docs/openwhisk/openwhisk_reference.html#openwhisk_ref_javascript_environments
- Credentials needed inside the action should be passed as params when creating the action object
- Watson Node.js package, including speec-to-text: https://www.npmjs.com/package/watson-developer-cloud#speech-to-text