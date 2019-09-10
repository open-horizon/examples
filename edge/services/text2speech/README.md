# Horizon Text to Speech Service

By default this service subscribes to the mqtt topic "ova/result" that processtext publishes. Once text2speech receives the text it will use the "espeak" command to play the text it received over a speaker if one is plugged into the Raspberry Pi. This audio clip will either be the IP address of the host machine running the service, or "I don't understand what you said."
