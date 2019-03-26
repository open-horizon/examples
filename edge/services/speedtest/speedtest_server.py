#
# Periodic speedtest client, and REST API for access to results
#
# Written by Glen Darling, February 2019.
# Copyright 2019, Glen Darling; all rights reserved.
#

from os import environ
from flask import Flask
import json
import subprocess
import threading
import time
import speedtest

# How long to pause between runs of the test (in seconds)
SPEEDTEST_PAUSE_SEC = environ.get('SPEEDTEST_PAUSE_SEC')
if SPEEDTEST_PAUSE_SEC is None:
  SPEEDTEST_PAUSE_SEC = (60 * 15)

# REST API details
REST_API_BIND_ADDRESS = '0.0.0.0'
REST_API_BIND_PORT = 80
webapp = Flask('speedtest')

# Global for the test results
last_test_results = None

# Run one speedtest (seems to take about 25 seconds on RPi3B)
def run_speedtest():
  servers = []
  s = speedtest.Speedtest()
  s.get_servers(servers)
  s.get_best_server()
  s.download()
  s.upload()
  s.results.share()
  return s.results.dict()
 
# Loop forever running the test
class SpeedTestThread(threading.Thread):
  def run(self):
    global last_test_results
    #print("\nSpeedTest thread started!")
    #t = 1
    while True:
      #print("\n\nRunning SpeedTest #" + str(t) + "...\n")
      #t += 1
      last_test_results = run_speedtest()
      #print(json.dumps(last_test_results))
      #print("\nSleeping for " + str(SPEEDTEST_PAUSE_SEC) + " seconds...\n")
      time.sleep(SPEEDTEST_PAUSE_SEC)

# A web server to make the speedtest results available on the LAN
@webapp.route("/v1/speedtest")
def get_speedtest():
  if None == last_test_results:
    return '{"error": "Patience, please. No speed test data received yet."}\n'
  else:
    return (json.dumps(last_test_results)) + '\n'

# Main program (instantiates and starts speedtest thread and then web server)
if __name__ == '__main__':
  tester = SpeedTestThread()
  tester.start()
  webapp.run(host=REST_API_BIND_ADDRESS, port=REST_API_BIND_PORT)
