# Creating Your Own Model Management Service (MMS) Example

Follow the steps in this page to create your own Horizon MMS example service.

## Preconditions for Developing Your Own Service

1. If you have not already done so, complete the steps in these sections:

  - [Preconditions for Using the MMSExample Edge Service](README.md#preconditions)
  - [Using the MMS Example Edge Service with Deployment Pattern](README.md#using-MMS-pattern)

2. If you are using macOS as your development host, configure Docker to store credentials in `~/.docker`:

  - Open the Docker **Preferences** dialog
  - Uncheck **Securely store Docker logins in macOS keychain**

3. If you do not already have a docker ID, obtain one at https://hub.docker.com/ . Log in to Docker Hub using your Docker Hub ID:

  ```bash
  export DOCKER_HUB_ID="<dockerhubid>"
  echo "<dockerhubpassword>" | docker login -u $DOCKER_HUB_ID --password-stdin
  ```

  Output example:

  ```bash
  WARNING! Your password will be stored unencrypted in /home/pi/.docker/config.json.
  Configure a credential helper to remove this warning. See
  https://docs.docker.com/engine/reference/commandline/login/#credentials-store

  Login Succeeded
  ```

4. Create a cryptographic signing key pair. This enables you to sign services when publishing them to the exchange. **This step only needs to be done once.**

  ```bash
  hzn key create "<x509-org>" "<x509-cn>"
  ```

  where `<x509-org>` is your company name, and `<x509-cn>` is typically set to your email address.

5. Install `git` and `jq`:

  On **Linux**:

  ```bash
  sudo apt install -y git jq
  ```

  On **macOS**:

  ```bash
  brew install git jq
  ```

## <a id=build-publish-your-hw> Building and Publishing Your Own Hello World Example Edge Service


## Building and Publishing Your Own Version of the CPU To IBM Event Streams Edge Service
- Use the developer tool to run the container with a local development instance of the Model Management Service (MMS). Normally, in production, you will use the MMS in the IBM Public Cloud, or ICP, but during development it is convenient to have a dedicated and private "dev MMS" instance you can use. So we will show that approach here first.

1. Clone this git repo:
```
cd ~   # or wherever you want
git clone git@github.com:open-horizon/examples.git
```

2. Copy the `mms` dir to where you will start development of your new service:
```
cp -a examples/edge/evtstreams/mms ~/myservice     # or wherever
cd ~/myservice
```

3. Set the values in `horizon/hzn.json` to your own values.

4. Edit `service.sh` however you want.
    - Note: this service is a shell script simply for brevity, but you can write your service in any language.

5. Build the service docker image:

```bash
make
```

6. Test the service by running it the simulated agent environment:

```bash
hzn dev service start
```

7. Check that the container is running:

```bash
sudo docker ps
```

8. Display the environment variables Horizon passes into your service container:

```bash
sudo docker inspect $(sudo docker ps -q --filter name=mms) | jq '.[0].Config.Env'
```

9. See the docker container running and look at the output:

  on **Linux**:

  ```bash
  sudo tail -f /var/log/syslog | grep mms[[]
  ```

  on **Mac**:

  ```bash
  sudo docker logs -f $(sudo docker ps -q --filter name=mms)
  ```

You should see something similar to the following, that is, the output should identify your Edge Node, and the message should be, "**Hello!**":

```bash
Jun  7 16:04:01 myedgenode0 workload-c9ef49dbf715f1477f72c001eb3933970690bea96c4d486a7fc60a686843fcd1_ibm.mms[823]: myedgenode0.dev.edge-fabric.com says: "Hello!"
Jun  7 16:04:04 myedgenode0 workload-c9ef49dbf715f1477f72c001eb3933970690bea96c4d486a7fc60a686843fcd1_ibm.mms[823]: myedgenode0.dev.edge-fabric.com says: "Hello!"
...
```

10. Use the "dev MMS" to send something through the MMS to the Service container running on the Edge Node. In a **host**  shell, run:

```
echo 'Goodbye!' | ./dev-css-write.sh example-type id-0
```

11. Observe the change in the mms container output:

  on **Linux**:

  ```bash
  sudo tail -f /var/log/syslog | grep mms[[]
  ```

  on **Mac**:

  ```bash
  sudo docker logs -f $(sudo docker ps -q --filter name=mms)
  ```

```bash
Jun  7 16:04:17 myedgenode0 workload-c9ef49dbf715f1477f72c001eb3933970690bea96c4d486a7fc60a686843fcd1_ibm.mms[823]: myedgenode0.dev.edge-fabric.com says: "Hello!"
Jun  7 16:04:20 myedgenode0 workload-c9ef49dbf715f1477f72c001eb3933970690bea96c4d486a7fc60a686843fcd1_ibm.mms[823]: myedgenode0.dev.edge-fabric.com says: "Hello!"
Jun  7 16:04:23 myedgenode0 workload-c9ef49dbf715f1477f72c001eb3933970690bea96c4d486a7fc60a686843fcd1_ibm.mms[823]: myedgenode0.dev.edge-fabric.com says: ""Goodbye!""
Jun  7 16:04:26 myedgenode0 workload-c9ef49dbf715f1477f72c001eb3933970690bea96c4d486a7fc60a686843fcd1_ibm.mms[823]: myedgenode0.dev.edge-fabric.com says: ""Goodbye!""
```

- Notice the the message changed to "**Goodbye!**".

12. You can send other messages and watch the updated versions being picked up. E.g.:

```bash
echo 'Something Random' | ./dev-css-write.sh example-type id-0
echo 'Rubber Duck' | ./dev-css-write.sh example-type whatever-you-like-here
```

13. Stop the service:

```bash
hzn dev service stop
```


14. Instruct Horizon to push your docker image to your registry and publish your service in the Horizon Exchange:

```bash
hzn exchange service publish -f horizon/service.definition.json
hzn exchange service list
```

15. Publish and view your edge node deployment pattern in the Horizon Exchange:

```bash
hzn exchange pattern publish -f horizon/pattern.json
hzn exchange pattern list
```

16. Register your edge node with Horizon to use your deployment pattern (substitute `<service-name>` for the `SERVICE_NAME` you specified in `horizon/hzn.json`):
```
hzn register -p pattern-<service-name>-$(hzn architecture) -f horizon/userinput.json
```

17. The edge device will make an agreement with one of the Horizon agreement bots (this typically takes about 15 seconds). Repeatedly query the agreements of this device until the `agreement_finalized_time` and `agreement_execution_start_time` fields are filled in:

  ```bash
  hzn agreement list
  ```

18. After the agreement is made, list the docker container edge service that has been started as a result:

  ```bash
  sudo docker ps
  ```

19. See the mms service output (you should see the, "**Hello!**" message as before):

  on **Linux**:

  ```bash
  sudo tail -f /var/log/syslog | grep mms[[]
  ```

  on **Mac**:

  ```bash
  sudo docker logs -f $(sudo docker ps -q --filter name=mms)
  ```
20. Use the production MMS to send a new message to your Service:

```bash
echo 'Goodbye!' | ./prod-css-write.sh example-type id-0
```

21. Again, observe the `mms` Service output (to see the message change to, "**Goodbye!**" as it did during development):

  on **Linux**:

  ```bash
  sudo tail -f /var/log/syslog | grep mms[[]
  ```

  on **Mac**:

  ```bash
  sudo docker logs -f $(sudo docker ps -q --filter name=mms)
  ```
- Be aware that if you send things in rapid succession using different IDs, they may arrive out of order.

22. Unregister your edge node, stopping the mms service:

```bash
hzn unregister -f
```


