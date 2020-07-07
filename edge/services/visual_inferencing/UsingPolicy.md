# Horizon Object Detection and Classification Example Edge Service

This is a simple example of using and creating an operator as edge service.

- [Preconditions for Using the Object Detection and Classification Example Edge Service](#preconditions)
- [Using the Object Detection and Classification Example Edge Service with Deployment Policy](#using-detect-policy)
- [Using the Object Detection and Classification Example Edge Service with Deployment Pattern](PatternRegister.md)
- [Creating Your Own Object Detection and Classification Edge Service](CreateService.md)
- Further Learning - to see more Horizon features demonstrated, continue on to the [cpu2evtstreams example](../../evtstreams/cpu2evtstreams).

## <a id=preconditions></a> Preconditions for Using the Object Detection and Classification Example Edge Service

If you haven't done so already, you must do these steps before proceeding with the object detection and classification example:

1. Install the Horizon management hub (exchange and agbot).

2. Install the Horizon agent on your edge device and configure it to point to your Horizon exchange.

3. As part of the management hub installation process a file called `agent-install.cfg` was created that contains the values for `HZN_ORG_ID` and the exchange and css URL values. Locate this file and set those environment variables in your shell now:

```bash
eval export $(cat agent-install.cfg)
```

 - **Note**: if for some reason you disconnected from ssh or your command line closes, run the above command again to set the required environment variables.

4. In addition to the file above, an API key associated with your Horizon instance would have been created, set the exchange user credentials, and verify them:

```bash
export HZN_EXCHANGE_USER_AUTH=iamapikey:<horizon-API-key>
hzn exchange user list
```

5. Choose an ID and token for your edge node, create it, and verify it:

```bash
export HZN_EXCHANGE_NODE_AUTH="<choose-any-node-id>:<choose-any-node-token>"
hzn exchange node create -n $HZN_EXCHANGE_NODE_AUTH
hzn exchange node confirm
```

6. If you have not done so already, unregister your node before moving on:

 ```bash
hzn unregister -f
```

## <a id=using-detect-policy></a> Using the Object Detection and Classification Example Edge Service with Deployment Policy

1. Get the required node policy and user input files on your edge device:

- if your edge device **does not** have a GPU, run the following commands:
  ```bash
  wget https://raw.githubusercontent.com/open-horizon/examples/master/edge/services/visual_inferencing/yolocpu/horizon/node.policy.json
  wget https://raw.githubusercontent.com/open-horizon/examples/master/edge/services/visual_inferencing/yolocpu/horizon/userinput.json
  ```
- if your edge device **does** have a GPU, run the following commands:
  ```bash
  wget https://raw.githubusercontent.com/open-horizon/examples/master/edge/services/visual_inferencing/yolocuda/horizon/node.policy.json
  wget https://raw.githubusercontent.com/open-horizon/examples/master/edge/services/visual_inferencing/yolocuda/horizon/userinput.json
  ```

- Below is the `node.policy.json` file you obtained in the step above if you do not have a GPU on your edge device:

```json
{
  "properties": [
    { "name": "openhorizon.example", "value": "visual_inferencing" },
    { "name": "GPUenabled", "value": "false" }
  ],
  "constraints": [
  ]
}
```

- Below is the `node.policy.json` file you obtained in the step above if you do not have a GPU on your edge device:

```json
{
  "properties": [
    { "name": "openhorizon.example", "value": "visual_inferencing" },
    { "name": "GPUenabled", "value": "true" }
  ],
  "constraints": [
  ]
}
```

- Both provide one value for `properties` (`openhorizon.example`), that will effect which services get deployed to this edge node, and state no `constraints`.

2. Register your Node Policy with this policy

```bash
hzn register --policy node.policy.json
```

4. When the registration completes, use the following command to review the Node Policy:

```bash
hzn policy list
```

5. The edge device will make an agreement with one of the Horizon agreement bots (this typically takes about 15 seconds). Repeatedly query the agreements of this device until the `agreement_finalized_time` and `agreement_execution_start_time` fields are filled in:

```bash
hzn agreement list
```

- Below is the `service.policy.json` that was published into the Exchange when the `yolocpu` example was published:

```json
{
    "properties": [],
    "constraints": [
        "GPUenabled == false"
    ]
}
```

- Below is the `service.policy.json` that was published into the Exchange when the `yolocuda` example was published:

```json
{
    "properties": [],
    "constraints": [
        "GPUenabled == true"
    ]
}
```
- Notice that the developer who wrote and published these object detection and classification services included a service policy constraint that will enable edge devices to run `yolocuda` only if they have the node policy states `GPUenabled` is `true`, otherwise it will run `yolocpu`.

- Below is the example `deployment.policy.json` that has been published into the Exchange along with the `yolocpu` service:

```json
{
  "IBM/policy-yolocpu_1.0.0": {
    "owner": "mycluster/operator1",
    "label": "yolocpu Deployment Policy",
    "description": "A super-simple sample Horizon Deployment Policy",
    "service": {
      "name": "yolocpu",
      "org": "IBM",
      "arch": "*",
      "serviceVersions": [
        {
          "version": "1.0.0",
          "priority": {},
          "upgradePolicy": {}
        }
      ],
      "nodeHealth": {}
    },
    "constraints": [
      "openhorizon.example == visual_inferencing"
    ],
    "userInput": [
      {
        "serviceOrgid": "IBM",
        "serviceUrl": "yolocpu",
        "serviceVersionRange": "[0.0.0,INFINITY)",
        "inputs": []
      }
    ],
    "created": "2020-06-23T06:58:16.964Z[UTC]",
    "lastUpdated": "2020-06-23T06:58:16.964Z[UTC]"
  }
}
```

- Note: the `deployment.policy.json` for `yolocuda` is identical to the deployment policy for `yolocpu` so it is not listed here.

6. Verify that the service is up and running:

```bash
sudo docker ps 
```

7. You can now navigate to `http://0.0.0.0:5200` to confirm the object detection and classification is working as intended.

8. Unregister your edge node (which will also stop the object detection service):

```bash
hzn unregister -f
```