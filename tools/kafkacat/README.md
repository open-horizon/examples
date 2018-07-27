# kafkacat

This directory is here because some of the services use kafkacat to publish messsages to IBM Message Hub.
So they need kafkacat on alpine. The repo https://github.com/sgerrand/alpine-pkg-kafkacat provides this for
x86, but not for arm or arm64. So we built it on arm and arm64 and put the (small) apk packages here.