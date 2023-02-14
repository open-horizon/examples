# Horizon Operator Example Edge Service

## Preconditions for Developing Your Own Service

1. If you have not already done so, complete the steps in this section:

   - [Preconditions for Using the Operator Example Edge Service](README.md#preconditions)

## <a id=using-operator-pattern></a> Using the Operator Example Edge Service with Deployment Pattern

1. Register your edge node with Horizon to use the `nginx-operator` pattern:

   ```bash
   hzn register -p IBM/pattern-ibm.nginx-operator-amd64 -u $HZN_EXCHANGE_USER_AUTH
   ```

   **Note**: using the `-s` flag with the `hzn register` command will cause Horizon to wait until agreements are formed and the service is running on your edge node to exit, or alert you of any errors encountered during the registration process.

2. Verify that the `nginx-operator` deployment is up and running:

   ```bash
   kubectl get pods -n openhorizon-agent
   ```

   If everything deployed correctly you you should see output similar to the following:

   ```text
   NAME                                   READY   STATUS    RESTARTS   AGE
   agent-dd984ff96-jmmdl                  1/1     Running   0          1d
   nginx-7d5598fb56-vw6lz                 1/1     Running   0          12s
   nginx-operator-55c6f56c47-b6p7c        1/1     Running   0          48s
   ```

3. Check that the service is up:

   ```bash
   kubectl get service -n openhorizon-agent
   ```

   If everything deployed correctly you should see output similar to the following after around 60 seconds:

   ```text
   NAME    TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
   nginx   NodePort    172.30.37.113    <none>        80:30080/TCP     45s
   ```

   If you are using an **OCP edge cluster** you will need to `curl` the service using the exposed `route`.

4. Get the exposed route name:

   ```bash
   kubectl get route -n openhorizon-agent
   ```

   If the route was exposed correctly you should see output similar to the following:

   ```text
   NAME          HOST/PORT                                                    PATH   SERVICES   PORT   TERMINATION   WILDCARD
   nginx-route   nginx-route-openhorizon-agent.apps.apollo5.cp.fyre.ibm.com          nginx      8080                 None
   ```

5. `curl` the service to test if it is functioning correctly:
   **OCP edge cluster** substitute the above `HOST/PORT` value:

   ```bash
   curl nginx-route-openhorizon-agent.apps.apollo5.cp.fyre.ibm.com
   ```

   **k3s or microk8s edge cluster**:

   ```bash
   curl <external-ip-address>:30080
   ```

   If the service is running you should see the following `Welcome to nginx!` output:

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

6. Unregister your edge node (which will also stop the operator and helloworld service):

   ```bash
   hzn unregister -f
   ```
