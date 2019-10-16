# Horizon IBM Watson Speech to Text to IBM Event Streams Service for Raspberry Pi

For details about using this service, see [watsons2text.md](watsons2text.md).

## Using the IBM Watson Speech to Text to IBM Event Streams Service

- Before following the steps in this section, install the Horizon agent on your edge device and point it to your Horizon exchange. Also get an API key that is associated with your Horizon instance.
1. Set your exchange org:
```
export HZN_ORG_ID="PUT-YOUR-CLUSTER-NAME-HERE"
```
2. Set your exchange user credentials in the Horizon-supported environment variable and verify it:
```
export HZN_EXCHANGE_USER_AUTH="iamapikey:PUT-YOUR-API-KEY-HERE"
hzn exchange user list
```
3. Choose a id and token for your edge node, create it, and verify it:
```
export HZN_EXCHANGE_NODE_AUTH="PUT-ANY-NODE-ID-HERE:PUT-ANY-NODE-TOKEN-HERE"
hzn exchange node create -n $HZN_EXCHANGE_NODE_AUTH
hzn exchange node confirm
```
4. Deploy (or get access to) an instance of IBM Event Streams that the watsons2text sample can send its data to. Ensure that the topic `watsons2text ` is created in Event Streams. Using information from the Event Streams UI, `export` these environment variables:
    - `MSGHUB_API_KEY`
    - `MSGHUB_BROKER_URL`
    - `MSGHUB_CERT_ENCODED` (if using IBM Event Streams in IBM Cloud Private) due to differences in the `base64` command set this variable as follows based on the platform you're using:
        - On Linux: `MSGHUB_CERT_ENCODED=“$(cat $MSGHUB_CERT_FILE| base64 -w 0)”`
	- On Mac: `MSGHUB_CERT_ENCODED="$(cat $MSGHUB_CERT_FILE| base64)"`
    - `MSGHUB_CERT_FILE` (if using IBM Event Streams in IBM Cloud Private)

- Deploy (or get access to) an instance of IBM Speech to Text that the watsons2text sample can send its data to. Ensure that the Speech to Text service is created. Using information from the Speech to Text UI, `export` these environment variables:
    - `STT_IAM_APIKEY`
    - `STT_URL`

5. Get the user input file for the watsons2text sample:
```
wget https://github.com/open-horizon/examples/raw/master/edge/evtstreams/watson_speech2text/horizon/userinput.json
```
6. Register your edge node with Horizon to use the watsons2text pattern:
```
hzn register -p IBM/pattern-ibm.watsons2text-arm -f userinput.json
```


7. The edge device will make an agreement with one of the Horizon agreement bots (this typically takes about 15 seconds). Repeatedly query the agreements of this device until the `agreement_finalized_time` and `agreement_execution_start_time` fields are filled in:
```
hzn agreement list
```

8. Once the agreement is made, list the docker container edge service that has been started as a result:
``` 
sudo docker ps
```


9. On any machine, install [kafkacat](https://github.com/edenhill/kafkacat#install), then subscribe to the Event Streams topic to see the json data that watsons2text is sending:
  - If using IBM Event Streams in IBM Cloud:
  ```
  kafkacat -C -q -o end -f "%t/%p/%o/%k: %s\n" -b $MSGHUB_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=${MSGHUB_API_KEY:0:16} -X sasl.password=${MSGHUB_API_KEY:16} -t $MSGHUB_TOPIC
  ```
  - If using IBM Event Streams in IBM Cloud Private:
  ```
  kafkacat -C -q -o end -f "%t/%p/%o/%k: %s\n" -b $MSGHUB_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=token -X sasl.password=$MSGHUB_API_KEY -X ssl.ca.location=$MSGHUB_CERT_FILE -t $MSGHUB_TOPIC
  ```
10 See the watsons2text service output:

	on **Linux**:
	```
	tail -f /var/log/syslog | grep watsons2text[[]
	```

11. Unregister your edge node, stopping the watsons2text service:
```
hzn unregister -f
```

## First-Time Edge Service Developer - Building and Publishing Your Own Version of the IBM Watson Speech to Text to IBM Event Streams Service

If you want to create your own Horizon edge service, based on this example, follow the next 2 sections to copy the watsons2text example and start modifying it.

### Preconditions for Developing Your Own Service

1. First, go through the steps in the section above to run the IBM watsons2text service on an edge node.
2. Get a docker hub id at https://hub.docker.com/ , if you don't already have one. (This example is set up to store the docker image in docker hub, but by modifying DOCKER_IMAGE_BASE you can store it in another registry.) Login to the docker registry using your id:
```
echo 'DOCKER-PASSWORD' | docker login -u DOCKER-HUB-ID --password-stdin
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
5. Copy the `watson_speech2text` dir to where you will start development of your new service:
```
cp -a examples/edge/evtstreams/watson_speech2text ~/myservice     # or wherever
cd ~/myservice
```
6. Set the values in `horizon/hzn.json` to your own values.
7. As part of the above section "Using the IBM Watson Speech to Text to IBM Event Streams Edge Service", you created your Exchange user credentials and edge node credentials. Ensure they are set and verify them:
```
export HZN_EXCHANGE_USER_AUTH="iamapikey:PUT-YOUR-API-KEY-HERE"
hzn exchange user list
export HZN_EXCHANGE_NODE_AUTH="PUT-ANY-NODE-ID-HERE:PUT-ANY-NODE-TOKEN-HERE"
hzn exchange node confirm
```
8. Verify that these environment variables are still set from when you used the existing watsons2text sample earlier in this document:
```
echo MSGHUB_API_KEY=$MSGHUB_API_KEY
echo MSGHUB_ADMIN_URL=$MSGHUB_ADMIN_URL
echo MSGHUB_BROKER_URL=$MSGHUB_BROKER_URL
```

### Building and Publishing Your Own Version of the CPU To IBM Event Streams Edge Service

1. Edit `service.sh` however you want.
    - Note: this service is a shell script simply for brevity, but you can write your service in any language.
2. Build the watsons2text docker image:
```
make
```
3. Test the service by having Horizon start it locally:
```
hzn dev service start -S
```
4. Check that the container is running:
```
sudo docker ps 
```

5. See the watsons2text service output:

	on **Linux**:
	```
	tail -f /var/log/syslog | grep watsons2text[[]
	```

6. See the environment variables Horizon passes into your service container:
```
docker inspect $(docker ps -q --filter name=watsons2text) | jq '.[0].Config.Env'
```

7. Stop the service:
```
hzn dev service stop
```
8. Create a service signing key pair in `~/.hzn/keys/` (if you haven't already done so):
```
hzn key create <my-company> <my-email>
```
9. Have Horizon push your docker image to your registry and publish your service in the Horizon Exchange and see it there:
```
hzn exchange service publish -f horizon/service.definition.json
hzn exchange service list
```
10 Publish your edge node deployment pattern in the Horizon Exchange and see it there:
```
hzn exchange pattern publish -f horizon/pattern.json
hzn exchange pattern list
```
11. Register your edge node with Horizon to use your deployment pattern (substitute for `SERVICE_NAME` the value you specified above for `hzn dev service new -s`):
```
hzn register -p pattern-SERVICE_NAME-$(hzn architecture) -f horizon/userinput.json
```


12. The edge device will make an agreement with one of the Horizon agreement bots (this typically takes about 15 seconds). Repeatedly query the agreements of this device until the `agreement_finalized_time` and `agreement_execution_start_time` fields are filled in:
```
hzn agreement list
```

13. Once the agreement is made, list the docker container edge service that has been started as a result:
``` 
sudo docker ps
```


14. On any machine, install [kafkacat](https://github.com/edenhill/kafkacat#install), then subscribe to the Event Streams topic to see the json data that watsons2text is sending:
  - If using IBM Event Streams in IBM Cloud:
  ```
  kafkacat -C -q -o end -f "%t/%p/%o/%k: %s\n" -b $MSGHUB_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=${MSGHUB_API_KEY:0:16} -X sasl.password=${MSGHUB_API_KEY:16} -t $MSGHUB_TOPIC
  ```
  - If using IBM Event Streams in IBM Cloud Private:
  ```
  kafkacat -C -q -o end -f "%t/%p/%o/%k: %s\n" -b $MSGHUB_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=token -X sasl.password=$MSGHUB_API_KEY -X ssl.ca.location=$MSGHUB_CERT_FILE -t $MSGHUB_TOPIC
  ```


15. See the watsons2text service output:
```
# soon you will use 'hzn service log ...' for all platforms
# For now on Linux:
tail -f /var/log/syslog | grep watsons2text[[]
``` 
16. Unregister your edge node, stopping the watsons2text service:
```
hzn unregister -f
```

## Process for the Horizon Development Team to Make Updates to the watsons2text Service

- Do the steps in the Preconditions section above, **except**:
    - export `HZN_EXCHANGE_URL` to the staging instance
    - Do **not** copy the watsons2text directory (use the git files in this directory instead)
    - export `HZN_EXCHANGE_USER_AUTH` to your credentials in the IBM org
- Make whatever code changes are necessary
- Increment `SERVICE_VERSION` in `horizon/hzn.json`
- Make `~/.hzn/keys/service.private.key` and `~/.hzn/keys/service.public.pem` actually be symbolic links to the common keys we use to sign all of our examples.
- Build, test, and publish service:
```
make publish
```
