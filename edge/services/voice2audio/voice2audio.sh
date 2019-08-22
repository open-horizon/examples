echo "[OVA]:STARTING VOICE 2 AUDIO SERVICE"
#recording voice command
echo "[OVA]:Say some command like whats your ip address?"
sleep 1
arecord -D sysdefault:CARD=0  -d 5 --format S16_LE --rate 16000 -c1 test.wav
echo "[OVA]:Test.wav file is created!!"
#converting to base64 encoding
base64 test.wav | tr -d \\n > encodeddata
echo "[OVA]:Wav file is converted to base64encoded file"

aplay test.wav

#publish to mqtt listenaudio topic
mosquitto_pub -h ibm.mqtt -t ova/audioheard -f encodeddata
echo "[OVA]:SENT AUDIO DATA TO AUDIO 2 TEXT SERVICE"

sleep 60
