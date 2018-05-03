# Start command for darknet yolov3, and a small HTTP REST service at port 8359
#

if [ -z "$YOLO_CAMERA" ]; then
        YOLO_CAMERA=0
fi

if [ -z "$YOLO_LOG_INTERVAL" ]; then
        YOLO_LOG_INTERVAL=15
fi

# Run the HTTP REST service process and VNC (via entrypoint script) with Darknet/YoloV3 args
socat TCP4-LISTEN:8359,fork EXEC:/service.sh &
/bin/bash /entrypoint.sh /darknet/darknet detector demo -c ${YOLO_CAMERA} cfg/coco.data cfg/yolov3.cfg yolov3.weights > /dev/null 2>&1 &

# Log yolo output every n seconds, just to show something (yolo by default logs way too fast... hence > dev null above)
while true; do 
  sleep $YOLO_LOG_INTERVAL
  echo -en "$(curl -sSL http://localhost:8359/v1/yolo)"
done
