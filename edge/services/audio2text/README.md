# Horizon Audio to Text Service

By default this service subscribes to the mqtt topic "ova/audioheard" that the voice2audio topic publishes to. Once audio2text receives the audio file, it will use the offline audio to text conversion tool Pocket Sphinx and publish the deciphered text to the mqtt topic "ova/textheard"

