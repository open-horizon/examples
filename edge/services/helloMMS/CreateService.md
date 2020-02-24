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

4. Edit `service.sh` however you want. For example, to make a simple change so you will be able to confirm that your new service is running, you could customize the `echo` statement near the end that says "Hello". While you are editing `service.sh`, read the comments and code to learn the basic pattern for using a MMS file in an edge service. This coding pattern will be the same, regardless of what language you implement your own edge services in.
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

  Once the `Object status` changes to `delivered` you will see the output of the hello-mms service (in the other terminal) change from **\<your-node-id\> says: Hello from the dockerfile!** to **\<your-node-id\> says: Hello from the MMS!**

13. Delete your MMS object and watch the service messages change back to the original value:

  ```bash
  HZN_FSS_CSSURL=http://localhost:8580  hzn mms object delete -t $HZN_DEVICE_ID.hello-mms -i config.json
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

20. Publish `config.json` as a new object in the cloud MMS:

  ```bash
  hzn mms object publish -m object.json -f config.json
  ```

21. View the published mms object:

  ```bash
  hzn mms object list -t $HZN_DEVICE_ID.hello-mms -i config.json -d
  ```

22. After approximately 15 seconds you should see the output of the service change to the value of `HW_WHO` set in the `config.json` file.

23. Clean up by deleting the published mms object and unregistering your edge node:

  ```bash
  hzn mms object delete -t $HZN_DEVICE_ID.hello-mms -i config.json
  hzn unregister -f
  ```

## More MMS Information

You can browse the [full MMS REST API](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/open-horizon/edge-sync-service/master/swagger.json) .

The ESS REST API (the APIs that an edge service uses) is a small subset of that. The most commonly used ESS REST APIs are:

- `GET /api/v1/objects/{objectType}` - Get metadata for objects of the specified type that have changed, but not yet been acknowledged by this edge service. (There is an optional URL parameter `?received=true` that will cause it to return all objects of this type, regardless of whether they've been acknowledged or not, but this is rarely needed.)
- `GET /api/v1/objects/{objectType}/{objectID}` - Get an object's metadata
- `PUT /api/v1/objects/{objectType}/{objectID}` - Create the metadata (specified in the body) for a new (or updated) object that this service is sending to MMS
- `GET /api/v1/objects/{objectType}/{objectID}/data` - Get the file associated with this object
- `PUT /api/v1/objects/{objectType}/{objectID}/data` - Send this file (specified in the body) to MMS
- `PUT /api/v1/objects/{objectType}/{objectID}/deleted` - Confirm that this service has seen that the object has been deleted
- `PUT /api/v1/objects/{objectType}/{objectID}/received` - Confirm that this service has seen that the object has been changed
