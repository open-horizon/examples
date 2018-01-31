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
hzn_organization = utils.check_env_var('HZN_ORGANIZATION')
wiotp_device_auth_token = utils.check_env_var('WIOTP_DEVICE_AUTH_TOKEN')
hzn_device_id = utils.check_env_var('HZN_DEVICE_ID', utils.get_serial())

# For device type and device id, there are 3 cases, checked for below:
#  1) The workload is run in wiotp/horizon and will be sending mqtt as the gw: 
#     HZN_DEVICE_ID contains class id, device type and id, and WIOTP_CLASS_ID, 
#     WIOTP_DEVICE_TYPE and WIOTP_DEVICE_ID are blank or '-'
#  2) The workload is run in wiotp/horizon and will be sending mqtt as a 
#     device: WIOTP_CLASS_ID, WIOTP_DEVICE_TYPE and/or WIOTP_DEVICE_ID have 
#     real values
#  3) The workload is run in a non-wiotp horizon instance: HZN_DEVICE_ID is 
#     the simple device id (use that), and WIOTP_CLASS_ID and WIOTP_DEVICE_TYPE
#     have values
# The way we will handle this cases is if WIOTP_CLASS_ID, WIOTP_DEVICE_TYPE 
#     and/or WIOTP_DEVICE_ID are set, they will override what we can parse from 
#     HZN_DEVICE_ID

# Case 2: WIOTP-specific variables (convenient to set first)
class_id = utils.check_env_var('WIOTP_CLASS_ID')
device_type = utils.check_env_var('WIOTP_DEVICE_TYPE')
device_id = utils.check_env_var('WIOTP_DEVICE_ID')

# Case 1: Workload deployed by WIoTP-Horizon; HZN_DEVICE_ID ~= 'g@mygwtype@mygw'.
if hzn_device_id.find('@'):
    ids = hzn_device_id.split('@')
    if len(ids) == 3:
        class_id, device_type, device_id = ids
    else:
        utils.print_("Workload config.py: class id, device_type, device_id could \
         not be set from HZN_DEVICE_ID. Possibly malformed env var.")
else: 
    # Case 3: When this workload is run in a non-wiotp horizon instance, 
    #  HZN_DEVICE_ID will be a simple device id
    if device_id != '' and device_id != '-':
        device_id = hzn_device_id

# Checks
if class_id == '' or class_id == '-':
    utils.print_("Workload config.py: class id could not be set in WIOTP_CLASS_ID \
     or HZN_DEVICE_ID.")
#    return 1
elif class_id != 'd' and class_id != 'g':
    utils.print_("Workload config.py: class id could not be set. Class ID can only \
        have value of 'g' or 'd'.")
#    return 1
if device_type == '' or class_id == '-':
    utils.print_("Workload config.py: device type not set in WIOTP_DEVICE_TYPE or \
        HZN_DEVICE_ID.")
#    return 1
if device_id == '' or device_id == '-':
    utils.print_("Workload config.py: device ID not set in WIOTP_DEVICE_ID or \
        HZN_DEVICE_ID.")
#    return 1

utils.print_("Optional override environment variables:")
utils.print_("  WIOTP_CLASS_ID=" + utils.check_env_var('WIOTP_CLASS_ID', '', False))
utils.print_("  WIOTP_DEVICE_TYPE=" + utils.check_env_var('WIOTP_DEVICE_TYPE', '', False))
utils.print_("  WIOTP_DEVICE_ID=" + utils.check_env_var('WIOTP_DEVICE_ID', '', False))
utils.print_("Derived variables:")
utils.print_("  CLASS_ID=" + class_id)
utils.print_("  DEVICE_TYPE=" + device_type)
utils.print_("  DEVICE_ID=" + device_id)

## Environment variables that can optionally be set, or default
# set in the pattern deployment_overrides field if you need to override
wiotp_domain = utils.check_env_var('WIOTP_DOMAIN', 'internetofthings.ibmcloud.com', False)    

# api key to create WIOTP device if it doesn't exist
wiotp_api_key = utils.check_env_var('WIOTP_API_KEY', '', False)                               

# api token to create WIOTP device if it doesn't exist
wiotp_api_auth_token = utils.check_env_var('WIOTP_API_AUTH_TOKEN', '', False)                 

# the cert to verify the WIoTP MQTT cloud broker
wiotp_pem_file = utils.check_env_var('WIOTP_PEM_FILE', '/messaging.pem', False)               

# local IP or hostname to send mqtt msgs via WIoTP Edge Connector ms
wiotp_edge_mqtt_ip = utils.check_env_var('WIOTP_EDGE_MQTT_IP', '')                            

# by default publish via MQTT to WioTP every n seconds
reporting_interval = utils.getEnvInt('REPORTING_INTERVAL', 10)                            
      
utils.print_("Optional environment variables (or default values):")
utils.print_("  WIOTP_DOMAIN=" + wiotp_domain)
utils.print_("  WIOTP_API_KEY=" + wiotp_api_key)
utils.print_("  WIOTP_API_AUTH_TOKEN=" + wiotp_api_auth_token)
utils.print_("  WIOTP_PEM_FILE=" + wiotp_pem_file)
utils.print_("  WIOTP_EDGE_MQTT_IP=" + wiotp_edge_mqtt_ip)
utils.print_("  REPORTING_INTERVAL=" + str(reporting_interval))     

# If Watson IoT Platform API credentials are not provided, assume device exists in WIoTP
if wiotp_api_key != '' or wiotp_api_auth_token != '':
    utils.print_("Workload_config.py: Watson IoT Platform REST API credentials not provided.")
    utils.print_(" assuming type %s with ID %s already exists in WIoTP." % (device_type, device_id))
else:  # TODO: both creds provided; prep for Watson IoT Platform REST API Calls
    pass

## Set up additional device metadata
# Get edge device latitude, longitude, and contract address / nonce
longitude = utils.check_env_var('HZN_LON', DEFAULT_LON, False)
latitude = utils.check_env_var('HZN_LAT', DEFAULT_LAT, False)
contract_id = utils.check_env_var('HZN_AGREEMENTID', DEFAULT_CONTRACT_ADDR, False)
contract_nonce = utils.check_env_var('HZN_HASH', '', False)

# Derive MQTT parameters for direct send to WIoTP
mqtt_port=8883
mqtt_broker = '.'.join([hzn_organization, "messaging", wiotp_domain])
mqtt_client_id  = ':'.join(['g', hzn_organization, device_type, device_id])
mqtt_ca_file = wiotp_pem_file
mqtt_auth = {'username': 'use-token-auth', 'password': wiotp_device_auth_token}
mqtt_tls = {'ca_certs': mqtt_ca_file}#, 'tls_version': ssl.PROTOCOL_TLSv1} # Do not spec TLS w/ WIOTP  

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