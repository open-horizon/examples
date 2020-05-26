ARG BUILD_IMAGE
ARG RUN_IMAGE
FROM $BUILD_IMAGE AS build

ENV GOOS=linux
ENV COMMIT_HASH=b89554f33b5516a9bbdb53167bc8504258b2ea2d

ARG GOARCH
ENV GOARCH=$GOARCH

RUN apk update && apk add git bash gcc g++ musl-dev binutils-gold && \
 	go get -v github.com/troian/govvv && \
 	GO111MODULE=off go get -v github.com/VolantMQ/volantmq/cmd/volantmq; exit 0

RUN cd /go/src/github.com/VolantMQ/volantmq/cmd/volantmq && \
	git reset --hard $COMMIT_HASH && \
 	go mod tidy && \
	govvv build -i -o /volantmq/volantmq
	
##################################################################################################

FROM $RUN_IMAGE

ENV VOLANTMQ_PLUGINS_DIR=/usr/local/bin/plugins/ \
	VOLANTMQ_CONFIG=/config/volantmq.yaml

COPY --from=build /volantmq/ /usr/local/bin/
COPY config.yaml /config/volantmq.yaml

RUN ln -s /usr/local/bin/volantmq /

EXPOSE 1883
CMD ["volantmq"]