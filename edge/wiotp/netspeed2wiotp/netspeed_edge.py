#!/usr/bin/env python
# -*- coding: utf-8 -*-
# Original Speedtest-cli code: Copyright 2012-2015 Matt Martz
# Modified with additions for edge compute demo
# All Rights Reserved.
#
#    Licensed under the Apache License, Version 2.0 (the "License"); you may
#    not use this file except in compliance with the License. You may obtain
#    a copy of the License at
#
#         http://www.apache.org/licenses/LICENSE-2.0
#
#    Unless required by applicable law or agreed to in writing, software
#    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
#    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
#    License for the specific language governing permissions and limitations
#    under the License.

import os
import re
import sys
import math
import signal
import socket
import timeit
import platform
import threading

import datetime
import calendar
from random import randint

import time
import string
import json
import csv
import subprocess

import paho.mqtt.publish as publish
import paho.mqtt.client as mqtt
import ssl

from workload_config import *   # read netspeed configuration
import mqtt_pub as mqttpub      # import mqtt pub supporting functions
import utils                    # utilities file in this dir (utils.py)

__version__ = '0.3.4'           # speedtest-cli version


# netspeed global variables
last_date = 0
last_hour = 0
netpoc_error = '0x0000'

# speedtest-cli global variables 
user_agent = None
source = None
shutdown_event = None
scheme = 'http'

# Used for bound_interface
socket_socket = socket.socket


# Speedtest-cli code starts here
try:
    import xml.etree.cElementTree as ET
except ImportError:
    try:
        import xml.etree.ElementTree as ET
    except ImportError:
        from xml.dom import minidom as DOM
        ET = None

# Begin import game to handle Python 2 and Python 3
try:
    from urllib2 import urlopen, Request, HTTPError, URLError
except ImportError:
    from urllib.request import urlopen, Request, HTTPError, URLError

try:
    from httplib import HTTPConnection, HTTPSConnection
except ImportError:
    e_http_py2 = sys.exc_info()
    try:
        from http.client import HTTPConnection, HTTPSConnection
    except ImportError:
        e_http_py3 = sys.exc_info()
        raise SystemExit('Your python installation is missing required HTTP '
                         'client classes:\n\n'
                         'Python 2: %s\n'
                         'Python 3: %s' % (e_http_py2[1], e_http_py3[1]))

try:
    from Queue import Queue
except ImportError:
    from queue import Queue

try:
    from urlparse import urlparse
except ImportError:
    from urllib.parse import urlparse

try:
    from urlparse import parse_qs
except ImportError:
    try:
        from urllib.parse import parse_qs
    except ImportError:
        from cgi import parse_qs

try:
    from hashlib import md5
except ImportError:
    from md5 import md5

try:
    from argparse import ArgumentParser as ArgParser
except ImportError:
    from optparse import OptionParser as ArgParser
'''
try:
    import builtins
except ImportError:
    def print_(*args, **kwargs):
        """The new-style print function taken from https://pypi.python.org/pypi/six/ """
        fp = kwargs.pop("file", sys.stdout)
        if fp is None:
            return

        def write(data):
            if not isinstance(data, basestring):
                data = str(data)
            fp.write(data)

        want_unicode = False
        sep = kwargs.pop("sep", None)
        if sep is not None:
            if isinstance(sep, unicode):
                want_unicode = True
            elif not isinstance(sep, str):
                raise TypeError("sep must be None or a string")
        end = kwargs.pop("end", None)
        if end is not None:
            if isinstance(end, unicode):
                want_unicode = True
            elif not isinstance(end, str):
                raise TypeError("end must be None or a string")
        if kwargs:
            raise TypeError("invalid keyword arguments to print()")
        if not want_unicode:
            for arg in args:
                if isinstance(arg, unicode):
                    want_unicode = True
                    break
        if want_unicode:
            newline = unicode("\n")
            space = unicode(" ")
        else:
            newline = "\n"
            space = " "
        if sep is None:
            sep = space
        if end is None:
            end = newline
        for i, arg in enumerate(args):
            if i:
                write(sep)
            write(arg)
        write(end)
else:
    print_ = getattr(builtins, 'print')
    del builtins
'''

class SpeedtestCliServerListError(Exception):
    """Internal Exception class used to indicate to move on to the next
    URL for retrieving speedtest.net server details

    """


def bound_socket(*args, **kwargs):
    """Bind socket to a specified source IP address"""

    global source
    sock = socket_socket(*args, **kwargs)
    sock.bind((source, 0))
    return sock


def distance(origin, destination):
    """Determine distance between 2 sets of [lat,lon] in km"""

    lat1, lon1 = origin
    lat2, lon2 = destination
    radius = 6371  # km

    dlat = math.radians(lat2 - lat1)
    dlon = math.radians(lon2 - lon1)
    a = (math.sin(dlat / 2) * math.sin(dlat / 2) +
         math.cos(math.radians(lat1)) *
         math.cos(math.radians(lat2)) * math.sin(dlon / 2) *
         math.sin(dlon / 2))
    c = 2 * math.atan2(math.sqrt(a), math.sqrt(1 - a))
    d = radius * c

    return d


def build_user_agent():
    """Build a Mozilla/5.0 compatible User-Agent string"""

    global user_agent
    if user_agent:
        return user_agent

    ua_tuple = (
        'Mozilla/5.0',
        '(%s; U; %s; en-us)' % (platform.system(), platform.architecture()[0]),
        'Python/%s' % platform.python_version(),
        '(KHTML, like Gecko)',
        'speedtest-net-poc/%s' % __version__
    )

    if debug_flag:
        utils.print_('Platform information:')
        utils.print_(platform.platform())
        utils.print_(ua_tuple)
      
    user_agent = ' '.join(ua_tuple)
    return user_agent


def build_request(url, data=None, headers={}):
    """Build a urllib2 request object. This function automatically adds a User-Agent header to all requests"""

    if url[0] == ':':
        schemed_url = '%s%s' % (scheme, url)
    else:
        schemed_url = url

    headers['User-Agent'] = user_agent
    return Request(schemed_url, data=data, headers=headers)


def catch_request(request):
    """Helper function to catch common exceptions encountered when establishing a connection with a HTTP/HTTPS request"""

    try:
        uh = urlopen(request)
        return uh, False
    except (HTTPError, URLError, socket.error):
        e = sys.exc_info()[1]
        return None, e


class FileGetter(threading.Thread):
    """Thread class for retrieving a URL"""

    def __init__(self, url, start):
        self.url = url
        self.result = None
        self.starttime = start
        threading.Thread.__init__(self)

    def run(self):
        self.result = [0]
        try:
            if (timeit.default_timer() - self.starttime) <= 10:
                request = build_request(self.url)
                f = urlopen(request)
                while 1 and not shutdown_event.isSet():
                    self.result.append(len(f.read(10240)))
                    if self.result[-1] == 0:
                        break
                f.close()
        except IOError:
            pass


def downloadSpeed(files, quiet=False):
    """Function to launch FileGetter threads and calculate download speeds"""

    download_metrics = {
        'speed_Mbs': 0.0,      # Mbits / sec
        'data_MB': 0.0,        # MBytes
        'time_s': 0.0          # seconds
    }

    start = timeit.default_timer()

    def producer(q, files):
        for file in files:
            thread = FileGetter(file, start)
            thread.start()
            q.put(thread, True)
            if not quiet and not shutdown_event.isSet():
                sys.stdout.write('.')
                sys.stdout.flush()

    finished = []

    def consumer(q, total_files):
        while len(finished) < total_files:
            thread = q.get(True)
            while thread.isAlive():
                thread.join(timeout=0.1)
            finished.append(sum(thread.result))
            del thread

    q = Queue(6)
    prod_thread = threading.Thread(target=producer, args=(q, files))
    cons_thread = threading.Thread(target=consumer, args=(q, len(files)))
    start = timeit.default_timer()
    prod_thread.start()
    cons_thread.start()
    while prod_thread.isAlive():
        prod_thread.join(timeout=0.1)
    while cons_thread.isAlive():
        cons_thread.join(timeout=0.1)

    time_s = (timeit.default_timer() - start)     # total download time in sec  
    data_B = sum(finished)                        # total download data in bytes

# netspeed edge modifications

    download_speed = data_B / time_s              # download speed in Bytes/sec
    download_metrics = {
        'speed_Mbs': ((download_speed * 8) / 1024 / 1024),   # speed in Mbits / sec 
        'data_MB': (data_B / 1024 / 1024),                   # data in MBytes
        'time_s': time_s                                     # time in seconds
    }
    if debug_flag:
        utils.print_('\nData (Bytes): %d  Time (sec): %0.3f  Speed(Bytes/sec): %0.3f' 
                  % (data_B, time_s, download_speed) )
        utils.print_('Download Volume(MB): %(data_MB)0.3f Speed(Mbps): %(speed_Mbs)0.3f' 
                  % download_metrics)   

    return download_metrics


class FilePutter(threading.Thread):
    """Thread class for putting a URL"""

    def __init__(self, url, start, size):
        self.url = url
        chars = '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ'
        data = chars * (int(round(int(size) / 36.0)))
        self.data = ('content1=%s' % data[0:int(size) - 9]).encode()
        del data
        self.result = None
        self.starttime = start
        threading.Thread.__init__(self)

    def run(self):
        try:
            if ((timeit.default_timer() - self.starttime) <= 10 and
                    not shutdown_event.isSet()):
                request = build_request(self.url, data=self.data)
                f = urlopen(request)
                f.read(11)
                f.close()
                self.result = len(self.data)
            else:
                self.result = 0
        except IOError:
            self.result = 0


def uploadSpeed(url, sizes, quiet=False):
    """Function to launch FilePutter threads and calculate upload speeds"""

    upload_metrics = {
        'speed_Mbs': 0.0,      # Mbits / sec
        'data_MB': 0.0,        # MBytes
        'time_s': 0.0          # seconds
    }

    start = timeit.default_timer()

    def producer(q, sizes):
        for size in sizes:
            thread = FilePutter(url, start, size)
            thread.start()
            q.put(thread, True)
            if not quiet and not shutdown_event.isSet():
                sys.stdout.write('.')
                sys.stdout.flush()

    finished = []

    def consumer(q, total_sizes):
        while len(finished) < total_sizes:
            thread = q.get(True)
            while thread.isAlive():
                thread.join(timeout=0.1)
            finished.append(thread.result)
            del thread

    q = Queue(6)
    prod_thread = threading.Thread(target=producer, args=(q, sizes))
    cons_thread = threading.Thread(target=consumer, args=(q, len(sizes)))
    start = timeit.default_timer()
    prod_thread.start()
    cons_thread.start()
    while prod_thread.isAlive():
        prod_thread.join(timeout=0.1)
    while cons_thread.isAlive():
        cons_thread.join(timeout=0.1)

    time_s = (timeit.default_timer() - start)          # total upload time in sec
    data_B = sum(finished)                             # total upload data in bytes

# netspeed edge modifications

    upload_speed = data_B / time_s                     # upload speed in Bytes/sec
    upload_metrics = {
        'speed_Mbs': ((upload_speed*8) / 1024 / 1024), # speed in Mbits/sec 
        'data_MB': (data_B / 1024 / 1024),             # data in MBytes
        'time_s': time_s                               # time in seconds
    }

    if debug_flag:
        utils.print_('\nData (Bytes): %d  Time(sec): %0.3f  Speed(Bytes/sec): %0.3f' 
                   % (data_B, time_s, upload_speed) )
        utils.print_('Upload Volume(MB): %(data_MB)0.3f Speed(Mbps): %(speed_Mbs)0.3f' 
                   % upload_metrics)   

    return upload_metrics


def getAttributesByTagName(dom, tagName):
    """Retrieve an attribute from an XML document and return it in a
    consistent format

    Only used with xml.dom.minidom, which is likely only to be used
    with python versions older than 2.5
    """
    elem = dom.getElementsByTagName(tagName)[0]
    return dict(list(elem.attributes.items()))


def getConfig():
    """Download the speedtest.net configuration and return only the data we are interested in. Returns None if error."""

    request = build_request('://www.speedtest.net/speedtest-config.php')
    uh, e = catch_request(request)
    if e:
        utils.print_('Error: could not retrieve speedtest.net configuration: %s' % e)
        return None
    configxml = []
    while 1:
        configxml.append(uh.read(10240))
        if len(configxml[-1]) == 0:
            break
    if int(uh.code) != 200:
        utils.print_('Error: got HTTP code %s when trying to retrieve speedtest.net configuration' % str(uh.code))
        return None
    uh.close()
    try:
        try:
            root = ET.fromstring(''.encode().join(configxml))
            config = {
                'client': root.find('client').attrib,
                'times': root.find('times').attrib,
                'download': root.find('download').attrib,
                'upload': root.find('upload').attrib}
        except AttributeError:  # Python3 branch
            root = DOM.parseString(''.join(configxml))
            config = {
                'client': getAttributesByTagName(root, 'client'),
                'times': getAttributesByTagName(root, 'times'),
                'download': getAttributesByTagName(root, 'download'),
                'upload': getAttributesByTagName(root, 'upload')}
    except SyntaxError as e:
        utils.print_('Failed to parse speedtest.net configuration: %s' % str(e))
        return None
    del root
    del configxml
    return config


def closestServers(client, all=False):
    """Determine the 5 closest speedtest.net servers based on geographic distance"""

    global netpoc_error
    urls = [
        '://www.speedtest.net/speedtest-servers-static.php',
        '://c.speedtest.net/speedtest-servers-static.php',
        '://www.speedtest.net/speedtest-servers.php',
        '://c.speedtest.net/speedtest-servers.php',
    ]
    errors = []
    servers = {}
    for url in urls:
        try:
            request = build_request(url)
            uh, e = catch_request(request)
            if e:
                errors.append('%s' % e)
                raise SpeedtestCliServerListError
            serversxml = []
            while 1:
                serversxml.append(uh.read(10240))
                if len(serversxml[-1]) == 0:
                    break
            if int(uh.code) != 200:
                uh.close()
                raise SpeedtestCliServerListError
            uh.close()
            try:
                try:
                    root = ET.fromstring(''.encode().join(serversxml))
                    elements = root.getiterator('server')
                except AttributeError:  # Python3 branch
                    root = DOM.parseString(''.join(serversxml))
                    elements = root.getElementsByTagName('server')
            except SyntaxError:
                raise SpeedtestCliServerListError
            for server in elements:
                try:
                    attrib = server.attrib
                except AttributeError:
                    attrib = dict(list(server.attributes.items()))
                d = distance([float(client['lat']),
                              float(client['lon'])],
                             [float(attrib.get('lat')),
                              float(attrib.get('lon'))])
                attrib['d'] = d
                if d not in servers:
                    servers[d] = [attrib]
                else:
                    servers[d].append(attrib)
            del root
            del serversxml
            del elements
        except SpeedtestCliServerListError:
            continue

        # We were able to fetch and parse the list of speedtest.net servers
        if servers:
            break

    if not servers:
        if debug_flag:
            utils.print_('Failed to retrieve list of speedtest.net servers:\n\n %s' %
               '\n'.join(errors))
        netpoc_error = 'netx0002' # cannot get list of closest servers
        sys.exit(1)

    closest = []
    for d in sorted(servers.keys()):          # sort servers based on distance
        for s in servers[d]:
            closest.append(s)
            if len(closest) == 5 and not all:
                break
        else:
            continue
        break

    del servers
    return closest


def getBestServer(servers):
    """Perform a speedtest.net latency request to determine which speedtest.net server has the lowest latency"""

    results = {}
    for server in servers:
        cum = []
        url = '%s/latency.txt' % os.path.dirname(server['url'])
        urlparts = urlparse(url)
        for i in range(0, 3):
            try:
                if urlparts[0] == 'https':
                    h = HTTPSConnection(urlparts[1])
                else:
                    h = HTTPConnection(urlparts[1])
                headers = {'User-Agent': user_agent}
                start = timeit.default_timer()
                h.request("GET", urlparts[2], headers=headers)
                r = h.getresponse()
                total = (timeit.default_timer() - start)
            except (HTTPError, URLError, socket.error):
                cum.append(3600)
                continue
            text = r.read(9)
            if int(r.status) == 200 and text == 'test=test'.encode():
                cum.append(total)
            else:
                cum.append(3600)
            h.close()
        avg = round((sum(cum) / 6) * 1000, 3)     # latency in ms
        results[avg] = server
    fastest = sorted(results.keys())[0]
    best = results[fastest]
    best['latency'] = fastest

    return best


def ctrl_c(signum, frame):
    """Catch Ctrl-C key sequence and set a shutdown_event for our threaded
    operations
    """

    global shutdown_event
    shutdown_event.set()
    raise SystemExit('\nCancelling...')


def version():
    """Print the version"""

    raise SystemExit(__version__)

# netspeed edge code starts here
def post_networkdata(jsonpayload, event_id, heart_beat=False):
    """Sends network data in json format to mqtt. Returns True if successful."""
    retries = 2
    if heart_beat:  retries = SEND_MAX_RETRIES
    for i in range(1,retries+1):
        result = mqttpub.post_networkdata_single_wiotp(jsonpayload, event_id, heart_beat=heart_beat)
        if result == 1:  return True        # success
        if result == -1:
            # We were not registered
            utils.print_('Send to mqtt failed. Not registered.')
        else:
            # The send failed for some reason other than not be registered
            time.sleep(SEND_RETRY_DELAY)

    return False

def myspeedtest():
    """ tests edge device's network speed using speedtest.net test """

    global total_volume_MB_month, debug_flag, mqtt_flag
    global latitude, longitude, contract_id, device_id, jsonfile, file_flag
    global netpoc_error

    if debug_flag:
        utils.print_('Retrieving speedtest.net configuration...')  
    try:
        config = getConfig()
        if not config:
            # When getConfig() hits an error, it prints it and returns None. Returning None here will cause our caller to retry.
            return None
    except URLError:
        # getConfig() catches this, but leaving it here for safety
        if debug_flag:
            utils.print_('Cannot retrieve speedtest configuration')
        return None

    if debug_flag:
        utils.print_('Retrieving speedtest.net server list...')

    if  (target_server_criteria == 'closest') | (target_server_criteria == 'fastest'):   
        if debug_flag:
            utils.print_('Testing from %(isp)s (%(ip)s)...' % config['client'])

        servers = closestServers(config['client'])    # get top 5 closest servers

        if (target_server_criteria == 'fastest'):
            if debug_flag:
                utils.print_('Selecting best server based on latency...')

            best = getBestServer(servers)  # get server with lowest latency from 5 closest servers

        else:
            if debug_flag:
                utils.print_('Selecting best server based on distance...')

            best = getBestServer(servers)  # looks like this is using same criteria as fastest??

    elif (target_server_criteria == 'random'):        # select random server
        if debug_flag:
            utils.print_('Selecting random server ...')

        servers = closestServers(config['client'], True)    # get full list of servers
        serverrange = len(servers)
        targetserver = randint(0, serverrange-1)
        serverid=servers[targetserver]['id']
        if debug_flag:
            utils.print_('server[%d] out of %d: %s %s \n' % (targetserver, serverrange,
                      serverid, servers[targetserver]['name']) )
            utils.print_('Testing from %(isp)s (%(ip)s)...' % config['client'])

        try:
            best = getBestServer(filter(lambda x: x['id'] == serverid, servers))   
        except IndexError as e:
            utils.print_('Invalid server ID: %s' % str(e))
            return None
    
    timestamp = datetime.datetime.now()   # get time of test

    if debug_flag:
        utils.print_(('Hosted by %(sponsor)s (%(name)s) [%(d)0.2f km]: '
               '%(latency)s ms' % best).encode('utf-8', 'ignore'))

    # sizes = [350, 500, 750, 1000, 1500, 2000, 2500, 3000, 3500, 4000]
    sizes = [500, 1000, 2000]
    urls = []
    for size in sizes:
        for i in range(0, 2):
            urls.append('%s/random%sx%s.jpg' %
                        (os.path.dirname(best['url']), size, size))

    if debug_flag:
        utils.print_('Testing download speed', end='')

    download_metrics = downloadSpeed(urls, not(debug_flag))   
    dlspeed = download_metrics['speed_Mbs']

    sizesizes = [int(.25 * 1000 * 1000), int(.5 * 1000 * 1000)]
    sizes = []
    for size in sizesizes:
        for i in range(0, 25):
            sizes.append(size)
    
    if debug_flag:    
        utils.print_()
        utils.print_('Testing upload speed', end='')

    upload_metrics = uploadSpeed(best['url'], sizes, not(debug_flag))   
    ulspeed = upload_metrics['speed_Mbs']

    """ if edge device lat and lon not available, use value computed by nettest """

    if (latitude == DEFAULT_LAT):
        latitude = config['client']['lat']

    if (longitude == DEFAULT_LON):
        longitude = config['client']['lon']

    ping_dict = {
        "host": best['url'],
        "min": float(best['latency']),
        "max": float(best['latency']),
        "avg": float(best['latency'])
    }

    netspeedresults = {
        'upload_Mbps': round(ulspeed,4),
        'download_Mbps': round(dlspeed,4),
        'ping_ms': ping_dict,
        'distance_km': round(best['d'],4),
        'targetserver': [ best['sponsor'], best['name'], best['country'],  best['url'] ],
        'latitude': float(latitude),
        'longitude': float(longitude),
        'device_id':device_id
    }

    networkdata = {
        't': calendar.timegm(timestamp.timetuple()),
        'r': netspeedresults
    }
    jsonpayload = json.dumps(networkdata)

    if (mqtt_flag):
        post_networkdata(jsonpayload, event_id='netspeed-speedtest')
    else:
        if debug_flag:
            utils.print_(jsonpayload)

    if (file_flag):
        jsonfile = open('./netspeedresults.json', 'w')
        json.dump(networkdata, jsonfile, sort_keys = True, indent = 4, ensure_ascii=False )
        jsonfile.write('\n')
        jsonfile.close()

    """ keep track of data usage per month """ 
    total_volume_MB_month += upload_metrics['data_MB']
    total_volume_MB_month += download_metrics['data_MB']

    if debug_flag:
        utils.print_('Total BW per month (MB): %0.3f MB \n'
               ' ' % (total_volume_MB_month) )   
    return netspeedresults


def speedtest_with_retry():
    """Calls myspeedtest(), retrying if it fails."""
    for i in range(1,SPEEDTEST_MAX_RETRIES+1):
        result = myspeedtest()
        if result:  return result
        if i < SPEEDTEST_MAX_RETRIES:
            utils.print_('Speed test was unsuccessful, sleeping %s seconds and then will try again...' % str(SPEEDTEST_RETRY_DELAY))
            time.sleep(SPEEDTEST_RETRY_DELAY)
    return None


def clear_monthly_data():
    global total_volume_MB_month
    global max_mbps_exceeded, max_volume_exceeded

    total_volume_MB_month = 0.0
    max_volume_exceeded = 0
    max_mbps_exceeded = 0

def pingstatus():
    """Gets ping latency info and schedules itself to run again at the next interval."""
    global contract_id, device_id, json_filename
    
    threading.Timer(PING_INTERVAL,pingstatus).start ()
    timestamp = datetime.datetime.now()

    uptime_raw = subprocess.check_output(["uptime"])
    m = re.search('(?<=up ).+',uptime_raw)
    uptime_m = m.group(0)
    uptime_db_col = uptime_m.strip()
    load_avg_list = uptime_db_col.split('average:')[-1].strip().split(',')
     
    uptime_dict = {
        'uptime': uptime_db_col.split(',')[0],
        'load_avg': [float(la.strip()) for la in load_avg_list]
    }
    
    free_raw = subprocess.check_output(["free","-mh"])
    free_raw_lines = free_raw.splitlines()
    free_db_col = ''
    if len(free_raw_lines) >= 1:
        mem_cols = free_raw_lines[1].split()
        if len(mem_cols) >= 6:
            free_db_col = 'total:' + str(mem_cols[1]) + ', free:' + str(mem_cols[2]) + ', shared:' + str(mem_cols[3]) + ', buffers:' + str(mem_cols[4]) + ', cached:' + str(mem_cols[5])
    
    ping_raw = subprocess.check_output(["ping", "-c", "1", "www.ibm.com"])
    ping_raw_lines = ping_raw.splitlines()
    ping_db_col = ''
    matching_lines = [line for line in ping_raw_lines if 'ping statistics ---' in line]     # make sure we got summary stats
    if len(matching_lines) > 0:
        ping_db_col = matching_lines[0].strip()
    matching_lines = [line for line in ping_raw_lines if 'round-trip min/avg/max' in line]     # get the stats
    ping_dict = {}
    if len(matching_lines) > 0:
        ping_db_col = ping_db_col + matching_lines[0].strip()
        ping_list = matching_lines[0].split('=')[-1].strip().split('/')[-3:]    # should be [min, avg, max]
        if len(ping_list) > 0:
            ping_dict = {
                'min': float(ping_list[0].strip()),
                'avg': float(ping_list[1].strip()),
                'max': float(ping_list[2].strip(' ms'))
            }

    df_raw = subprocess.check_output(["df","-h"])
    df_raw_lines = df_raw.splitlines()
    df_raw_lines.pop(0)
    df_db_col = ''
    df_lines = ''
    for l in df_raw_lines:
        df_cols = l.split()
        if len(df_cols) >= 6:
            df_line = '[' + df_cols[0] + ' size:' + df_cols[1] + ', used:' + df_cols[2] + ', avail:' + df_cols[3] + ', use%:' + df_cols[4] + ', mnt:' + df_cols[5] + ']'
            df_lines += df_line
    df_db_col = df_lines.strip()
            
    netspeedping = {
        'uptime': uptime_dict,
        'ping_ms': ping_dict,
        'contract_id': contract_id,
        'device_id': device_id
    }

    networkdata = {
        't': calendar.timegm(timestamp.timetuple()),
        'r': netspeedping
    }

    jsonpayload = json.dumps(networkdata)

    if (mqtt_flag):
        post_networkdata(jsonpayload, event_id='netspeed-ping')
    
    if (file_flag):
        jsonfile = open(json_filename, 'w')
        json.dump(jsonpayload, jsonfile, sort_keys = True, indent = 4, ensure_ascii=False )
        jsonfile.write('\n')
        jsonfile.close()  
  
    if debug_flag:
        utils.print_(jsonpayload)

def speedtestscheduler():
    """Gets speed data and schedules itself to run again at the next interval."""
    global total_volume_MB_month
    global max_volume_exceeded, max_mbps_exceeded
    global last_date
    global policy_flag

    threading.Timer(run_interval, speedtestscheduler).start ()

    current_date = datetime.datetime.now()     

    if debug_flag:
        utils.print_('\nCurrent date: ', current_date.strftime('%Y-%m-%d %H:%M:%S'))
        utils.print_('Last date: ', last_date.strftime('%Y-%m-%d %H:%M:%S'))

    if not(policy_flag):            # do not check policy, perform netspeed test
        testresults = speedtest_with_retry()
        if not testresults:
            utils.print_('Error: speed test failed after maximum retries.')
        return
 
    """ if new month, clear all cumulative values and policy exceeded flags """
    if (current_date.month != last_date.month):
        if (max_volume_exceeded | max_mbps_exceeded):
           if debug_flag:
                utils.print_('resuming network tests.') 

        """ Clear monthly data and perform first network test of the month """
        clear_monthly_data()
        last_date = current_date
    else:
        """ Check policy exceeded flags """
        if (max_volume_exceeded | max_mbps_exceeded):
            return

    """ Perform netspeed test """
    testresults = speedtest_with_retry()
    if not testresults:
        utils.print_('Error: speed test failed after maximum retries.')
        return

    """ Check data usage """
    if (total_volume_MB_month > max_volume_MB_month):
        if debug_flag:
            utils.print_('Send volume exceeded. Total upload %0.3f MB > %0.3f'
                     'Suspending network tests...'
                     % (total_volume_MB_month, max_volume_MB_month) )
        max_volume_exceeded = 1

def netpoc_init():
    global last_date, target_server_criteria, run_interval
    global latitude, longitude, contract_id, contract_nonce, device_id, debug_flag
    global total_volume_MB_month, max_volume_MB_month
    global mqtt_broker, mqtt_tls, mqtt_port, mqtt_client_id
    #global wiotp_domain, wiotp_org_id, wiotp_device_auth_token, wiotp_device_type
  
    clear_monthly_data()
 
    socket.setdefaulttimeout(10)

    # Pre-cache the user agent string
    build_user_agent()

    last_date = datetime.datetime.now()

    # Log the settings we are running with
    utils.print_('Running with these settings:')
    utils.print_('  Target network speed test server: %s' % target_server_criteria)
    utils.print_('  Run interval: %d' % run_interval)
    utils.print_('  Monthly bandwidth cap: %d' % max_volume_MB_month)
    utils.print_('  Ping interval: %d' % PING_INTERVAL)
    utils.print_('  Latitude: %s' % latitude)
    utils.print_('  Longitude: %s' % longitude)
    utils.print_('  Horizon agreement id: %s' % contract_id)
    utils.print_('  Horizon hash: %s' % contract_nonce)
    utils.print_('  Horizon device id: %s' % device_id)
    if 'HZN_EXCHANGE_URL' in os.environ: utils.print_('  Horizon exchange URL: %s' % os.environ['HZN_EXCHANGE_URL'])
    utils.print_('  MQTT broker hostname: %s' % mqtt_broker)
    utils.print_('  MQTT broker port: %s' % mqtt_port)
    utils.print_('  MQTT broker PEM file: %s' % mqtt_ca_file)
    utils.print_('  REG_MAX_RETRIES: %d' % REG_MAX_RETRIES)
    utils.print_('  REG_RETRY_DELAY: %d' % REG_RETRY_DELAY)
    utils.print_('  REG_SUCCESS_SLEEP: %d' % REG_SUCCESS_SLEEP)
    utils.print_('  SEND_MAX_RETRIES: %d' % SEND_MAX_RETRIES)
    utils.print_('  SEND_RETRY_DELAY: %d' % SEND_RETRY_DELAY)
    utils.print_('  SPEEDTEST_MAX_RETRIES: %d' % SPEEDTEST_MAX_RETRIES)
    utils.print_('  SPEEDTEST_RETRY_DELAY: %d' % SPEEDTEST_RETRY_DELAY)


def main():
 
    global shutdown_event, target_server_criteria, run_interval
    global send_policy_MB_month, receive_policy_MB_month 
    global policy_flag, mqtt_flag, debug_flag, file_flag, json_filename
    global netpoc_error

    description = (
        'summit poc network test to measure network bandwidth for edge device \n '
        'based on speedtest-cli: https://github.com/sivel/speedtest-cli.\n'
        '---------------------------------------------------------------------\n'
        'ssh://https://github.com/open-horizon/examples ... netspeed...' )

    parser = ArgParser(description=description)
    # Give optparse.OptionParser an `add_argument` method for
    # compatibility with argparse.ArgumentParser
    try:
        parser.add_argument = parser.add_option
    except AttributeError:
        pass

    parser.add_argument('--policy', action='store_true',
                        help='apply service policy')
    parser.add_argument('--target', default='closest', type=str,
                        help='override server criteria: closest, fastest, random. Default closest')
    parser.add_argument('--verbose', action='store_true',
                        help='verbose output; send test results to std output')
    parser.add_argument('--file', action='store_true',
                        help='write test results to json file')
    parser.add_argument('--mqtt', action='store_true',
                        help='send test results to mqtt')

    options = parser.parse_args()
    if isinstance(options, tuple):
        args = options[0]
    else:
        args = options
    del options

    if (args.verbose):
        debug_flag = 1

    if (args.policy):
        policy_flag = 1

    if (args.mqtt):
        mqtt_flag = 1

    if args.target:
        target_server_criteria = args.target

    shutdown_event = threading.Event()
    signal.signal(signal.SIGINT, ctrl_c)  

    netpoc_init()
    
    if args.file:
        """ values for testing purposes """
        file_flag = 1
        try:
            json_filename = './netspeedresults.json'
            # Ensure that we can open the json dump file
            jsonfile = open(json_filename, 'w')
            jsonfile.close()
        except IOError:
            utils.print_('Could not open file... writing results to std output\n')
            file_flag = 0
            debug_flag = 1

    try:
        # Every time these run, they schedule a timer to run themselves again at the next interval
        pingstatus()
        speedtestscheduler()
    except KeyboardInterrupt:
        if debug_flag:
            utils.print_('\nCancelling...')
        netpoc_error = 'netx0005'  # unexpected interrupt

if __name__ == '__main__':
    main()

# vim:ts=4:sw=4:expandtab
