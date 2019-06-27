# Horizon Policy Example (Using the Hello World Example Edge Service)

## The Hello World Example Edge Service is a Prerequisite

- This example makes use of the Horizon "helloworld" example. If you have not already built and published that Horizon Service to your Exchange, you need to do that first by following the instructions in the [helloworld example](https://github.com/open-horizon/examples/tree/master/edge/services/helloworld).
- First, go through the "Try It" page "Installing Horizon Software On Your Edge Machine" to set up your edge node.

- You also need to have a Horizon Edge Node set up with Docker, and the Horizon software installed.

- On that node, the following command can be used to verify that the `hzn` CLI command and the Horizon Agent are functioning appropriately:

```
hzn version
```

- You will also need to configure your Horizon Exchange credentials in your working shell.

- You can verify your Horizon Exchange user credentials are correctly configured, by running this command:

```
hzn exchange user list
```

## An Introduction to Policy

- Most of the other examples provided in this repository make use of Horizon Deployment Patterns to register Edge Nodes. Deployment Patterns are often the simplest and easiest way to use Horizon.  

- The Horizon Policy mechanism offers an alternative to using Deployment Patterns. Policies provide much finer control over agreement forming between Horizon Agents on Edge Nodes, and the Horizon AgBots. It also provides a greater separation of concerns, allowing Edge Nodes owners, Service code developers, and Business owners to each independently articulate their own Policies. There are therefore three types of Horizon Policies:

1. Node Policy (provided at registration time by the node owner)

2. Service Policy (may be applied to a published Service in the Exchange)

3. Business Policy (which approximately corresponds to a Deployment Pattern)

- Each of these policy types is described in detail later in this document.

- Policies may define attributes (`properties`) as key/value pairs. They may also assert specific requirements (`constraints`) that must be satisfied by other applicable Policies before agreements can be completed. In general, using the Policy mechanism provides a superset of the capabilities in Deployment Patterns.

- Go ahead and clone the Horizon Examples repo, and change into the `hellopolicy` directory:

```
git clone https://github.com/open-horizon/examples.git
cd examples/policy/hellopolicy
```

- This example will show you how to orchestrate the deployment of the "helloworld" example onto your Edge Nodes using the Policy mechanism instead of the Deployment Pattern mechanism using the 3 policy files in this directory:

```
 $ ls
README.md  business_policy.json  node_policy.json  service_policy.json
 $ 
```

- Let's begin by looking at Node Policy...

## Node Policy

- As an alternative to specifying a Deployment Pattern when you register your Edge Node, you may register with a Node Policy.

- Like the other two Policy types, Node Policy contains a set of `properties` and a set of `constraints`. The `properties` of a Node Policy could state characteristics of the Edge Node, such as the product model, or attributes of the node hardware, attached devices, software configuration, or anything else deemed relevant.  The `constraints` of the Node Policy can be used to restrict what Services the Horizon Agent will permit to be run on this node.

- Make sure your Edge Node is not registered by running:

```
hzn unregister -f
```

- Now let's register using the `node_policy.json` example in this example. It contains:

```
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

- Register using this command (which requires the usual node authorization, and your organization):

```
hzn register -n $EXCHANGE_NODE_AUTH $HZN_ORG_ID --policy node_policy.json
```

- Note that this is very similar to the registration command you used previously for Deployment Patterns, except that here you explicitly state `--policy`.

- When the registration command completes, use the `hzn policy list` command to review the Node Policy:

```
 $ hzn policy  list
{
  "properties": [
    {
      "name": "model",
      "value": "Thingamajig ULTRA"
    },
    {
      "name": "serial",
      "value": 9123456
    },
    {
      "name": "configuration",
      "value": "Mark-II-PRO"
    },
    {
      "name": "openhorizon.cpu",
      "value": 8
    },
    {
      "name": "openhorizon.arch",
      "value": "amd64"
    },
    {
      "name": "openhorizon.memory",
      "value": 16039
    }
  ]
}
 $ 
```

- Notice that in addition to the three `properties` stated in the `node_policy.json` file, Horizon has added a few more (`openhorizon.cpu`, `openhorizon.arch`, and `openhorizon.memory`). Horizon provides this additional information automatically and these properties may be used in any of your Policy `constraints`.

- Now let's take a look at Service Policy...

## Service Policy

- Like the other two Policy types, Service Policy contains a set of `properties` and a set of `constraints`. The `properties` of a Service Policy could state characteristics of the Service code that Node Policy authors or Business Policy authors may find relevant.  The `constraints` of a Service Policy can be used to restrict where this Service can be run. The Service developer could, for example,  assert that this Service requires a particular hardware setup such as CPU/GPU constraints, memory constraints, specific sensors, actuators or other peripheral devices required, etc.

- This example directory contains an example Service Policy (in `service_policy.json` in this directory). It looks like this:

```
{
  "properties": [
  ],
  "constraints": [
    "model == \"Whatsit ULTRA\" OR model == \"Thingamajig ULTRA\""
  ]
}
```

- This simple Service Policy doesn't provide any `properties`, but it does have a `constraint`. This example constraint is one that a Service developer might add, stating that their Service must only run on the `models` named `Whatsit ULTRA` or `Thingamajig ULTRA`. If you recall the Node Policy we used above, the `model` property was set to `Thingamajig ULTRA`, so this Service should be compatible with our Edge Node.

- Note that the constraint language is capable of representing many complex conditions based on the values of the `properties` set in the other Policies.

- Now let's attach this Service Policy to the `helloworld` Service you previously published.

- Use the `hzn exchange service list` command to review the list of Services you have published and save the full name of your `hellowworld` Service.

```
 $ hzn exchange service list
[
  "major-peacock-icp-cluster/helloworld_1.0.0_amd64"
]
 $ 
```

- To attach the example Service policy to this service, use the following command (substituting your service name in place of the one I used here):

```
hzn exchange service addpolicy -f service_policy.json "major-peacock-icp-cluster/helloworld_1.0.0_amd64"
```

- Once that competes, you can look at the results with the `hzn exchange service listpolicy <service>` command, e.g.:

```
 $ hzn exchange service listpolicy "major-peacock-icp-cluster/helloworld_1.0.0_amd64"
{
  "properties": [
    {
      "name": "openhorizon.service.url",
      "value": "helloworld"
    },
    {
      "name": "openhorizon.service.name",
      "value": "helloworld"
    },
    {
      "name": "openhorizon.service.org",
      "value": "major-peacock-icp-cluster"
    },
    {
      "name": "openhorizon.service.version",
      "value": "1.0.0"
    },
    {
      "name": "openhorizon.service.arch",
      "value": "amd64"
    }
  ],
  "constraints": [
    "model == \"Whatsit ULTRA\" OR model == \"Thingamajig ULTRA\""
  ]
}
 $ 
```

- Notice that Horizon has again automatically added some additional `properties` to your Policy. These generated property values can be used in `constraints` in Node Policies and Business Policies.

- Now that we have set up the Policies for an Edge Node and the Policies for a published Service, we can move on to the final step of defining a **Business Policy** to tie them all together and cause software to be automatically deployed on your Edge Node.

## Business Policy
 
- Business Policy is what ties together Edge Nodes, and published Services, and the Policies defined for each of those, making it roughly analogous to the Deployment Patterns you have previously worked with.

- Business Policy, like the other two Policy types, contains a set of `properties` and a set of `constraints`, but it contains other things as well. For example, it explicitly identifies the Service it applies to, and which it will cause to be deployed onto Edge Nodes if negotiation is successful. It also contains configuration variable values, performing the equivalent function to the `-f horizon/userinput.json` clause of a Deployment Pattern `hzn register ...` command. The Business Policy approach for configuration values is more powerful because this operation can be performed centrally (no need to connect directly to the Edge Node).

- This example directory contains an example Business Policy (in `business_policy.json` in this directory). It looks like this:

```
{
  "service": {
    "name": "helloworld",
    "org": "major-peacock-icp-cluster",
    "arch": "*",
    "serviceVersions": [
      {
        "version": "1.0.0",
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
      "serviceOrgid": "major-peacock-icp-cluster",
      "serviceUrl": "helloworld",
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

- This simple Business Policy doesn't provide any `properties`, but it does have two `constraints`. These constraints require that any Edge Nodes must have a property 'serial` with a value greater than or equal to 9,000,000 *and* the `model` value must be exactly `Thingamajig ULTRA`. If you recall the `properties` set in the Node Policy above, it satisfies these constraints, so this Business Policy should successfully deploy our Service onto the Edge Node.

- Recall that in the Service Policy there was a single constraint, using an `OR` to test two conditions. Here two constraints are used to test two conditions. When multiple constraints are provided, they are always `AND`-ed together (i.e., *all* of the constraints must be satisfied before an agreement can be formed).

- This example Business Policy identifies the Service to which it applies. It also itemizes the architectures and versions of this Service that it applies to.

- At the bottom, the `userInput` section has the same purpose as the `horizon/userinput.json` files provided for other examples (which are needed when using Deployment Patterns instead of Policy). That is, this section provides values for the configuration variables exposed by the Service specified by the Business Policy, or any of the other transitively `requiredServices` it includes. In this case there is only one Service (`helloworld`) and that Service defines only one configuration variable, `HW_WHO`. So correspondingly the `userInput` section here provides a value for `HW_WHO` (i.e., `Valued Customer`). Later we will be able to see this value in the logs of this Service when it is run on the Edge Node through these Policies.

- Now let's publish this Business Policy to the Exchange to start the ball rolling and get this Service running on the Edge Node. To do this, edit the `business_policy.json` file to correctly identify your specific Service name, org, version, arch, etc.  When your Business Policy is ready, run the following command to publish it, giving it a memorable name (`biz1` in this example:

```
hzn exchange business addpolicy -f ~/businesspolicy.json biz1
```

- Once that competes, you can look at the results with the `hzn exchange business listpolicy` command, e.g.:

```
hzn exchange business listpolicy "major-peacock-icp-cluster/biz1"
{
  "major-peacock-icp-cluster/biz1": {
    "owner": "major-peacock-icp-cluster/user1",
    "label": "",
    "description": "",
    "service": {
      "name": "helloworld",
      "org": "major-peacock-icp-cluster",
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
      "serial >= 9000000",
      "model == \"Thingamajig ULTRA\""
    ],
    "userInput": [
      {
        "serviceOrgid": "major-peacock-icp-cluster",
        "serviceUrl": "helloworld",
        "serviceVersionRange": "[0.0.0,INFINITY)",
        "inputs": [
          {
            "name": "HW_WHO",
            "value": "Valued Customer"
          }
        ]
      }
    ],
    "created": "2019-06-26T19:50:49.542Z[UTC]",
    "lastUpdated": "2019-06-26T21:28:04.273Z[UTC]"
  }
}
```

- The results above should look very similar to your original `business_policy.json` file, except that `owner`, `created` and `lastUpdated` and a few other fields have been added.

- Now let's go back to the Edge Node and see what's happening there...

## Results

- Return to a terminal on the Edge Node and see if an Agreement is in force for your `helloworld` Service, as a result of your work with the Policies:

```
hzn agreement list
```

- An agreement should appear there within a matter of seconds, or at most a couple of minutes.

- It will typically take less than a few minutes for an agreement to be created, accepted and finalized, then execution of the corresponding Docker container should begin.

- Check that the `helloworld` container is running:

```
 $ docker ps
CONTAINER ID        IMAGE                         COMMAND                  CREATED             STATUS              PORTS               NAMES
030e0c497dcb        ibmosquito/helloworld_amd64   "/bin/sh -c /serviceâ¦"   38 seconds ago      Up 36 seconds                           1173d19d419657fd7927ff3dc22c81bfd673139184bea31bc3b7b3826a18d9cb-helloworld
 $ 
```

- If everything has gone as expected, and the agreement is finalized, and your container is executing, then check the container output to see that the `userInput` configuration was applied.

- These commands will show the `helloworld` Service output on various platforms:

```
# soon you will use 'hzn service log ...' for all platforms
# For now on Linux:
tail -f /var/log/syslog | grep helloworld[[]
# For now on Mac:
docker logs -f $(docker ps -q --filter name=helloworld)
``` 

- You should see output similar to the lines below, showing `Hello Valued Customer`:

```
Jun 26 14:49:03 glen-dev1 workload-1173d19d419657fd7927ff3dc22c81bfd673139184bea31bc3b7b3826a18d9cb_helloworld[865]: major-peacock-icp-cluster says: Hello Valued Customer!!
Jun 26 14:49:09 glen-dev1 workload-1173d19d419657fd7927ff3dc22c81bfd673139184bea31bc3b7b3826a18d9cb_helloworld[865]: message repeated 2 times: [ major-peacock-icp-cluster says: Hello Valued Customer!!]
Jun 26 14:49:12 glen-dev1 workload-1173d19d419657fd7927ff3dc22c81bfd673139184bea31bc3b7b3826a18d9cb_helloworld[865]: major-peacock-icp-cluster says: Hello Valued Customer!!
```

- Before moving on to others examples, unregister your Edge Node, stopping the `helloworld` Service:

```
hzn unregister -f
```


