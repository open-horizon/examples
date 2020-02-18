# Using the CPU To IBM Event Streams Service with Deployment Policy

Follow the steps on this page to deploy CPU To IBM Event Streams Edge Service using deployment policy.

1. If you have not already done so, complete the steps in these sections:

  - [Preconditions for Using the CPU To IBM Event Streams Example Edge Service](README.md#preconditions)

## <a id=using-cpu2evtstreams-policy></a> Using the CPU To IBM Event Streams Service with Deployment Policy

- The Horizon Policy mechanism offers an alternative to using Deployment Patterns. Policies provide much finer control over agreement forming between Horizon Agents on Edge Nodes, and the Horizon AgBots. It also provides a greater separation of concerns, allowing Edge Nodes owners, Service code developers, and Business owners to each independently articulate their own Policies. There are therefore three types of Horizon Policies:

1. Node Policy (provided at registration time by the node owner)

2. Service Policy (may be applied to a published Service in the Exchange)

3. Business Policy (which approximately corresponds to a Deployment Pattern)

### Node Policy 

- As an alternative to specifying a Deployment Pattern when you register your Edge Node, you may register with a Node Policy.

1. Get the required `cpu2evtstreams` node and business policy files and the user input file:
```bash
wget https://raw.githubusercontent.com/open-horizon/examples/master/edge/evtstreams/cpu2evtstreams/horizon/node_policy.json
wget https://raw.githubusercontent.com/open-horizon/examples/master/edge/evtstreams/cpu2evtstreams/horizon/business_policy.json
wget https://raw.githubusercontent.com/open-horizon/examples/master/edge/evtstreams/cpu2evtstreams/horizon/use/userinput.json
```

- Below is the `node_policy.json` file:

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

2. Register your Node Policy using this command:

```bash
hzn register --policy node_policy.json -f userinput.json
```

3. When the registration completes, use the following command to review the Node Policy:

```bash
hzn policy list
```

- Notice that in addition to the three `properties` stated in the node_policy.json file, Horizon has added a few more (openhorizon.cpu, openhorizon.arch, and openhorizon.memory). Horizon provides this additional information automatically and these `properties` may be used in any of your Policy `constraints`.

### Service Policy 

- Like the other two Policy types, Service Policy contains a set of `properties` and a set of `constraints`. The `properties` of a Service Policy could state characteristics of the Service code that Node Policy authors or Business Policy authors may find relevant. The `constraints` of a Service Policy can be used to restrict where this Service can be run. The Service developer could, for example, assert that this Service requires a particular hardware setup such as CPU/GPU constraints, memory constraints, specific sensors, actuators or other peripheral devices required, etc.

- Below is the Service Policy for `ibm.cpu2evtstreams` published by the cpu2evtstreams service developer that is in the exchange:

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

1. View the pubished service policy for the `amd64` architecture:
```bash
hzn exchange service listpolicy IBM/ibm.cpu2evtstreams_1.4.3_amd64
```
- Notice that Horizon has again automatically added some additional `properties` to your Policy. These generated property values can be used in `constraints` in Node Policies and Business Policies.

- Now that we have set up the Policies for an Edge Node and the Policies for a published Service, we can move on to the final step of defining a Business Policy to tie them all together and cause software to be automatically deployed on your Edge Node.

### Business Policy 

- Business Policy (sometimes called Deployment Policy) is what ties together Edge Nodes, Published Services, and the Policies defined for each of those, making it roughly analogous to the Deployment Patterns you have previously worked with.

- Business Policy, like the other two Policy types, contains a set of `properties` and a set of `constraints`, but it contains other things as well. For example, it explicitly identifies the Service it will cause to be deployed onto Edge Nodes if negotiation is successful, in addition to configuration variable values, performing the equivalent function to the `-f horizon/userinput.json` clause of a Deployment Pattern `hzn register ...` command. The Business Policy approach for configuration values is more powerful because this operation can be performed centrally (no need to connect directly to the Edge Node).

- Below is the `business_policy.json` file you grabbed in step one:

```bash
{
  "label": "ibm.cpu2evtstreams Business Policy for $ARCH",
  "description": "A Horizon Business Policy example to run cpu2evtstreams",
  "service": {
    "name": "ibm.cpu2evtstreams",
    "org": "IBM",
    "arch": "$ARCH",
    "serviceVersions": [
      {
        "version": "1.4.3",
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

1. Set the architecture of the Edge Node you want to deploy the `ibm.cpu2evtstreams` service to `(arm, arm64, amd64)`:
```bash
export ARCH=<edge-node-arch>
```

2. Publish this Business Policy to the Exchange and get this Service running on the Edge Node and give it a memorable name:

```bash
hzn exchange business addpolicy -f business_policy.json <choose-any-policy-name>
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
 - **Note**: Press **Ctrl C** to stop the command output.

6. Unregister your edge node, stopping the cpu2evtstreams service:
```bash
hzn unregister -f
```
