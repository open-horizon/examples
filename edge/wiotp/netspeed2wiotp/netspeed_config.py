""" netspeed configuration file """
import os
import ssl
import string

# These are potentially overridden by cmd line args
debug_flag = 0
policy_flag = 0
file_flag = 0
mqtt_flag = 0

# netspeed config constants
NETSPEED_MQTT_ID = 1

DEFAULT_LAT = '0.0001'
DEFAULT_LON = '0.0001'
DEFAULT_CONTRACT_ADDR = '0x0000000000000000000000000000000000000000'

# Get edge device env var's for WIoTP publish
if 'WIOTP_DOMAIN' in os.environ:
   wiotp_domain = os.getenv('WIOTP_DOMAIN')
else:
   print("Environment variable 'WIOTP_DOMAIN' not found.")

if 'WIOTP_ORG_ID' in os.environ:
   wiotp_org_id = os.getenv('WIOTP_ORG_ID')
else:
   print("Environment variable 'WIOTP_DOMAIN' not found.")

if 'WIOTP_DEVICE_TYPE' in os.environ:
   wiotp_device_type = os.getenv('WIOTP_DEVICE_TYPE')
else:
   print("Environment variable 'WIOTP_DEVICE_TYPE' not found.")

if 'WIOTP_DEVICE_AUTH_TOKEN' in os.environ:
   wiotp_device_auth_token = os.getenv('WIOTP_DEVICE_AUTH_TOKEN')
else:
   print("Environment variable 'WIOTP_DEVICE_AUTH_TOKEN' not found.")

# Get edge device latitude, longitude, and contract address
if 'HZN_LAT' in os.environ:
   latitude = os.getenv('HZN_LAT')
else:
   latitude = DEFAULT_LAT

if 'HZN_LON' in os.environ:
   longitude = os.getenv('HZN_LON')
else:
   longitude = DEFAULT_LON

if 'HZN_AGREEMENTID' in os.environ:
   contract_id = os.getenv('HZN_AGREEMENTID')  
else:
   contract_id = DEFAULT_CONTRACT_ADDR

if 'HZN_HASH' in os.environ:
   contract_nonce = os.getenv('HZN_HASH')
else:
   contract_nonce = ''

if 'HZN_DEVICE_ID' in os.environ:       
    device_id = os.getenv('HZN_DEVICE_ID')
else:
    device_id = get_serial()

def get_serial():
    # Extract serial from cpuinfo file
    cpuserial = "0000000000000000"
    try:
        f = open('/proc/cpuinfo','r')
        for line in f:
            if line[0:6]=='Serial':
               cpuserial = line[10:26]
        f.close()
    except:
        cpuserial = "ERROR000000000"
    
    if debug_flag:
        print_('processor s/n: %s' % cpuserial)

    return cpuserial

# MQTT parameters
mqtt_port=8883
mqtt_broker = '.'.join([wiotp_org_id, "messaging", wiotp_domain])
mqtt_client_id  = ':'.join(['d', wiotp_org_id, wiotp_device_type, device_id])
mqtt_ca_file = '/messaging.pem'
mqtt_auth = {'username': 'use-token-auth', 'password': wiotp_device_auth_token}
mqtt_tls = {'ca_certs': mqtt_ca_file}  #'tls_version': ssl.PROTOCOL_TLSv1 # Do not spec TLS w/ WIOTP  

def getEnvInt(name, default):
    """Return the named env var value as an int, or the default value."""
    if name in os.environ:
        strVal = os.getenv(name)
        try:
            return int(strVal)
        except ValueError as e:
            print_('Error: invalid value for environment variable %s: %s. Using default value %d.' % (name, str(e), default) )
            return default
    else:
        return default


if 'HZN_TARGET_SERVER' in os.environ:
    target_server_criteria = os.getenv('HZN_TARGET_SERVER')
else:
    target_server_criteria = 'closest'
target_server_criteria = string.lower(target_server_criteria)

run_interval = getEnvInt('RUN_INTERVAL', 25200)           # by default do speed test every 7 hours
max_volume_MB_month = getEnvInt('MONTHLY_CAP_MB', 30000)  # default max bandwidth used of 30 GB
PING_INTERVAL = getEnvInt('PING_INTERVAL', 300)           # seconds
REG_MAX_RETRIES = getEnvInt('REG_MAX_RETRIES', 4)         # constants for registration and send
REG_RETRY_DELAY = getEnvInt('REG_RETRY_DELAY', 3)         # seconds
REG_SUCCESS_SLEEP = getEnvInt('REG_SUCCESS_SLEEP', 20)    # number of seconds to sleep to allow the registration to take hold
SEND_MAX_RETRIES = getEnvInt('SEND_MAX_RETRIES', 5)
SEND_RETRY_DELAY = getEnvInt('SEND_RETRY_DELAY', 3)
SPEEDTEST_MAX_RETRIES = getEnvInt('SPEEDTEST_MAX_RETRIES', 5)
SPEEDTEST_RETRY_DELAY = getEnvInt('SPEEDTEST_RETRY_DELAY', 5)

# Initialize some global vars
total_ul_MB_month = 0.0
total_dl_MB_month = 0.0
total_bw_hour = 0.0
max_bw_hour = 0.0

receive_policy_exceeded = 0
send_policy_exceeded = 0
bandwidth_per_hour_exceeded = 0
max_mbps_exceeded = 0
