/*eslint-env node*/
const Promise = require('promise');
const SpeechToTextV1 = require('watson-developer-cloud/speech-to-text/v1');
const NaturalLanguageUnderstandingV1 = require('watson-developer-cloud/natural-language-understanding/v1.js');
const stream = require('stream');
// const asyncmod = require('async');   // https://caolan.github.io/async/
const { Pool } = require('pg');
const minConfidence = 0.5;

// These needs to be global vars, so the sttTranscribe() and nluSentiment() functions can access them
var db = null;
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

    db = connect2db(params.postgresUrl);
    /* db.query('select noun, sentiment, numberofmentions from globalnouns').then((res) => {
      console.log('globalnouns table:');
      console.log(res.rows);
    }).catch((e) => console.error(e)); */
  
    speechToText = new SpeechToTextV1({
      username: params.watsonSttUsername,
      password: params.watsonSttPassword
    });
    naturalLangUnderstanding = new NaturalLanguageUnderstandingV1({
      username: params.watsonNluUsername,
      password: params.watsonNluPassword,
      version: '2018-04-05'
    });

    // Process each msg. The processing of each msg is asynchronous (a promise), so use Promise.all() to wait for them all before we complete our promise.
    //todo: Not sure if we should process them in paralell or sequentially. If there are a lot of msgs, paralell may
    //  exceed the max number of simultaneous calls to STT, but sequential may exceed the timeout of our action.
    var promises = [];
    params.messages.forEach(m => {
      promises.push(sttTranscribe(m));
    });

    // Promise.all(params.messages.map(sttTranscribe));
    Promise.all(promises)
      .then(() => { resolve({ "result": params.messages.length+" Message(s) from IBM Message Hub processed successfully" }) })
      .catch((err) => { reject(err) });
  
    // Process each msg. The processing of each msg is asynchronous, so use the async module to wait for them all before we complete our promise.
    /* asyncmod.eachOfSeries(params.messages, sttTranscribe, function (error) {
      if (error) {
        console.log(error);
        reject(error);
      } else {
        resolve({ "result": params.messages.length+" Message(s) from IBM Message Hub processed successfully" });
      }
    }); */

  });   // close of the promise
}   // end of main


// sttTranscribe converts 1 audio clip to text. It is called by asyncmod.eachOfSeries() so its function signature is determined by that.
function sttTranscribe(message) {
  return new Promise((resolve, reject) => {
    const msg = message.value
    const audioDecoded = Buffer.from(msg.audio, 'base64');
    console.log("Transcribing msg from " + msg.devID + ', audio length ' + audioDecoded.length + ", at " + new Date(msg.ts*1000) + '...');

    const recognizeParams = {
      audio: buffer2stream(audioDecoded),
      content_type: 'audio/ogg',
      max_alternatives: 1,  // default is 1
    };

    speechToText.recognize(recognizeParams, function(error, sttResults) {
      if (error) {
        console.log(error);
        return reject(error);
      } else {
        console.log("Processing Watson STT results...");
        /* const resultStr = JSON.stringify(sttResults, null, 2);
        console.log(resultStr);
        return callback(); */
        // console.log(sttResults);

        // Process each result. The processing of each result is asynchronous, so use Promise.all() to wait for them all before we resolve our promise.
        //todo: Not sure if we should process them in paralell (eachOf()) or sequentially (eachOfSeries()), but parallel is probably fine for now.
        var promises = [];
        sttResults.results.forEach(r => {
          promises.push(nluSentiment(r, msg.ts, msg.devID));
        });
        Promise.all(promises)
          .then(() => {return addNodeStationToDB(db, msg.ts, msg.devID, msg.freq, msg.lat, msg.lon, msg.expectedValue)})
          .then(() => { resolve() })
          .catch((err) => { reject(err) });

        /* asyncmod.eachOfSeries(sttResults.results, nluSentiment, function (error) {
          if (error) {
            console.log(error);
            return reject(error);
          } else {
            return resolve();
          }
        }); */

      }
    });
  });
}


// nluSentiment does sentiment analysis on text. It is called by asyncmod.eachOfSeries() so its function signature is determined by that.
function nluSentiment(sttResult, timeStamp, nodeID) {
  return new Promise((resolve, reject) => {
    const r = sttResult;
    // Note: we only ask for 1 alternative
    if (r.final && r.alternatives[0].confidence > minConfidence) {
      // Run sentiment analysis on this text
      console.log('Analyzing sentiment of result, alternative: '+r.alternatives[0].transcript);
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
          return reject(error);
        } else {
          console.log('Processing sentiment analysis results...');
          // const resultStr = JSON.stringify(nluResults, null, 2);
          // console.log(resultStr);
          addSentimentsToDB(db, nluResults, timeStamp, nodeID)
            .then(() => { resolve() })
            .catch((err) => { reject(err) });
        }
      });
    } else {
      console.log('Skipping alternative: Final: '+r.final+', Confidence: '+r.alternatives[0].confidence+', Text: '+r.alternatives[0].transcript);
      return resolve();
    }
  });
}

// addSentimentsToDB adds the nouns, sentiments, and node id to DB tables, and returns the combined db promise
function addSentimentsToDB(db, sentiments, timeStamp, nodeID) {
  console.log(`adding the nouns and sentiments from node id ${nodeID} to DB tables...`);
  const ts = seconds2timestamp(timeStamp);
  const ents = sentiments.entities.concat(sentiments.keywords);  // concat the 2 lists
  var dups = {};  // we have to filter out duplicates
  console.log('length of ents: ' + ents.length);
  //for (let e of ents) {}  // use lodash.uniqueBy()
  const entities = ents.filter((e) => {
    if (dups[e.text]) {return false} dups[e.text] = true; return true
  });
  console.log('length of entities: ' + entities.length);
  // console.log(`timestamp is ${ts}`);

  // For all of the entities and keywords add their nouns and sentiments to the DB
  var promises = [];
  entities.forEach(e => {
    promises.push(addOneSentimentToDB(e, db, ts, nodeID));
  });
  return Promise.all(promises);
}

// addOneSentimentToDB adds 1 sentiment to the globalnouns and nodenouns tables, and returns the combined db promise
function addOneSentimentToDB(entity, db, ts, nodeID) {
  const e = entity;
  const noun = e.text;
  const sentiment = e.sentiment.score;
  console.log(`adding noun ${noun} with sentiment score ${sentiment} to db...`);

  // Add this noun/sentiment to the globalnouns table and the noun/sentiment to the nodenouns table
  // This is the postgres way to upsert a row (insert if not there, update if there)
  // stmt, err := db.Prepare("INSERT INTO globalnouns VALUES ($1, $2, 1, $3) ON CONFLICT (noun) DO UPDATE SET sentiment = ((globalnouns.sentiment * globalnouns.numberofmentions) + $2) / (globalnouns.numberofmentions + 1), numberofmentions = globalnouns.numberofmentions + 1, timeupdated = $3"); _, err = stmt.Exec(noun, sentiment, ts)
  // stmt, err = db.Prepare("INSERT INTO nodenouns VALUES ($1, $4, $2, 1, $3) ON CONFLICT ON CONSTRAINT nodenouns_pkey DO UPDATE SET sentiment = ((nodenouns.sentiment * nodenouns.numberofmentions) + $2) / (nodenouns.numberofmentions + 1), numberofmentions = nodenouns.numberofmentions + 1, timeupdated = $3"); _, err = stmt.Exec(noun, sentiment, ts, nodeID)

  return db.query("INSERT INTO globalnouns VALUES ($1, $2, 1, $3) ON CONFLICT (noun) DO UPDATE SET sentiment = ((globalnouns.sentiment * globalnouns.numberofmentions) + $2) / (globalnouns.numberofmentions + 1), numberofmentions = globalnouns.numberofmentions + 1, timeupdated = $3", [noun, sentiment, ts])
  .then((result) => {
    console.log(`globalnouns table: inserted/updated ${result.rowCount} rows`);
    return db.query("INSERT INTO nodenouns VALUES ($1, $4, $2, 1, $3) ON CONFLICT ON CONSTRAINT nodenouns_pkey DO UPDATE SET sentiment = ((nodenouns.sentiment * nodenouns.numberofmentions) + $2) / (nodenouns.numberofmentions + 1), numberofmentions = nodenouns.numberofmentions + 1, timeupdated = $3", [noun, sentiment, ts, nodeID])
  }, () => console.error('reject!!!!!!!'))
  .then((result) => {console.log(`nodenouns table: inserted/updated ${result.rowCount} rows`)})
  .catch((err) => {console.error(err)} );
  //return new Promise(() => { console.log('here!!!!!!!!!!!!!!!!!!!!!')});
}

// addNodeStationToDB adds the node and station info to DB tables. This is only called once per msg hub msg. Returns the db promise.
function addNodeStationToDB(db, timeStamp, nodeID, stationFreq, latitude, longitude, expectedValue) {
  console.log("adding the node and station info to DB tables...")
  const ts = seconds2timestamp(timeStamp);
  // console.log(`timestamp is ${ts}`);

	// Add station and node info to the db tables. Chain the promises together and return the chain.
	// from the golang version: stmt, err := db.Prepare("INSERT INTO stations VALUES ($1, $2, 1, $3, $4) ON CONFLICT ON CONSTRAINT stations_pkey DO UPDATE SET numberofclips = stations.numberofclips + 1, dataqualitymetric =$3, timeupdated = $4"); _, err = stmt.Exec(nodeID, stationFreq, expectedValue, ts) //todo: not sure what to do with expectedValue
	//      stmt, err = db.Prepare("INSERT INTO edgenodes VALUES ($1, $2, $3, $4) ON CONFLICT (edgenode) DO UPDATE SET latitude = $2, longitude = $3, timeupdated = $4");  _, err = stmt.Exec(nodeID, latitude, longitude, ts)
  return db.query("INSERT INTO stations VALUES ($1, $2, 1, $3, $4) ON CONFLICT ON CONSTRAINT stations_pkey DO UPDATE SET numberofclips = stations.numberofclips + 1, dataqualitymetric =$3, timeupdated = $4", [nodeID, stationFreq, expectedValue, ts])
    .then((result) => {
      console.log(`stations table: inserted/updated ${result.rowCount} rows`);
      return db.query("INSERT INTO edgenodes VALUES ($1, $2, $3, $4) ON CONFLICT (edgenode) DO UPDATE SET latitude = $2, longitude = $3, timeupdated = $4", [nodeID, latitude, longitude, ts])
    })
    .then((result) => {console.log(`edgenodes table: inserted/updated ${result.rowCount} rows`)})
    .catch((err) => {console.error(err)});
}

// connect2db connects to the postgre db, tests the connection, and returns the handle
function connect2db(connectionString) {
  console.log(`Connecting (eventually) to ${connectionString}`);
  const db = new Pool({ connectionString: connectionString });
  // verify the connection to the db
  db.query('SELECT NOW()')
      .then((res) => console.log('Connected to db: ' + res.rows[0].now))
      .catch((e) => setImmediate(() => { throw e; }));
  return db
}

// seconds2timestamp converts unix seconds to a data/time formatted in the way postgres expects it for a timestamp type
function seconds2timestamp(unixSeconds) {
  const d = new Date(unixSeconds*1000);
  return d.getUTCFullYear()+'-'+(d.getUTCMonth()+1)+'-'+d.getUTCDate()+' '+d.getUTCHours()+':'+d.getUTCMinutes()+':'+d.getUTCSeconds()+'.'+d.getUTCMilliseconds()+'+00';
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

/* how to use Promises in various ways...
function sleep(seconds) {
  return new Promise((resolve) => setTimeout(resolve, (seconds*1000)));
}

function myindex(index) {
  // this is the promise w/o a sleep
  return new Promise((resolve, reject) => {
    console.log(`my index is ${index}`);
    if (index > 3) return reject("index can not > 3 !");
    else return resolve(index);
  });

  //await sleep(2);  // must be inside an async function. This sleeps outside the promise

  // this makes the sleep promise resolve inside our promise (works, but longer)
  return new Promise((resolve, reject) => {
    sleep(2).then(() => {
      console.log(`my index is ${index}`);
      return resolve();
    });
  });

  // this just chains the sleep promise with our promise content
  return sleep(2).then(() => {
    return console.log(`my index is ${index}`);
  });
}

// chain the myindex calls. The return stmt is critical.
myindex(2).then(() => {return myindex(4)}).then(() => {console.log('done')});

const nums = [1, 2, 3];
// run an arbitrary number of myindex calls serially. Note: i'm not sure how to handle rejects
var result = Promise.resolve();
nums.forEach(n => {
  result = result.then(() => myindex(n));
});
//result; <- not necessary?

// run promises in parallel and wait for them all to finish
var promises = [];
nums.forEach(n => {
  promises.push(myindex(n));
});
Promise.all(promises)
  .then((results) => { console.log('results:'); console.log(results) })
  .catch((err) => { console.error(err); });
*/

/* using the async module to call this...
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
      return callback(error);
    } else {
      console.log("Processing Watson STT results...");
      // const resultStr = JSON.stringify(sttResults, null, 2);
      // console.log(resultStr);
      // return callback();
      // console.log(sttResults);

      // Process each result. The processing of each result is asynchronous, so use the async module to wait for them all before we call our callback.
      //todo: Not sure if we should process them in paralell (eachOf()) or sequentially (eachOfSeries()). These calls are quick, so sequential is probably fine for now.
      asyncmod.eachOfSeries(sttResults.results, nluSentiment, function (error) {
        if (error) {
          console.log(error);
          return callback(error);
        } else {
          return callback();
        }
      });

    }
  });
} */

/* using the async module to call this...
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
        return callback(error);
      } else {
        console.log('Sentiment analysis results:');
        const resultStr = JSON.stringify(nluResults, null, 2);
        console.log(resultStr);
        return callback();    //todo: call addSentimentsToDB() passing it the callback
      }
    });
  } else {
    console.log('Skipping alternative: Final: '+r.final+', Confidence: '+r.alternatives[0].confidence+', Text: '+r.alternatives[0].transcript);
    return callback();
  }
} */

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
  const Client = require('node-rest-client').Client;
  const options_auth = { user: params.watsonSttUsername, password: params.watsonSttPassword };
  const client = new Client(options_auth);
  client.get("https://stream.watsonplatform.net/speech-to-text/api/v1/models", function (data, response) {
    console.log(data);     // parsed response body as js object
    console.log(response);   // raw response
  });
  */
