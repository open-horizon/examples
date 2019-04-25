from watson_developer_cloud import SpeechToTextV1, WatsonApiException
from os.path import join, dirname
import speech_recognition as sr
import signal
import json
import os
import sys
import tempfile
import snowboydecoder

mqtt_host = os.environ['MQTT_HOST']   #"zhangl.mqtt"
mqtt_port = os.environ['MQTT_PORT'] #1883
mqtt_topic = os.environ['MQTT_TOPIC'] #"zhangl_mqtt_topic"

remove_stopword = False
remove_stopword_string = os.environ['REMOVE_SW']
if remove_stopword_string.lower() == "true":
    remove_stopword = True

stt_iam_apikey = os.environ['STT_IAM_APIKEY']
stt_url = os.environ['STT_URL']

print("envs: ")
print("mqtt_host: %s, mqtt_port: %s, mqtt_topic: %s, remove_stopword: %s, stt_iam_apikey: %s, stt_url: %s" %(mqtt_host, mqtt_port, mqtt_topic, remove_stopword, stt_iam_apikey, stt_url) )

interrupted = False
stopwordlist = ['the', 'a', 'for', 'of', 'on', 'at', 'as', 'are', 'am', 'is', 'before', 'but', 'do', 'by', 'in', 'with']

speech_to_text = SpeechToTextV1(
    iam_apikey=stt_iam_apikey,
    url=stt_url
)
print("start...", flush=True) 

def call_speech_to_text(name):
    fn = None
    try:
        #Invoke a Speech to Text method
        # files = ['audio-file1.flac', 'audio-file2.flac']
        # files = ['audio-file.flac']
        #for file in files:
            #with open(join(dirname(__file__), './audio-files', file),'rb') as audio_file:
        print("speech_to_text...", flush=True) 
        recognizer = sr.Recognizer()
        with sr.AudioFile(name) as src:
            recognizer.adjust_for_ambient_noise(src)
            
            aud = recognizer.listen(src)
            speech_recognition_results = speech_to_text.recognize(
                audio=aud.get_flac_data(),
                content_type='audio/flac'
            ).get_result()

            results = speech_recognition_results.get('results', [])

            if len(results) != 0:
                transcript = speech_recognition_results['results'][0]['alternatives'][0]['transcript']
                print("transcript: " + transcript, flush=True)

                # stopword removal
                if remove_stopword == True:
                    text = ' '.join(w for w in transcript.split() if w not in stopwordlist)
                else:
                    text = transcript
                print(text, flush=True)
                
                #sending to mqtt
                cmd = "mosquitto_pub -d -h %s -p %s -t %s -m \"%s\""%(mqtt_host, mqtt_port, mqtt_topic, text)
                print("publish to mqtt using command: " + cmd, flush=True)
                os.system(cmd)
    except WatsonApiException as ex:
        print("Method failed with status code " + str(ex.code) + ": " + ex.message, flush=True)
    finally:
        if fn:
            os.remove(fn)
        os.remove(name)


def interrupt_callback():
    global interrupted
    return interrupted

def signal_handler(signal, frame):
    global interrupted
    sys.exit(0)
    interrupted = True

signal.signal(signal.SIGINT, signal_handler)

model = "./model/Watson.pmdl"
detector = snowboydecoder.HotwordDetector(model, sensitivity=0.5)
print("Listening...", flush=True)


detector.start(audio_recorder_callback=call_speech_to_text,
            interrupt_check=interrupt_callback,
            sleep_time=0.03)

detector.terminate()