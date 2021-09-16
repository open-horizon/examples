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

3. As part of the infrastructure installation process for Horizon, a file called `agent-install.cfg` was created that contains the values for `HZN_ORG_ID` and the exchange and css URLs. Locate this file and set those environment variables in your shell now:

  ```bash
  eval export $(cat agent-install.cfg)
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

6. If you have not done so already, unregister your node before moving on:

  ```bash
  hzn unregister -f
  ```

## <a id=using-hello-mms-pattern></a> Using the Hello MMS Example Edge Service with Deployment Pattern

1. Register your edge node with Horizon to use the hello-mms pattern:

  ```bash
  export SERVICE_NAME=ibm.hello-mms
  export HZN_PATTERN=pattern-$SERVICE_NAME
  hzn register -p IBM/$HZN_PATTERN -s $SERVICE_NAME --serviceorg IBM
  ```

2. After the service has started, list the docker containers to see it:

  ``` bash
  sudo docker ps
  ```

3. **Open another terminal** and view the hello-mms service output with the "follow" flag. This sample service repeatedly checks for an updated config file (`config.json`) that contains a parameter of who it should say "hello" to. The initial config file is built into the docker image. Updated config files can come via MMS. Initially you should see the message like: **<your-node-id> says: Hello from the dockerfile!** .

  ```bash
  export SERVICE_NAME=ibm.hello-mms
  hzn service log -f $SERVICE_NAME
  ```

4. While observing the service output in the other terminal, get the sample files that will be needed to create and publish a new config file in MMS:

  ```bash
  wget -q --show-progress https://github.com/open-horizon/examples/raw/master/edge/services/helloMMS/object.json
  wget -q --show-progress https://github.com/open-horizon/examples/raw/master/edge/services/helloMMS/config.json
  ```

5. Modify `config.json` and publish it (along with its metadata `object.json`) as a new mms object:

  ```bash
  jq '.HW_WHO = "from the MMS"' config.json > config.tmp && mv config.tmp config.json
  export HZN_DEVICE_ID="${HZN_EXCHANGE_NODE_AUTH%%:*}"   # this env var is referenced in object.json
  hzn mms object publish -m object.json -f config.json
  ```

6. View the published mms object:

  ```bash
  hzn mms object list -t $HZN_DEVICE_ID.hello-mms -i config.json -d
  ```

  Once the `Object status` changes to `delivered` you will see the output of the hello-mms service (in the other terminal) change from **\<your-node-id\> says: Hello from the dockerfile!** to **\<your-node-id\> says: Hello from the MMS!**

7. Delete the published mms object:

  ```bash
  hzn mms object delete -t $HZN_DEVICE_ID.hello-mms --id config.json
  ```

  **Note**: in the service output in the other terminal that this will cause the service to revert to the original config file, and therefore the original "hello" message.

8. Unregister your edge node (which will also stop the hello-mms service):

  ```bash
  hzn unregister -f
  ```

## <a id=mms-deets></a> More MSS Details

The `hzn mms ...` command provides additional tooling for working with the MMS. Get  help for this command with:

```bash
hzn mms --help
```

A good place to start is with the `hzn mms object new` command, which will emit an MMS object metadata template. You can take this template and fill in the fields that are relevant to your use case and then use it to publish your MMS object.

You can view all of the MMS objects that are used with a particular pattern like this:

```bash
hzn mms object list --destinationType pattern-ibm.hello-mms -d
```

To view the current MMS status, use: `hzn mms status`
