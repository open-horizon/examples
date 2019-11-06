# Horizon SDR To IBM Event Streams Service

This example illustrates a more realistic Horizon edge service by including additional aspects of typical edge services.

- [Preconditions for Using the SDR To IBM Event Streams Example Edge Service](#preconditions)

- [Using the SDR To IBM Event Streams Edge Service with Deployment Pattern](#using-sdr2evtstreams-pattern)

- [Creating Your Own SDR to IBM Event Streams Edge Service](CreateService.md)

- For details about using this service, see [sdr2evtstreams.md](sdr2evtstreams.md).


## <a id=preconditions></a> Preconditions for Using the SDR To IBM Event Streams Example Edge Service

If you haven't done so already, you must do these steps before proceeding with the sdr2evtstreams example:

1. Install the Horizon management infrastructure (exchange and agbot).

2. Install the Horizon agent on your edge device and configure it to point to your Horizon exchange.

3. Set your exchange org:

```bash
export HZN_ORG_ID=<your-cluster-name>
```

4. Create a cloud API key that is associated with your Horizon instance, set your exchange user credentials, and verify them:

```bash
export HZN_EXCHANGE_USER_AUTH=iamapikey:<your-API-key>
hzn exchange user list
```

5. Choose an ID and token for your edge node, create it, and verify it:

```bash
export HZN_EXCHANGE_NODE_AUTH=<choose-any-node-id>:<choose-any-node-token>
hzn exchange node create -n $HZN_EXCHANGE_NODE_AUTH
hzn exchange node confirm
```
6. Deploy (or get access to) an instance of IBM Event Streams that the sdr2evtstreams sample can send its data to. Ensure that the topic `sdr-audio` is created in Event Streams. Using information from the Event Streams UI, `export` these environment variables:
    - `EVTSTREAMS_API_KEY`
    - `EVTSTREAMS_BROKER_URL`
    - `EVTSTREAMS_CERT_ENCODED` **(if using IBM Event Streams in IBM Cloud Private)** due to differences in the base64 command set this variable as follows based on the platform you're using:
        - on **Linux**: `EVTSTREAMS_CERT_ENCODED=“$(cat $EVTSTREAMS_CERT_FILE | base64 -w 0)”`
        - on **Mac**: `EVTSTREAMS_CERT_ENCODED=“$(cat $EVTSTREAMS_CERT_FILE | base64)”`
    - `EVTSTREAMS_CERT_FILE` **(if using IBM Event Streams in IBM Cloud Private)**


## <a id=using-sdr2evtstreams-pattern></a> Using the SDR To IBM Event Streams Edge Service with Deployment Pattern

1. Get the user input file for the sdr2evtstreams sample:
```
wget https://github.com/open-horizon/examples/raw/master/edge/evtstreams/sdr2evtstreams/horizon/userinput.json
```
2. Register your edge node with Horizon to use the sdr2evtstreams pattern:
```
hzn register -p IBM/pattern-ibm.sdr2evtstreams -f userinput.json
```
3. The edge device will make an agreement with one of the Horizon agreement bots (this typically takes about 15 seconds). Repeatedly query the agreements of this device until the `agreement_finalized_time` and `agreement_execution_start_time` fields are filled in:
```
hzn agreement list
```
4. Once the agreement is made, list the docker container edge service that has been started as a result:
```
sudo docker ps
```

5. On any machine, install `kafkacat`, then subscribe to the Event Streams topic to see the json data that sdr2evtstreams is sending:
```
kafkacat -C -q -o end -f "%t/%p/%o/%k: %s\n" -b $EVTSTREAMS_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=token -X sasl.password=$EVTSTREAMS_API_KEY -X ssl.ca.location=$EVTSTREAMS_CERT_FILE -t $EVTSTREAMS_TOPIC
```

6. See the sdr2evtstreams service output:

	on **Linux**:
	```
	tail -f /var/log/syslog | grep sdr2evtstreams[[]
	```

	on **Mac**:
	```
	docker logs -f $(docker ps -q --filter name=sdr2evtstreams)
	``` 

7. Unregister your edge node, stopping the sdr2evtstreams service:
```
hzn unregister -f
```

