# Horizon SDR To IBM Event Streams Service

For details about using this service, see [sdr2evtstreams.md](sdr2evtstreams.md).

## Using the SDR To IBM Event Streams Edge Service

- First, go through the "Try It" page "Installing Horizon Software On Your Edge Machine" to set up your edge node.
- Get an IBM cloud account (and for now have your org created in the exchange)
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
- Deploy (or get access to) an instance of IBM Event Streams in the IBM Cloud that the sdr2evtstreams sample can send its data to. In the Event Streams UI, go to the `Service credentials` tab, create new credentials, and use the following values to `export` these environment variables:
    - Set `EVTSTREAMS_API_KEY` to the value of `api_key`
    - Set `EVTSTREAMS_ADMIN_URL` to the value of `kafka_admin_url`
    - Set `EVTSTREAMS_BROKER_URL` to all of the values in `kafka_brokers_sasl` separated by commas
- Create the `sdr2evtstreams` topic (sdr-audio) in your event streams instance:
```
export EVTSTREAMS_TOPIC=sdr-audio
curl -sS -w %{http_code} -H 'Content-Type: application/json' -H "X-Auth-Token: $EVTSTREAMS_API_KEY" -d "{ \"name\": \"$EVTSTREAMS_TOPIC\", \"partitions\": 2 }" $EVTSTREAMS_ADMIN_URL/admin/topics
```
- Verify the `sdr-audio` topic is now in your event streams instance:
```
curl -sS -H "X-Auth-Token: $EVTSTREAMS_API_KEY" $EVTSTREAMS_ADMIN_URL/admin/topics | jq -r ".[] | .name"
```
- Get the user input file for the sdr2evtstreams sample:
```
wget https://github.com/open-horizon/examples/raw/master/edge/evtstreams/sdr2evtstreams/horizon/use/userinput.json
```
- Register your edge node with Horizon to use the sdr2evtstreams pattern:
```
hzn register -p IBM/pattern-ibm.sdr2evtstreams -f userinput.json
```
- Look at the Horizon agreement until it is finalized and then see the running container:
```
hzn agreement list
docker ps
```
- On any machine, install `kafkacat`, then subscribe to the Event Streams topic to see the json data that sdr2evtstreams is sending:
```
kafkacat -C -q -o end -f "%t/%p/%o/%k: %s\n" -b $EVTSTREAMS_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=${EVTSTREAMS_API_KEY:0:16} -X sasl.password=${EVTSTREAMS_API_KEY:16} -t $EVTSTREAMS_TOPIC
```
- (Optional) To see the sdr2evtstreams service output:
```
# soon you will use 'hzn service log ...' for all platforms
# For now on Linux:
tail -f /var/log/syslog | grep sdr2evtstreams[[]
# For now on Mac:
docker logs -f $(docker ps -q --filter name=sdr2evtstreams)
``` 
- Unregister your edge node, stopping the sdr2evtstreams service:
```
hzn unregister -f
```

## First-Time Edge Service Developer - Building and Publishing Your Own Version of the SDR To IBM Event Streams Edge Service

If you want to create your own Horizon edge service, based on this example, follow the next 2 sections to copy the sdr2evtstreams example and start modifying it.

### Preconditions for Developing Your Own Service

- First, go through the steps in the section above to run the IBM sdr2evtstreams service on an edge node.
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
- Copy the `sdr2evtstreams` dir to where you will start development of your new service:
```
cp -a examples/edge/evtstreams/sdr2evtstreams ~/myservice     # or wherever
cd ~/myservice
```
- Set the values in `horizon/hzn.json` to your own values.
- As part of the above section "Using the SDR To IBM Event Streams Edge Service", you created your Exchange user credentials and edge node credentials. Ensure they are set and verify them:
```
export HZN_EXCHANGE_USER_AUTH="iamapikey:<myapikey>"
hzn exchange user list
export HZN_EXCHANGE_NODE_AUTH="<mynodeid>:<mynodetoken>"
hzn exchange node confirm
```
- Verify that these environment variables are still set from when you used the existing sdr2evtstreams sample earlier in this document:
```
echo EVTSTREAMS_API_KEY=$EVTSTREAMS_API_KEY
echo EVTSTREAMS_ADMIN_URL=$EVTSTREAMS_ADMIN_URL
echo EVTSTREAMS_BROKER_URL=$EVTSTREAMS_BROKER_URL
```
- Verify the `sdr2evtstreams` topic is now in your event streams instance:
```
make evtstreams-topic-list
```

### Building and Publishing Your Own Version of the SDR To IBM Event Streams Edge Service

- Edit `service.sh` however you want.
    - Note: this service is a shell script simply for brevity, but you can write your service in any language.
- Build the sdr2evtstreams docker image:
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
tail -f /var/log/syslog | grep sdr2evtstreams[[]
# For now on Mac:
docker logs -f $(docker ps -q --filter name=sdr2evtstreams)
```
- See the environment variables Horizon passes into your service container:
```
docker inspect $(docker ps -q --filter name=sdr2evtstreams) | jq '.[0].Config.Env'
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
- On any machine, subscribe to the Event Streams topic to see the json data that sdr2evtstreams is sending:
```
kafkacat -C -q -o end -f "%t/%p/%o/%k: %s\n" -b $EVTSTREAMS_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=${EVTSTREAMS_API_KEY:0:16} -X sasl.password=${EVTSTREAMS_API_KEY:16} -t $EVTSTREAMS_TOPIC
```
- See the sdr2evtstreams service output:
```
# soon you will use 'hzn service log ...' for all platforms
# For now on Linux:
tail -f /var/log/syslog | grep sdr2evtstreams[[]
# For now on Mac:
docker logs -f $(docker ps -q --filter name=sdr2evtstreams)
``` 
- Unregister your edge node, stopping the sdr2evtstreams service:
```
hzn unregister -f
```

## Process for the Horizon Development Team to Make Updates to the sdr2evtstreams Service

- Do the steps in the Preconditions section above, **except**:
    - export `HZN_EXCHANGE_URL` to the staging instance
    - Do **not** copy the sdr2evtstreams directory (use the git files in this directory instead)
    - export `HZN_EXCHANGE_USER_AUTH` to your credentials in the IBM org
- Make whatever code changes are necessary
- Increment `SERVICE_VERSION` in `horizon/hzn.json`
- Make `~/.hzn/keys/service.private.key` and `~/.hzn/keys/service.public.pem` actually be symbolic links to the common keys we use to sign all of our examples.
- Build, test, and publish for all architectures:
```
make publish-all-arches
```
Note: building all architectures works on mac os x, and can be made to work on ubuntu via: http://wiki.micromint.com/index.php/Debian_ARM_Cross-compile , https://wiki.debian.org/QemuUserEmulation

