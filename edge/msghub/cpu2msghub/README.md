# Horizon CPU To IBM Message Hub Service

For details about using this service, see [cpu2msghub.md](cpu2msghub.md).

## Building and Publishing

Set environment variables as recommended in `horizon/envvars.sh.sample` and then:

```
make build
make docker-push
make exchange-publish
```

## Building kafkacat for arm and arm64

This only needs to be done once for each new version of kakfacat. This uses https://github.com/sgerrand/alpine-pkg-kafkacat/ and
roughly follows https://wiki.alpinelinux.org/wiki/Creating_an_Alpine_package for building apk packages.

```
# On the arm or arm64 machine you are building on:
mkdir -p ~/tmp/sdr/apkbuild
docker run --name alpine -d -t --privileged -v $HOME:/outside alpine

# Go into the alpine container:
docker exec -it alpine sh

apk update
apk add alpine-sdk

adduser -s /bin/ash bp   # or use another username
addgroup bp abuild
echo 'bp ALL=(ALL) ALL' >> /etc/sudoers
su - bp
sudo whoami

abuild-keygen -a  # puts private signing key in file like: ~/.abuild/-5b5afcec.rsa
sudo cp ~/.abuild/*.rsa.pub /etc/apk/keys/

mkdir apkbuild; cd apkbuild
wget https://raw.githubusercontent.com/sgerrand/alpine-pkg-kafkacat/master/APKBUILD
abuild -r

sudo apk add ~/packages/bp/*/kafkacat-1.3.1-r0.apk   # test installing it
kafkacat   # test running it

# Copy the keys and pkg out of the container
sudo cp ~/.abuild/*.rsa* /outside/tmp/sdr/apkbuild/
sudo cp ~/packages/bp/*/kafkacat-*.apk /outside/tmp/sdr/apkbuild/

exit  # get out of the container

# then put the pkg and public key into git:
cp ~/tmp/sdr/apkbuild/*.rsa.pub ~/tmp/sdr/apkbuild/kafkacat-*.apk ~/src/github.com/open-horizon/examples/tools/kafkacat/$ARCH/
# or scp them to a machine you have the repo on
```
