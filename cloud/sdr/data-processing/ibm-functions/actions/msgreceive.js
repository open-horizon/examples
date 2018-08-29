/*eslint-env node*/
//const Client = require('node-rest-client').Client;  // <- don't need this because using the Nodejs SDK
const Promise = require('promise');
const SpeechToTextV1 = require('watson-developer-cloud/speech-to-text/v1');
const NaturalLanguageUnderstandingV1 = require('watson-developer-cloud/natural-language-understanding/v1.js');
const stream = require('stream');
const async = require('async');
const minConfidence = 0.5;

// These needs to be global vars, so the sttTranscribe() and nluSentiment() functions can access them
var speechToText = null;
var naturalLangUnderstanding = null;

function main(params) {
  const date = new Date();
  console.log("msgreceive.js invoked at: "+date.toLocaleString());
  // console.log(params);
  if (params.messages && params.messages[0] && params.messages[0].value) {
    console.log(params.messages.length + " msgs received");
  } else {
    console.log("did not receive any messages!");
  }
  return new Promise((resolve, reject) => {
    if (!params.messages || !params.messages[0] || !params.messages[0].value) {
      reject("Invalid arguments. In the params object you must include 'messages' JSON array with 'value' field");
    }

    speechToText = new SpeechToTextV1({
      username: params.watsonSttUsername,
      password: params.watsonSttPassword
    });
    naturalLangUnderstanding = new NaturalLanguageUnderstandingV1({
      username: params.watsonNluUsername,
      password: params.watsonNluPassword,
      version: '2018-03-16'
    });

    // Process each msg. The processing of each msg is asynchronous, so use the async module to wait for them all before we complete our promise.
    //todo: Not sure if we should process them in paralell (eachOf()) or sequentially (eachOfSeries()). If there are a lot of msgs, paralell may
    //  exceed the max number of simultaneous calls to STT, but sequential may exceed the timeout of our action.
    async.eachOfSeries(params.messages, sttTranscribe, function (error) {
      if (error) {
        console.log(error);
        reject(error);
      } else {
        resolve({ "result": params.messages.length+" Message(s) from IBM Message Hub processed successfully" });
      }
    });
    
  });   // close of the promise
}   // end of main


// sttTranscribe converts 1 audio clip to text. It is called by async.eachOfSeries() so its function signature is determined by that.
function sttTranscribe(message, index, callback) {
  const msg = message.value
  const audioDecoded = Buffer.from(msg.audio, 'base64');
  console.log("Transcribing msg " + (index+1) + " from " + msg.devID + ', audio length ' + audioDecoded.length + ", at " + new Date(msg.ts*1000) + '...');

  const recognizeParams = {
    audio: buffer2stream(audioDecoded),
    content_type: 'audio/ogg',
    max_alternatives: 1,  // default is 1
  };

  speechToText.recognize(recognizeParams, function(error, sttResults) {
    if (error) {
      console.log(error);
      callback(error);
    } else {
      // const resultStr = JSON.stringify(sttResults, null, 2);
      console.log("Processing Watson STT results...");
      // console.log(resultStr);
      // console.log(sttResults);

      // Process each result. The processing of each result is asynchronous, so use the async module to wait for them all before we call our callback.
      //todo: Not sure if we should process them in paralell (eachOf()) or sequentially (eachOfSeries()). These calls are quick, so sequential is probably fine for now.
      async.eachOfSeries(sttResults.results, nluSentiment, function (error) {
        if (error) {
          console.log(error);
          callback(error);
        } else {
          callback();
        }
      });

    }
  });
}


// nluSentiment does sentiment analysis on text. It is called by async.eachOfSeries() so its function signature is determined by that.
function nluSentiment(sttResult, index, callback) {
    const r = sttResult;
    // Note: we only ask for 1 alternative
    if (r.final && r.alternatives[0].confidence > minConfidence) {
      // Run sentiment analysis on this text
      console.log('Analyzing sentiment of result '+(index+1)+', alternative: '+r.alternatives[0].transcript);
      const analyzeParams = {
        text: r.alternatives[0].transcript,
        features: {
          entities: { sentiment: true },
          keywords: { sentiment: true }
        }
      };
      naturalLangUnderstanding.analyze(analyzeParams, function(error, nluResults) {
        if (error) {
          console.log(error);
          callback(error);
        } else {
          const resultStr = JSON.stringify(nluResults, null, 2);
          console.log('Sentiment analysis results:');
          console.log(resultStr);
          callback();
        }
      });
    } else {
      console.log('Skipping alternative: Final: '+r.final+', Confidence: '+r.alternatives[0].confidence+', Text: '+r.alternatives[0].transcript);
      callback();
    }
}

// Convert the buffer to a stream
function buffer2stream(buffer) {
  const audioStream = new stream.Readable();
  audioStream._read = () => {}; // _read is required but we noop it because we will push data into it
  audioStream.push(buffer);
  audioStream.push(null);
  return audioStream
}

exports.main = main;

      /* just kept for reference....
      speechToText.listModels(null, function(error, speechModels) {
        if (error) {
          console.log(error);
          reject(error);
        } else {
          // const resolveStr = JSON.stringify(speechModels.models[0], null, 2);
          console.log("Result from Watson Speech to Text Service:")
          console.log(speechModels.models[0])
          resolve({ "result": "Message from IBM Message Hub processed successfully" });
        }
      }); */

      /* This is how we would use the watson stt rest api...
      const options_auth = { user: params.watsonSttUsername, password: params.watsonSttPassword };
      const client = new Client(options_auth);
      client.get("https://stream.watsonplatform.net/speech-to-text/api/v1/models", function (data, response) {
        console.log(data);     // parsed response body as js object
        console.log(response);   // raw response
      });
      */
