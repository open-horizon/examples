/*eslint-env node*/
// Locally test msg-receive.js

const msgreceive = require('./msgreceive.js');
// var protobuf = require("protobufjs");	// Use this if using audiomsg.proto
var root = require("./audiomsg.js");

// Marshal the message value to protobufs format before giving it to msgreceive.js

/* Use this if using audiomsg.proto
protobuf.load("../../../../../edge/msghub/sdr2msghub/audiomsg.proto").then(function(root) {
 	var AudioMsg = root.lookupType("audiolib.AudioMsg");	// get our PB type */

// This is our object that will be PB marshaled and put in the message.value field
// console.log(Math.floor(Date.now()/1000))
var payload = {
	devID: "bpmac",
	lat: 42.1,
	lon: -73.0,
	freq: 97.8,
	expectedValue: 0.8,
	audio: Buffer.from("this is some audio data that could end up being pretty big"),
	ts: {seconds: Math.floor(Date.now()/1000)}
};
var errMsg = root.audiolib.AudioMsg.verify(payload);	// check that it conforms to our PB type
if (errMsg) throw Error(errMsg);
console.log("payload verified");

// Marshal our payload
var msg = root.audiolib.AudioMsg.create(payload);
var buffer = root.audiolib.AudioMsg.encode(msg).finish();

const params = {
	"messages": [
		// { "value": Buffer.from("this is my first msg").toString('base64') },
		// { "value": buffer.toString('base64') }
		{ "value": buffer }
	],
	"watsonSttUsername": process.env.STT_USERNAME,
	"watsonSttPassword": process.env.STT_PASSWORD
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
