# Horizon Hello World Example Edge Service

This is a simple example of using and creating a Horizon edge service.

- [Preconditions for Using the Hello World Example Edge Service](#preconditions)
- [Using the Hello World Example Edge Service with Deployment Pattern](#using-helloworld-pattern)
- [Using the Hello World Example Edge Service with Deployment Policy](#using-helloworld-policy)
- [Creating Your Own Hello World Edge Service](CreateService.md)
- Further Learning - to see more Horizon features demonstrated, continue on to the [cpu2evtstreams example](../../evtstreams/cpu2evtstreams).

## <a id=preconditions></a> Preconditions for Using the Hello World Example Edge Service

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

## <a id=using-helloworld-pattern></a> Using the Hello World Example Edge Service with Deployment Pattern

1. Register your edge node with Horizon to use the helloworld pattern:

```bash
hzn register -p IBM/pattern-ibm.helloworld
```

2. The edge device will make an agreement with one of the Horizon agreement bots (this typically takes about 15 seconds). Repeatedly query the agreements of this device until the `agreement_finalized_time` and `agreement_execution_start_time` fields are filled in:

```bash
hzn agreement list
```

3. After the agreement is made, list the docker container edge service that has been started as a result:

``` bash
sudo docker ps
```

4. See the helloworld service output:

``` bash
hzn service log -f ibm.helloworld
```

5. Unregister your edge node (which will also stop the myhelloworld service):

```bash
hzn unregister -f
```
