# Horizon CPU To IBM Event Streams Service

This example illustrates a more realistic Horizon edge service by including additional aspects of typical edge services. 

- [Preconditions for Using the CPU To IBM Event Streams Example Edge Service](#preconditions)

- [Using the CPU To IBM Event Streams Example Edge Service with Deployment Pattern](#using-cpu2evtstreams-pattern)

- [Using the CPU To IBM Event Streams Example Edge Service with Deployment Policy](PolicyRegister.md)

- [Creating Your Own CPU To IBM Event Streams Example Edge Service](CreateService.md)

- For details about using this service, see [cpu2evtstreams.md](cpu2evtstreams.md).


## <a id=preconditions></a> Preconditions for Using the CPU To IBM Event Streams Example Edge Service

If you haven't done so already, you must do these steps before proceeding with the cpu2evtstreams example:

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
export HZN_EXCHANGE_NODE_AUTH="<choose-any-node-id>:<choose-any-node-token>"
hzn exchange node create -n $HZN_EXCHANGE_NODE_AUTH
hzn exchange node confirm
```

6. While this service can be used with any kafka based message brokers, if you are using IBM Event Streams and an instance has already been deployed for you, obtain the `evtstreams.cfg` file that was created during this process. This file contains all the necessary environment variables for `cpu2evtstreams` to publish data to IBM Event Streams. Set these environment variables in your shell now:
```bash
eval export $(cat agent-install.cfg)
```

## <a id=using-cpu2evtstreams-pattern></a> Using the CPU To IBM Event Streams Edge Service with Deployment Pattern

1. Get the user input file for the cpu2evtstreams sample:
```bash
wget https://github.com/open-horizon/examples/raw/master/edge/evtstreams/cpu2evtstreams/horizon/use/userinput.json
```
2. Register your edge node with Horizon to use the cpu2evtstreams pattern:
```bash
hzn register -p IBM/pattern-ibm.cpu2evtstreams -f userinput.json
```

The cpu2evtstreams sample has a dependent service called gps. This service requires a privileged node to run due to the hardware access it requires to get your coordinates.
3. Create a `policy.json` file to update your node policy:
```bash
cat << EndOfContent > policy.json
{
  "properties": [
    {
      "name": "openhorizon.allowPrivileged",
      "value": true
    }
  ]
}
EndOfContent
```
4. Update your node policy to authorize the node:
```bash
hzn hzn policy update --input-file=policy.json
```

5. The edge device will make an agreement with one of the Horizon agreement bots (this typically takes about 15 seconds). Repeatedly query the agreements of this device until the `agreement_finalized_time` and `agreement_execution_start_time` fields are filled in:
```bash
hzn agreement list
```

6. Once the agreement is made, list the docker container edge service that has been started as a result:
```bash
sudo docker ps
```

7. On any machine, install [kafkacat](https://github.com/edenhill/kafkacat#install), then subscribe to the Event Streams topic to see the json data that cpu2evtstreams is sending:
  ```bash
  kafkacat -C -q -o end -f "%t/%p/%o/%k: %s\n" -b $EVTSTREAMS_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=token -X sasl.password=$EVTSTREAMS_API_KEY -X ssl.ca.location=$EVTSTREAMS_CERT_FILE -t cpu2evtstreams
  ```
8. See the cpu2evtstreams service output:

```bash
hzn service log -f ibm.cpu2evtstreams
```

9. Unregister your edge node, stopping the cpu2evtstreams service:
```bash
hzn unregister -f
```
