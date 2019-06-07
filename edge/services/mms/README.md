# Horizon Model Management Service (MMS) Example

## Getting Ready

- It is assumed you have gone through the developer workflow for at least one of the other Horizon examples on this host, so you have verified your credentials are configured in the environment of this shell, you are logged-in to DockerHub, and you have created your cryptographic signing key pair, etc.

- Clone this git repo:

```
git clone git@github.com:open-horizon/examples.git
```

- Enter the `mms` example directory

```
cd  examples/edge/services/mms
```

- Set the values in `horizon/hzn.json` to your own values.
- Add those values into your shell environment.

```
`hzn util configconv -f horizon/hzn.json`
```

- Build and appropriately tag this example Docker container

```
docker build -t "${DOCKER_IMAGE_BASE}_$(hzn architecture):${SERVICE_VERSION}" .
```

## Running In Development Mode

- Use the developer tool to run the container with a local development instance of the Cloud Sync Server (CSS). Normally, in production, you will use the CSS in the IBM Cloud, but during development it is convenient to have a dedicated and private "dev CSS" instance you can use. So we will show that approach here first.

```
hzn dev service start
```

- Observe the `mms` Service example output (keep this running in a separate terminal so you can watch it as it changes a bit later):

```
# Soon you will use 'hzn service log ...' for all platforms
# For now on Linux:
tail -f /var/log/syslog | grep mms[[]
# For now on Mac:
docker logs -f $(docker ps -q --filter name=mms)
``` 

You should see something similar to this:

```
Jun  7 16:04:01 myedgenode0 workload-c9ef49dbf715f1477f72c001eb3933970690bea96c4d486a7fc60a686843fcd1_ibm.mms[823]: myedgenode0.dev.edge-fabric.com says: "Hello!"
Jun  7 16:04:04 myedgenode0 workload-c9ef49dbf715f1477f72c001eb3933970690bea96c4d486a7fc60a686843fcd1_ibm.mms[823]: myedgenode0.dev.edge-fabric.com says: "Hello!"
...
```

That is, the output should identify your Edge Node, and the message should be, "**Hello!**". This is how the Service is initially configured. Now let's use the "dev CSS" to send something to the Edge Sync Service (ESS) that the container is monitoring. In a **host**  shell, run this:

```
echo 'Goodbye!' | ./dev-css-write.sh example-type id-0
```

- Observe the change in the mms Service example output:

```
Jun  7 16:04:17 myedgenode0 workload-c9ef49dbf715f1477f72c001eb3933970690bea96c4d486a7fc60a686843fcd1_ibm.mms[823]: myedgenode0.dev.edge-fabric.com says: "Hello!"
Jun  7 16:04:20 myedgenode0 workload-c9ef49dbf715f1477f72c001eb3933970690bea96c4d486a7fc60a686843fcd1_ibm.mms[823]: myedgenode0.dev.edge-fabric.com says: "Hello!"
Jun  7 16:04:23 myedgenode0 workload-c9ef49dbf715f1477f72c001eb3933970690bea96c4d486a7fc60a686843fcd1_ibm.mms[823]: myedgenode0.dev.edge-fabric.com says: ""Goodbye!""
Jun  7 16:04:26 myedgenode0 workload-c9ef49dbf715f1477f72c001eb3933970690bea96c4d486a7fc60a686843fcd1_ibm.mms[823]: myedgenode0.dev.edge-fabric.com says: ""Goodbye!""
```

- Notice the the message changed to "**Goodbye!**".
- You can send other messages and watch them being picked up. Just make be to change the object id from `id-0` above to something else each time you send. E.g.:

```
echo 'Something Random' | ./dev-css-write.sh example-type id-1
echo 'Rubber Duck' | ./dev-css-write.sh example-type whatever-you-like-here
```

- Stop the service:

```
hzn dev service stop
```

## Publishing In Preparation For Registration

- Have Horizon push your docker image to your registry and use your signing key to publish your service in the Horizon Exchange and see it there:

```
hzn exchange service publish -f horizon/service.definition.json
hzn exchange service list
```

- Publish your edge node deployment pattern in the Horizon Exchange and see it there:

```
hzn exchange pattern publish -f horizon/pattern.json
hzn exchange pattern list
```

## Running In Production Mode

- Register your edge node with Horizon to use your deployment pattern:

```
hzn register -p pattern-${SERVICE_NAME}-$(hzn architecture)
```

- Repeatedly list the horizon agreements until one is finalized and then list the running containers and see the `mms` example container running:

```
hzn agreement list
docker ps
```

- Monitor the `mms` Service output (you should see the, "**Hello!**" message as before):

```
# soon you will use 'hzn service log ...' for all platforms
# for now on linux:
tail -f /var/log/syslog | grep mms[[]
# for now on mac:
docker logs -f $(docker ps -q --filter name=mms)
``` 

- While watching the output logs from the container, use the production CSS to send a new message to your Service:

```
echo 'Goodbye!' | ./prod-css-write.sh example-type id-0
```

- Again, observe the `mms` Service output (to see the message change to, "**Goodbye!**" as it did during development):
- As before, you can send again, just be sure to change the object id from `id-0` above to something new each time you send.
- Be aware that if you send things in rapid succession, they may arrive out of order.
- Unregister your edge node, stopping the mms service:

```
hzn unregister -f
```

## Some MSS Usage Tips

Cover stuff here, such as:
- `hzn mms object list -t <TYPE> -i <ID>`
- go over the code in the scripts


