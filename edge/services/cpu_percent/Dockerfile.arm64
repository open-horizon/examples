FROM aarch64/alpine:latest
RUN apk --no-cache --update add gawk bc socat
COPY *.sh /
WORKDIR /
CMD /start.sh
