# Horizon CPU To IBM Event Streams Service

For details about using this service, see [cpu2evtstreams.md](cpu2evtstreams.md).

## Using the CPU To IBM Event Streams Edge Service

- Before following the steps in this section, install the Horizon agent on your edge device and point it to your Horizon exchange. Also get an API key that is associated with your Horizon instance.
1. Set your exchange org:
```
export HZN_ORG_ID="<yourorg>"
```
2. Set your exchange user credentials in the Horizon-supported environment variable and verify it:
```
export HZN_EXCHANGE_USER_AUTH="iamapikey:<myapikey>"
hzn exchange user list
```
3. Choose a id and token for your edge node, create it, and verify it:
```
export HZN_EXCHANGE_NODE_AUTH="<mynodeid>:<mynodetoken>"
hzn exchange node create -n $HZN_EXCHANGE_NODE_AUTH
hzn exchange node confirm
```
4. Deploy (or get access to) an instance of IBM Event Streams that the cpu2evtstreams sample can send its data to. Ensure that the topic `cpu2evtstreams` is created in Event Streams. Using information from the Event Streams UI, `export` these environment variables:
    - `EVTSTREAMS_API_KEY`
    - `EVTSTREAMS_BROKER_URL`
    - `EVTSTREAMS_CERT_ENCODED` (if using IBM Event Streams in IBM Cloud Private) due to differences in the base64 command set this variable as follows based on the platform you're using:
        - On Linux: `EVTSTREAMS_CERT_ENCODED=“$(cat $EVTSTREAMS_CERT_FILE| base64 -w 0)”`
        - On Mac: `EVTSTREAMS_CERT_ENCODED=“$(cat $EVTSTREAMS_CERT_FILE| base64)”`
    - `EVTSTREAMS_CERT_FILE` (if using IBM Event Streams in IBM Cloud Private)
5. Get the user input file for the cpu2evtstreams sample:
```
wget https://github.com/open-horizon/examples/raw/master/edge/evtstreams/cpu2evtstreams/horizon/use/userinput.json
```
6. Register your edge node with Horizon to use the cpu2evtstreams pattern:
```
hzn register -p IBM/pattern-ibm.cpu2evtstreams -f userinput.json
```
7. Look at the Horizon agreement until it is finalized and then see the running container:
```
hzn agreement list
docker ps
```
8. On any machine, install [kafkacat](https://github.com/edenhill/kafkacat#install), then subscribe to the Event Streams topic to see the json data that cpu2evtstreams is sending:
  - If using IBM Event Streams in IBM Cloud:
  ```
  kafkacat -C -q -o end -f "%t/%p/%o/%k: %s\n" -b $EVTSTREAMS_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=${EVTSTREAMS_API_KEY:0:16} -X sasl.password=${EVTSTREAMS_API_KEY:16} -t $EVTSTREAMS_TOPIC
  ```
  - If using IBM Event Streams in IBM Cloud Private:
  ```
  kafkacat -C -q -o end -f "%t/%p/%o/%k: %s\n" -b $EVTSTREAMS_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=token -X sasl.password=$EVTSTREAMS_API_KEY -X ssl.ca.location=$EVTSTREAMS_CERT_FILE -t cpu2evtstreams
  ```
9. (Optional) To see the cpu2evtstreams service output:
```
# On Linux:
tail -f /var/log/syslog | grep cpu2evtstreams[[]
# On Mac:
docker logs -f $(docker ps -q --filter name=cpu2evtstreams)
``` 
10. Unregister your edge node, stopping the cpu2evtstreams service:
```
hzn unregister -f
```

## Using the CPU To IBM Event Streams Service as a Policy

- The Horizon Policy mechanism offers an alternative to using Deployment Patterns. Policies provide much finer control over agreement forming between Horizon Agents on Edge Nodes, and the Horizon AgBots. It also provides a greater separation of concerns, allowing Edge Nodes owners, Service code developers, and Business owners to each independently articulate their own Policies. There are therefore three types of Horizon Policies:

1. Node Policy (provided at registration time by the node owner)

2. Service Policy (may be applied to a published Service in the Exchange)

3. Business Policy (which approximately corresponds to a Deployment Pattern)

### Node Policy 

- As an alternative to specifying a Deployment Pattern when you register your Edge Node, you may register with a Node Policy.

1. Make sure your Edge Node is not registered by running:

```
hzn unregister -f
```

- Now let's register using the `horizon/node_policy.json` file:

```
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

```
hzn register --policy horizon/node_policy.json
```

3. When the registration completes, use the following command to review the Node Policy:

```
hzn policy list
```

- Notice that in addition to the three `properties` stated in the node_policy.json file, Horizon has added a few more (openhorizon.cpu, openhorizon.arch, and openhorizon.memory). Horizon provides this additional information automatically and these `properties` may be used in any of your Policy `constraints`.

### Service Policy 

- Like the other two Policy types, Service Policy contains a set of `properties` and a set of `constraints`. The `properties` of a Service Policy could state characteristics of the Service code that Node Policy authors or Business Policy authors may find relevant. The `constraints` of a Service Policy can be used to restrict where this Service can be run. The Service developer could, for example, assert that this Service requires a particular hardware setup such as CPU/GPU constraints, memory constraints, specific sensors, actuators or other peripheral devices required, etc.

- Now let's attach this Service Policy to the cpu2evtstreams Service previously published using the `horizon/service_policy.json` file:

```
{
    "properties": [],
    "constraints": [
        "model == \"Mac\" OR model == \"Pi3B\"",
        "os == \"Mojave\""
    ]
}
```

- Note this simple Service Policy doesn't provide any `properties`, but it does have a `constraint`. This example `constraint` is one that a Service developer might add, stating that their Service must only run on the models named `Mac` or `Pi3B `. If you recall the Node Policy we used above, the model `property` was set to `Mac`, so this Service should be compatible with our Edge Node.

1. To attach the example Service policy to this service, use the following command (substituting your service name):

```
hzn exchange service addpolicy -f horizon/service_policy.json <published-cpu2evtstreams-service-name>
```

2. Once that completes, you can look at the results with the following command:

```
hzn exchange service listpolicy <published-cpu2evtstreams-service-name>
```
- Notice that Horizon has again automatically added some additional `properties` to your Policy. These generated property values can be used in `constraints` in Node Policies and Business Policies.

- Now that we have set up the Policies for an Edge Node and the Policies for a published Service, we can move on to the final step of defining a Business Policy to tie them all together and cause software to be automatically deployed on your Edge Node.

### Business Policy 

- Business Policy is what ties together Edge Nodes, Published Services, and the Policies defined for each of those, making it roughly analogous to the Deployment Patterns you have previously worked with.

- Business Policy, like the other two Policy types, contains a set of `properties` and a set of `constraints`, but it contains other things as well. For example, it explicitly identifies the Service it will cause to be deployed onto Edge Nodes if negotiation is successful, in addition to configuration variable values, performing the equivalent function to the `-f horizon/userinput.json` clause of a Deployment Pattern `hzn register ...` command. The Business Policy approach for configuration values is more powerful because this operation can be performed centrally (no need to connect directly to the Edge Node).

- Below is the `horizon/business_policy.json` file used for this example:

```
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
  "userInput": [
    {
      "serviceOrgid": "$HZN_ORG_ID",
      "serviceUrl": "$SERVICE_NAME",
      "serviceVersionRange": "[0.0.0,INFINITY)",
      "inputs": [
        {
          "name": "EVTSTREAMS_API_KEY",
          "value": "$EVTSTREAMS_API_KEY"
        },
        {
          "name": "EVTSTREAMS_BROKER_URL",
          "value": "$EVTSTREAMS_BROKER_URL"
        },
        {
          "name": "EVTSTREAMS_CERT_ENCODED",
          "value": "$EVTSTREAMS_CERT_ENCODED"
        }
      ]
    }
  ]
}
```
- This simple example of a Business Policy doesn't provide any `properties`, but it does have two `constraints` that are satisfied by the `properties` set in the `horizon/node_policy.json` file, so this Business Policy should successfully deploy our Service onto the Edge Node.

- At the bottom, the userInput section has the same purpose as the horizon/userinput.json files provided for other examples if the given services requires them. In this case the cpu2evtstreams service defines the configuration variables needed to send the data to IBM Event Streams, which will by default be taken from the environment variables themselves.

1. To publish this Business Policy to the Exchange and get this Service running on the Edge Node edit the `horizon/business_policy.json` file to correctly identify your specific Service name, org, version, arch, etc. When your Business Policy is ready, run the following command to publish it, giving it a memorable name (cpu2evtstreamsPolicy in this example):

```
hzn exchange business addpolicy -f horizon/business_policy.json cpu2evtstreamsPolicy
```

2. Once that competes, you can look at the results with the following command, substituting your own org id:

```
hzn exchange business listpolicy major-peacock-icp-cluster/cpu2evtstreamsPolicy
```

- The results should look very similar to your original `horizon/business_policy.json` file, except that `owner`, `created`, and `lastUpdated` and a few other fields have been added.

3. Look at the Horizon agreement until it is finalized and then see the running container:
```
hzn agreement list
docker ps
```

4. See the cpu2evtstreams service output:
```
# soon you will use 'hzn service log ...' for all platforms
# For now on Linux:
tail -f /var/log/syslog | grep cpu2evtstreams[[]
# For now on Mac:
docker logs -f $(docker ps -q --filter name=cpu2evtstreams)
```
5. Unregister your edge node, stopping the cpu2evtstreams service:
```
hzn unregister -f
```



## First-Time Edge Service Developer - Building and Publishing Your Own Version of the CPU To IBM Event Streams Edge Service

If you want to create your own Horizon edge service, based on this example, follow the next 2 sections to copy the cpu2evtstreams example and start modifying it.

### Preconditions for Developing Your Own Service

1. First, go through the steps in the section above to run the IBM cpu2evtstreams service on an edge node.
2. Get a docker hub id at https://hub.docker.com/ , if you don't already have one. (This example is set up to store the docker image in docker hub, but by modifying DOCKER_IMAGE_BASE you can store it in another registry.) Login to the docker registry using your id:
```
echo 'mydockerpw' | docker login -u mydockehubid --password-stdin
```
3. If you have the HZN_ORG_ID environment variable set from previous work, unset it (in a moment this value will now come from `horizon/hzn.json`):
```
unset HZN_ORG_ID
```
4. Clone this git repo:
```
cd ~   # or wherever you want
git clone git@github.com:open-horizon/examples.git
```
5. Copy the `cpu2evtstreams` dir to where you will start development of your new service:
```
cp -a examples/edge/evtstreams/cpu2evtstreams ~/myservice     # or wherever
cd ~/myservice
```
6. Set the values in `horizon/hzn.json` to your own values.
7. As part of the above section "Using the CPU To IBM Event Streams Edge Service", you created your Exchange user credentials and edge node credentials. Ensure they are set and verify them:
```
export HZN_EXCHANGE_USER_AUTH="iamapikey:<myapikey>"
hzn exchange user list
export HZN_EXCHANGE_NODE_AUTH="<mynodeid>:<mynodetoken>"
hzn exchange node confirm
```
8. Verify that these environment variables are still set from when you used the existing cpu2evtstreams sample earlier in this document:
```
echo EVTSTREAMS_API_KEY=$EVTSTREAMS_API_KEY
echo EVTSTREAMS_ADMIN_URL=$EVTSTREAMS_ADMIN_URL
echo EVTSTREAMS_BROKER_URL=$EVTSTREAMS_BROKER_URL
```
9. Verify the `cpu2evtstreams ` topic is now in your event streams instance:
```
make evtstreams-topic-list
```

### Building and Publishing Your Own Version of the CPU To IBM Event Streams Edge Service

1. Edit `service.sh` however you want.
    - Note: this service is a shell script simply for brevity, but you can write your service in any language.
2. Build the cpu2evtstreams docker image:
```
make
```
3. Test the service by having Horizon start it locally:
```
hzn dev service start -S
```
4. See the docker container running and look at the output:
```
docker ps
# soon you will use 'hzn service log ...' for all platforms
# For now on Linux:
tail -f /var/log/syslog | grep cpu2evtstreams[[]
# For now on Mac:
docker logs -f $(docker ps -q --filter name=cpu2evtstreams)
```
5. See the environment variables Horizon passes into your service container:
```
docker inspect $(docker ps -q --filter name=cpu2evtstreams) | jq '.[0].Config.Env'
```
6. Stop the service:
```
hzn dev service stop
```
7. Create a service signing key pair in `~/.hzn/keys/` (if you haven't already done so):
```
hzn key create <my-company> <my-email>
```
8. Have Horizon push your docker image to your registry and publish your service in the Horizon Exchange and see it there:
```
hzn exchange service publish -f horizon/service.definition.json
hzn exchange service list
```
9. Publish your edge node deployment pattern in the Horizon Exchange and see it there:
```
hzn exchange pattern publish -f horizon/pattern.json
hzn exchange pattern list
```
10. Register your edge node with Horizon to use your deployment pattern (substitute for `SERVICE_NAME` the value you specified above for `hzn dev service new -s`):
```
hzn register -p pattern-SERVICE_NAME-$(hzn architecture) -f horizon/userinput.json
```
11. Look at the Horizon agreement until it is finalized and then see the running container:
```
hzn agreement list
docker ps
```
12. On any machine, subscribe to the Event Streams topic to see the json data that cpu2evtstreams is sending:
```
kafkacat -C -q -o end -f "%t/%p/%o/%k: %s\n" -b $EVTSTREAMS_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=${EVTSTREAMS_API_KEY:0:16} -X sasl.password=${EVTSTREAMS_API_KEY:16} -t $EVTSTREAMS_TOPIC
```
13. See the cpu2evtstreams service output:
```
# soon you will use 'hzn service log ...' for all platforms
# For now on Linux:
tail -f /var/log/syslog | grep cpu2evtstreams[[]
# For now on Mac:
docker logs -f $(docker ps -q --filter name=cpu2evtstreams)
``` 
14. Unregister your edge node, stopping the cpu2evtstreams service:
```
hzn unregister -f
```

## Process for the Horizon Development Team to Make Updates to the Cpu2evtstreams Service

- Do the steps in the Preconditions section above, **except**:
    - export `HZN_EXCHANGE_URL` to the staging instance
    - Do **not** copy the cpu2evtstreams directory (use the git files in this directory instead)
    - export `HZN_EXCHANGE_USER_AUTH` to your credentials in the IBM org
- Make whatever code changes are necessary
- Increment `SERVICE_VERSION` in `horizon/hzn.json`
- Make `~/.hzn/keys/service.private.key` and `~/.hzn/keys/service.public.pem` actually be symbolic links to the common keys we use to sign all of our examples.
- Build, test, and publish for all architectures:
```
make publish-all-arches
```
Note: building all architectures works on mac os x, and can be made to work on ubuntu via: http://wiki.micromint.com/index.php/Debian_ARM_Cross-compile , https://wiki.debian.org/QemuUserEmulation
