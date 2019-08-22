# Horizon Offline Voice Assistant Service

Every 60 seconds this service records a five second audio clip, converts the audio clip to text, finally having the host machine execute the command and speak the output.

This services depends on four lower level services: mqtt, voice2audio, audio2text, and text2speech. The voice2audio service records the five second audio clip and publishes it to the mqtt broker to the audio2text service that takes the audio clip and converts it to text offline using pocket sphinx, which is then sent to the processtext service. This service will take the text and attempt to execute the originally recorded command and send it to the text2speech service which will play the output through a speaker. 

As a proof of concept, this service can only speak the output to the command "What is your IP address?" Otherwise it will say "I don't understand what you said."

#### Processtext parameters:

Receives text over mqtt on topic "ova/textheard" 

Sends either the text output of the executed command or "I don't understand what you said" over the mqtt topic "ova/result"

## Using the Horizon Offline Voice Assistant Service 

In order to get the most out of this service you will need be watching the output using the command:
`tail -f /var/log/syslog | grep OVA` 

Once the service begins to run it will prompt you to say a command with the following message:

```[OVA]:STARTING VOICE 2 AUDIO SERVICE
[OVA]:Say some command like whats your ip address?```

If the audio clip is successfully recorded and sent to the audio2text service you will see:
```[OVA]:STARTING AUDIO 2 TEXT SERVICE
[OVA]:Test.wav file is created!!
[OVA]:Wav file is converted to base64encoded file
[OVA]:SENT AUDIO DATA TO AUDIO 2 TEXT SERVICE```

You will then hear the audio clip you recorded played back to you and the service will begin using Pocket Sphinx to convert the audio to text and display the following message to the terminal:

`[OVA]:Converting Audio 2 Text OFFLINE using Poecket_Sphinx. Wait for few seconds.`

Once that has completed you will see the converted text in the terminal as well:

```[OVA]: what is your ip address
[OVA]:SENDING TEXT TO PROCESSTEXT SERVICE```

Finally, if all the services executed successfully, you will hear your machine say your IP address thru the speaker. 



