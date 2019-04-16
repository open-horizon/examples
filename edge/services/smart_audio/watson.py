from watson_developer_cloud import SpeechToTextV1, WatsonApiException
from os.path import join, dirname
import speech_recognition as sr
import signal
import json
import os
import sys
import tempfile
import snowboydecoder

msghub_username = "*****"
msghub_password = "*****"
msghub_broker_url = "kafka03-prod02.messagehub.services.us-south.bluemix.net:9093,kafka02-prod02.messagehub.services.us-south.bluemix.net:9093,kafka04-prod02.messagehub.services.us-south.bluemix.net:9093,kafka05-prod02.messagehub.services.us-south.bluemix.net:9093,kafka01-prod02.messagehub.services.us-south.bluemix.net:9093"
msghub_topic = "zhangl_us.ibm.com.IBM_cpu2msghub"

remove_stopword = True
interrupted = False

stopwordlist = ['the', 'a', 'for', 'of', 'on', 'at', 'as', 'are', 'am', 'is', 'before', 'but', 'do', 'by', 'in', 'with']

speech_to_text = SpeechToTextV1(
    iam_apikey='*****',
    url='https://stream.watsonplatform.net/speech-to-text/api'
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
                
                # sending to kafka
                j = json.dumps({"text": text})
                print("json: " + j + ", sending json to kafka topic: " + msghub_topic, flush=True)
                _, fn = tempfile.mkstemp(dir='.')
                with open(fn, 'w') as f:
                    json.dump(j, f)
                cmd = "cat %s | kafkacat -P -b %s -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=%s -X sasl.password=%s -t %s"%(fn, msghub_broker_url, msghub_username, msghub_password, msghub_topic)
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