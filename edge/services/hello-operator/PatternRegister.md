# Horizon Operator Example Edge Service

## Preconditions for Developing Your Own Service

1. If you have not already done so, complete the steps in this section:

  - [Preconditions for Using the Operator Example Edge Service](README.md#preconditions)
  
## <a id=using-operator-pattern></a> Using the Operator Example Edge Service with Deployment Pattern

1. Register your edge node with Horizon to use the `ibm.operator` pattern:
  ```bash
  hzn register -p IBM/pattern-hello-operator-amd64 -s hello-operator --serviceorg IBM -u $HZN_EXCHANGE_USER_AUTH
  ```
 - **Note**: using the `-s` flag with the `hzn register` command will cause Horizon to wait until agreements are formed and the service is running on your edge node to exit, or alert you of any errors encountered during the registration process. 

2. Verify that the `hello-operator` deployment is up and running:
  ```bash
  kubectl get pods -n openhorizon-agent
  ```

- If everything deployed correctly you you should see output similar to the following:

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

3. Verify that the operator is running successfully by curl-ing the service using one of the following methods:
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

4. Unregister your edge node (which will also stop the operator and helloworld service):
  ```bash
  hzn unregister -f
  ```
