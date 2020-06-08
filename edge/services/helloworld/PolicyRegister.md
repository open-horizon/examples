# Using the Hello World Example Edge Service with Deployment Policy

## Preconditions for Developing Your Own Service

1. If you have not already done so, complete the steps in this section:

  - [Preconditions for Using the Hello World Example Edge Service](README.md#preconditions)
  
## <a id=using-helloworld-policy></a> Using the Hello World Example Edge Service with Deployment Policy

The Horizon Policy mechanism offers an alternative to using Deployment Patterns. Policies provide much finer control over the deployment placement of edge services. It also provides a greater separation of concerns, allowing Edge Nodes owners, Service code developers, and Business owners to each independently articulate their own Policies. There are therefore three types of Horizon Policies:

1. Node Policy (provided at registration time by the node owner)

2. Service Policy (may be applied to a published Service in the Exchange)

3. Deployment Policy (which approximately corresponds to a Deployment Pattern)

### Node Policy

- As an alternative to specifying a Deployment Pattern when you register your Edge Node, you may register with a Node Policy.

1. Get the required helloworld node policy file:
```bash
wget https://raw.githubusercontent.com/open-horizon/examples/master/edge/services/helloworld/horizon/node.policy.json
```

- Below is the `node.policy.json` file you obtained in step one:

```json
{
  "properties": [
    { "name": "openhorizon.example", "value": "helloworld" }
  ],
  "constraints": [
  ]
}
```

- It provides values for three `properties` (`model`, `serial`, and `configuration`), that will effect which services get deployed to this edge node, and states no `constraints`.

2. Register your Node Policy with this policy

```bash
hzn register --policy node.policy.json
```

3. When the registration completes, use the following command to review the Node Policy:

```bash
hzn policy list
```

- Notice that in addition to the three `properties` stated in the node_policy.json file, Horizon has added a few more (openhorizon.cpu, openhorizon.arch, and openhorizon.memory). Horizon provides this additional information automatically and these `properties` may be used in any of your Policy `constraints`.

### Service Policy

- Like the other two Policy types, Service Policy contains a set of `properties` and a set of `constraints`. The `properties` of a Service Policy could state characteristics of the Service code that Node Policy authors or Deployment Policy authors may find relevant. The `constraints` of a Service Policy can be used to restrict where this Service can be run. The Service developer could, for example, assert that this Service requires a particular hardware setup such as CPU/GPU constraints, memory constraints, specific sensors, actuators or other peripheral devices required, etc.

- Below is the `service.policy.json` file the service developer attached to `ibm.helloworld` when it was published:

```json
{
  "properties": [
  ],
  "constraints": [
    "openhorizon.memory >= 100"
  ]
}
```

- Note this simple Service Policy doesn't provide any `properties`, but it does have a `constraint`. This example `constraint` is one that a Service developer might add, stating that their Service must only run on the models named `Whatsit ULTRA` or `Thingamajig ULTRA`. If you recall the Node Policy we used above, the model `property` was set to `Thingamajig ULTRA`, so this Service should be compatible with our Edge Node.

1. View the pubished service policy attached to `ibm.helloworld`:

```bash
hzn exchange service listpolicy IBM/ibm.helloworld_1.0.0_amd64
```

- Notice that Horizon has again automatically added some additional `properties` to your Policy. These generated property values can be used in `constraints` in Node Policies and Deployment Policies.

- Now that you have set up the Policy for your Edge Node and the published Service policy is in the exchange, we can move on to the final step of defining a Deployment Policy to tie them all together and cause software to be automatically deployed on your Edge Node.

### Deployment Policy

- Deployment Policy is what ties together Edge Nodes, Published Services, and the Policies defined for each of those, making it roughly analogous to the Deployment Patterns you have previously worked with.

- Deployment Policy, like the other two Policy types, contains a set of `properties` and a set of `constraints`, but it contains other things as well. For example, it explicitly identifies the Service it will cause to be deployed onto Edge Nodes if negotiation is successful, in addition to configuration variable values, performing the equivalent function to the `-f horizon/userinput.json` clause of a Deployment Pattern `hzn register ...` command. The Deployment Policy approach for configuration values is more powerful because this operation can be performed centrally (no need to connect directly to the Edge Node).

1. Get the required `helloworld` deployment policy file and the `hzn.json` file:
```bash
wget https://raw.githubusercontent.com/open-horizon/examples/master/edge/services/helloworld/horizon/deployment.policy.json
wget https://raw.githubusercontent.com/open-horizon/examples/master/edge/services/helloworld/horizon/hzn.json
```
- Below is the `deployment,policy.json` file you just grabbed in step one:

```json
{
  "label": "$SERVICE_NAME Deployment Policy",
  "description": "A super-simple sample Horizon Deployment Policy",
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
    "openhorizon.example == helloworld"
  ],
  "userInput": [
    {
      "serviceOrgid": "IBM",
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

- This simple example of a Deployment Policy doesn't provide any `properties`, but it does have two `constraints` that are satisfied by the `properties` set in the `node.policy.json` file, so this Deployment Policy should successfully deploy our Service onto the Edge Node.

- At the end, the userInput section has the same purpose as the `horizon/userinput.json` files provided for other examples if the given services requires them. In this case the helloworld service defines only one configuration variable, HW_WHO, and the userInput section here provides a value for HW_WHO (i.e., Valued Customer).

2. Run the following commands to set the environment variables needed by the `deployment.policy.json` file in your shell:
```bash
export ARCH=$(hzn architecture)
eval $(hzn util configconv -f hzn.json)
eval export $(cat agent-install.cfg)
```

3. Publish this Deployment Policy to the Exchange to deploy the `ibm.helloworld` service to the Edge Node (give it a memorable name):

```bash
hzn exchange deployment addpolicy -f deployment.policy.json <choose-any-policy-name>
```

- The results should look very similar to your original `deployment.policy.json` file, except that `owner`, `created`, and `lastUpdated` and a few other fields have been added.

4. The edge device will make an agreement with one of the Horizon agreement bots (this typically takes about 15 seconds). Repeatedly query the agreements of this device until the `agreement_finalized_time` and `agreement_execution_start_time` fields are filled in:

```bash
hzn agreement list
```

5. After the agreement is made, list the edge service docker container that has been started as a result:

```bash
sudo docker ps
```

6. See the `ibm.helloworld` service output:

``` bash
hzn service log -f ibm.helloworld
```
 - **Note**: Press **Ctrl C** to stop the command output.

7. Unregister your edge node:

```bash
hzn unregister -f
```
