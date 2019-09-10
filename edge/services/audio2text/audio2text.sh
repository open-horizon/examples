#Audio 2 Text
echo "[OVA]:STARTING AUDIO 2 TEXT SERVICE"


#get base64 encoded output from listenaudio topic
mosquitto_sub -h ibm.mqtt -t ova/audioheard -p 1883 | while read -r line
#| encoded
do
if [ ! -z "$line" ]; then
	echo $line > encoded

	base64 -d encoded > new.wav
	aplay -D plughw:0,0 new.wav
	echo "[OVA]:Converting Audio 2 Text OFFLINE using Poecket_Sphinx. Wait for few seconds."
	pocketsphinx_continuous -infile new.wav -logfn /dev/null > textcommand
	sleep 30
	echo "[OVA]: $(cat textcommand)"
	mosquitto_pub -h ibm.mqtt -t ova/textheard -p 1883 -f textcommand
	echo "[OVA]:SENDING TEXT TO PROCESSTEXT SERVICE"
fi
done
