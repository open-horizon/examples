## Docker build for cogwerx-mjpg-streamer-pi3 on IBM edge/Horizon, as an edge microservice
#  Runs a streaming camera feed, accessible on LAN at <Rpi3 IP>:8080
#  Ref: https://github.com/open-horizon/cogwerx-mjpg-streamer-pi3

## Start from cogwerx's image    ~200MB
FROM openhorizon/mjpg-streamer-pi3:20180306

ENV ARCH armhf

# Copy in bash script -- simple run file
COPY start.sh start.sh
CMD ["./start.sh", "&"]

# For docker --squash
RUN apt-get -y autoremove && apt-get clean

# Manual start command (for dev/test): 
#./mjpg_streamer -o "output_http.so -w ./www" -i "input_raspicam.so -x 640 -y 480 -fps 20 -ex night -vf"
