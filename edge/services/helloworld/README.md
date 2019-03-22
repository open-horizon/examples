# Horizon Hello World Service

## Building and Publishing

Set environment variables as recommended in `horizon/envvars.sh.sample` and then:

```
make build
make hznstart
docker logs -f $(docker ps -q)
make hznstop
make publish-service-only
make publish-pattern
```
