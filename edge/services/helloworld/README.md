# Horizon Hello World Example Edge Service

## Using the Hello World Example Edge Service

- Before following the steps in this section, install the Horizon agent on your edge device and point it to your Horizon exchange. Also get an API key that is associated with your Horizon instance.
- Set your exchange org:
```
export HZN_ORG_ID="<yourorg>"
```
- Set your exchange user credentials in the Horizon-supported environment variable and verify it:
```
export HZN_EXCHANGE_USER_AUTH="iamapikey:<myapikey>"
hzn exchange user list
```
- Choose a id and token for your edge node, create it, and verify it:
```
export HZN_EXCHANGE_NODE_AUTH="<mynodeid>:<mynodetoken>"
hzn exchange node create -n $HZN_EXCHANGE_NODE_AUTH
hzn exchange node confirm
```
- Register your edge node with Horizon to use the helloworld pattern:
```
hzn register -p IBM/pattern-ibm.helloworld
```
- Look at the Horizon agreement until it is finalized and then see the running container:
```
hzn agreement list
docker ps
```
- See the helloworld service output:
```
# soon you will use 'hzn service log ...' for all platforms
# For now on Linux:
tail -f /var/log/syslog | grep helloworld[[]
# For now on Mac:
docker logs -f $(docker ps -q --filter name=helloworld)
``` 
- Unregister your edge node, stopping the helloworld service:
```
hzn unregister -f
```

## Using the Hello World Example Edge Service as a Policy

- The Horizon Policy mechanism offers an alternative to using Deployment Patterns. Policies provide much finer control over agreement forming between Horizon Agents on Edge Nodes, and the Horizon AgBots. It also provides a greater separation of concerns, allowing Edge Nodes owners, Service code developers, and Business owners to each independently articulate their own Policies. There are therefore three types of Horizon Policies:

1. Node Policy (provided at registration time by the node owner)

2. Service Policy (may be applied to a published Service in the Exchange)

3. Business Policy (which approximately corresponds to a Deployment Pattern)

### Node Policy 

- As an alternative to specifying a Deployment Pattern when you register your Edge Node, you may register with a Node Policy.

- Make sure your Edge Node is not registered by running:

```
hzn unregister -f
```

- Now let's register using the `horizon/node_policy.json` file:

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

- Register using this command:

```
hzn register --policy node_policy.json
```

- When the registration completes, use the following command to review the Node Policy:
```
hzn policy list
```

- Notice that in addition to the three properties stated in the node_policy.json file, Horizon has added a few more (openhorizon.cpu, openhorizon.arch, and openhorizon.memory). Horizon provides this additional information automatically and these properties may be used in any of your Policy constraints.

### Service Policy 

- Like the other two Policy types, Service Policy contains a set of properties and a set of constraints. The properties of a Service Policy could state characteristics of the Service code that Node Policy authors or Business Policy authors may find relevant. The constraints of a Service Policy can be used to restrict where this Service can be run. The Service developer could, for example, assert that this Service requires a particular hardware setup such as CPU/GPU constraints, memory constraints, specific sensors, actuators or other peripheral devices required, etc.

- Now let's attach this Service Policy to the helloworld Service previously published using the `horizon/service_policy.json` file:

```
{
  "properties": [
  ],
  "constraints": [
    "model == \"Whatsit ULTRA\" OR model == \"Thingamajig ULTRA\""
  ]
}
```

- Note this simple Service Policy doesn't provide any properties, but it does have a constraint. This example constraint is one that a Service developer might add, stating that their Service must only run on the models named Whatsit ULTRA or Thingamajig ULTRA. If you recall the Node Policy we used above, the model property was set to Thingamajig ULTRA, so this Service should be compatible with our Edge Node.

- To attach the example Service policy to this service, use the following command (substituting your service name):

```
hzn exchange service addpolicy -f service_policy.json <published-helloworld-service-name>
```

Once that completes, you can look at the results with the following command:

```
hzn exchange service listpolicy <published-helloworld-service-name>
```
- Notice that Horizon has again automatically added some additional properties to your Policy. These generated property values can be used in constraints in Node Policies and Business Policies.

- Now that we have set up the Policies for an Edge Node and the Policies for a published Service, we can move on to the final step of defining a Business Policy to tie them all together and cause software to be automatically deployed on your Edge Node.

### Business Policy 




## First-Time Edge Service Developer - Building and Publishing Your Own Version of the Hello World Edge Service

If you want to create your own Horizon edge service, follow the next 2 sections to copy the hello world example and start modifying it.

### Preconditions for Developing Your Own Service

- First, go through the steps in the section above to run the IBM helloworld service on an edge node.
- Get a docker hub id at https://hub.docker.com/ , if you don't already have one. (This example is set up to store the docker image in docker hub, but by modifying DOCKER_IMAGE_BASE you can store it in another registry.) Login to the docker registry using your id:
```
echo 'mydockerpw' | docker login -u mydockehubid --password-stdin
```
- If you have the HZN_ORG_ID environment variable set from previous work, unset it (in a moment this value will now come from `horizon/hzn.json`):
```
unset HZN_ORG_ID
```
- Cd to the directory in which you want to create your new service and then:
```
hzn dev service new -o <org-id> -s <service-name> -i <docker-image-base>
# E.g.:  hzn dev service new -o bp@someemail.com -s bp.helloworld -i brucemp/bp.helloworld
```
- As part of the above section "Using the Hello World Example Edge Service", you created your Exchange user credentials and edge node credentials. Ensure they are set and verify them:
```
export HZN_EXCHANGE_USER_AUTH="iamapikey:<myapikey>"
hzn exchange user list
export HZN_EXCHANGE_NODE_AUTH="<mynodeid>:<mynodetoken>"
hzn exchange node confirm
```

### Building and Publishing Your Own Version of the Hello World Example Edge Service

- Edit `service.sh`, for example changing "Hello" to "Hey there"
    - Note: this service is a shell script simply for brevity, but you can write your service in any language.
- Build the hello world docker image:
```
make
```
- Test the service by having Horizon start it locally:
```
hzn dev service start -S
```
- See the docker container running and look at the output:
```
docker ps
# soon you will use 'hzn service log ...' for all platforms
# For now on Linux:
tail -f /var/log/syslog | grep helloworld[[]
# For now on Mac:
docker logs -f $(docker ps -q --filter name=helloworld)
```
- See the environment variables Horizon passes into your service container:
```
docker inspect $(docker ps -q --filter name=helloworld) | jq '.[0].Config.Env'
```
- Stop the service:
```
hzn dev service stop
```
- Create a service signing key pair in `~/.hzn/keys/` (if you haven't already done so):
```
hzn key create <my-company> <my-email>
```
- Have Horizon push your docker image to your registry and publish your service in the Horizon Exchange and see it there:
```
hzn exchange service publish -f horizon/service.definition.json
hzn exchange service list
```
- Publish your edge node deployment pattern in the Horizon Exchange and see it there:
```
hzn exchange pattern publish -f horizon/pattern.json
hzn exchange pattern list
```
- Register your edge node with Horizon to use your deployment pattern (substitute for `SERVICE_NAME` the value you specified above for `hzn dev service new -s`):
```
hzn register -p pattern-SERVICE_NAME-$(hzn architecture)
```
- Look at the Horizon agreement until it is finalized and then see the running container:
```
hzn agreement list
docker ps
```
- See the helloworld service output:
```
# soon you will use 'hzn service log ...' for all platforms
# For now on Linux:
tail -f /var/log/syslog | grep helloworld[[]
# For now on Mac:
docker logs -f $(docker ps -q --filter name=helloworld)
``` 
- Unregister your edge node, stopping the helloworld service:
```
hzn unregister -f
```

## Further Learning

To see more Horizon features demonstrated, continue on to the [cpu2msghub example](../../msghub/cpu2msghub).

## Process for the Horizon Development Team to Make Updates to the Helloworld Service

- Do the steps in the Preconditions section above, **except**:
    - export `HZN_EXCHANGE_URL` to the staging instance
    - Do **not** run `hzn dev service new ...` (use the git files in this directory instead)
    - export `HZN_EXCHANGE_USER_AUTH` to your credentials in the IBM org
- Make whatever code changes are necessary
- Increment `SERVICE_VERSION` in `horizon/hzn.json`
- Make `~/.hzn/keys/service.private.key` and `~/.hzn/keys/service.public.pem` actually be symbolic links to the common keys we use to sign all of our examples.
- Build, test, and publish for all architectures:
```
make publish-all-arches
```
Note: building all architectures works on mac os x, and can be made to work on ubuntu via: http://wiki.micromint.com/index.php/Debian_ARM_Cross-compile , https://wiki.debian.org/QemuUserEmulation
