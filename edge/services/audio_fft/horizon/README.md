## Publishing horizon services

Generate MQTT users file using 

```
echo username: $(echo -n "password" | openssl dgst -sha256 ) >> vmq.users 
```

e.g.

```
echo fft-server: $(echo -n "server-pass" | openssl dgst -sha256 ) >> vmq.users
echo fft-client: $(echo -n "client-pass" | openssl dgst -sha256 ) >> vmq.users
```

Publish service

```
hzn exchange service publish -f service.definition.json -I
```


Publish pattern 

```
hzn exchange pattern publish -f pattern.json
```

Change volumes' binding in `service.json` to target your `vmq.users` file.

Register node

```
 hzn register -p pattern-fft --policy node.policy 
```

Check that agreement created 

```
 hzn agreement list   
```

And containers are up

```
docker ps
```