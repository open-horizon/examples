/*eslint-env node*/
// Locally test msg-receive.js

const msgreceive = require('./msgreceive.js');
// var protobuf = require("protobufjs");	// Use this if using audiomsg.proto
// var root = require("./audiomsg.js");
var fs = require('fs');		// to read the sample ogg file

// Read the sample ogg file and base64 encode it
const mockAudioFile = "../../../../edge/services/sdr/mock_audio.ogg" // if running it from the Makefile in cloud/sdr/data-processing/ibm-functions
var audioEncoded = fs.readFileSync(mockAudioFile, { encoding: 'base64' });
if (!audioEncoded) { audioEncoded = fs.readFileSync('../'+mockAudioFile, { encoding: 'base64' }); } // if running it directly from cloud/sdr/data-processing/ibm-functions/actions

/* Use this if using audiomsg.proto
protobuf.load("../../../../../edge/msghub/sdr2msghub/audiomsg.proto").then(function(root) {
 	var AudioMsg = root.lookupType("audiolib.AudioMsg");	// get our PB type */

// This is our object that will be put in the messages.value field
var payload = {
	devID: "IBM/bpmac",
	lat: 42.1,
	lon: -73.0,
	freq: 97.8,
	expectedValue: 0.8,
	audio: audioEncoded,
	// audio: Buffer.from("this is some audio data that could end up being pretty big").toString('base64'),
	// audio: "this is some audio data that could end up being pretty big",
	ts: Math.floor(Date.now()/1000)
};
/* var payload2 = {
	devID: "IBM/bpmac2",
	lat: 42.1,
	lon: -73.0,
	freq: 101.5,
	expectedValue: 0.8,
	audio: audioEncoded,
	ts: Math.floor(Date.now()/1000)
}; */
/* var errMsg = root.audiolib.AudioMsg.verify(payload);	// check that it conforms to our PB type
if (errMsg) throw Error(errMsg);
console.log("payload verified"); */

// Marshal our payload
/* var msg = root.audiolib.AudioMsg.create(payload);
var buffer = root.audiolib.AudioMsg.encode(msg).finish(); */
const params = {
	"messages": [
		// { "value": Buffer.from("this is my first msg").toString('base64') },
		// { "value": buffer.toString('base64') }
		// { "value": buffer }
		{ "value": payload },
		// { "value": payload2 },
	],
	"watsonSttUsername": process.env.STT_USERNAME,
	"watsonSttPassword": process.env.STT_PASSWORD,
	"watsonNluUsername": process.env.NLU_USERNAME,
	"watsonNluPassword": process.env.NLU_PASSWORD,
	"postgresUrl": process.env.SDR_DB_URL,
}

// const result = msgreceive.main(params)
// console.log("msgreceive.main() result:")
// console.log(result)

msgreceive.main(params).then(function(response){
	console.log("msgreceive.main() result:");
	console.log(response);
}, function(error) {
	console.log(error.message);
});

/* Use this if using audiomsg.proto
}, function(err) {
	if (err) {
		console.log(err);
	}
}); */
