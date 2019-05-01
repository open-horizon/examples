from watson_developer_cloud import SpeechToTextV1, WatsonApiException
import json
import urllib.request
import os
import shlex
import sys
import base64

from subprocess import Popen, PIPE

mqtt_host = os.environ['MQTT_HOST']   #"zhangl.mqtt"
mqtt_port = os.environ['MQTT_PORT'] #1883
mqtt_receive_topic = os.environ['MQTT_RECEIVE_TOPIC'] #"zhangl_mqtt_flac_audio_topic"
mqtt_send_topic = os.environ['MQTT_SEND_TOPIC'] #"zhangl_mqtt_send_topic"

remove_stopword = False
remove_stopword_string = os.environ['REMOVE_SW']
if remove_stopword_string.lower() == "true":
    remove_stopword = True
stopword_removal_host = os.environ['SW_HOST'] #127.0.0.1
stopword_removal_port = os.environ['SW_PORT'] #5002

stt_iam_apikey = os.environ['STT_IAM_APIKEY']
stt_url = os.environ['STT_URL']

print("envs: ")
print("mqtt_host: %s, mqtt_port: %s, mqtt_receive_topic: %s, mqtt_send_topic: %s, stt_iam_apikey: %s, stt_url: %s, remove_stopword: %s, stopword_removal_host: %s, stopword_removal_port: %s" %(mqtt_host, mqtt_port, mqtt_receive_topic, mqtt_send_topic, stt_iam_apikey, stt_url, remove_stopword, stopword_removal_host, stopword_removal_port) )

speech_to_text = SpeechToTextV1(
    iam_apikey=stt_iam_apikey,
    url=stt_url
)
print("start...", flush=True) 

def call_speech_to_text():
    try:
        #Invoke a Speech to Text method
        # files = ['audio-file1.flac', 'audio-file2.flac']
        # files = ['audio-file.flac']
        #for file in files:
            #with open(join(dirname(__file__), './audio-files', file),'rb') as audio_file:

        cmd = "mosquitto_sub -h %s -p %s -t %s"%(mqtt_host, mqtt_port, mqtt_receive_topic)
        print("receive from mqtt using command: %s"%cmd, flush=True)
        p = Popen(shlex.split(cmd), stdout=PIPE, stderr=PIPE)
        for l in iter(p.stdout.readline, ''):
            print("speech_to_text...", flush=True) 

            audio_json_string = l.decode('utf-8').strip('\n')
            try:
                audio_json_string = json.loads(audio_json_string)
            except:
                continue
            print("audio_json_string: %s"%audio_json_string, flush=True)

            audio_content_string = audio_json_string['content'].encode('utf-8')
            flac_audio=base64.b64decode(audio_content_string)

            speech_recognition_results = speech_to_text.recognize(
                audio=flac_audio,
                content_type='audio/flac'
            ).get_result()

            results = speech_recognition_results.get('results', [])

            if len(results) != 0:
                transcript = speech_recognition_results['results'][0]['alternatives'][0]['transcript']
                print("transcript: %s"%transcript, flush=True)

                # stopword removal
                if remove_stopword == True:
                    print("remove stopword: ", flush=True)
                    # curl -X POST http://127.0.0.1:5002/remove_stopword -H 'Content-Type: application/json' -H 'cache-control: no-cache' -d '{"text": "how are you today"}'
                    sw_url = "http://%s:%s/remove_stopword" %(stopword_removal_host, stopword_removal_port) #
                    post_data = {"text": transcript}
                    params = json.dumps(post_data).encode('utf-8')
                    print("calling out to remove stopword: post_url: %s, post_data: %s" %(sw_url, post_data))
                    
                    req = urllib.request.Request(sw_url, data=params, headers={'content-type': 'application/json'})
                    response = urllib.request.urlopen(req).read().decode('utf-8')
                    if response:
                        response_json = json.loads(response)
                        text = response_json['result']
                    else:
                        text = ''

                else:
                    text = transcript
                print(text, flush=True)

                #sending to mqtt
                cmd = "mosquitto_pub -h %s -p %s -t %s -m \"%s\""%(mqtt_host, mqtt_port, mqtt_send_topic, text)
                print("publish to mqtt using command: %s"%cmd, flush=True)
                os.system(cmd)
        
    except WatsonApiException as ex:
        print("Method failed with status code %d: %s"%(ex.code, ex.message), flush=True)

if __name__ == "__main__":
    call_speech_to_text()