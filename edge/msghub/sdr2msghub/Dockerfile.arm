FROM arm32v7/golang:1.10 as go_build

RUN apt-get update && apt-get install -y \
  curl \
  git \
  lame

RUN curl -L \
  "http://1dd40.http.dal05.cdn.softlayer.net/arm_artifacts/libtensorflow-cpu-linux-armv7l-1.9.0-rc0.tar.gz" | \
  tar -C /usr/local -xz

ENV LD_LIBRARY_PATH=/lib:/usr/lib:/usr/local/lib

# go deps for tf
RUN mkdir -p ${GOPATH}/src/github.com/tensorflow/ && \
  cd ${GOPATH}/src/github.com/tensorflow/ && \
  git clone https://github.com/tensorflow/tensorflow.git

RUN cd ${GOPATH}/src/github.com/tensorflow/tensorflow && git checkout v1.9.0-rc0

RUN go get github.com/Shopify/sarama
RUN apt-get install libmp3lame-dev
RUN go get github.com/viert/lame




COPY msghub/sdr2msghub/main.go /
COPY msghub/sdr2msghub/audiolib/audiolib.go /go/src/github.com/open-horizon/examples/edge/msghub/sdr2msghub/audiolib/audiolib.go
COPY services/sdr/rtlsdrclientlib/clientlib.go /go/src/github.com/open-horizon/examples/edge/services/sdr/rtlsdrclientlib/clientlib.go
RUN go build -o /bin/data_broker /main.go

FROM arm32v7/ubuntu:18.04
RUN apt-get update && apt-get install -y --no-install-recommends \
  ca-certificates \
  tar \
  ncdu \
  curl \
  lame

RUN curl -L \
"http://1dd40.http.dal05.cdn.softlayer.net/arm_artifacts/libtensorflow-cpu-linux-armv7l-1.9.0-rc0.tar.gz" | \
  tar -C /usr/local -xz

ENV LD_LIBRARY_PATH=/usr/local/lib

COPY msghub/sdr2msghub/model.pb /model.pb

COPY --from=go_build /bin/data_broker /bin/data_broker

CMD ["data_broker"]
