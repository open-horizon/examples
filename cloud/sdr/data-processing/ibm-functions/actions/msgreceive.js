/*eslint-env node*/
//var Client = require('node-rest-client').Client;  // <- don't need this because using the Nodejs SDK
var Promise = require('promise');
// var protobuf = require("protobufjs");	// Use this if using audiomsg.proto
var root = require("./audiomsg.js");
// var SpeechToTextV1 = require('watson-developer-cloud/speech-to-text/v1');

function main(params) {
  var date = new Date();
  console.log("msgreceive.js invoked at: "+date.toLocaleString());
  // console.log(params);
  if (params.messages && params.messages[0] && params.messages[0].value) {
    console.log(params.messages.length + " msgs received, 1st msg length: " + params.messages[0].value.length);
  } else {
    console.log("did not receive any messages!");
  }
  return new Promise((resolve, reject) => {
    /* if (!params) { reject("no params given"); }
    else { resolve({ "result": "got a params" }); } */

    if (!params.messages || !params.messages[0] || !params.messages[0].value) {
      reject("Invalid arguments. In the params object you must include 'messages' JSON array with 'value' field");
    }

    /* Use this if using audiomsg.proto
    protobuf.load("../../../../../edge/msghub/sdr2msghub/audiomsg.proto").then(function(root) {
      var AudioMsg = root.lookupType("audiolib.AudioMsg");	// get our PB type */
    
      //todo: support multiple msgs
      // var msgBuf = Buffer.from(params.messages[0].value, 'base64')
      // var msg = AudioMsg.decode(msgBuf);
      console.log("decoding msg...")
      var msg = root.audiolib.AudioMsg.decode(params.messages[0].value);
      console.log("got msg from " + msg.devID + " at " + new Date(msg.ts.seconds*1000));
      console.log(msg)
      console.log(msg.audio.toString())
      resolve({"result": "got msg and PB decoded it"})

      /* const msgs = params.messages;
      for (let i = 0; i < msgs.length; i++) {
        const msg = msgs[i];
        console.log("msg "+i+": "+msg.value);
      } */

      /* var speechToText = new SpeechToTextV1({
        username: params.watsonSttUsername,
        password: params.watsonSttPassword
      });

      speechToText.listModels(null, function(error, speechModels) {
        if (error) {
          console.log(error);
          reject(error);
        } else {
          // const resolveStr = JSON.stringify(speechModels.models[0], null, 2);
          console.log("Result from Watson Speech to Text Service:")
          console.log(speechModels.models[0])
          resolve({ "result": "Message from IBM Message Hub processed successfully processed" });
        }
      }); */

      /* This is how we would use the watson stt rest api...
      var options_auth = { user: params.watsonSttUsername, password: params.watsonSttPassword };
      var client = new Client(options_auth);
      client.get("https://stream.watsonplatform.net/speech-to-text/api/v1/models", function (data, response) {
        console.log(data);     // parsed response body as js object
        console.log(response);   // raw response
      });
      */

    /* Use this if using audiomsg.proto
    }, function(err) {
      if (err) {
        console.log(err);
      }
    }); */
    
  });
}

exports.main = main;
