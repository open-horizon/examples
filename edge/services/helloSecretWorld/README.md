# OpenHorizon Hello Secret World Example Edge Service

This is a simple example of using and creating an OpenHorizon edge service that uses secrets in the secrets manager.

- [Introduction to the OpenHorizon Secrets Manager](#introduction)
- [Preconditions for Using the Hello Secret World Example Edge Service](#preconditions)
- [Using the Hello Secret World Example Edge Service with a Deployment Pattern](#using-hello-secret-pattern)
- [Secrets Manager Access Control and CLI](#secret-info)
- [Creating Your Own Hello Secret World Edge Service](CreateService.md)

## <a id=introduction></a> Introduction

Services deployed to the edge often require access to cloud services, which means the edge service needs credentials to authenticate to the cloud service. OpenHorizon provides secrure deployment of these credentials without exposing them within application metadata (e.g. service definitions, policies, configuration files), or to other users in the system that should not have access to the credentials.

The OpenHorizon secrets manager provides service developers with access to protected resources, such as an RSA key, certificate, or user ID/password, for use in their services. It prevents protected resources from being accessed by other users in the organization who should not be able to read them directly, but still allows service developers to use them in their services.

A secret consists of a key and value pair. They are stored in the secrets manager on the management hub and have a unique name that is used to identify the secret in the manager. This name is org-unique and should give no hints about the contents of the secret. The secrets manager is administered by an organization administrator authenticated as such by the Exchange.

When a service developer needs a secret, they can come up with their own name to refer to the secret they chose. This name is used exclusively by the service and does not need to match the name of the secret in the secrets manager. The **service definition** defines the secret names, along with their descriptions, used by the service. When the service is ready to be deployed, a **secret association** is added to the pattern or deployment policy that associates the name of a secret used by the service with the name of a secret in the secrets manager. More details on how to do this can be found in the instructions for [creating your own service](CreateService.md).

This document will walk you through the process of deploying this example service that uses a secret called `hw-secret-name` in the secrets manager.

## <a id=preconditions></a> Preconditions for Using the Hello Secret World Example Edge Service

If you haven't done so already, you must do these steps before proceeding with the Hello Secret World example:

1. Install the OpenHorizon management hub infrastructure (exchange and agbot) by following the directions [here](https://github.com/open-horizon/devops/tree/master/mgmt-hub).

2. Save the exchange user credentials output by the management hub installation and set the `HZN_EXCHNAGE_USER_AUTH` environment variable accordingly:

  ```bash
  export HZN_EXCHANGE_USER_AUTH="<your-exchange-user>:<your-exchange-password>"
  hzn exchange user list
  ```

3. Install the OpenHorizon agent on your edge device and configure it to point to your OpenHorizon exchange.

4. As part of the infrastructure installation process for OpenHorizon, a file called `agent-install.cfg` was created that contains the values for `HZN_ORG_ID` and the exchange and css URLs. Locate this file and set those environment variables in your shell now:

  ```bash
  eval export $(cat agent-install.cfg)
  ```

5. If you have not done so already, unregister your node before moving on. This allows you to re-register with a deployment pattern in the next section.

  ```bash
  hzn unregister -f
  ```

## <a id=using-hello-secret-pattern></a> Using the Hello Secret World Example Edge Service with a Deployment Pattern

The following steps will show you how to use the OpenHorizon CLI to create a secret, deploy it to your edge node, and consume it within an edge service. In addition, since it is common for secrets to change more often than the services themselves, the following steps will also show you how to use the OpenHorizon CLI to update an existing secret and consume the updated secret in an edge service.

1. (This step must be done by an organization admin) Create an organization-wide secret called `hw-secret-name` in your organization's secrets manager. The flags `--secretKey` and `--secretDetail` provide the key and value for the secret, respectively. The key is ignored by the service, but the value is used and will be displayed in the service's log output like this: `"Hello <secret-value>!"`.

  ```bash
  hzn sm secret add --secretKey <secret-key> --secretDetail <secret-value> hw-secret-name
  ```

2. Register your edge node with OpenHorizon to use the hello-secret pattern:

  ```bash
  hzn register -p IBM/pattern-ibm.hello-secret -s ibm.hello-secret --serviceorg IBM
  ```

3. After the service has started, list the docker containers to see it:

  ``` bash
  sudo docker ps
  ```

4. **Open another terminal** and view the hello-secret service output with the "follow" flag. This sample service repeatedly checks for updates to the secret `hw-secret-name` in the secrets manager and uses its value as a parameter of who it should say "hello" to. Initially you should see the message like: **\<your-node-id> says: Hello \<secret-value>!**, where \<secret-value> is the value of the secret as passed to the creation command in Step 1.

  ```bash
  hzn service log -f ibm.hello-secret
  ```

5. (This step must be done by an organization admin) Update the secret `hw-secret-name` in the secrets manager using the OpenHorizon CLI. To see any change reflected in the output, the argument passed to `--secretDetail` should be different from the initial secret creation command, so that the value of the secret can change.

  ```
  hzn sm secret add --secretKey <secret-key> --secretDetail <new-secret-value> hw-secret-name
  ```

A prompt will pop up in the terminal asking you to confirm that you want to overwrite the existing `hw-secret-name` in the secrets manager. Enter `y` to confirm.

6. After some time you should see the output of the service change to **\<your-node-id> says: Hello \<new-secret-value>!**, where `<new-secret-value>` is the new value of the secret as updated in the secrets manager in Step 5.

7. Unregister your edge node (which will also stop the hello-secret service):

  ```bash
  hzn unregister -f
  ```

## <a id=secret-info></a> Secrets Manager Access Control and CLI

### Access Control

To access secrets using the OpenHorizon secrets manager CLI, a user must provide their Exchange credentials through the environment variable `HZN_EXCHANGE_USER_AUTH` or can specify their credentials through flags passed to the command. The secrets manager in the OpenHorizon management hub authenticates users by delegating the authentication to the exchange. The secrets manager can also be accessed directly, using the secrets manager specific CLI or API, but this is outside the scope of OpenHorizon.

Secrets are created within two different scopes: organization-wide secrets and user private secrets. Organization wide secrets are the most common type of secret because any admin in the organization can work with them and deploy services which refer to these secrets. In a production environment, with rare exception, organization-wide secrets will be the only kinds of secrets in use. User private secrets are intended for special use cases where the secret is managed by a single authenticated user. Organization admins can delete user private secrets, but are unable to read the secret details.

Any authenticated user in the organization can list organization-wide secrets, meaning they can see the names of secrets in their organization. However, only organization admins can read the details of organization-wide secrets, meaning they can access the key/pair of a secret using the CLI. Organization admins are also able to add, remove, and update organization-wide secrets.

User private secrets are owned by individual users in the organization. These secrets will not be accessible by anyone other than the user that owns them. User private secrets follow a naming convention; they take the form of `user/<user>/<name>`, where `<user>` is the Exchange ID of the user who owns that secret and `<name>` is the user-unique name of the secret. This convention is enforced by OpenHorizon when referring to user private secrets in the secrets manager.

### Secret Names 

OpenHorizon allows the use of special characters `/`, `-`, and `_` when creating secret names. These characters can be used to create hierarchical secret names, such as `test/password` and `test/certificate`. Any secret name that does not follow the naming convention for user private secrets is considered an organization-wide secret by the secrets manager. Numbers can also be used in secret names, but no other non-alphanumeric characters can be used other than the three specified. 

As stated previously, user private secrets also follow a naming convention that utilizes `/`. User private secret names take the form of `user/<user>/<name>`, where `<user>` is the owner of the secret and `<name>` can be any secret name that follows the rules above. 

### Secrets Manager OpenHorizon CLI

The OpenHorizon CLI includes a subcommand `secretsmanager` (or the short form `sm`) for organization administrators and users to access the resources in the secrets manager, with access restricted as described above.

For organization administrators looking to manage the resources in their organization's secrets manager, the CLI has commands to `add`, `remove`, and `read` secrets in the secrets manager. These commands can also be used by users on their own private secrets. 

For service developers who need to specify secrets in their patterns and deployment policies, the `list` command allows them to list all the secrets in their organization, as well as check whether an organization-wide secret exists. Of course, this command can also be used by users on their own private secrets. To list a user's private secrets, the name `user/<user>` should be passed to `list`, where `<user>` is the owner of the secrets. Only the user who owns those secrets is able to do this.

The full documentation of these commands is available through `hzn sm secret --help`.


