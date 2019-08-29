# Horizon CPU To IBM Event Streams Service

For details about using this service, see [cpu2evtstreams.md](cpu2evtstreams.md).

## Using the CPU To IBM Event Streams Edge Service

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
- Deploy (or get access to) an instance of IBM Event Streams that the cpu2evtstreams sample can send its data to. Ensure that the topic `cpu2evtstreams` is created in Event Streams. Using information from the Event Streams UI, `export` these environment variables:
    - `EVTSTREAMS_API_KEY`
    - `EVTSTREAMS_BROKER_URL`
    - `EVTSTREAMS_CERT_ENCODED` (if using IBM Event Streams in IBM Cloud Private)
    - `EVTSTREAMS_CERT_FILE` (if using IBM Event Streams in IBM Cloud Private)
- Get the user input file for the cpu2evtstreams sample:
```
wget https://github.com/open-horizon/examples/raw/master/edge/evtstreams/cpu2evtstreams/horizon/use/userinput.json
```
- Register your edge node with Horizon to use the cpu2evtstreams pattern:
```
hzn register -p IBM/pattern-ibm.cpu2evtstreams -f userinput.json
```
- Look at the Horizon agreement until it is finalized and then see the running container:
```
hzn agreement list
docker ps
```
- On any machine, install [kafkacat](https://github.com/edenhill/kafkacat#install), then subscribe to the Event Streams topic to see the json data that cpu2evtstreams is sending:
  - If using IBM Event Streams in IBM Cloud:
  ```
  kafkacat -C -q -o end -f "%t/%p/%o/%k: %s\n" -b $EVTSTREAMS_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=${EVTSTREAMS_API_KEY:0:16} -X sasl.password=${EVTSTREAMS_API_KEY:16} -t $EVTSTREAMS_TOPIC
  ```
  - If using IBM Event Streams in IBM Cloud Private:
  ```
  kafkacat -C -q -o end -f "%t/%p/%o/%k: %s\n" -b $MEVTSTREAMS_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=token -X sasl.password=$EVTSTREAMS_API_KEY -X ssl.ca.location=$EVTSTREAMS_CERT_FILE -t cpu2evtstreams
  ```
- (Optional) To see the cpu2evtstreams service output:
```
# On Linux:
tail -f /var/log/syslog | grep cpu2evtstreams[[]
# On Mac:
docker logs -f $(docker ps -q --filter name=cpu2evtstreams)
``` 
- Unregister your edge node, stopping the cpu2evtstreams service:
```
hzn unregister -f
```

## First-Time Edge Service Developer - Building and Publishing Your Own Version of the CPU To IBM Event Streams Edge Service

If you want to create your own Horizon edge service, based on this example, follow the next 2 sections to copy the cpu2evtstreams example and start modifying it.

### Preconditions for Developing Your Own Service

- First, go through the steps in the section above to run the IBM cpu2evtstreams service on an edge node.
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
- Copy the `cpu2evtstreams` dir to where you will start development of your new service:
```
cp -a examples/edge/evtstreams/cpu2evtstreams ~/myservice     # or wherever
cd ~/myservice
```
- Set the values in `horizon/hzn.json` to your own values.
- As part of the above section "Using the CPU To IBM Event Streams Edge Service", you created your Exchange user credentials and edge node credentials. Ensure they are set and verify them:
```
export HZN_EXCHANGE_USER_AUTH="iamapikey:<myapikey>"
hzn exchange user list
export HZN_EXCHANGE_NODE_AUTH="<mynodeid>:<mynodetoken>"
hzn exchange node confirm
```
- Verify that these environment variables are still set from when you used the existing cpu2evtstreams sample earlier in this document:
```
echo EVTSTREAMS_API_KEY=$EVTSTREAMS_API_KEY
echo EVTSTREAMS_ADMIN_URL=$EVTSTREAMS_ADMIN_URL
echo EVTSTREAMS_BROKER_URL=$EVTSTREAMS_BROKER_URL
```
- Verify the `cpu2evtstreams ` topic is now in your event streams instance:
```
make evtstreams-topic-list
```

### Building and Publishing Your Own Version of the CPU To IBM Event Streams Edge Service

- Edit `service.sh` however you want.
    - Note: this service is a shell script simply for brevity, but you can write your service in any language.
- Build the cpu2evtstreams docker image:
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
tail -f /var/log/syslog | grep cpu2evtstreams[[]
# For now on Mac:
docker logs -f $(docker ps -q --filter name=cpu2evtstreams)
```
- See the environment variables Horizon passes into your service container:
```
docker inspect $(docker ps -q --filter name=cpu2evtstreams) | jq '.[0].Config.Env'
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
- On any machine, subscribe to the Event Streams topic to see the json data that cpu2evtstreams is sending:
```
kafkacat -C -q -o end -f "%t/%p/%o/%k: %s\n" -b $EVTSTREAMS_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=${EVTSTREAMS_API_KEY:0:16} -X sasl.password=${EVTSTREAMS_API_KEY:16} -t $EVTSTREAMS_TOPIC
```
- See the cpu2evtstreams service output:
```
# soon you will use 'hzn service log ...' for all platforms
# For now on Linux:
tail -f /var/log/syslog | grep cpu2evtstreams[[]
# For now on Mac:
docker logs -f $(docker ps -q --filter name=cpu2evtstreams)
``` 
- Unregister your edge node, stopping the cpu2evtstreams service:
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
