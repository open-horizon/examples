# Horizon Hello Model Management Service (MMS) Example Edge Service

This is a simple example of using and creating a Horizon edge service.

- [Introduction to the Horizon Model Management Service](#introduction)
- [Preconditions for Using the Hello MMS Example Edge Service](#preconditions)
- [Using the Hello MMS Example Edge Service with Deployment Pattern](#using-hello-mms-pattern)
- [More MMS Details](#mms-deets)
- [Creating Your Own Hello MMS Edge Service](CreateService.md)

## <a id=introduction></a> Introduction

The Horizon Model Management Service (MMS) enables you to have independent lifecycles for your code and for your data. While Horizon Services, Patterns, and Policies enable you to manage the lifecycles of your code components, the MMS performs an analogous service for your data files.  This can be useful for remotely updating the configuration of your Services in the field. It can also enable you to continuously train and update of your neural network models in powerful central data centers, then dynamically push new versions of the models to your small edge machines in the field. The MMS enables you to manage the lifecycle of data files on your edge node, remotely and independently from your code updates. In general the MMS provides facilities for you to securely send any data files to and from your edge nodes.

This document will walk you through the process of using the Model Management Service to send a file to your edge nodes. It also shows how your nodes can detect the arrival of a new version of the file, and then consume the contents of the file.

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

4. See the hello-mms service output (you should see the message **<your-node-id> says: Hello World!**:

  on **Linux**:

  ```bash
  sudo tail -f /var/log/syslog | grep hello-mms[[]
  ```

  on **Mac**:

  ```bash
  sudo docker logs -f $(sudo docker ps -q --filter name=hello-mms)
  ```

5. While observing the output, in another terminal, open the `object.json` file and change the `destinationID` value to your node id.

6. Publish the `input.json` file as a new mms object:
```bash
hzn mms object publish -m object.json -f input.json
```

7. View the published mms object:
```bash
hzn mms object list -t json -i input.json -d
```

Once the `Object status` changes to `delivered` you will see the output of the hello-mms service change from **\<your-node-id\> says: Hello World!** to **\<your-node-id\> says: Hello Everyone!**

8. Delete the published mms object:
```bash
hzn mms object delete -t json --id input.json
```

9. Unregister your edge node (which will also stop the hello-mms service):

```bash
hzn unregister -f
```

## <a id=mms-deets></a> More MSS Details

The `hzn mms ...` command provides additional tooling for working with the MMS. Get  help for this command with:

```bash
hzn mms --help
```

A good place to start is with the `hzn mms object new` command, which will emit an MMS object metadata template. You can take this template, fill in the fields that are relevant to your use case, and remove all of the "comments" wrapped in `/* ... */`. Then you can pass it to the `hzn mms object publish -m <my-metadata-file` (as your `<my-metadata-file>`).

To publish an object with the MMS, you can use the scripts you used above, or the `hzn mms object publish ...` command. For the latter you need to provide `-t <my-type>` and `-i <my-id>` (passing your own type, `<my-type>`, and ID, `<my-id>`). This command also takes a `-p <my-pattern>` flag that you can use to tell the MMS to deliver this object only to Edge Nodes that are registered with Deployment Pattern `<my-pattern>.

The `hzn mms object list -t <my-type>` can be used to list all the MMS objects of type, `<my-type>`.

To delete a specific object, of type `<my-type>` with ID `<my-id>` you can use, `hzn mms object delete -t <my-type> -i <my-id>`.

To view the current MMS status, use, `hzn mms status`.
