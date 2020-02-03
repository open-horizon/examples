# Horizon CPU To IBM Event Streams Service

This example illustrates a more realistic Horizon edge service by including additional aspects of typical edge services. 

- [Preconditions for Using the CPU To IBM Event Streams Example Edge Service](#preconditions)

- [Using the CPU To IBM Event Streams Example Edge Service with Deployment Pattern](#using-cpu2evtstreams-pattern)

- [Using the CPU To IBM Event Streams Example Edge Service with Deployment Policy](#using-cpu2evtstreams-policy)

- [Creating Your Own CPU To IBM Event Streams Example Edge Service](CreateService.md)

- For details about using this service, see [cpu2evtstreams.md](cpu2evtstreams.md).


## <a id=preconditions></a> Preconditions for Using the CPU To IBM Event Streams Example Edge Service

If you haven't done so already, you must do these steps before proceeding with the cpu2evtstreams example:

1. Install the Horizon management infrastructure (exchange and agbot).

2. Install the Horizon agent on your edge device and configure it to point to your Horizon exchange.

3. Set your exchange org:

```bash
export HZN_ORG_ID=<your-cluster-name>
```

4. Create a cloud API key that is associated with your Horizon instance, set your exchange user credentials, and verify them:

```bash
export HZN_EXCHANGE_USER_AUTH=iamapikey:<your-API-key>
hzn exchange user list
```

5. Choose an ID and token for your edge node, create it, and verify it:

```bash
export HZN_EXCHANGE_NODE_AUTH="<choose-any-node-id>:<choose-any-node-token>"
hzn exchange node create -n $HZN_EXCHANGE_NODE_AUTH
hzn exchange node confirm
```

6. Deploy (or get access to) an instance of IBM Event Streams that the cpu2evtstreams sample can send its data to. Ensure that the topic `cpu2evtstreams` is created in Event Streams. Using information from the Event Streams UI, `export` these environment variables:
    - `EVTSTREAMS_API_KEY`
    - `EVTSTREAMS_BROKER_URL`
    - `EVTSTREAMS_CERT_ENCODED` **(if using IBM Event Streams in IBM Cloud Private)** due to differences in the base64 command set this variable as follows based on the platform you're using:
        - on **Linux**: `EVTSTREAMS_CERT_ENCODED=“$(cat $EVTSTREAMS_CERT_FILE | base64 -w 0)”`
        - on **Mac**: `EVTSTREAMS_CERT_ENCODED=“$(cat $EVTSTREAMS_CERT_FILE | base64)”`
    - `EVTSTREAMS_CERT_FILE` **(if using IBM Event Streams in IBM Cloud Private)**


## <a id=using-cpu2evtstreams-pattern></a> Using the CPU To IBM Event Streams Edge Service with Deployment Pattern

1. Get the user input file for the cpu2evtstreams sample:
```bash
wget https://github.com/open-horizon/examples/raw/master/edge/evtstreams/cpu2evtstreams/horizon/use/userinput.json
```
2. Register your edge node with Horizon to use the cpu2evtstreams pattern:
```bash
hzn register -p IBM/pattern-ibm.cpu2evtstreams -f userinput.json
```

3. The edge device will make an agreement with one of the Horizon agreement bots (this typically takes about 15 seconds). Repeatedly query the agreements of this device until the `agreement_finalized_time` and `agreement_execution_start_time` fields are filled in:
```bash
hzn agreement list
```

4. Once the agreement is made, list the docker container edge service that has been started as a result:
```bash
sudo docker ps
```

5. On any machine, install [kafkacat](https://github.com/edenhill/kafkacat#install), then subscribe to the Event Streams topic to see the json data that cpu2evtstreams is sending:
  ```bash
  kafkacat -C -q -o end -f "%t/%p/%o/%k: %s\n" -b $EVTSTREAMS_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=token -X sasl.password=$EVTSTREAMS_API_KEY -X ssl.ca.location=$EVTSTREAMS_CERT_FILE -t cpu2evtstreams
  ```
6. See the cpu2evtstreams service output:

```bash
hzn service log -f ibm.cpu2evtstreams
```

7. Unregister your edge node, stopping the cpu2evtstreams service:
```bash
hzn unregister -f
```

## <a id=using-cpu2evtstreams-policy></a> Using the CPU To IBM Event Streams Service with Deployment Policy

- The Horizon Policy mechanism offers an alternative to using Deployment Patterns. Policies provide much finer control over agreement forming between Horizon Agents on Edge Nodes, and the Horizon AgBots. It also provides a greater separation of concerns, allowing Edge Nodes owners, Service code developers, and Business owners to each independently articulate their own Policies. There are therefore three types of Horizon Policies:

1. Node Policy (provided at registration time by the node owner)

2. Service Policy (may be applied to a published Service in the Exchange)

3. Business Policy (which approximately corresponds to a Deployment Pattern)

### Node Policy 

- As an alternative to specifying a Deployment Pattern when you register your Edge Node, you may register with a Node Policy.

1. Install `git`:

On **Linux**:

```bash
sudo apt install -y git
```

On **macOS**:

```bash
brew install git
```

2. If you have not done so already, clone this git repo:

```bash
git clone git@github.com:open-horizon/examples.git
```

3. Go to the `cpu2evtstreams` directory:

```bash
cd examples/edge/evtstreams/cpu2evtstreams
```

4. Make sure your Edge Node is not registered by running:

```bash
hzn unregister -f
```

- Now let's register using the `horizon/node_policy.json` file:

```bash
{
    "properties": [
        { "name": "model", "value": "Mac" },
        { "name": "year", "value": "2018" },
        { "name": "os", "value": "Mojave" }
    ],
    "constraints": []
}
```

- It provides values for three `properties` (`model`, `year`, and `os`). It states no `constraints`, so any appropriately signed and authorized code can be deployed on this Edge Node,

5. Get the user input file for the cpu2evtstreams sample:
```bash
wget https://github.com/open-horizon/examples/raw/master/edge/evtstreams/cpu2evtstreams/horizon/use/userinput.json
```

6. Register your Node Policy using this command:

```bash
hzn register --policy horizon/node_policy.json -f userinput.json
```

7. When the registration completes, use the following command to review the Node Policy:

```bash
hzn policy list
```

- Notice that in addition to the three `properties` stated in the node_policy.json file, Horizon has added a few more (openhorizon.cpu, openhorizon.arch, and openhorizon.memory). Horizon provides this additional information automatically and these `properties` may be used in any of your Policy `constraints`.

### Service Policy 

- Like the other two Policy types, Service Policy contains a set of `properties` and a set of `constraints`. The `properties` of a Service Policy could state characteristics of the Service code that Node Policy authors or Business Policy authors may find relevant. The `constraints` of a Service Policy can be used to restrict where this Service can be run. The Service developer could, for example, assert that this Service requires a particular hardware setup such as CPU/GPU constraints, memory constraints, specific sensors, actuators or other peripheral devices required, etc.

- Now let's attach this Service Policy to the cpu2evtstreams Service previously published using the `horizon/service_policy.json` file:

```bash
{
    "properties": [],
    "constraints": [
        "model == \"Mac\" OR model == \"Pi3B\"",
        "os == \"Mojave\""
    ]
}
```

- Note this simple Service Policy doesn't provide any `properties`, but it does have a `constraint`. This example `constraint` is one that a Service developer might add, stating that their Service must only run on the models named `Mac` or `Pi3B `. If you recall the Node Policy we used above, the model `property` was set to `Mac`, so this Service should be compatible with our Edge Node.

1. List the services in your org:
```bash
hzn exchange service list
```
2. To attach the example Service policy to this service, use the following command (substituting your service name):

```bash
hzn exchange service addpolicy -f horizon/service_policy.json <published-cpu2evtstreams-service-name>
```

3. Once that completes, you can look at the results with the following command:

```bash
hzn exchange service listpolicy <published-cpu2evtstreams-service-name>
```
- Notice that Horizon has again automatically added some additional `properties` to your Policy. These generated property values can be used in `constraints` in Node Policies and Business Policies.

- Now that we have set up the Policies for an Edge Node and the Policies for a published Service, we can move on to the final step of defining a Business Policy to tie them all together and cause software to be automatically deployed on your Edge Node.

### Business Policy 

- Business Policy (sometimes called Deployment Policy) is what ties together Edge Nodes, Published Services, and the Policies defined for each of those, making it roughly analogous to the Deployment Patterns you have previously worked with.

- Business Policy, like the other two Policy types, contains a set of `properties` and a set of `constraints`, but it contains other things as well. For example, it explicitly identifies the Service it will cause to be deployed onto Edge Nodes if negotiation is successful, in addition to configuration variable values, performing the equivalent function to the `-f horizon/userinput.json` clause of a Deployment Pattern `hzn register ...` command. The Business Policy approach for configuration values is more powerful because this operation can be performed centrally (no need to connect directly to the Edge Node).

- Below is the `horizon/business_policy.json` file used for this example:

```bash
{
  "label": "$SERVICE_NAME Business Policy for $ARCH",
  "description": "A Horizon Business Policy example to run cpu2evtstreams",
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
    "os == \"Mojave\"",
    "model == \"Mac\" OR model == \"Pi3B\""
  ],
  "userInput": []
    }
  ]
}
```
- This simple example of a Business Policy doesn't provide any `properties`, but it does have two `constraints` that are satisfied by the `properties` set in the `horizon/node_policy.json` file, so this Business Policy should successfully deploy our Service onto the Edge Node.

- At the bottom, the userInput section has the same purpose as the `horizon/userinput.json` files provided for other examples if the given services requires them. In this case the cpu2evtstreams service defines the configuration variables needed to send the data to IBM Event Streams. Though, for this example we have left the `userInput` section blank and will use the same `userinput.json` file we used before during pattern registration. 

1. To publish this Business Policy to the Exchange and get this Service running on the Edge Node edit the `horizon/business_policy.json` file to correctly identify your specific Service name, org, version, arch, etc. When your Business Policy is ready, run the following command to publish it, giving it a memorable name (cpu2evtstreamsPolicy in this example):

```bash
hzn exchange business addpolicy -f horizon/business_policy.json cpu2evtstreamsPolicy
```

2. Once that competes, you can look at the results with the following command, substituting your own org id:

```bash
hzn exchange business listpolicy major-peacock-icp-cluster/cpu2evtstreamsPolicy
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

5. See the cpu2evtstreams service output:

```bash
hzn service log -f ibm.cpu2evtstreams
```


6. Unregister your edge node, stopping the cpu2evtstreams service:
```bash
hzn unregister -f
```
