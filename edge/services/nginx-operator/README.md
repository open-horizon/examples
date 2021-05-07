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
  hzn exchange node create -n $HZN_EXCHANGE_NODE_AUTH -u $HZN_EXCHANGE_USER_AUTH -T "<device/cluster>"
  hzn exchange node confirm -n $HZN_EXCHANGE_NODE_AUTH -u $HZN_EXCHANGE_USER_AUTH
  ```

7. If you have not done so already, unregister your node before moving on:
  ```bash
  hzn unregister -f
  ```

## <a id=using-operator-policy></a> Using the Operator Example Edge Service with Deployment Policy

In the following steps you will deploy the `nginx-operator` to your edge cluster. This operator will then create a pod running a hello world service you can `curl` externally or interally.

1. Get the required node policy file on your edge cluster host:
  ```bash
  wget https://raw.githubusercontent.com/open-horizon/examples/master/edge/services/nginx-operator/horizon/node.policy.json
  ```

- Below is the `node_policy.json` file you obtained in the step above:

  ```json
  {
    "properties": [
      { "name": "openhorizon.example", "value": "nginx-operator" }
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
          "openhorizon.arch in \"'ppc64le','amd64'\""
      ]
  }
  ```

- Below is the example `deployment.policy.json` that has been published into the Exchange as part of the example operator:

  ```json
  {
    "mycluster/policy-hello-operator": {
      "owner": "root/root",
      "label": "nginx-operator Deployment Policy",
      "description": "A super-simple sample Horizon Deployment Policy",
      "service": {
        "name": "nginx-operator",
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
        "openhorizon.example == nginx-operator"
      ],
      "created": "2020-11-05T19:18:17.722Z[UTC]",
      "lastUpdated": "2020-11-05T19:18:17.722Z[UTC]"
    }
  }
  ```

5. Verify that the `nginx-operator` deployment is up and running:
  ```bash
  kubectl get pods -n openhorizon-agent
  ```

If everything deployed correctly you should see an output similar to the following:
  ```
  NAME                                   READY   STATUS    RESTARTS   AGE
  agent-dd984ff96-jmmdl                  1/1     Running   0          1d
  nginx-8699f45b-pw7dj                   1/1     Running   0          23s
  nginx-operator-898999564-nnzb5         1/1     Running   0          50s
  ```

6. Check that the service is up:
  ```bash
  kubectl get service -n openhorizon-agent
  ```

If everything deployed correctly you should see an output similar to the following after around 60 seconds:
  ```
  NAME                  TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)             AGE
  nginx                 NodePort    172.30.37.113    <none>        80:30080/TCP        45s
  ```

If you are using an **OCP edge cluster** you will need to `curl` the service using the exposed `route`.
7. Get the exposed route name:
  ```bash
  kubectl get route -n openhorizon-agent
  ```
  
If the route was exposed correctly you should see an output similar to the following:
  ```bash
  NAME          HOST/PORT                                                    PATH   SERVICES   PORT   TERMINATION   WILDCARD
  nginx-route   nginx-route-openhorizon-agent.apps.apollo5.cp.fyre.ibm.com          nginx      8080                 None
  ```

8. `curl` the service to test if it is functioning correctly:
   **OCP edge cluster** substitute the above `HOST/PORT` value:
      ```bash
      curl nginx-route-openhorizon-agent.apps.apollo5.cp.fyre.ibm.com
      ```
   
   **k3s or microk8s edge cluster**:
      ```bash
      curl <external-ip-address>:30080
      ```

If the service is running you should see following `Welcome to nginx!` output:
   ```
   <!DOCTYPE html>
   <html>
   <head>
   <title>Welcome to nginx!</title>
   <style>
       body {
           width: 35em;
         margin: 0 auto;
         font-family: Tahoma, Verdana, Arial, sans-serif;
      }
   </style>
   </head>
   <body>
   <h1>Welcome to nginx!</h1>
   <p>If you see this page, the nginx web server is successfully installed and
   working. Further configuration is required.</p>

   <p>For online documentation and support please refer to
   <a href="http://nginx.org/">nginx.org</a>.<br/>
   Commercial support is available at
   <a href="http://nginx.com/">nginx.com</a>.</p>

   <p><em>Thank you for using nginx.</em></p>
   </body>
   </html>
   ```

9. Unregister your edge cluster:
   ```
   hzn unregister -f
   ```
