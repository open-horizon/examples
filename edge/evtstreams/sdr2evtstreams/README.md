# Horizon SDR To IBM Event Streams Service

This is a simple example of using and creating a Horizon edge service.

- [Preconditions for Using the SDR To IBM Event Streams Example Edge Service](#preconditions)

- [Using the SDR To IBM Event Streams Edge Service with Deployment Pattern](#using-sdr2evtstreams-pattern)

- [Building and Publishing Your Own Version of the SDR To IBM Event Streams Edge Service](#building-your-own-sdr2evtstreams-pattern)

- [Process for the Horizon Development Team to Make Updates to the SDR To IBM Event Streams Edge Service](#dev-team-updates-sdr2evtstreams)

- For details about using this service, see [sdr2evtstreams.md](sdr2evtstreams.md).


## <a id=preconditions></a> Preconditions for Using the SDR To IBM Event Streams Example Edge Service

If you haven't done so already, you must do these steps before proceeding with the helloworld example:

1. Install the Horizon management infrastructure (exchange and agbot).

2. Install the Horizon agent on your edge device and configure it to point to your Horizon exchange.

3. Set your exchange org:

```bash
export HZN_ORG_ID="<your-cluster-name>"
```

4. Create a cloud API key that is associated with your Horizon instance, set your exchange user credentials, and verify them:

```bash
export HZN_EXCHANGE_USER_AUTH="iamapikey:<your-API-key>"
hzn exchange user list
```

5. Choose an ID and token for your edge node, create it, and verify it:

```bash
export HZN_EXCHANGE_NODE_AUTH="<choose-any-node-id>:<choose-any-node-token>"
hzn exchange node create -n $HZN_EXCHANGE_NODE_AUTH
hzn exchange node confirm
```


## <a id=using-sdr2evtstreams-pattern></a> Using the SDR To IBM Event Streams Edge Service

- First, go through the "Try It" page "Installing Horizon Software On Your Edge Machine" to set up your edge node.
- Get an IBM cloud account (and for now have your org created in the exchange)

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

4. Deploy (or get access to) an instance of IBM Event Streams that the sdr2evtstreams sample can send its data to. Ensure that the topic `sdr-audio` is created in Event Streams. Using information from the Event Streams UI, `export` these environment variables:
    - `EVTSTREAMS_API_KEY`
    - `EVTSTREAMS_BROKER_URL`
    - `EVTSTREAMS_CERT_ENCODED` (if using IBM Event Streams in IBM Cloud Private) due to differences in the base64 command set this variable as follows based on the platform you're using:
        - on **Linux**: `EVTSTREAMS_CERT_ENCODED=“$(cat $EVTSTREAMS_CERT_FILE | base64 -w 0)”`
        - on **Mac**: `EVTSTREAMS_CERT_ENCODED=“$(cat $EVTSTREAMS_CERT_FILE | base64)”`
    - `EVTSTREAMS_CERT_FILE` (if using IBM Event Streams in IBM Cloud Private)

5. Get the user input file for the sdr2evtstreams sample:
```
wget https://github.com/open-horizon/examples/raw/master/edge/evtstreams/sdr2evtstreams/horizon/userinput.json
```
6. Register your edge node with Horizon to use the sdr2evtstreams pattern:
```
hzn register -p IBM/pattern-ibm.sdr2evtstreams -f userinput.json
```
7. The edge device will make an agreement with one of the Horizon agreement bots (this typically takes about 15 seconds). Repeatedly query the agreements of this device until the `agreement_finalized_time` and `agreement_execution_start_time` fields are filled in:
```
hzn agreement list
```
8. Once the agreement is made, list the docker container edge service that has been started as a result:
```
sudo docker ps
```

9. On any machine, install `kafkacat`, then subscribe to the Event Streams topic to see the json data that sdr2evtstreams is sending:
```
kafkacat -C -q -o end -f "%t/%p/%o/%k: %s\n" -b $EVTSTREAMS_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=${EVTSTREAMS_API_KEY:0:16} -X sasl.password=${EVTSTREAMS_API_KEY:16} -t $EVTSTREAMS_TOPIC
```

10. See the sdr2evtstreams service output:

	on **Linux**:
	```
	tail -f /var/log/syslog | grep sdr2evtstreams[[]
	```

	on **Mac**:
	```
	docker logs -f $(docker ps -q --filter name=sdr2evtstreams)
	``` 

11. Unregister your edge node, stopping the sdr2evtstreams service:
```
hzn unregister -f
```


## <a id=building-your-own-sdr2evtstreams-pattern></a> First-Time Edge Service Developer - Building and Publishing Your Own Version of the SDR To IBM Event Streams Edge Service

If you want to create your own Horizon edge service, based on this example, follow the next 2 sections to copy the sdr2evtstreams example and start modifying it.

### Preconditions for Developing Your Own Service

1. First, go through the steps in the section above to run the IBM sdr2evtstreams service on an edge node.
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
5. Copy the `sdr2evtstreams` dir to where you will start development of your new service:
```
cp -a examples/edge/evtstreams/sdr2evtstreams ~/myservice     # or wherever
cd ~/myservice
```
6. Set the values in `horizon/hzn.json` to your own values.
7. As part of the above section "Using the SDR To IBM Event Streams Edge Service", you created your Exchange user credentials and edge node credentials. Ensure they are set and verify them:
```
export HZN_EXCHANGE_USER_AUTH="iamapikey:<myapikey>"
hzn exchange user list
export HZN_EXCHANGE_NODE_AUTH="<mynodeid>:<mynodetoken>"
hzn exchange node confirm
```
8. Verify that these environment variables are still set from when you used the existing sdr2evtstreams sample earlier in this document:
```
echo EVTSTREAMS_API_KEY=$EVTSTREAMS_API_KEY
echo EVTSTREAMS_ADMIN_URL=$EVTSTREAMS_ADMIN_URL
echo EVTSTREAMS_BROKER_URL=$EVTSTREAMS_BROKER_URL
```
9. Verify the `sdr2evtstreams` topic is now in your event streams instance:
```
make evtstreams-topic-list
```

### Building and Publishing Your Own Version of the SDR To IBM Event Streams Edge Service

1. Edit `main.go` however you want.
    - Note: this service is written in go, but you can write your service in any language.
2. Build the sdr2evtstreams docker image:
```
make
```
3. Test the service by having Horizon start it locally:
```
hzn dev service start -S
```
4. Check that the containers are running:
```
sudo docker ps
```

5. See the sdr2evtstreams service output:

	on **Linux**:
	```
	tail -f /var/log/syslog | grep sdr2evtstreams[[]
	```

	on **Mac**:
	```
	docker logs -f $(docker ps -q --filter name=sdr2evtstreams)
	```

6. See the environment variables Horizon passes into your service container:
```
docker inspect $(docker ps -q --filter name=sdr2evtstreams) | jq '.[0].Config.Env'
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
10. Publish your edge node deployment pattern in the Horizon Exchange and see it there:
```
hzn exchange pattern publish -f horizon/pattern.json
hzn exchange pattern list
```
11. Register your edge node with Horizon to use your deployment pattern (substitute for `SERVICE_NAME` the value you specified above for `hzn dev service new -s`):
```
hzn register -p pattern-SERVICE_NAME-$(hzn architecture) -f horizon/userinput.json
```
12. Look at the Horizon agreement until it is finalized and then see the running container:
```
hzn agreement list
docker ps
```
13. On any machine, subscribe to the Event Streams topic to see the json data that sdr2evtstreams is sending:
```
kafkacat -C -q -o end -f "%t/%p/%o/%k: %s\n" -b $EVTSTREAMS_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=token -X sasl.password=$EVTSTREAMS_API_KEY -X ssl.ca.location=$EVTSTREAMS_CERT_FILE -t $EVTSTREAMS_TOPIC
```

14. See the sdr2evtstreams service output:

	on **Linux**:
	```
	tail -f /var/log/syslog | grep sdr2evtstreams[[]
	```

	on **Mac**:
	```
	docker logs -f $(docker ps -q --filter name=sdr2evtstreams)
	``` 

15. Unregister your edge node, stopping the sdr2evtstreams service:
```
hzn unregister -f
```

## <a id=dev-team-updates-sdr2evtstreams></a> Process for the Horizon Development Team to Make Updates to the sdr2evtstreams Service

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

