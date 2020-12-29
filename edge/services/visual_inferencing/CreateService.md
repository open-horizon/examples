# Creating Your Own Object Detection and Classification Edge Service

Follow the steps in this page to create your own CPU To IBM Event Streams Edge Service.

## Preconditions for Developing Your Own Service

1. If you have not already done so, complete the steps in these sections:

  - [Preconditions for Using the Object Detection and Classification Example Edge Service](UsingPolicy.md#preconditions)
  - [Using the Object Detection and Classification Example Edge Service with Deployment Policy](UsingPolicy.md#-using-the-object-detection-and-classification-example-edge-service-with-deployment-policy)

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

## <a id=build-publish-your-yolo> Building and Publishing Your Own Version of the Object Detection and Classification Edge Service


1. Clone this git repo:
  ```bash
  cd ~   # or wherever you want
  git clone git@github.com:open-horizon/examples.git
  ```

2. Copy either the `yolocpu` or `yolocuda` dir to where you will start development of your new service:
- if your edge device **does not** have a GPU, run the following commands:
  ```bash
  cp -a examples/edge/services/visual_inferencing/yolocpu ~/myservice     # or wherever
  cd ~/myservice
  ```
- if your edge device **does** have a GPU, run the following commands:
  ```bash
  cp -a examples/edge/services/visual_inferencing/yolocuda ~/myservice     # or wherever
  cd ~/myservice
  ```

3. Set the values in `horizon/hzn.json` to your own values.

4. Edit the service however you want.

5. Build the yolo docker image:
  ```bash
  make
  ```

6. Test the service by having Horizon start it locally:
  ```bash
  hzn dev service start -S
  ```

7. Check that the container is running:
  ```bash
  sudo docker ps 
  ```

8. Navigate to http://0.0.0.0:5200 to confirm the object detection and classification is working as intended (this can take a couple minutes).

9. See the environment variables Horizon passes into your service container:
  ```bash
  docker inspect $(docker ps -q --filter name=cpu2evtstreams) | jq '.[0].Config.Env'
  ```

10. Stop the service:
  ```bash
  hzn dev service stop
  ```

11. Have Horizon push your docker image to your registry and publish your service in the Horizon Exchange and see it there:
  ```bash
  hzn exchange service publish -f horizon/service.definition.json
  hzn exchange service list
  ```

## Publishing A Policy For Your Object Detection and Classification Example Edge Service

1. Set the required environment variables:
  ```bash
  eval $(hzn util configconv -f horizon/hzn.json)
  export ARCH=$(hzn architecture)
  ```

2. Publish and view your service policy in the Horizon Exchange:
  ```bash
  hzn exchange service addpolicy -f horizon/service.policy.json ${HZN_ORG_ID}/${SERVICE_NAME}_${SERVICE_VERSION}_${ARCH}
  hzn exchange service listpolicy ${HZN_ORG_ID}/${SERVICE_NAME}_${SERVICE_VERSION}_${ARCH}
  ```

3. Publish and view your deployment policy in the Horizon Exchange:
  ```bash
  hzn exchange deployment addpolicy -f horizon/deployment.policy.json ${HZN_ORG_ID}/policy-${SERVICE_NAME}_${SERVICE_VERSION}
  hzn exchange deployment listpolicy ${HZN_ORG_ID}/policy-${SERVICE_NAME}_${SERVICE_VERSION}
  ```

4. Register your edge device with the node policy:
  ```bash
  hzn register --policy horizon/node.policy.json
  ```
  
5. The edge device will make an agreement with one of the Horizon agreement bots (this typically takes about 15 seconds). Repeatedly query the agreements of this device until the `agreement_finalized_time` and `agreement_execution_start_time` fields are filled in:

  ```bash
  hzn agreement list
  ```
  
6. After the agreement is made, list the docker container edge service that has been started as a result:

  ```bash
  sudo docker ps
  ```

7. Navigate to http://0.0.0.0:5200 to confirm the object detection and classification is working as intended (this can take a couple minutes).

8. Unregister your edge device (which will also stop the myhelloworld service):

  ```bash
  hzn unregister -f
  ```
  
## Publishing A Pattern For Your Object Detection and Classification Example Edge Service

1. Publish and view your edge node deployment pattern in the Horizon Exchange:

  ```bash
  hzn exchange pattern publish -f horizon/pattern.json
  hzn exchange pattern list
  ```

2. Register your edge node with Horizon to use your deployment pattern (substitute `<service-name>` for the `SERVICE_NAME` you specified in `horizon/hzn.json`):
  ```bash
  hzn register -p pattern-<service-name>-$(hzn architecture) -f horizon/userinput.json
  ```

3. The edge device will make an agreement with one of the Horizon agreement bots (this typically takes about 15 seconds). Repeatedly query the agreements of this device until the `agreement_finalized_time` and `agreement_execution_start_time` fields are filled in:
  ```bash
  hzn agreement list
  ```

4. Once the agreement is made, list the docker container edge service that has been started as a result:
  ```bash
  sudo docker ps
  ```
5. Navigate to http://0.0.0.0:5200 to confirm the object detection and classification is working as intended (this can take a couple minutes).

6. Unregister your edge node, stopping the yolo service:
```bash
hzn unregister -f
```
