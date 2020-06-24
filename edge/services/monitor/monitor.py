#
# Debug client for monitoring/debugging the achatina examples
#
# Written by Glen Darling, October 2019.
#

import json
import os
import subprocess
import threading
import time
from datetime import datetime
import base64

# Configuration constants
MQTT_SUB_COMMAND = 'mosquitto_sub -h mqtt -p 1883 -C 1 '
MQTT_DETECT_TOPIC = '/detect'
FLASK_BIND_ADDRESS = '0.0.0.0'
FLASK_PORT = 5200
DUMMY_DETECT_IMAGE='/dummy_detect.jpg'

# Globals for the cached JSON data (last messages on these MQTT topics)
last_detect = None

if __name__ == '__main__':

  from io import BytesIO
  from flask import Flask
  from flask import send_file
  webapp = Flask('monitor')                             
  webapp.config['SEND_FILE_MAX_AGE_DEFAULT'] = 0

  # Loop forever collecting object detection / classification data from MQTT
  class DetectThread(threading.Thread):
    def run(self):
      global last_detect
      # print("\nMQTT \"" + MQTT_DETECT_TOPIC + "\" topic monitor thread started!")
      DETECT_COMMAND = MQTT_SUB_COMMAND + '-t ' + MQTT_DETECT_TOPIC
      while True:
        last_detect = subprocess.check_output(DETECT_COMMAND, shell=True)
        # print("\n\nMessage received on detect topic...\n")
        # print(last_detect)

  @webapp.route("/images/detect.jpg")
  def get_detect_image():
    if last_detect:
      j = json.loads(last_detect)
      i = base64.b64decode(j['detect']['image'])
      buffer = BytesIO()
      buffer.write(i)
      buffer.seek(0)
      return send_file(buffer, mimetype='image/jpg')
    else:
      return send_file(DUMMY_DETECT_IMAGE)

  @webapp.route("/json")
  def get_json():
    if last_detect:
      return last_detect.decode("utf-8") + '\n'
    else:
      return '{}\n'

  @webapp.route("/")
  def get_results():
    if None == last_detect: return '{"error":"Server not ready."}'
    j = json.loads(last_detect)
    n = j['deviceid']
    c = len(j['detect']['entities'])
    ct = j['detect']['camtime']
    it = j['detect']['time']
    s = j['source']
    u = j['source-url']
    # print(s, u)
    kafka_msg = '<p> &nbsp; <em>NOTE:</em> Nothing is being published to EventStreams (kafka)!</p>\n'
    if 'kafka-sub' in j:
      sub = j['kafka-sub']
      kafka_msg = '<p> &nbsp; <em>NOTE:</em> This data is also being published to EventStreams (kafka). Subscribe with:</p>\n' + \
        '<p style="font-family:monospace;">' + sub + '</p>\n'
    OUT = \
      '<html>\n' + \
      ' <head>\n' + \
      '   <title>achatina monitor</title>\n' + \
      '   <style>body { width: 650px; }</style>\n' + \
      ' </head>\n' + \
      ' <body>\n' + \
      '   <div>\n' + \
      '   <table>\n' + \
      '     <tr>\n' + \
      '       <th><h2 style="color:blue;">' + s + '</h2></th>\n' + \
      '     </tr>\n' + \
      '     <tr>\n' + \
      '       <th>Device ID: ' + n + '</th>\n' + \
      '     </tr>\n' + \
      '     <tr>\n' + \
      '       <th><span id="when">&nbsp;</span></th>\n' + \
      '     </tr>\n' + \
      '     <tr>\n' + \
      '       <td><img id="detect" height="480px" width="640px" src="/images/detect.jpg" alt="Prediction Image" /></td>\n' + \
      '     </tr>\n' + \
      '     <tr><td> &nbsp; </td></tr>\n' + \
      '     <tr>\n' + \
      '       <td> &nbsp; Found entities in <span id="classes">' + str(c) + '</span> classes.</td>\n' + \
      '     </tr>\n' + \
      '     <tr>\n' + \
      '       <td> &nbsp; Camera time: <span id="camtime">' + str(ct) + '</span> seconds.</td>\n' + \
      '     </tr>\n' + \
      '     <tr>\n' + \
      '       <td> &nbsp; Inferencing time: <span id="inftime" style="font-weight:bold;color:blue;">' + str(it) + '</span> seconds.</td>\n' + \
      '     </tr>\n' + \
      '     <tr><td> &nbsp; </td></tr>\n' + \
      '     <tr>\n' + \
      '       <td> &nbsp; More information: <a href="' + u + '">' + u + '</a></td>\n' + \
      '     </tr>\n' + \
      '     <tr><td> &nbsp; </td></tr>\n' + \
      '     <tr>\n' + \
      '       <td>\n' + \
      '         ' + kafka_msg + \
      '       </td>\n' + \
      '     </tr>\n' + \
      '   </table>\n' + \
      '   </div>\n' + \
      '   <script>\n' + \
      '     function refresh(d_image, d_date, d_classes, d_camtime, d_inftime) {\n' + \
      '       var t = 500;\n' + \
      '       (async function startRefresh() {\n' + \
      '         var address;\n' + \
      '         if(d_image.src.indexOf("?")>-1)\n' + \
      '           address = d_image.src.split("?")[0];\n' + \
      '         else\n' + \
      '           address = d_image.src;\n' + \
      '         d_image.src = address+"?time="+new Date().getTime();\n' + \
      '         const response = await fetch("/json");\n' + \
      '         const j = await response.json();\n' + \
      '         var when = new Date(j.detect.date * 1000);\n' +\
      '         var c = j.detect.entities.length;\n' +\
      '         var ct = j.detect.camtime;\n' +\
      '         var it = j.detect.time;\n' +\
      '         d_date.innerHTML = when;\n' +\
      '         d_classes.innerHTML = c;\n' +\
      '         d_camtime.innerHTML = ct;\n' +\
      '         d_inftime.innerHTML = it;\n' +\
      '         setTimeout(startRefresh, t);\n' + \
      '       })();\n' + \
      '     }\n' + \
      '     window.onload = function() {\n' + \
      '       var d_image = document.getElementById("detect");\n' + \
      '       var d_date = document.getElementById("when");\n' + \
      '       var d_classes = document.getElementById("classes");\n' + \
      '       var d_camtime = document.getElementById("camtime");\n' + \
      '       var d_inftime = document.getElementById("inftime");\n' + \
      '       refresh(d_image, d_date, d_classes, d_camtime, d_inftime);\n' + \
      '     }\n' + \
      '   </script>\n' + \
      ' </body>\n' + \
      '</html>\n'
    return (OUT)

  # Prevent caching everywhere
  @webapp.after_request
  def add_header(r):
    r.headers["Cache-Control"] = "no-cache, no-store, must-revalidate"
    r.headers["Pragma"] = "no-cache"
    r.headers["Expires"] = "0"
    r.headers['Cache-Control'] = 'public, max-age=0'
    return r

  # Main program (instantiates and starts monitor threads and then web server)
  monitor_detect = DetectThread()
  monitor_detect.start()
  webapp.run(host=FLASK_BIND_ADDRESS, port=FLASK_PORT)

