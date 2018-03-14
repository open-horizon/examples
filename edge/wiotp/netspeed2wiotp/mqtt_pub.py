import time
import string
import json
import sys

import paho.mqtt.publish as publish
import paho.mqtt.client as mqtt

from workload_config import *   # Read configuration
import utils                    # Utilities file in this dir (utils.py) 

def post_networkdata_single_wiotp(jsonpayload, event_id, heart_beat=False):
    """Tries once to send network data in json format to WIoTP via mqtt. 
       Returns 1 if successful, 0 if not, -1 if failed because not registered.
    """

    try:
        retain = True
        qos = 2   # since speed data is sent so infrequently we can afford to make sure it gets there exactly once
        
        if debug_flag:
            utils.print_("mqtt_pub.py: Sending data to mqtt... \
                mqtt_topic=%s, mqtt_broker=%s, client_id=%s" % (mqtt_topic, mqtt_broker, mqtt_client_id))

        # Publish to MQTT
        publish.single(topic=mqtt_topic, payload=jsonpayload, qos=qos, hostname=mqtt_broker,
            protocol=mqtt.MQTTv311, client_id=mqtt_client_id, port=mqtt_port, #auth=mqtt_auth,
            tls=mqtt_tls, retain=retain)
        if debug_flag: utils.print_('mqtt_pub.py: Send to mqtt successful')
        return 1
    except:
        e = sys.exc_info()[1]
        if 'not authori' in str(e).lower() or 'bad user name or password' in str(e).lower():
            # The data send failed because we are not successfully registered
            return -1
        else:
            utils.print_('Send to mqtt failed: %s' % e)
            return 0