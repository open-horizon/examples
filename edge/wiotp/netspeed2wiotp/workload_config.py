""" Workload configuration file """
import os
import ssl
import string

import utils        # Utilities file in this dir (utils.py)

# These are potentially overridden by cmd line args
debug_flag = 0
mqtt_flag = 0
file_flag = 0

# Edge parameters, environment variable based
DEFAULT_LAT = '0.0001'
DEFAULT_LON = '0.0001'
DEFAULT_CONTRACT_ADDR = '0x0000000000000000000000000000000000000000'


## Get edge device env var's for WIoTP publish
# Check primary env vars first
hzn_organization = utils.check_env_var('HZN_ORGANIZATION')     # automatically passed in by Horizon
#wiotp_device_auth_token = utils.check_env_var('WIOTP_DEVICE_AUTH_TOKEN')  # note: this is no longer needed because we can now send msgs an an app to edge-connector unauthenticated, as long as we are local.
hzn_device_id = utils.check_env_var('HZN_DEVICE_ID', utils.get_serial())      # automatically passed in by Horizon. Wiotp automatically gives this a value of: g@mygwtype@mygw

# When the Workload is deployed by WIoTP-Horizon; HZN_DEVICE_ID ~= 'g@mygwtype@mygw'.
ids = hzn_device_id.split('@')
if len(ids) == 3:
    class_id, device_type, device_id = ids     # the class id is not actually used anymore
else:
    utils.print_("Error: HZN_DEVICE_ID must have the format: g@mygwtype@mygw")

#utils.print_("Workload config.py: Optional override environment variables:")
#utils.print_("Workload config.py:   WIOTP_CLASS_ID=" + utils.check_env_var('WIOTP_CLASS_ID', '', False))
#utils.print_("Workload config.py:   WIOTP_DEVICE_TYPE=" + utils.check_env_var('WIOTP_DEVICE_TYPE', '', False))
#utils.print_("Workload config.py:   WIOTP_DEVICE_ID=" + utils.check_env_var('WIOTP_DEVICE_ID', '', False))
utils.print_("Workload config.py: Derived variables:")
#utils.print_("Workload config.py:   CLASS_ID=" + class_id)
utils.print_("Workload config.py:   DEVICE_TYPE=" + device_type)
utils.print_("Workload config.py:   DEVICE_ID=" + device_id)

## Environment variables that can optionally be set, or default
# set in the pattern deployment_overrides field if you need to override
wiotp_domain = utils.check_env_var('WIOTP_DOMAIN', 'internetofthings.ibmcloud.com', False)    

# the cert to verify the WIoTP MQTT cloud broker
wiotp_pem_file = utils.check_env_var('WIOTP_PEM_FILE', '/messaging.pem', False)               

# local IP or hostname to send mqtt msgs via WIoTP Edge Connector ms
wiotp_edge_mqtt_ip = utils.check_env_var('WIOTP_EDGE_MQTT_IP', '')                            

# by default publish via MQTT to WioTP every n seconds
reporting_interval = utils.getEnvInt('REPORTING_INTERVAL', 10)                            
      
utils.print_("Workload config.py: Optional environment variables (or default values):")
utils.print_("Workload config.py:   WIOTP_DOMAIN=" + wiotp_domain)
utils.print_("Workload config.py:   WIOTP_PEM_FILE=" + wiotp_pem_file)
utils.print_("Workload config.py:   WIOTP_EDGE_MQTT_IP=" + wiotp_edge_mqtt_ip)
utils.print_("Workload config.py:   REPORTING_INTERVAL=" + str(reporting_interval))     

## Set up additional device metadata
# Get edge device latitude, longitude, and contract address / nonce
longitude = utils.check_env_var('HZN_LON', DEFAULT_LON, False)
latitude = utils.check_env_var('HZN_LAT', DEFAULT_LAT, False)
contract_id = utils.check_env_var('HZN_AGREEMENTID', DEFAULT_CONTRACT_ADDR, False)
contract_nonce = utils.check_env_var('HZN_HASH', '', False)

# Derive MQTT parameters for direct send to WIoTP
mqtt_port=8883
if wiotp_edge_mqtt_ip != '':        # Case 1: Send via local broker
    mqtt_broker = wiotp_edge_mqtt_ip
else:                               # Case 2: Send directly to WIoTP
    mqtt_broker = '.'.join([hzn_organization, "messaging", wiotp_domain])
#mqtt_client_id  = ':'.join(['g', hzn_organization, device_type, device_id])
mqtt_client_id  = ':'.join(['a', hzn_organization, device_type+device_id])
mqtt_ca_file = wiotp_pem_file
#mqtt_auth = {'username': 'use-token-auth', 'password': wiotp_device_auth_token}
mqtt_tls = {'ca_certs': mqtt_ca_file}#, 'tls_version': ssl.PROTOCOL_TLSv1} # Do not spec TLS w/ WIOTP  
#mqtt_topic = '/'.join(['iot-2/type', device_type, 'id', device_id, 'evt', event_id, 'fmt/json'])
mqtt_topic = 'iot-2/evt/status/fmt/json'

# Optional Netspeed-specific settings. These define the run "policy", and may be set via env var's
if 'HZN_TARGET_SERVER' in os.environ:
    target_server_criteria = os.getenv('HZN_TARGET_SERVER')
else:
    target_server_criteria = 'closest'
target_server_criteria = string.lower(target_server_criteria)

PING_INTERVAL = reporting_interval                          # seconds, for ping only
run_interval = utils.getEnvInt('RUN_INTERVAL', 25200)             # by default do speed test every 7 hours
max_volume_MB_month = utils.getEnvInt('MONTHLY_CAP_MB', 30000)    # default max bandwidth used of 30 GB
REG_MAX_RETRIES = utils.getEnvInt('REG_MAX_RETRIES', 4)           # constants for registration and send
REG_RETRY_DELAY = utils.getEnvInt('REG_RETRY_DELAY', 3)           # seconds
REG_SUCCESS_SLEEP = utils.getEnvInt('REG_SUCCESS_SLEEP', 20)      # number of seconds to sleep to allow the registration to take hold
SEND_MAX_RETRIES = utils.getEnvInt('SEND_MAX_RETRIES', 5)
SEND_RETRY_DELAY = utils.getEnvInt('SEND_RETRY_DELAY', 3)
SPEEDTEST_MAX_RETRIES = utils.getEnvInt('SPEEDTEST_MAX_RETRIES', 5)
SPEEDTEST_RETRY_DELAY = utils.getEnvInt('SPEEDTEST_RETRY_DELAY', 5)

# Initialize some global vars
total_ul_MB_month = 0.0
total_dl_MB_month = 0.0
total_bw_hour = 0.0
max_bw_hour = 0.0

receive_policy_exceeded = 0
send_policy_exceeded = 0
bandwidth_per_hour_exceeded = 0
max_mbps_exceeded = 0