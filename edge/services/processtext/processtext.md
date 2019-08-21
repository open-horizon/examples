# Horizon Offline Voice Assistant Service

Every 60 seconds this service records a five second audio clip, converts the audio clip to text, finally having the host machine execute the command and speak the output.

This services depends on four lower level services: mqtt, voice2audio, audio2text, and text2speech. The voice2audio service records the five second audio clip and publishes it to the mqtt broker to the audio2text service that takes the audio clip and converts it to text offline using pocket sphinx, which is then sent to the processtext service. This service will take the text and attempt to execute the originally recorded command and send it to the text2speech service which will play the output through a speaker. 

As a proof of concept, this service can only speak the output to the command "What is your IP address?" Otherwise it will say "I don't understand what you said."

#### Processtext parameters:

Receives text over mqtt on topic "ova/textheard" 

Sends either the text output of the executed command or "I don't understand what you said" over the mqtt topic "ova/result"