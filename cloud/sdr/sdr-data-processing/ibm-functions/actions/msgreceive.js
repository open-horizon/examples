//var Client = require('node-rest-client').Client;
var SpeechToTextV1 = require('watson-developer-cloud/speech-to-text/v1');

function main(params) {
  var date = new Date();
  if (params.messages && params.messages[0] && params.messages[0].value) {
    console.log(params.messages.length + " messages received at: "+date.toLocaleString());
  } else {
    console.log("received message(s) at: "+date.toLocaleString());
  }
  return new Promise((resolve, reject) => {
    if (!params.messages || !params.messages[0] || !params.messages[0].value) {
      reject("Invalid arguments. Must include 'messages' JSON array with 'value' field");
    }
    const msgs = params.messages;
    //const cats = [];
    for (let i = 0; i < msgs.length; i++) {
      const msg = msgs[i];
      console.log("msg "+i+": "+msg.value);
      /* parse the msg json
      for (let j = 0; j < msg.value.cats.length; j++) {
        const cat = msg.value.cats[j];
        console.log(`A ${cat.color} cat named ${cat.name} was received.`);
        cats.push(cat);
      }
      */
    }

    /*
    var options_auth = { user: "8281900f-8621-43ae-b8a6-4656420bef9c", password: "OIZKxcVxxqF6" };
    var client = new Client(options_auth);
    client.get("https://stream.watsonplatform.net/speech-to-text/api/v1/models", function (data, response) {
      console.log(data);     // parsed response body as js object
      console.log(response);   // raw response
    });
    */

    var speechToText = new SpeechToTextV1({
      username: '8281900f-8621-43ae-b8a6-4656420bef9c',
      password: 'OIZKxcVxxqF6'
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
    });

  });
}

exports.main = main;
