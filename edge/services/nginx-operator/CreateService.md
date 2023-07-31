# Creating Your Own Operator Edge Service

Follow the steps on this page to create your first ansible operator that deploys an edge service to a cluster. With this operator we will expose the nginx service so you can `curl` it externally using a route on OCP, or thru the cluster's IP address if you are using a microk8s or k3s edge cluster.

## Preconditions for developing your own service

1. If you have not already done so, complete the steps in these sections:

   - [Preconditions for Using the Operator Example Edge Service](README.md#preconditions)
   - [Using the Operator Example Edge Service with Deployment Policy](README.md#using-operator-policy)

2. If you are using macOS as your development host, configure Docker to store credentials in `~/.docker`:

   - Open the Docker **Preferences** dialog
   - Uncheck **Securely store Docker logins in macOS keychain**

3. Install [operator-sdk](https://v1-27-x.sdk.operatorframework.io/docs/installation/) and all its prerequisites.

4. Install the Kubenetes CLI [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

5. If you do not already have a docker ID, obtain one at [https://hub.docker.com](https://hub.docker.com/), and log in to Docker Hub using your Docker Hub ID:

   ```bash
   export DOCKER_HUB_ID="<dockerhubid>"
   echo "<dockerhubpassword>" | docker login -u $DOCKER_HUB_ID --password-stdin
   ```

   Output example:

   ```text
   WARNING! Your password will be stored unencrypted in /home/pi/.docker/config.json.
   Configure a credential helper to remove this warning. See
   https://docs.docker.com/engine/reference/commandline/login/#credentials-store

   Login Succeeded
   ```

6. Create a cryptographic signing key pair. This enables you to sign services when publishing them to the exchange. **This step only needs to be done once.**

   ```bash
   hzn key create "<x509-org>" "<x509-cn>"
   ```

   where `<x509-org>` is your company name, and `<x509-cn>` is typically set to your email address.

7. Install `git` and `jq`:

   - On **Linux**:

     ```bash
     sudo apt install -y git jq
     ```

   - On **macOS**:

     ```bash
     brew install git jq
     ```

## <a id=build-publish-your-op></a> Building and Publishing Your Own Nginx Ansible Operator Example Edge Service

In order to deploy a containerized edge service to an edge cluster, a software developer first has to build a Kubernetes Operator that deploys the containerized edge service in a Kubernetes cluster. There are several options when writing a Kubernetes operator. This example will guide through creating an ansible operator. These steps are based on the [Ansible Operator Tutorial](https://v1-27-x.sdk.operatorframework.io/docs/building-operators/ansible/tutorial/) on the official `operator-sdk` website. If you have never created an operator before, I highly suggest skimming over the information there as well. The following steps were originally performed on the cluster host machine. 

1. Creating yourself a base working directory, and grab the Makefile int his repo:
   ```bash
   mkdir operator-example-service
   cd operator-example-service/
   wget https://raw.githubusercontent.com/open-horizon/examples/master/edge/services/nginx-operator/Makefile
   ```

2. Export the following environment variables to customize the `operator-sdk` init:
   ```bash
   export OPERATOR_GROUP_NAME=my-nginx-ansible-operator
   export OPERATOR_TYPE=ansible
   export OPERATOR_API_VERSION=v1alpha1
   export OPERATOR_DOMAIN=$DOCKER_HUB_ID
   export OPERATOR_KIND=MyNginxAnsibleOperator
   export OPERATOR_NAMESPACE=operator-project
   ```

3. Create `my-nginx-ansible-operator`:
   ```bash
   make init
   ```

The above Makefile command will:
- Create a new `my-nginx-ansible-operator` project and generate the entire operator structure and empty `roles/mynginxansibleoperator` 
- Create a `MyNginxAnsibleOperator` API
- Added `services` to the operator RBAC
- Changed the default namespace to the value set with `OPERATOR_NAMESPACE`
- Added `size: 1` to the operators custom resource

4. Gather the necessary nginx deployment, service, and task files for the operator to deploy the nginx service:
   ```bash
   make nginx-files
   ```

5. Build and push the operator:
   ```bash
   make build push
   ```

6. Deploy the operator locally:
   ```
   make deploy
   ```

7. After a few moments you should see your pods begin to start and your service get created:
   ```bash
   kubectl get pods -n $OPERATOR_NAMESPACE
   ```

   If everything deployed correctly, you should see output similar to the following after around 60 seconds:

   ```text
   NAME                                                            READY   STATUS    RESTARTS   AGE
   nginx-6ccdb77fd6-6s59q                                          1/1     Running   0          2s
   my-nginx-ansible-operator-controller-manager-7b65577c94-smsm6   2/2     Running   0          10s
   ```

8. Check that the service is up:

   ```bash
   kubectl get service -n $OPERATOR_NAMESPACE
   ```

   If everything deployed correctly, you should see output similar to the following after around 60 seconds:

   ```text
   NAME                                                           TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
   my-nginx-ansible-operator-controller-manager-metrics-service   ClusterIP   10.43.46.134    <none>        8443/TCP       23s
   nginx                                                          NodePort    10.43.193.244   <none>        80:30080/TCP   12s
   ```

   If you are using an **OCP edge cluster** you will need to `curl` the service using the exposed `route`.

9. Get the exposed route name:

    ```bash
    kubectl get route -n $OPERATOR_NAMESPACE
    ```

    If the route was exposed correctly you should see output similar to the following:

    ```text
    NAME          HOST/PORT                                                    PATH   SERVICES   PORT   TERMINATION   WILDCARD
    nginx-route   nginx-route-openhorizon-agent.apps.apollo5.cp.fyre.ibm.com          nginx      8080                 None
    ```

10. `curl` the service to test if it is functioning correctly:

     **OCP edge cluster** substitute the above `HOST/PORT` value:

     ```bash
     curl nginx-route-openhorizon-agent.apps.apollo5.cp.fyre.ibm.com
     ```

     **k3s or microk8s edge cluster**:

     ```bash
     curl <ip-address>:30080
     ```

     If the service is running you should see following `Welcome to nginx!` output:

     ```html
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

11. Delete the operator:
   ```bash
   make undeploy
   ```

12. Create an `operator.tar.gz` file before moving onto the next section.
   ```bash
   make tar
   ```

## <a id=publish-op-service></a> Publish Your Operator Example Edge Service

**Note:** The following commands were performed from a separate machine, **not** the cluster host.

1. Create a new working directory for a new horizon project:

   ```bash
   hzn dev service new -V 1.0.0 -s $OPERATOR_GROUP_NAME -c cluster
   ```

2. Transfer the `operator.tar.gz` archive you created in the previous section to your `my-operator/` directory.

3. Edit the `horizon/service.definition.json` file to point to the operator's yaml archive created in the previous step. Assuming it is in the `my-operator/` directory, you can make it the following:

   ```json
   "operatorYamlArchive": "../operator.tar.gz"
   ```

4. Publish your operator service:

   ```bash
   hzn exchange service publish -f horizon/service.definition.json
   ```

## <a id=publish-op-policy></a> Create a Deployment Policy For Your Operator Example Edge Service

1. Create a `deployment.policy.json` file to deploy your operator service to an edge cluster:

   ```bash
   cat << 'EOF' > horizon/deployment.policy.json
   {
   "label": "$SERVICE_NAME Deployment Policy",
   "description": "A super-simple sample Horizon Deployment Policy",
   "service": {
      "name": "$SERVICE_NAME",
      "org": "$HZN_ORG_ID",
      "arch": "*",
      "serviceVersions": [
         {
         "version": "$SERVICE_VERSION",
         "priority":{}
         }
      ]
   },
   "properties": [
   ],
   "constraints": [
      "example == my-operator"
   ],
   "userInput": [
   ]
   }
   EOF
   ```

   Notice we have given this deployment policy the following constraint: `"example == my-operator"`

2. Publish your deployment policy:

   ```bash
   hzn exchange deployment addpolicy -f horizon/deployment.policy.json policy-my-operator
   ```

3. Back on your cluster host, create a `node.policy.json` file:

   ```bash
   cat << 'EOF' > node.policy.json
   {
   "properties": [
      { "name": "example", "value": "my-operator" }
   ]
   }
   EOF
   ```

4. Register your edge cluster with your new node policy:

   ```bash
   hzn register -u $HZN_EXCHANGE_USER_AUTH
   cat node.policy.json | hzn policy update -f-
   hzn policy list
   ```

   Notice the node policy contains the property that matches the constraint specified by the deployment policy we created, which means your operator service will begin deploying to your edge cluster momentarily.

5. Check to see the agreement has been created (this can take approximately 15 seconds):

   ```bash
   hzn agreement list
   ```

6. Once you see an `agreement_execution_start_time`, you should start to see the operator pod and service deployment pod begin to start up:

   ```bash
   kubectl get pods -n openhorizon-agent
   ```

   If everything deployed correctly, you should see output similar to the following after around 60 seconds:

   ```text
   NAME                                                            READY   STATUS    RESTARTS   AGE
   nginx-6ccdb77fd6-6s59q                                          1/1     Running   0          2s
   my-nginx-ansible-operator-controller-manager-7b65577c94-smsm6   2/2     Running   0          10s
   ```

7. Check that the service is up:

   ```bash
   kubectl get service -n openhorizon-agent
   ```

   If everything deployed correctly, you should see output similar to the following after around 60 seconds:

   ```text
   NAME                                                           TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
   my-nginx-ansible-operator-controller-manager-metrics-service   ClusterIP   10.43.46.134    <none>        8443/TCP       23s
   nginx                                                          NodePort    10.43.193.244   <none>        80:30080/TCP   12s
   ```

   If you are using an **OCP edge cluster** you will need to `curl` the service using the exposed `route`.

8. Get the exposed route name:

    ```bash
    kubectl get route -n openhorizon-agent
    ```

    If the route was exposed correctly you should see output similar to the following:

    ```text
    NAME          HOST/PORT                                                    PATH   SERVICES   PORT   TERMINATION   WILDCARD
    nginx-route   nginx-route-openhorizon-agent.apps.apollo5.cp.fyre.ibm.com          nginx      8080                 None
    ```

9. `curl` the service to test if it is functioning correctly:

     **OCP edge cluster** substitute the above `HOST/PORT` value:

     ```bash
     curl nginx-route-openhorizon-agent.apps.apollo5.cp.fyre.ibm.com
     ```

     **k3s or microk8s edge cluster**:

     ```bash
     curl <ip-address>:30080
     ```

     If the service is running you should see following `Welcome to nginx!` output:

     ```html
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

10. Unregister your edge cluster:

    ```bash
    hzn unregister -f
    ```
