#!/bin/bash

# Example IBM Edge "Microservice" that returns results of Yolo inference (objects detected by camera)
# Note: requires gawk and bc to be installed.

getObjsFromFile() {
        # Get saw JSON content to be passed along in request
        filename=/darknet/output.json  # must be at this location
        $(echo cat "$filename")
}

# Get the currect YOLO object detection result, then construct the HTTP response message
OBJS=$(getObjsFromFile)
HEADERS="Content-Type: text/html; charset=ISO-8859-1"
BODY="{\"yolo\":${OBJS}}"
HTTP="HTTP/1.1 200 OK\r\n${HEADERS}\r\n\r\n${BODY}\r\n"

# Emit the HTTP response
echo -en $HTTP

