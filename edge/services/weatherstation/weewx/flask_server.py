from flask import Flask
#from multiprocess import Process, Manager
import json
import copy
app = Flask(__name__)
 
@app.route("/")
def hello():
    return "weewx service: queries must be directed to /v1/weather"

@app.route("/v1/weather")
def weather():
    global _data
    return json.dumps(copy.deepcopy(_data))
 
def run_server(host, port, data):
    global _data
    _data = data    
    app.run(host=host, port=port)

if __name__ == "__main__":
    global _data
    _data = argv[1]
    # (In a docker container on linux, 127.0.0.1 can't connect)
    app.run(host="0.0.0.0", port=5000)
