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

3. As part of the infrasctucture installation process for IBM Edge Computing Manager a file called `agent-install.cfg` was created that contains the values for `HZN_ORG_ID` and the exchange and css url values. Locate this file and set those environment variables in your shell now:

```bash
eval export $(cat agent-install.cfg)
```

 - Note: if for some reason you are disconnected from ssh or your command line closes, run the above command again to set the required environment variables.

4. In addition to the file above, an API key associated with your Horizon instance would have been created, set the exchange user credentials, and verify them:

```bash
export HZN_EXCHANGE_USER_AUTH=iamapikey:<horizon-API-key>
hzn exchange user list
```

5. Choose an ID and token for your edge node, create it, and verify it:

```bash
export HZN_EXCHANGE_NODE_AUTH="<choose-any-node-id>:<choose-any-node-token>"
hzn exchange node create -n $HZN_EXCHANGE_NODE_AUTH
hzn exchange node confirm
```

6. While this service can be used with any kafka based message brokers, if you are using IBM Event Streams and an instance has already been deployed for you, obtain the `event-streams.cfg` file that was created during this process. This file contains all the necessary environment variables for `cpu2evtstreams` to publish data to IBM Event Streams. Set these environment variables in your shell now:
```bash
eval export $(cat event-streams.cfg)
```

7. If you have not done so already, unregister your node before moving on:
 ```bash
hzn unregister -f
```


## <a id=using-cpu2evtstreams-pattern></a> Using the CPU To IBM Event Streams Edge Service with Deployment Pattern

1. Get the user input file for the cpu2evtstreams sample and the policy file for the gps service to run privileged:
```bash
wget https://raw.githubusercontent.com/open-horizon/examples/master/edge/evtstreams/cpu2evtstreams/horizon/use/userinput.json
wget https://raw.githubusercontent.com/open-horizon/examples/master/edge/services/gps/horizon/node_policy_privileged.json
```
2. Register your edge node with Horizon to use the cpu2evtstreams pattern:
```bash
hzn register -p IBM/pattern-ibm.cpu2evtstreams -f userinput.json -s ibm.cpu2evtstreams --serviceorg IBM -t 120 --policy=node_policy_privileged.json
```
 - **Note**: using the `-s` flag with the `hzn register` command will cause Horizon to wait until agreements are formed and the service is running on your edge node to exit, or alert you of any errors encountered during the registration process. 

3. View the formed agreement:
```bash
hzn agreement list
```

4. Once the agreement is made, list the docker container edge service that has been started as a result:
```bash
sudo docker ps
```

5. On any machine, install [kafkacat](https://github.com/edenhill/kafkacat#install), then subscribe to the Event Streams topic to see the json data that cpu2evtstreams is sending:
  ```bash
  kafkacat -C -q -o end -f "%t/%p/%o/%k: %s\n" -b $EVTSTREAMS_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=token -X sasl.password=$EVTSTREAMS_API_KEY -X ssl.ca.location=$EVTSTREAMS_CERT_FILE -t cpu2evtstreams
  ```
6. See the cpu2evtstreams service output:

```bash
hzn service log -f ibm.cpu2evtstreams
```
 - **Note**: Press **Ctrl C** to stop the command output.

7. Unregister your edge node, stopping the cpu2evtstreams service:
```bash
hzn unregister -f
```
