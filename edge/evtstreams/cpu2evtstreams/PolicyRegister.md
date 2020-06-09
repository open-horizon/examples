# Using the CPU To IBM Event Streams Service with Deployment Policy

Follow the steps on this page to deploy CPU To IBM Event Streams Edge Service using deployment policy.

1. If you have not already done so, complete the steps in this section:

  - [Preconditions for Using the CPU To IBM Event Streams Example Edge Service](README.md#preconditions)

## <a id=using-cpu2evtstreams-policy></a> Using the CPU To IBM Event Streams Service with Deployment Policy

- The Horizon Policy mechanism offers an alternative to using Deployment Patterns. Policies provide much finer control over agreement forming between Horizon Agents on Edge Nodes, and the Horizon AgBots. It also provides a greater separation of concerns, allowing Edge Nodes owners, Service code developers, and Business owners to each independently articulate their own Policies. There are therefore three types of Horizon Policies:

1. Node Policy (provided at registration time by the node owner)

2. Service Policy (may be applied to a published Service in the Exchange)

3. Deployment Policy (which approximately corresponds to a Deployment Pattern)

### Node Policy 

- As an alternative to specifying a Deployment Pattern when you register your Edge Node, you may register with a Node Policy.

1. Get the required `cpu2evtstreams` node policy file and the use input file:
```bash
wget https://raw.githubusercontent.com/open-horizon/examples/master/edge/evtstreams/cpu2evtstreams/horizon/node.policy.json
wget https://raw.githubusercontent.com/open-horizon/examples/master/edge/evtstreams/cpu2evtstreams/horizon/use/userinput.json
```

- Below is the `node_policy.json` file you just grabbed in step one:

```bash
{
    "properties": [
        {
            "name": "openhorizon.example", "value": "cpu2evtstreams"
        }
    ],
    "constraints": []
}
```

- It provides values for one just properties (`openhorizon.example`), that should restrict the services that get deployed to this edge node to the `cpu2evtstreams` service.

2. Register your Node Policy using this command:

```bash
hzn register --policy node.policy.json -f userinput.json
```

3. When the registration completes, use the following command to review the Node Policy:

```bash
hzn policy list
```

- Notice that in addition to the `properties` stated in the `node.policy.json` file, Horizon has added a few more (openhorizon.cpu, openhorizon.arch, and openhorizon.memory). Horizon provides this additional information automatically and these `properties` may be used in any of your Policy `constraints`.

### Service Policy 

- Like the other two Policy types, Service Policy contains a set of `properties` and a set of `constraints`. The `properties` of a Service Policy could state characteristics of the Service code that Node Policy authors or Deployment Policy authors may find relevant. The `constraints` of a Service Policy can be used to restrict where this Service can be run. The Service developer could, for example, assert that this Service requires a particular hardware setup such as CPU/GPU constraints, memory constraints, specific sensors, actuators or other peripheral devices required, etc.

- Below is the `service.policy.json` file the service developer attached to the `ibm.cpu2evtstreams` service when it was published to the exchange:

```bash
{
    "properties": [],
    "constraints": [
        "openhorizon.memory >= 100"
    ]
}
```

- Note this simple Service Policy doesn't provide any `properties`, but it does have a `constraint`. This example constraint is one that a Service developer might add, stating that their Service must only run on edge nodes with 100 MB or more memory.

1. View the pubished service policy attached to `ibm.cpu2evtstreams`:
```bash
hzn exchange service listpolicy IBM/ibm.cpu2evtstreams_1.4.3_amd64
```
- Notice that Horizon has again automatically added some additional `properties` to your Policy. These generated property values can be used in `constraints` in Node Policies and Deployment Policies.

- Now that we have set up the Policies for an Edge Node and the Policies for a published Service, we can move on to the final step of defining a Deployment Policy to tie them all together and cause software to be automatically deployed on your Edge Node.

### Deployment Policy 

- Deployment Policy (sometimes called Deployment Policy) is what ties together Edge Nodes, Published Services, and the Policies defined for each of those, making it roughly analogous to the Deployment Patterns you have previously worked with.

- Deployment Policy, like the other two Policy types, contains a set of `properties` and a set of `constraints`, but it contains other things as well. For example, it explicitly identifies the Service it will cause to be deployed onto Edge Nodes if negotiation is successful, in addition to configuration variable values, performing the equivalent function to the `-f horizon/userinput.json` clause of a Deployment Pattern `hzn register ...` command. The Deployment Policy approach for configuration values is more powerful because this operation can be performed centrally (no need to connect directly to the Edge Node).

- Below is the `deployment.policy.json` for `cpu2evtstreams` that was published to the Exchange when it was created:

```bash
{
  "label": "$SERVICE_NAME Deployment Policy",
  "description": "A Horizon Deployment Policy example to run cpu2evtstreams",
  "service": {
    "name": "$SERVICE_NAME",
    "org": "IBM",
    "arch": "*",
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
    "openhorizon.example == cpu2evtstreams"
  ],
  "userInput": []
    }
  ]
}
```
- This simple example of a Deployment Policy doesn't provide any `properties`, but it does have one `constraints` value that are satisfied by the `properties` set in the `horizon/node.policy.json` file, so this Deployment Policy should successfully deploy our Service onto the Edge Node.

2. Run the following commands to set the environment variables needed by the `deployment.policy.json` file in your shell:
```bash
export ARCH=$(hzn architecture)
eval $(hzn util configconv -f hzn.json)
eval export $(cat agent-install.cfg)
```

3. Publish this Deployment Policy to the Exchange and get this Service running on the Edge Node and give it a memorable name:

```bash
hzn exchange deployment addpolicy -f deployment.policy.json <choose-any-policy-name>
```

- The results should look very similar to your original `horizon/deployment.policy.json` file, except that `owner`, `created`, and `lastUpdated` and a few other fields have been added.


4. The edge device will make an agreement with one of the Horizon agreement bots (this typically takes about 15 seconds). Repeatedly query the agreements of this device until the `agreement_finalized_time` and `agreement_execution_start_time` fields are filled in:
```bash
hzn agreement list
```

5. Once the agreement is made, list the docker container edge service that has been started as a result:
```bash
sudo docker ps
```

6. On any machine, install [kafkacat](https://github.com/edenhill/kafkacat#install), then subscribe to the Event Streams topic to see the json data that `cpu2evtstreams` is sending:
  ```bash
  kafkacat -C -q -o end -f "%t/%p/%o/%k: %s\n" -b $EVTSTREAMS_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=token -X sasl.password=$EVTSTREAMS_API_KEY -X ssl.ca.location=$EVTSTREAMS_CERT_FILE -t cpu2evtstreams
  ```
 - **Note**: Press **Ctrl C** to stop the command output.
  
7. See the cpu2evtstreams service output:

```bash
hzn service log -f ibm.cpu2evtstreams
```
 - **Note**: Press **Ctrl C** to stop the command output.

8. Unregister your edge node, stopping the cpu2evtstreams service:
```bash
hzn unregister -f
```
