# Horizon Hello Model Management Service (MMS) Example Edge Service

This is a simple example of using and creating a Horizon edge service.

- [Preconditions for Using the Hello MMS Example Edge Service](#preconditions)
- [Using the Hello MMS Example Edge Service with Deployment Pattern](#using-hello-mms-pattern)
- [Creating Your Own Hello MMS Edge Service](CreateService.md)
- Further Learning - to see more Horizon features demonstrated, continue on to the [cpu2evtstreams example](../../evtstreams/cpu2evtstreams).

## <a id=preconditions></a> Preconditions for Using the Hello MMS Example Edge Service

If you haven't done so already, you must do these steps before proceeding with the hello-mms example:

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

## <a id=using-hello-mms-pattern></a> Using the Hello MMS Example Edge Service with Deployment Pattern

1. Register your edge node with Horizon to use the hello-mms pattern:

```bash
hzn register -p IBM/pattern-ibm.hello-mms
```

2. The edge device will make an agreement with one of the Horizon agreement bots (this typically takes about 15 seconds). Repeatedly query the agreements of this device until the `agreement_finalized_time` and `agreement_execution_start_time` fields are filled in:

```bash
hzn agreement list
```

3. After the agreement is made, list the docker container edge service that has been started as a result:

``` bash
sudo docker ps
```

4. See the hello-mms service output:

  on **Linux**:

  ```bash
  sudo tail -f /var/log/syslog | grep hello-mms[[]
  ```

  on **Mac**:

  ```bash
  sudo docker logs -f $(sudo docker ps -q --filter name=hello-mms)
  ```

5. While observing the output, in another terminal open the `object.json` file and change the `destinationID` value to your node id.

6. Publish the `input.json` file as a new mms object:
```
make publish-mms-object
```
7. View the published mms object:
```bash
make list-mms-object
```

8. You should now see the output of the hello-mms service change from `<your-node-id> says: Hello World!!` to `<your-node-id> says: Hello Everyone!!`


9. Unregister your edge node (which will also stop the hello-mms service):

```bash
hzn unregister -f
```

10. Delete the published mms object:
```bash
make delete-mms-object
```
