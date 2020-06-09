This service starts a [VolantMQ](https://volantmq.io) MQTT server. 

## Building

You will need a quemu to build a multi-arch image. For Mac just update to edge version.

```
gmake build
```

### Generate users file

Before starting the service you need to create a users file in `username: sha256(password)` format.

On Mac you can use 

```
echo username: $(echo -n "password" | openssl dgst -sha256 | sed 's/^.* //') >> vmq.users 
```

e.g.

```
echo fft-server: $(echo -n "server-pass" | openssl dgst -sha256 | sed 's/^.* //') > /tmp/hzn/vmq.users
echo fft-client: $(echo -n "client-pass" | openssl dgst -sha256 | sed 's/^.* //') >> /tmp/hzn/vmq.users
```

Should result in this

```
cat /tmp/hzn/vmq.users
fft-server: 7b1bf1e4f9535de960093f1c303fe35f49167bdc103ba99ad7dc9d62e2807a1d
fft-client: fbfc2da74af1af1945ba7bf403cde789091e39b13c420170080872323dd2d148
```