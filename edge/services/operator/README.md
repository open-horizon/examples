# Horizon Operator Example Edge Service

This is a simple example of using and creating an operator as edge service.

- [Preconditions for Using the Operator Example Edge Service](#preconditions)
- [Using the Operator Example Edge Service with Deployment Pattern](PatternRegister.md)
- [Creating Your Own Operator Edge Service](CreateService.md)


- [Using the Operator Example Edge Service with Deployment Policy](#using-operator-pattern)

- Further Learning - to see more Horizon features demonstrated, continue on to the [cpu2evtstreams example](../../evtstreams/cpu2evtstreams).

## <a id=preconditions></a> Preconditions for Using the Operator Example Edge Service

If you haven't done so already, you must do these steps before proceeding with the operator example:

1. Install the Horizon management infrastructure (exchange and agbot).

2. Install the Horizon agent on your edge cluster and configure it to point to your Horizon exchange.

3. As part of the infrasctucture installation process for IBM Edge Computing Manager a file called `agent-install.cfg` was created that contains the values for `HZN_ORG_ID` and the exchange and css url values. Locate this file and set those environment variables in your shell now:

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

## <a id=using-helloworld-pattern></a> Using the Operator Example Edge Service with Deployment Policy

In the following steps you will deploy the `ibm.operator` to your edge cluster. This operator will then create three pods running the `ibm.helloworld` service. 

1. Get the required node policy file:
```bash
wget https://raw.githubusercontent.com/open-horizon/examples/master/edge/services/operator/horizon/node.policy.json
```

- Below is the `node_policy.json` file you obtained in step one:

```json
{
  "properties": [
    { "name": "openhorizon.service", "value": "ibm.operator" }
  ],
  "constraints": [
  ]
}
```

- It provides values for one `properties` (`openhorizon.service`), that will effect which services get deployed to this edge node, and states no `constraints`.

2. Register your Node Policy with this policy
```bash
hzn register --policy node.policy.json
```

3. When the registration completes, use the following command to review the Node Policy:
```bash
hzn policy list
```

4. The edge device will make an agreement with one of the Horizon agreement bots (this typically takes about 15 seconds). Repeatedly query the agreements of this device until the `agreement_finalized_time` and `agreement_execution_start_time` fields are filled in:

```bash
hzn agreement list
```

- Below is the `business_policy.json` file you just grabbed in step one:

```json
{
  "mycluster/ibm.operator_1.0.0": {
    "owner": "mycluster/operator1",
    "label": "ibm.operator Business Policy for amd64",
    "description": "A super-simple sample Horizon Business Policy",
    "service": {
      "name": "ibm.operator",
      "org": "mycluster",
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
      "openhorizon.service == ibm.operator"
    ],
    "userInput": [
      {
        "serviceOrgid": "mycluster",
        "serviceUrl": "ibm.operator",
        "serviceVersionRange": "[0.0.0,INFINITY)",
        "inputs": [
          {
            "name": "HW_WHO",
            "value": "from the cluster!"
          }
        ]
      }
    ],
    "created": "2020-05-21T17:18:34.956Z[UTC]",
    "lastUpdated": "2020-05-21T19:50:56.937Z[UTC]"
  }
}
```







1. Register your edge cluster with Horizon to use the operator pattern:

```bash
hzn register -p IBM/pattern-ibm.helloworld -s ibm.helloworld --serviceorg IBM
```
 - **Note**: using the `-s` flag with the `hzn register` command will cause Horizon to wait until agreements are formed and the service is running on your edge node to exit, or alert you of any errors encountered during the registration process. 

2. After the agreement is made, list the docker container edge service that has been started as a result:

``` bash
sudo docker ps
```

3. See the helloworld service output:

``` bash
hzn service log -f ibm.helloworld
```
 - **Note**: Press **Ctrl C** to stop the command output.

4. Unregister your edge node (which will also stop the myhelloworld service):

```bash
hzn unregister -f
```
