# Creating Your Own Hello Secret World Edge Service

Follow the steps in this page to create your first OpenHorizon edge service that exploits secrets.

## Preconditions for Developing Your Own Service

1. If you have not already done so, complete the steps in these sections:

  - [Preconditions for Using the Hello Secret World Example Edge Service](README.md#preconditions)
  - [Using the Hello Secret World Example Edge Service with Deployment Pattern](README.md#using-hello-secret-pattern)

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

## <a id=build-publish-your-hw></a> Building and Publishing Your Own Hello Secret World Example Edge Service

1. Clone this git repo:

  ```bash
  cd ~   # or wherever you want
  git clone git@github.com:open-horizon/examples.git
  ```

2. Check your Horizon CLI version:

  ```bash
  hzn version
  ```

3. Starting with Horizon version `v2.29.0-595` you can `checkout` to a version of the example services that directly corresponds to your Horizon CLI version with these commands: 

  ```bash
  export EXAMPLES_REPO_TAG="v$(hzn version 2>/dev/null | grep 'Horizon CLI' | awk '{print $4}')"
  git checkout tags/$EXAMPLES_REPO_TAG -b $EXAMPLES_REPO_TAG
  ```

4. Copy the `helloSecretWorld` dir to where you will start development of your new service:

  ```bash
  cp -a examples/edge/services/helloSecretWorld ~/myservice     # or wherever
  cd ~/myservice
  ```

5. Set the values in `horizon/hzn.json` to your liking. These variables are used in the service and pattern files in `horizon`. They are also used in some of the commands in this procedure. After editing `horizon/hzn.json`, set the variables in your environment:

  ```bash
  eval $(hzn util configconv -f horizon/hzn.json)
  ```

6. Edit `service.sh` however you want. For example, to make a simple change so you will be able to confirm that your new service is running, you could customize the `echo` statement near the end that says "Hello". While you are editing `service.sh`, read the comments and code to learn the basic pattern for using a secret in an edge service. This coding pattern will be the same, regardless of what language you implement your own edge services in.
    - Note: this service is a shell script simply for brevity, but you can write your service in any language.

7. Build the service docker image.

  ```bash
  make
  ```

8. Edit the secret value in `horizon/hw_who` to whatever you want. The secret key is ignored by the code in `service.sh`, but the code will read the secret value from the file and output it using the `echo` statement. Because we will be testing our service with `hzn dev service start` without a deployment pattern or policy, the agent cannot get the secret details of any secrets in the secrets manager. Instead, it will read this file as the secret details of the secret used by the service. 

9. Test the service by running it with the simulated agent environment. Pass in the **absolute path** to the `hw_who` file, which is in the `horizon` subdirectory. 

  ```bash
  hzn dev service start --secret=<abs-path-to-hw_who>
  ```

Upon starting the service, this command will mount the local `hw_who` file containing the contents of the secret to the service container in the `open-horizon-secrets` folder. If this file is changed while the service is running, the change will be reflected in the container as well.

10. Check that the container is running:

  ```bash
  sudo docker ps
  ```

11. Display the environment variables OpenHorizon passes into your service container. Note the variables that start with `HZN_ESS_`. These are used by the service to query the secrets API for updating secrets.

  ```bash
  sudo docker inspect $(sudo docker ps -q --filter name=$SERVICE_NAME) | jq '.[0].Config.Env'
  ```

12. View the service output (you should see messages like **\<your-node-id\> says: Hello \<secret-value\>!**, where \
\<secret-value\> is the value of the secret as stated in the `hw_who` file.

  ```bash
  hzn dev service log $SERVICE_NAME
  ```

13. Stop the service:

  ```bash
  hzn dev service stop
  ```

14. You are now ready to publish your edge service and pattern, so that they can be deployed to real edge nodes. Instruct OpenHorizon to push your docker image to your registry and publish your service in the OpenHorizon Exchange:

  ```bash
  hzn exchange service publish -f horizon/service.definition.json
  hzn exchange service list
  ```

15. (This step must be done by an organization admin) Create an organization-wide secret called `hw-secret-name` in your organization's secrets manager. The flags `--secretKey` and `--secretDetail` provide the key and value for the secret, respectively.

  ```bash
  hzn sm secret add --secretKey <secret-key> --secretDetail <secret-value> hw-secret-name
  ```

16. Edit your pattern definition file to make the pattern not public, then publish your edge node deployment pattern in the OpenHorizon Exchange:

  ```bash
  jq '.public = false' horizon/pattern.json > horizon/pattern.tmp && mv horizon/pattern.tmp horizon/pattern.json
  hzn exchange pattern publish -f horizon/pattern.json
  hzn exchange pattern list
  ```
OpenHorizon will check for the existence of `hw-secret-name` in your organization's secrets manager. If it is not found in the secrets manager, OpenHorizon will not be able to publish the pattern to the exchange.

17. Register your edge node with OpenHorizon to use your deployment pattern:

  ```bash
  hzn register -p pattern-${SERVICE_NAME}-$(hzn architecture) -s $SERVICE_NAME --serviceorg $HZN_ORG_ID
  ```

18. **Open another terminal** and view the service output with the "follow" flag:

  ```bash
  hzn service log -f $SERVICE_NAME
  ```
  
19. (This step must be done by an organization admin) Update the secret `hw-secret-name` in the secrets manager using the OpenHorizon CLI. To see any change reflected in the output, the argument passed to `--secretDetail` should be different from the initial secret creation command, so that the value of the secret can change.

  ```
  hzn sm secret add --secretKey <secret-key> --secretDetail <new-secret-value> hw-secret-name
  ```

A prompt will pop up in the terminal asking you to confirm that you want to overwrite the existing `hw-secret-name` in the secrets manager. Enter `y` to confirm.

20. After some time your service should be able to observe the change in the secret as made in Step 17.

21. Clean up by unregistering your edge node:

  ```bash
  hzn unregister -f
  ```

## Adding Secrets to a Service Definition

The service definition is updated to include a `"secrets"` field within the deployment configuration of a container, as follows:

```
"deployment": {
        "services": {
            "<container-name>": {
                "image": "<image-name>",
                "secrets": {
                    "<secret-name>": {
                        "description": "<optional-description>"
                    }
                }
            }
        }
    }
```

Any secret name used by a service must be unique within the service definition. When a secret name is referred to more than once within a service definition, all of those usages refer to the same secret in the service. In addition, all secrets must be associated with at least 1 container in the service definition. 

## Adding Secret Associations to a Deployment Pattern/Policy

Secret associations can be made in patterns and deployment policies. When these are used to make an agreement with a node, the secret name used by the service is associated with the name used by the secrets manager that is provided in the pattern or deployment policy. 

To make these secret associations, the following JSON snippet is added to the patterns and deployment policies:
```
"secretBinding": [
    {
        "serviceUrl": "<service-URL>",
        "serviceOrgid": "<service-org-id>",
        "serviceArch": "<service-hardware-architecture>,
        "serviceVersionRange": "x.y.z", 
        "secrets": [
            {"<service-secret-name>": "<secrets-manager-secret-name>"}
        ]
    }
]
```

An example of a service definition, pattern, and deployment policy with secrets are provided in the `horizon` folder of the example service. 

When a deployment pattern or policy that contains a secret association is published to the Exchange, OpenHorizon will check for the existence of all secret names stated in the pattern/policy to be associated to a service's secret name. If a secret name is not found in the secrets manager, OpenHorizon will not be able to publish the pattern or policy to the exchange.

## Using Secrets in the Service Code

### Accessing Secret Details

When an agreement is formed and a service is deployed to a node, a folder called `open-horizon-secrets` is mounted to the service's container. This folder will contain one file for each secret used by a service, and the names of these files with match the name of the secret as specified in the service definition. For example, the secret details for a service secret `password` will be in `open-horizon-secrets/password`. These files will contain the secret details encoded in JSON format as follows:
```
{
    "key": "<secret-key>",
    "value: "<secret-value>"
}
```
This file is automatically updated if the associated secret in the secrets manager, as specified in the pattern or deployment policy, has its contents updated.

### Updating Secrets

The agent hosts a REST API to notify services when a secret is updated in the secrets manager by an administrator. This API allows the service implementation to query for updated secrets and receive the new details when it is ready to do so. 

The full documentation of this API can be accessed [here](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/open-horizon/anax/master/resource/swagger.json). The example service also includes examples of calls to this API for reference.

Each service is provided with an authentication ID and token to access the ESS API. These same credentials are used to authenticate to this API. 

Note: When testing with `hzn dev service start`, this API will always return 404 since secret updates will not be received after the service has started. Instead, the agent will use the local secret file passed to the `start` command (e.g. `hzn dev service start --secret=<path-to-secret>`) to mount the secret to `open-horizon-secrets` in the service container. When the local file is updated, the file in the container is updated as well.


