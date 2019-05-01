import speech_recognition as sr
import signal
import json
import os
import sys
import snowboydecoder
import base64

mqtt_host = os.environ['MQTT_HOST']   #"zhangl.mqtt"
mqtt_port = os.environ['MQTT_PORT'] #1883
mqtt_topic = os.environ['MQTT_AUDIO_TOPIC'] #"zhangl_mqtt_flac_audio_topic" for flac format
audio_output_format = os.environ['AUDIO_FORMAT'] #flac or wav

print("envs: ")
print("mqtt_host: %s, mqtt_port: %s, mqtt_topic: %s, audio_output_format: %s" %(mqtt_host, mqtt_port, mqtt_topic, audio_output_format) )

interrupted = False
print("start...", flush=True) 

def publish_to_mqtt(name):
    fn = None
    try:

        print("convert audio to audio data: ", flush=True) 
        recognizer = sr.Recognizer()
        with sr.AudioFile(name) as src:
            recognizer.adjust_for_ambient_noise(src)
            
            aud = recognizer.listen(src)

            audio_data = {}
            if audio_output_format == "flac":
                byte_content = base64.b64encode(aud.get_flac_data())
                string_content = byte_content.decode('utf-8')
                audio_data = {
                    "format":"flac",
                    "content": string_content
                }

            if audio_output_format == "wav":
                byte_content = base64.b64encode(aud.get_wav_data())
                string_content = byte_content.decode('utf-8')
                audio_data = {
                    "format":"wav",
                    "content": string_content
                }

            audio_json = json.dumps(audio_data)

            #sending to mqtt
            cmd = "mosquitto_pub -d -h %s -p %s -t %s -m %s"%(mqtt_host, mqtt_port, mqtt_topic, json.dumps(audio_json))
            print("publish to mqtt using command: " + cmd, flush=True)
            os.system(cmd)
            
    except Exception as ex:
        print("Method failed with trace: " + ex.__traceback__, flush=True)
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

detector.start(audio_recorder_callback=publish_to_mqtt,
            interrupt_check=interrupt_callback,
            sleep_time=0.03)

detector.terminate()