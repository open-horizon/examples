# Horizon Voice to Audio Service

Every 60 seconds the Voice to Audio service records a five second audio clip, such as "What is your IP address?", and sends it via MQTT broker by default to the Audio to Text service. 

By default this service uses audio card 0 on the raspberry pi to record the five second audio clip that it saves as a file called test.wav. That file is then decoded using base64 and published to the mqtt topic "ova/audioheard" 
