# Creating Your Own Hello MMS Edge Service

Follow the steps in this page to create your first Horizon edge service that uses the Model Management Service.

## Preconditions for Developing Your Own Service

1. If you have not already done so, complete the steps in these sections:

  - [Preconditions for Using the Hello MMS Example Edge Service](README.md#preconditions)
  - [Using the Hello MMS Example Edge Service with Deployment Pattern](README.md#using-hello-mms-pattern)

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

## <a id=build-publish-your-hw> Building and Publishing Your Own Hello MMS Example Edge Service

1. Clone this git repo:
```bash
cd ~   # or wherever you want
git clone git@github.com:open-horizon/examples.git
```

2. Copy the `hello-mms` dir to where you will start development of your new service:
```bash
cp -a examples/edge/services/helloMMS ~/myservice     # or wherever
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
  sudo docker inspect $(sudo docker ps -q --filter name=hello-mms) | jq '.[0].Config.Env'
  ```

9. See the hello-mms service output (you should see the message **$<your-node-id$> says: Hello World!**:

  on **Linux**:

  ```bash
  sudo tail -f /var/log/syslog | grep hello-mms[[]
  ```

  on **Mac**:

  ```bash
  sudo docker logs -f $(sudo docker ps -q --filter name=hello-mms)
  ```

10. While observing the output, in another terminal, open the `object.json` file and change the `destinationID` value to your node id.

11. Publish the `input.json` file as a new mms object:
```bash
make publish-mms-object
```

12. View the published mms object:
```bash
make list-mms-object
```


13. After approximately 15 seconds you should see the output of the hello-mms service change to the value of `HW_WHO` that is set in the `input.json` file. For instance, you may see the message change from **<your-node-id> says: Hello World!** to **<your-node-id> says: Hello Everyone!**

14. Stop the service:

  ```bash
  hzn dev service stop
  ```

15. Instruct Horizon to push your docker image to your registry and publish your service in the Horizon Exchange:

  ```bash
  hzn exchange service publish -f horizon/service.definition.json
  hzn exchange service list
  ```

16. Publish and view your edge node deployment pattern in the Horizon Exchange:

  ```bash
  hzn exchange pattern publish -f horizon/pattern.json
  hzn exchange pattern list
  ```

17. Register your edge node with Horizon to use your deployment pattern (substitute `<service-name>` for the `SERVICE_NAME` you specified in `horizon/hzn.json`):
```bash
hzn register -p pattern-<service-name>-$(hzn architecture) -f horizon/userinput.json
```

18. The edge device will make an agreement with one of the Horizon agreement bots (this typically takes about 15 seconds). Repeatedly query the agreements of this device until the `agreement_finalized_time` and `agreement_execution_start_time` fields are filled in:

  ```bash
  hzn agreement list
  ```

19. After the agreement is made, list the docker container edge service that has been started as a result:

  ```bash
  sudo docker ps
  ```

20. See the hello-mms service output:

  on **Linux**:

  ```bash
  sudo tail -f /var/log/syslog | grep hello-mms[[]
  ```

  on **Mac**:

  ```bash
  sudo docker logs -f $(sudo docker ps -q --filter name=hello-mms)
  ```

21. While observing the output, in another terminal open the `input.json` file and change the `"value":` field to whatever you want.

22. Publish the `input.json` file as a new mms object:
```bash
make publish-mms-object
```

23. After approximately 15 seconds you should see the output of the hello-mms service change to the value of `HW_WHO` set in the `input.json` file.

24. Delete the published mms object:
```bash
make delete-mms-object
```

25. Unregister your edge node (which will also stop the hello-mms service):

  ```bash
  hzn unregister -f
  ```

