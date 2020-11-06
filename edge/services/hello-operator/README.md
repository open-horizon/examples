# Horizon Operator Example Edge Service

This is a simple example of using and creating an operator as edge service.

- [Preconditions for Using the Operator Example Edge Service](#preconditions)
- [Using the Operator Example Edge Service with Deployment Policy](#using-operator-policy)
- [Using the Operator Example Edge Service with Deployment Pattern](PatternRegister.md)
- [Creating Your Own Operator Edge Service](CreateService.md)
- Further Learning - to see more Horizon features demonstrated, continue on to the [cpu2evtstreams example](../../evtstreams/cpu2evtstreams).

## <a id=preconditions></a> Preconditions for Using the Operator Example Edge Service

If you haven't done so already, you must do these steps before proceeding with the operator example:

1. Install the Horizon management hub (exchange and agbot).

2. Install the Horizon agent on your edge cluster and configure it to point to your Horizon exchange.

3. As part of the management hub installation process a file called `agent-install.cfg` was created that contains the values for `HZN_ORG_ID` and the exchange and css URL values. Locate this file and set those environment variables in your shell now:

  ```bash
  eval export $(cat agent-install.cfg)
  ```
 - **Note**: if for some reason you disconnected from ssh or your command line closes, run the above command again to set the required environment variables.
 
4. The `hzn` command is inside the agent container, but you can set some aliases to make it possible to run `hzn` from the cluster host with the following commands:
  ```bash
  cat << 'END_ALIASES' >> ~/.bash_aliases
  alias getagentpod='kubectl -n openhorizon-agent get pods --selector=app=agent -o jsonpath={.items[*].metadata.name}'
  alias hzn='kubectl -n openhorizon-agent exec -i $(getagentpod) -- hzn'
  END_ALIASES
  source ~/.bash_aliases
  ```

5. In addition to the file above, an API key associated with your Horizon instance would have been created, set the exchange user credentials, and verify them:
  ```bash
  export HZN_EXCHANGE_USER_AUTH=iamapikey:<horizon-API-key>
  hzn exchange user list -u $HZN_EXCHANGE_USER_AUTH
  ```

6. Choose an ID and token for your edge node, create it, and verify it:
  ```bash
  export HZN_EXCHANGE_NODE_AUTH="<choose-any-node-id>:<choose-any-node-token>"
  hzn exchange node create -n $HZN_EXCHANGE_NODE_AUTH -u $HZN_EXCHANGE_USER_AUTH
  hzn exchange node confirm -n $HZN_EXCHANGE_NODE_AUTH -u $HZN_EXCHANGE_USER_AUTH
  ```

7. If you have not done so already, unregister your node before moving on:
  ```bash
  hzn unregister -f
  ```

## <a id=using-operator-policy></a> Using the Operator Example Edge Service with Deployment Policy

In the following steps you will deploy the `hello-operator` to your edge cluster. This operator will then create a pod running a hello world service you can `curl` externally or interally.

1. Get the required node policy file on your edge cluster host:
  ```bash
  wget https://raw.githubusercontent.com/open-horizon/examples/master/edge/services/hello-operator/horizon/node.policy.json
  ```

- Below is the `node_policy.json` file you obtained in the step above:

  ```json
  {
    "properties": [
      { "name": "openhorizon.example", "value": "hello-operator" }
    ],
    "constraints": [
    ]
  }
  ```

- It provides one value for `properties` (`openhorizon.example`), that will effect which services get deployed to this edge node, and states no `constraints`.

2. Register your Node Policy with this policy
  ```bash
  hzn register -u $HZN_EXCHANGE_USER_AUTH
  cat node.policy.json | hzn policy update -f-
  ```

3. When the registration completes, use the following command to review the Node Policy:
  ```bash
  hzn policy list
  ```

4. The edge device will make an agreement with one of the Horizon agreement bots (this typically takes about 15 seconds). Repeatedly query the agreements of this device until the `agreement_finalized_time` and `agreement_execution_start_time` fields are filled in:
  ```bash
  hzn agreement list
  ```

- Below is the `service.policy.json` that was published into the Exchange when the example operator was published:

  ```json
  {
      "properties": [],
      "constraints": [
          "openhorizon.arch == amd64"
      ]
  }
  ```

- Below is the example `deployment.policy.json` that has been published into the Exchange as part of the example operator:

  ```json
  {
    "mycluster/policy-hello-operator": {
      "owner": "root/root",
      "label": "hello-operator Deployment Policy",
      "description": "A super-simple sample Horizon Deployment Policy",
      "service": {
        "name": "hello-operator",
        "org": "IBM",
        "arch": "*",
        "serviceVersions": [
          {
            "version": "1.0.0",
            "priority": {},
            "upgradePolicy": {}
          }
        ],
        "nodeHealth": {}
      },
      "constraints": [
        "openhorizon.example == hello-operator"
      ],
      "created": "2020-11-05T19:18:17.722Z[UTC]",
      "lastUpdated": "2020-11-05T19:18:17.722Z[UTC]"
    }
  }
  ```

5. Verify that the `hello-operator` deployment is up and running:
  ```bash
  kubectl get pods -n openhorizon-agent
  ```

- If everything deployed correctly you should see output similar to the following:

  ```
   NAME                                   READY   STATUS    RESTARTS   AGE
   agent-dd984ff96-jmmdl                  1/1     Running   0          1d
   hello-operator-6c5f8c4458-6ggwx        1/1     Running   0          24s
   mosquito-helloworld-7bccc7668c-x9qf7   1/1     Running   0          7s
   ```

**Note:** If you are attempting to run this service on an **OCP edge cluster** and the operator does not start you may have to grant the operator the privileges it requires to execute with the following command:
  ```bash
  oc adm policy add-scc-to-user privileged -z hello-operator -n openhorizon-agent
  ```

6. Verify that the operator is running successfully by `curl`-ing the service using one of the following methods:
  ```bash
   curl -sS <INTERNAL_IP>:8000 | jq .

   - or externally - 

   curl -sS <NODE_IP>:30007 | jq .
   ```

If the service is running you should see output similar to the following:
   ```json
   {
     "Hello": "10.22.29.174"
   }
   ```

7. Unregister your edge node (which will also stop the operator and helloworld service):

  ```bash
  hzn unregister -f
  ```
