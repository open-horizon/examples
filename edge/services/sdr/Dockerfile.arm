FROM arm32v6/alpine:latest as rtl_build
RUN apk --no-cache add git cmake libusb-dev make gcc g++ alpine-sdk
COPY never_close.patch /tmp/
COPY signal_hack.patch /tmp/
WORKDIR /tmp
RUN git clone https://github.com/texane/librtlsdr
WORKDIR /tmp/librtlsdr
RUN git checkout rpc
RUN patch -p1 </tmp/never_close.patch
RUN patch -p1 </tmp/signal_hack.patch

RUN mkdir build && cd build && cmake ../ && make && make install
RUN ls /usr/local/bin/rtl_*

FROM arm32v6/golang:1.10.0-alpine as go_build
RUN apk --no-cache add git
RUN go get github.com/hajimehoshi/go-mp3
COPY main.go /
COPY rtlsdrclientlib/clientlib.go /go/src/github.com/open-horizon/examples/edge/services/sdr/rtlsdrclientlib/clientlib.go
COPY bbcfake/bbcfake.go /go/src/github.com/open-horizon/examples/edge/services/sdr/bbcfake/bbcfake.go
ARG version=0.0.2
ENV MIC_VERSION $version
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo --ldflags "-X main.version=${MIC_VERSION}" -o /bin/rtlsdrd /main.go

FROM arm32v6/alpine:latest
RUN apk --no-cache add alsa-utils libusb ca-certificates
COPY --from=go_build /bin/rtlsdrd /bin/rtlsdrd
COPY --from=rtl_build /usr/local/bin/rtl_rpcd /bin/rtl_rpcd
COPY --from=rtl_build /usr/local/bin/rtl_fm /bin/rtl_fm
COPY --from=rtl_build /usr/local/bin/rtl_power /bin/rtl_power
COPY --from=rtl_build /usr/local/lib/librtlsdr.so.0 /usr/local/lib/librtlsdr.so.0
WORKDIR /
CMD ["/bin/rtlsdrd"]
