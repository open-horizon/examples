# Horizon Hello World Example Edge Service

A simple example of using and creating a Horizon edge service.

- [Preconditions for Using the Hello World Example Edge Service](#preconditions)
- [Using the Hello World Example Edge Service with Deployment Pattern](#using-helloworld-pattern)
- [Using the Hello World Example Edge Service with Deployment Policy](#using-helloworld-policy)
- [Creating Your Own Hello World Edge Service](CreateService.md)
- Further Learning - to see more Horizon features demonstrated, continue on to the [cpu2msghub example](../../evtstreams/cpu2evtstreams).

## <a id=preconditions></a> Preconditions for Using the Hello World Example Edge Service

If you haven't done so already, you must do these steps before proceeding with the helloworld example:

1. Install the Horizon management infrastructure (exchange and agbot).

2. Install the Horizon agent on your edge device and configure it to point to your Horizon exchange.

3. Set your exchange org:

```bash
export HZN_ORG_ID="<your-cluster-name>"
```

4. Create a cloud API key that is associated with your Horizon instance, set your exchange user credentials, and verify them:

```bash
export HZN_EXCHANGE_USER_AUTH="iamapikey:<your-API-key>"
hzn exchange user list
```

5. Choose an id and token for your edge node, create it, and verify it:

```bash
export HZN_EXCHANGE_NODE_AUTH="<choose-any-node-id>:<choose-any-node-token>"
hzn exchange node create -n $HZN_EXCHANGE_NODE_AUTH
hzn exchange node confirm
```

## <a id=using-helloworld-pattern></a> Using the Hello World Example Edge Service with Deployment Pattern

1. Register your edge node with Horizon to use the helloworld pattern:

```bash
hzn register -p IBM/pattern-ibm.helloworld
```

2. The edge device will make an agreement with one of the Horizon agreement bots (this typically takes about 15 seconds). Repeatedly query the agreements of this device until the `agreement_finalized_time` and `agreement_execution_start_time` fields are filled in:

```bash
hzn agreement list
```

3. Once the agreement is made, list the docker container edge service that has been started as a result:

``` bash
sudo docker ps
```

4. See the helloworld service output:

  on **Linux**:

  ```bash
  sudo tail -f /var/log/syslog | grep helloworld[[]
  ```

  on **Mac**:

  ```bash
  sudo docker logs -f $(sudo docker ps -q --filter name=helloworld)
  ```

5. Unregister your edge node, stopping the helloworld service:

```bash
hzn unregister -f
```

## <a id=using-helloworld-policy></a> Using the Hello World Example Edge Service with Deployment Policy

The Horizon Policy mechanism offers an alternative to using Deployment Patterns. Policies provide much finer control over the deployment placement of edge services. It also provides a greater separation of concerns, allowing Edge Nodes owners, Service code developers, and Business owners to each independently articulate their own Policies. There are therefore three types of Horizon Policies:

1. Node Policy (provided at registration time by the node owner)

2. Service Policy (may be applied to a published Service in the Exchange)

3. Business Policy (which approximately corresponds to a Deployment Pattern)

### Node Policy

- As an alternative to specifying a Deployment Pattern when you register your Edge Node, you may register with a Node Policy.

1. Make sure your Edge Node is not registered by running:

```bash
hzn unregister -f
```

- Now let's register using the `horizon/node_policy.json` file:

```json
{
  "properties": [
    { "name": "model", "value": "Thingamajig ULTRA" },
    { "name": "serial", "value": 9123456 },
    { "name": "configuration", "value": "Mark-II-PRO" }
  ],
  "constraints": [
  ]
}
```

- It provides values for three `properties` (`model`, `serial`, and `configuration`). It states no `constraints`, so any appropriately signed and authorized code can be deployed on this Edge Node,

2. Register your Node Policy using this command:

```bash
hzn register --policy horizon/node_policy.json
```

3. When the registration completes, use the following command to review the Node Policy:

```bash
hzn policy list
```

- Notice that in addition to the three `properties` stated in the node_policy.json file, Horizon has added a few more (openhorizon.cpu, openhorizon.arch, and openhorizon.memory). Horizon provides this additional information automatically and these `properties` may be used in any of your Policy `constraints`.

### Service Policy

- Like the other two Policy types, Service Policy contains a set of `properties` and a set of `constraints`. The `properties` of a Service Policy could state characteristics of the Service code that Node Policy authors or Business Policy authors may find relevant. The `constraints` of a Service Policy can be used to restrict where this Service can be run. The Service developer could, for example, assert that this Service requires a particular hardware setup such as CPU/GPU constraints, memory constraints, specific sensors, actuators or other peripheral devices required, etc.

- Now let's attach this Service Policy to the helloworld Service previously published using the `horizon/service_policy.json` file:

```json
{
  "properties": [
  ],
  "constraints": [
    "model == \"Whatsit ULTRA\" OR model == \"Thingamajig ULTRA\""
  ]
}
```

- Note this simple Service Policy doesn't provide any `properties`, but it does have a `constraint`. This example `constraint` is one that a Service developer might add, stating that their Service must only run on the models named `Whatsit ULTRA` or `Thingamajig ULTRA`. If you recall the Node Policy we used above, the model `property` was set to `Thingamajig ULTRA`, so this Service should be compatible with our Edge Node.

1. To attach the example Service policy to this service, use the following command (substituting your service name):

```bash
hzn exchange service addpolicy -f horizon/service_policy.json <published-helloworld-service-name>
```

2. Once that completes, you can look at the results with the following command:

```bash
hzn exchange service listpolicy <published-helloworld-service-name>
```

- Notice that Horizon has again automatically added some additional `properties` to your Policy. These generated property values can be used in `constraints` in Node Policies and Business Policies.

- Now that we have set up the Policies for an Edge Node and the Policies for a published Service, we can move on to the final step of defining a Business Policy to tie them all together and cause software to be automatically deployed on your Edge Node.

### Business Policy

- Business Policy (sometimes called Deployment Policy) is what ties together Edge Nodes, Published Services, and the Policies defined for each of those, making it roughly analogous to the Deployment Patterns you have previously worked with.

- Business Policy, like the other two Policy types, contains a set of `properties` and a set of `constraints`, but it contains other things as well. For example, it explicitly identifies the Service it will cause to be deployed onto Edge Nodes if negotiation is successful, in addition to configuration variable values, performing the equivalent function to the `-f horizon/userinput.json` clause of a Deployment Pattern `hzn register ...` command. The Business Policy approach for configuration values is more powerful because this operation can be performed centrally (no need to connect directly to the Edge Node).

- Below is the `horizon/business_policy.json` file used for this example:

```json
{
  "label": "$SERVICE_NAME Business Policy for $ARCH",
  "description": "A super-simple sample Horizon Business Policy",
  "service": {
    "name": "$SERVICE_NAME",
    "org": "$HZN_ORG_ID",
    "arch": "$ARCH",
    "serviceVersions": [
      {
        "version": "$SERVICE_VERSION",
        "priority":{}
      }
    ]
  },
  "properties": [
  ],
  "constraints": [
    "serial >= 9000000",
    "model == \"Thingamajig ULTRA\""
  ],
  "userInput": [
    {
      "serviceOrgid": "$HZN_ORG_ID",
      "serviceUrl": "$SERVICE_NAME",
      "serviceVersionRange": "[0.0.0,INFINITY)",
      "inputs": [
        {
          "name": "HW_WHO",
          "value": "Valued Customer"
        }
      ]
    }
  ]
}
```

- This simple example of a Business Policy doesn't provide any `properties`, but it does have two `constraints` that are satisfied by the `properties` set in the `horizon/node_policy.json` file, so this Business Policy should successfully deploy our Service onto the Edge Node.

- At the bottom, the userInput section has the same purpose as the horizon/userinput.json files provided for other examples if the given services requires them. In this case the helloworld service defines only one configuration variable, HW_WHO, and the userInput section here provides a value for HW_WHO (i.e., Valued Customer).

1. To publish this Business Policy to the Exchange and get this Service running on the Edge Node edit the `horizon/business_policy.json` file to correctly identify your specific Service name, org, version, arch, etc. When your Business Policy is ready, run the following command to publish it, giving it a memorable name (bizPolicy1 in this example):

```bash
hzn exchange business addpolicy -f horizon/business_policy.json bizPolicy1
```

2. Once that competes, you can look at the results with the following command, substituting your own org id:

```bash
hzn exchange business listpolicy major-peacock-icp-cluster/bizPolicy1
```

- The results should look very similar to your original `horizon/business_policy.json` file, except that `owner`, `created`, and `lastUpdated` and a few other fields have been added.

3. The edge device will make an agreement with one of the Horizon agreement bots (this typically takes about 15 seconds). Repeatedly query the agreements of this device until the `agreement_finalized_time` and `agreement_execution_start_time` fields are filled in:

```bash
hzn agreement list
```

4. Once the agreement is made, list the docker container edge service that has been started as a result:

```bash
sudo docker ps
```

5. See the helloworld service output:

  on **Linux**:

  ```bash
  sudo tail -f /var/log/syslog | grep helloworld[[]
  ```

  on **Mac**:

  ```bash
  sudo docker logs -f $(sudo docker ps -q --filter name=helloworld)
  ```

6. Unregister your edge node, stopping the helloworld service:

```bash
hzn unregister -f
```
