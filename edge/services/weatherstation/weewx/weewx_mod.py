##
#  weewx_mod.py
#  Script to change weewx.conf parameters to values set by the user
#  (wee_config does not offer all options via command line)
#  Author / Maintainer: dyec@us.ibm.com
#

import os
import subprocess
import sys
import logging

class weewx_mod:
    ''' weewx config file modifier class '''

    config_dict = {"AcuRite":          "weewx.drivers.acurite",
            "CC3000":           "weewx.drivers.cc3000",
            "FineOffsetUSB":    "weewx.drivers.fousb",
            "Simulator":        "weewx.drivers.simulator",
            "TE923":            "weewx.drivers.te923",
            "Ultimeter":        "weewx.drivers.ultimeter",
            "Vantage":          "weewx.drivers.vantage",
            "WMR100":           "weewx.drivers.wmr100",
            "WMR200":           "weewx.drivers.wmr200",
            "WMR9x8":           "weewx.drivers.wmr9x8",
            "WS1":              "weewx.drivers.ws1",
            "WS23xx":           "weewx.drivers.ws23xx",
            "WS28xx":           "weewx.drivers.ws28xx" }
  
    # Default constructor
    def __init__(self, config_file=None, station_type=None):
        self.config_file = config_file
        self.wee_config_script = "wee_config"
        logging.basicConfig(level=logging.INFO)
        self.logger = logging.getLogger(__name__)

    # Station params
        self.station_type = station_type
        self.driver = self.weewx_driver_lookup(self.station_type)
        self.pws_model = "WS2080"      		# (just a default)
        self.pws_lat = 0.000
        self.pws_lon = 0.000

    # WU-specific params
        self.wu_station_id = None 
	self.wu_pwd = None
        self.wu_location = None
        self.wu_rapidfire = 'False'     	# config file expects String

    ## Member Functions
    def splitup(self, line):
        indent = ''
        param = ''
        value = ''
        comment = ''
        comm_buf = ''
        
        if line.find(' ') == 0:
            indent = ' '*(len(line) - len(line.lstrip(' ')))
        if line.find('#') >= 0:
            line, comment = line.split('#', 1)
            comm_buf = ' '*(len(line) - len(line.lstrip(' ')))
            comment = comm_buf + '#' + comment.rstrip('\n')
        if line.find(' = ') >= 0:
            param, value = line.strip().split(' = ')
        
        return indent, param, value, comment
        
    
    def replace_param(self, lines, keywd, p, old_val, new_val, limit=1):
        ''' Replace a parameter's value with a new one, using optional match '''

        def assemble(indent, param, value, comment):
            ''' Helper function for line reassembly '''
            return indent + param + ' = ' + value + comment + '\n'

        replaced = 0
        i = 0
        kwd_found = False
        while replaced < limit and i < len(lines):
            line = lines[i]
            if kwd_found == True:   # Found the section with our keyword
                if line.find(p) >= 0 and line.find('=') >= 0:
                    indent, param, value, comment = self.splitup(line)
                    if old_val != '':
                        if value.find(old_val) >= 0:
                            lines[i] = assemble(indent, param, str(new_val), comment)
                            replaced += 1
                    else:
                        lines[i] = assemble(indent, param, str(new_val), comment)
                        replaced += 1
            elif line.find(keywd) >= 0:
                kwd_found = True
            i += 1
        return replaced

    # Weewx station type lookup from driver
    def weewx_driver_lookup(self, station_type):
        # Find station driver by key (made case-insensitive)
        types = [d.upper() for d in self.config_dict.keys()]
        if station_type.upper() in types:
           idx = types.index(station_type.upper())
           return self.config_dict.values()[idx]
        else:
           return None

    def set_latlon(self, lat, lon):
        ''' Set values for lat/lon in class obj '''

        def check_val(val, extents):
            if val >= extents[0] and val <= extents[1]:
                return True
            return False

        # Check values prior to setting (ignore if out of range)
        if lat != "" and lon != "":
            if check_val(float(lat), [-90.0, 90.0]):
                self.pws_lat = float(lat)
                    
            if check_val(float(lon), [-180.0, 180.0]):
                self.pws_lon = float(lon)
      

    def set_wu_cfg(self, wu_id, wu_pwd, rpdf='False', loc_str=''):
        ''' ''' 
        self.wu_station_id = wu_id
        self.wu_pwd = wu_pwd
        if rpdf == 'True':    	self.wu_rapidfire = 'True'
        else:             	self.wu_rapidfire = 'False'
        self.wu_location = loc_str

    def run_wee_config(self):
        ''' Update weewx config file using wee_conf utility (run 1st)'''
        
        cmd = [self.wee_config_script, 
            "--reconfigure", 
            "--driver=%s" % str(self.driver),
            "--longitude=%s" % str(self.pws_lon),
            "--latitude=%s" % str(self.pws_lat),
            "--location=%s" % str(self.wu_location),
            "--no-prompt",
            "--no-backup"]

        retcode = subprocess.call(cmd)
        if retcode != 0:
            self.logger.error("subprocess cmd '%s' returned %s\n" % (cmd,retcode))

    def update_config_file(self):
        ''' Update the weewx config file manually (run 2nd) ''' 

        # Grab all lines from input file and store in list
        lines = []
        with open(self.config_file, 'r') as infile:
            line = None
            while line != "":
                line = infile.readline()
                lines.append(line)

        # Modify other weewx station params manually (wee_config doesn't do these)
        self.replace_param(lines, '[Station]', 'model', '', self.pws_model)
        self.replace_param(lines, '[[Wunderground]]', 'enable', '', 'true')
        self.replace_param(lines, '[[Wunderground]]', 'station', '', self.wu_station_id)
        self.replace_param(lines, '[[Wunderground]]', 'password', '', '"'+self.wu_pwd+'"')
        self.replace_param(lines, '[[Wunderground]]', 'rapidfire', '', self.wu_rapidfire)
        self.replace_param(lines, '[[Wunderground]]', 'location', '', self.wu_location)

        ## Write out new config file (replace the input copy)
        with open(self.config_file, 'w') as outfile:
            for line in lines:
                outfile.write(line)
        self.logger.debug("%s Station params modified\n" % (self.config_file))

    def update_all(self):
        ''' '''

        # Run wee_config to modify the config file for us, where it can
        self.run_wee_config()

        # With the config file properly updated, we can modify it with the 
        # remaining WU-specific parameters
        self.update_config_file()

       
#### Main
###
def main(args=None):
    """The main routine."""
    if args is None:
        args = sys.argv[1:]

    if len(args) < 2:
        usage_string = ("%s Usage: %s <weewx_config_file> <PWS Station Type>" %
         (sys.argv[0], sys.argv[0]))
        print(usage_string)
        print("%s weewx Station types: %s" % (sys.argv[0], weewx_mod.config_dict.keys()))
        sys.exit(1)

    # Input args
    config_file = args[1]
    station_type = args[2]
    
    # Set up class instance and run mod routine
    weemod = weewx_mod(config_file, station_type)
    weemod.run()

