#!/usr/bin/env python
# Entry point to the PWS edge Microservice
# Uses weewx weather utility to extract weather condition data
#  from a variety of PWS (Personal Weather Stations)
#  and provide a local HTTP endpoint which can be queried
#
# Usage: start.py <options> <weewx_config_file>
#
# Author: Chris Dye (dyec@us.ibm.com)
# Some of this code borrowed from "weewxd"
#

import os, sys
import time
import copy
import logging

from optparse import OptionParser
from multiprocess import Process, Manager
import subprocess

# Local utilities
import flask_server as fl
from weewx_mod import weewx_mod

# weewx itself
import weewx.engine

#===============================================================================
#                 weewx function parseArgs()
#===============================================================================

usagestr = """Usage: weewxd --help
       weewxd --version
       weewxd config_file [--daemon] [--pidfile=PIDFILE] 
                          [--exit]   [--loop-on-init]
                          [--log-label=LABEL]
           
  Entry point to the weewx weather program. Can be run directly, or as a daemon
  by specifying the '--daemon' option.

Arguments:
    config_file: The weewx configuration file to be used.
"""

def parseArgs():
    """Parse any command line options."""

    parser = OptionParser(usage=usagestr)
    parser.add_option("-d", "--daemon",  action="store_true", dest="daemon",  help="Run as a daemon")
    parser.add_option("-p", "--pidfile", type="string",       dest="pidfile", help="Store the process ID in PIDFILE", default="/var/run/weewx.pid", metavar="PIDFILE")     
    parser.add_option("-v", "--version", action="store_true", dest="version", help="Display version number then exit")
    parser.add_option("-x", "--exit",    action="store_true", dest="exit"   , help="Exit on I/O and database errors instead of restarting")
    parser.add_option("-r", "--loop-on-init", action="store_true", dest="loop_on_init"  , help="Retry forever if device is not ready on startup")
    parser.add_option("-n", "--log-label", type="string", dest="log_label", help="Label to use in syslog entries", default="weewx", metavar="LABEL")
    (options, args) = parser.parse_args()
    
    if options.version:
        logger.info(weewx.__version__ + '\n')
        sys.exit(0)
        
    if len(args) < 1:
        logger.error("Missing argument(s).\n")
        logger.error(parser.parse_args(["--help"])+'\n')
        sys.exit(weewx.CMD_ERROR)

    return options, args

## Supporting sys functions
def check_env_var(envname, default='', printerr=True):
    ''' Checks linux environment variable and uses default if present'''
    if envname in os.environ:
       val = os.getenv(envname)
       if val == '' or val == '-':
           if printerr:
               logger.info("start.py: Environment variable" + envname + " value is '%s'\n" % val)
           return default
       return val
    
    else:
       if printerr:
           logger.info("start.py: Environment variable " + envname + " not found.\n")
       return default

## Supporting data transform function
def format_weather_data(data_str):
    """ Take weather data in weewx format and transform to param/value dict"""     
    ## Data sample direct from weewx (shortened)
	# "altimeter: 72.317316, ... maxSolarRad: None, ... windGustDir: 359.99994, windSpeed: 5.1645e-09"

    # Replace "None" values with 0's
    data_str = data_str.replace("None", "0.0")

    # Grab the list of param/values
    pairs_list = [p.strip() for p in data_str.strip().split(',')]
    
    # Capture each param/value in a dict
    pairs_dict = {}
    for p in pairs_list:
        k,v = p.split(':')
        pairs_dict[k.strip()] = v.strip()

    return pairs_dict
       

def clear_station_mem():
    p_wee_clearmem = subprocess.check_call(["/home/weewx/bin/wee_device", "--clear-memory", "-y"])
    if p_wee_clearmem != 0:
        logger.error("start.py: PWS Station memory could not be cleared. Exiting.\n")
        sys.exit(1)
    return 0


#==============================================================
#                  Main entry point
#==============================================================
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(sys.argv[0] + " " + __name__)

## Prep  --------------------------------------------------------------------
# Get the command line options and arguments:
(options, args) = parseArgs()
weewx_config_file = args[0]

## Get PWS config (via environment variables)
# Station type, Model, WU ID (required, or exit)
pws_station_type = check_env_var("PWS_ST_TYPE", printerr=True)
pws_model = check_env_var("PWS_MODEL", printerr=True)
pws_wu_id = check_env_var("PWS_WU_ID", printerr=True)
pws_wu_pwd = check_env_var("PWS_WU_KEY", printerr=True)
if pws_station_type == '' or pws_model == '' or pws_wu_id == '' or pws_wu_pwd == '':
    sys.exit(1)

# Settings: Location, Units, and rapidfire (optional)
latitude = check_env_var("HZN_LAT", printerr=True)
longitude = check_env_var("HZN_LON", printerr=True)
pws_units = check_env_var("PWS_UNITS", default='us', printerr=True)    # weewx recommends only using 'us'
pws_wu_loc = check_env_var("PWS_WU_LOC", default='', printerr=True)
pws_wu_rapidfire = check_env_var("PWS_WU_RPDF", default='False', printerr=True)

# Deal with a potential lower-case (boolean value from Horizon) or erroneous value
if pws_wu_rapidfire == "true" or pws_wu_rapidfire == "True": 
    pws_wu_rapidfire = "True"
else: 
    pws_wu_rapidfire = "False"


## Shared data structure (dict for flask server to read & serve)
manager = Manager()
sdata = manager.dict()
standard_params = ["wu_id", "stationtype", "model", "latitude", "longitude", "units", "location"]
standard_values = [pws_wu_id, pws_station_type, pws_model, latitude, longitude, pws_units, pws_wu_loc]
sdata["r"] = dict(zip(["status"], ["Station initializing..."]))
sdata["t"] = str(int(time.time()))                                      # Timestamp
sdata["i"] = dict(zip(standard_params, standard_values))                # Station Info

## Flask HTTPserver ----------------------------------------------------------
## Start simple flask server at localhost:port and pass in shared data dict
p_flask = Process(target=fl.run_server, args=('0.0.0.0', 8357, sdata))
p_flask.start()

## Weewx service -------------------------------------------------------------
# Modify the weewx configuration file with our env var settings
weemod = weewx_mod(weewx_config_file, pws_station_type)
weemod.wee_config_script = "/home/weewx/bin/wee_config"
weemod.set_latlon(lat=latitude, lon=longitude)
weemod.set_wu_cfg(wu_id=pws_wu_id, 
        wu_pwd=pws_wu_pwd, 
        rpdf=pws_wu_rapidfire, 
        loc_str=pws_wu_loc)
weemod.update_all()

# Clear PWS station memory
retcode = clear_station_mem()

## Start weewx engine with options/args and pass in shared data dict
p_weewx = Process(target=weewx.engine.main, args=(options, args, sdata))
p_weewx.start()

## -------------------------------------------------------------
#  Forever loop for data transform, check process health
last_ts = sdata["t"]    # Timestamp
last_clear_time = last_ts
while True:
    time.sleep(0.1)	# Weather data post interval >= 2.5s, so take it easy.
    # Check for new weather data from weewx (will be a string, not dict)
    if sdata["t"] != last_ts and type(sdata["r"]) == str:	# Transform weather data to dict
    	sdata["r"] = format_weather_data(copy.deepcopy(sdata["r"]))
        sdata["r"]["status"] = "Station running..."
        last_ts = sdata["t"]

    # Check processes and kill this script if one dies
    if not p_weewx.is_alive():
        logger.error("start.py: weewx process died: exit code = %s\n" % (p_weewx.exitcode))
        p_flask.terminate()
        sys.exit(1)

    if not p_flask.is_alive():
        logger.error("start.py: flask process died: exit code = %s\n" % (p_flask.exitcode))
        p_weewx.terminate()
        sys.exit(1)

logger.error("start.py: Script exited from run loop unexpectedly")
sys.exit(1)

