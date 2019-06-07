FROM ubuntu:bionic
RUN apt update && apt install -y bash curl jq vim
WORKDIR /
COPY ess-read-all.sh /
COPY service.sh /
#CMD tail -f /dev/null
CMD /service.sh
