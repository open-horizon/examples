
# Horizon Offline Voice Assistant Example Edge Service for Raspberry Pi

## Using the Offline Voice Assistant Example Edge Service

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
4. Register your edge node with Horizon to use the helloworld pattern:
```
hzn register -p IBM/pattern-ibm.processtext-arm
```


5. The edge device will make an agreement with one of the Horizon agreement bots (this typically takes about 15 seconds). Repeatedly query the agreements of this device until the `agreement_finalized_time` and `agreement_execution_start_time` fields are filled in:
```
hzn agreement list
```

6. Once the agreement is made, list the docker container edge service that has been started as a result:
``` 
sudo docker ps
```

7. See the processtext service output:

	on **Linux**:

	```
	tail -f /var/log/syslog | grep OVA
	``` 


8. Unregister your edge node, stopping the processtext service:
```
hzn unregister -f
```

## First-Time Edge Service Developer - Building and Publishing Your Own Version of the Offline Voice Assistant Edge Service

If you want to create your own Horizon edge service, follow the next 2 sections to copy the Offline Voice Assistant Example Edge Service and start modifying it.

### Preconditions for Developing Your Own Service

1. First, go through the steps in the section above to run the IBM processtext service on an edge node.
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
5. Copy the `processtext` dir to where you will start development of your new service:
```
cp -a examples/edge/services/processtext ~/myservice     # or wherever
cd ~/myservice
```

6. Set the values in `horizon/hzn.json` to your own values.


7. As part of the above section "Using the Offline Voice Assistant Example Edge Service", you created your Exchange user credentials and edge node credentials. Ensure they are set and verify them:
```
export HZN_EXCHANGE_USER_AUTH="iamapikey:PUT-YOUR-API-KEY-HERE"
hzn exchange user list
export HZN_EXCHANGE_NODE_AUTH="PUT-ANY-NODE-ID-HERE:PUT-ANY-NODE-TOKEN-HERE"
hzn exchange node confirm
```

### Building and Publishing Your Own Version of the Offline Voice Assistant Example Edge Service

1. Build the processtext docker image:
```
make
```
2. Test the service by having Horizon start it locally:
```
hzn dev service start -S
```
3.. Check that the container is running:
```
sudo docker ps 
```

4. See the processtext service output:

	on **Linux**:

	```
	tail -f /var/log/syslog | grep OVA
	```

5. See the environment variables Horizon passes into your service container:
```
docker inspect $(docker ps -q --filter name=ibm.processtext) | jq '.[0].Config.Env'
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
hzn register -p pattern-SERVICE_NAME-$(hzn architecture)
```

11. The edge device will make an agreement with one of the Horizon agreement bots (this typically takes about 15 seconds). Repeatedly query the agreements of this device until the `agreement_finalized_time` and `agreement_execution_start_time` fields are filled in:
```
hzn agreement list
```

12. Once the agreement is made, list the docker container edge service that has been started as a result:
``` 
sudo docker ps
```


13. See the processtext service output:

	on **Linux**:

	```
	tail -f /var/log/syslog | grep OVA
	```

14. Unregister your edge node, stopping the processtext service:
```
hzn unregister -f
```

## Further Learning

To see more Horizon features demonstrated, continue on to the [cpu2evtstreams example](../../evtstreams/cpu2evtstreams).

## Process for the Horizon Development Team to Make Updates to the Offline Voice Assistant Service

- Do the steps in the Preconditions section above, **except**:
    - export `HZN_EXCHANGE_URL` to the staging instance
    - Do **not** run `hzn dev service new ...` (use the git files in this directory instead)
    - export `HZN_EXCHANGE_USER_AUTH` to your credentials in the IBM org
- Make whatever code changes are necessary
- Increment `SERVICE_VERSION` in `horizon/hzn.json`
- Make `~/.hzn/keys/service.private.key` and `~/.hzn/keys/service.public.pem` actually be symbolic links to the common keys we use to sign all of our examples.
- Build, test, and publish service:
```
make publish 
```