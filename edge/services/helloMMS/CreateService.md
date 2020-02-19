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

5. Install `git`, `jq`, and `make`:

  On **Linux**:

  ```bash
  sudo apt install -y git jq make
  ```

  On **macOS**:

  ```bash
  brew install git jq make
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

3. Set the values in `horizon/hzn.json` to your liking. These variables are used in the service and pattern files in `horizon` and in the MMS metadata file `object.json`. They are also used in some of the commands in this procedure. After editing `horizon/hzn.json`, set the variables in your environment:

  ```bash
  eval $(hzn util configconv -f horizon/hzn.json)
  ```

4. Edit `service.sh` however you want. For example, for now to be able to confirm that you are running your own service, you could customize the `echo` statement near the end that says "Hello".
    - Note: this service is a shell script simply for brevity, but you can write your service in any language.

5. Build the service docker image. Note that the Dockerfiles copy `config.json` into the service image for it to initially use.

  ```bash
  make
  ```

6. Test the service by running it the simulated agent environment. (`HZN_PATTERN` is set so the simulated environment can find MMS object in subsequent steps.)

  ```bash
  export HZN_PATTERN=pattern-${SERVICE_NAME}-$(hzn architecture)
  hzn dev service start
  ```

7. Check that the container is running:

  ```bash
  sudo docker ps
  ```

8. Display the environment variables Horizon passes into your service container. Note the variables that start with `HZN_ESS_`. These are used by the service to contact the local MMS proxy.

  ```bash
  sudo docker inspect $(sudo docker ps -q --filter name=$SERVICE_NAME) | jq '.[0].Config.Env'
  ```

9. View the service output (you should see messages like **\<your-node-id\> says: Hello from the dockerfile!**:

  on **Linux**:

  ```bash
  sudo tail -f /var/log/syslog | grep ${SERVICE_NAME}[[]
  ```

  on **Mac**:

  ```bash
  sudo docker logs -f $(sudo docker ps -q --filter name=$SERVICE_NAME)
  ```

10. While observing the output in this terminal, **open another terminal** in the same directory to perform the next several steps. First, set the `horizon/hzn.json` variable values in this environment too:

  ```bash
  eval $(hzn util configconv -f horizon/hzn.json)
  ```

11. Modify `config.json` and publish it as a new mms object, using the provided `object.json` metadata. Since you are running in the local simulated agent environment right now, the `hzn mms ...` commands must be directed to the local MMS.

  ```bash
  jq '.HW_WHO = "from the MMS"' config.json > config.tmp && mv config.tmp config.json
  export HZN_DEVICE_ID="${HZN_EXCHANGE_NODE_AUTH%%:*}"   # this env var is referenced in object.json
  HZN_FSS_CSSURL=http://localhost:8580  hzn mms object publish -m object.json -f config.json
  ```

12. View the published mms object:

  ```bash
  HZN_FSS_CSSURL=http://localhost:8580  hzn mms object list -d
  ```

13. After approximately 15 seconds you should see the output of the service change to the value of `HW_WHO` that is now set in the `config.json` file, for example **\<your-node-id\> says: Hello from the MMS!** Now delete your MMS object and watch the service messages change back to the original value:

  ```bash
  HZN_FSS_CSSURL=http://localhost:8580  hzn mms object delete -t $SERVICE_NAME -i config.json
  ```

14. Stop the service:

  ```bash
  hzn dev service stop
  ```

15. You are now ready to publish your edge service and pattern, so that they can be deployed to real edge nodes. Instruct Horizon to push your docker image to your registry and publish your service in the Horizon Exchange:

  ```bash
  hzn exchange service publish -f horizon/service.definition.json
  hzn exchange service list
  ```

16. Edit your pattern definition file to make the pattern not public, then publish your edge node deployment pattern in the Horizon Exchange:

  ```bash
  jq '.public = false' horizon/pattern.json > horizon/pattern.tmp && mv horizon/pattern.tmp horizon/pattern.json
  hzn exchange pattern publish -f horizon/pattern.json
  hzn exchange pattern list
  ```

17. Register your edge node with Horizon to use your deployment pattern:

  ```bash
  hzn register -p pattern-${SERVICE_NAME}-$(hzn architecture) -s $SERVICE_NAME --serviceorg $HZN_ORG_ID
  ```

18. View the service output with the "follow" flag:

  ```bash
  hzn service log -f $SERVICE_NAME
  ```

19. While watching the output, switch back to your **other terminal** for the remainder of the steps.

20. Edit `config.json` and change the value associated with the `HW_WHO` field to some other value, for example **"from the MMS"**. Then publish it as a new object in the cloud MMS:

  ```bash
  hzn mms object publish -m object.json -f config.json
  ```

21. After approximately 15 seconds you should see the output of the service change to the value of `HW_WHO` set in the `config.json` file.

22. Clean up by deleting the published mms object and unregistering your edge node:

  ```bash
  hzn mms object delete -t $SERVICE_NAME -i config.json
  hzn unregister -f
  ```
