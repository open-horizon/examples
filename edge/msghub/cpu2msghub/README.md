# Horizon CPU To IBM Message Hub Service

For details about using this service, see [cpu2msghub.md](cpu2msghub.md).

## Building and Publishing

**Note:** if you are building for arm or arm64, first ensure that a recent version of kafkacat has
been built for those architectures on alpine. See our [kafkacat README](../../../tools/kafkacat/README.md).

Set environment variables as recommended in `horizon/envvars.sh.sample` and then:

```
make build
make docker-push
make exchange-publish
```
