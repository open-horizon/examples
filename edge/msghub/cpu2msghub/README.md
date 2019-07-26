# Horizon CPU To IBM Message Hub Service

For details about using this service, see [cpu2msghub.md](cpu2msghub.md).

## Using the CPU To IBM Message Hub Edge Service

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
- Deploy (or get access to) an instance of IBM Event Streams that the cpu2msghub sample can send its data to. Ensure that the topic `cpu2msghub` is created in Event Streams. Using information from the Event Streams UI, `export` these environment variables:
    - `MSGHUB_API_KEY`
    - `MSGHUB_BROKER_URL`
    - `MSGHUB_CERT_ENCODED` (if using IBM Event Streams in IBM Cloud Private)
    - `MSGHUB_CERT_FILE` (if using IBM Event Streams in IBM Cloud Private)
- Get the user input file for the cpu2msghub sample:
```
wget https://github.com/open-horizon/examples/raw/master/edge/msghub/cpu2msghub/horizon/use/userinput.json
```
- Register your edge node with Horizon to use the cpu2msghub pattern:
```
hzn register -p IBM/pattern-ibm.cpu2msghub -f userinput.json
```
- Look at the Horizon agreement until it is finalized and then see the running container:
```
hzn agreement list
docker ps
```
- On any machine, install [kafkacat](https://github.com/edenhill/kafkacat#install), then subscribe to the msg hub topic to see the json data that cpu2msghub is sending:
  - If using IBM Event Streams in IBM Cloud:
  ```
  kafkacat -C -q -o end -f "%t/%p/%o/%k: %s\n" -b $MSGHUB_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=${MSGHUB_API_KEY:0:16} -X sasl.password=${MSGHUB_API_KEY:16} -t $MSGHUB_TOPIC
  ```
  - If using IBM Event Streams in IBM Cloud Private:
  ```
  kafkacat -C -q -o end -f "%t/%p/%o/%k: %s\n" -b $MSGHUB_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=token -X sasl.password=$MSGHUB_API_KEY -X ssl.ca.location=$MSGHUB_CERT_FILE -t cpu2msghub
  ```
- (Optional) To see the cpu2msghub service output:
```
# On Linux:
tail -f /var/log/syslog | grep cpu2msghub[[]
# On Mac:
docker logs -f $(docker ps -q --filter name=cpu2msghub)
``` 
- Unregister your edge node, stopping the cpu2msghub service:
```
hzn unregister -f
```

## First-Time Edge Service Developer - Building and Publishing Your Own Version of the CPU To IBM Message Hub Edge Service

If you want to create your own Horizon edge service, based on this example, follow the next 2 sections to copy the cpu2msghub example and start modifying it.

### Preconditions for Developing Your Own Service

- First, go through the steps in the section above to run the IBM cpu2msghub service on an edge node.
- Get a docker hub id at https://hub.docker.com/ , if you don't already have one. (This example is set up to store the docker image in docker hub, but by modifying DOCKER_IMAGE_BASE you can store it in another registry.) Login to the docker registry using your id:
```
echo 'mydockerpw' | docker login -u mydockehubid --password-stdin
```
- If you have the HZN_ORG_ID environment variable set from previous work, unset it (in a moment this value will now come from `horizon/hzn.json`):
```
unset HZN_ORG_ID
```
- Clone this git repo:
```
cd ~   # or wherever you want
git clone git@github.com:open-horizon/examples.git
```
- Copy the `cpu2msghub` dir to where you will start development of your new service:
```
cp -a examples/edge/msghub/cpu2msghub ~/myservice     # or wherever
cd ~/myservice
```
- Set the values in `horizon/hzn.json` to your own values.
- As part of the above section "Using the CPU To IBM Message Hub Edge Service", you created your Exchange user credentials and edge node credentials. Ensure they are set and verify them:
```
export HZN_EXCHANGE_USER_AUTH="iamapikey:<myapikey>"
hzn exchange user list
export HZN_EXCHANGE_NODE_AUTH="<mynodeid>:<mynodetoken>"
hzn exchange node confirm
```
- Verify that these environment variables are still set from when you used the existing cpu2msghub sample earlier in this document:
```
echo MSGHUB_API_KEY=$MSGHUB_API_KEY
echo MSGHUB_ADMIN_URL=$MSGHUB_ADMIN_URL
echo MSGHUB_BROKER_URL=$MSGHUB_BROKER_URL
```
- Verify the `cpu2msghub` topic is now in your event streams instance:
```
make msghub-topic-list
```

### Building and Publishing Your Own Version of the CPU To IBM Message Hub Edge Service

- Edit `service.sh` however you want.
    - Note: this service is a shell script simply for brevity, but you can write your service in any language.
- Build the cpu2msghub docker image:
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
tail -f /var/log/syslog | grep cpu2msghub[[]
# For now on Mac:
docker logs -f $(docker ps -q --filter name=cpu2msghub)
```
- See the environment variables Horizon passes into your service container:
```
docker inspect $(docker ps -q --filter name=cpu2msghub) | jq '.[0].Config.Env'
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
hzn register -p pattern-SERVICE_NAME-$(hzn architecture) -f horizon/userinput.json
```
- Look at the Horizon agreement until it is finalized and then see the running container:
```
hzn agreement list
docker ps
```
- On any machine, subscribe to the msg hub topic to see the json data that cpu2msghub is sending:
```
kafkacat -C -q -o end -f "%t/%p/%o/%k: %s\n" -b $MSGHUB_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=${MSGHUB_API_KEY:0:16} -X sasl.password=${MSGHUB_API_KEY:16} -t $MSGHUB_TOPIC
```
- See the cpu2msghub service output:
```
# soon you will use 'hzn service log ...' for all platforms
# For now on Linux:
tail -f /var/log/syslog | grep cpu2msghub[[]
# For now on Mac:
docker logs -f $(docker ps -q --filter name=cpu2msghub)
``` 
- Unregister your edge node, stopping the cpu2msghub service:
```
hzn unregister -f
```

## Process for the Horizon Development Team to Make Updates to the Cpu2msghub Service

- Do the steps in the Preconditions section above, **except**:
    - export `HZN_EXCHANGE_URL` to the staging instance
    - Do **not** copy the cpu2msghub directory (use the git files in this directory instead)
    - export `HZN_EXCHANGE_USER_AUTH` to your credentials in the IBM org
- Make whatever code changes are necessary
- Increment `SERVICE_VERSION` in `horizon/hzn.json`
- Make `~/.hzn/keys/service.private.key` and `~/.hzn/keys/service.public.pem` actually be symbolic links to the common keys we use to sign all of our examples.
- Build, test, and publish for all architectures:
```
make publish-all-arches
```
Note: building all architectures works on mac os x, and can be made to work on ubuntu via: http://wiki.micromint.com/index.php/Debian_ARM_Cross-compile , https://wiki.debian.org/QemuUserEmulation
