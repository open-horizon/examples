/*eslint-env node*/

// node.js ibm cloud starter application for for the SDR app

// This application uses express as its web server, for more info, see: http://expressjs.com
var express = require('express');

// cfenv provides access to your Cloud Foundry environment, for more info, see: https://www.npmjs.com/package/cfenv
var cfenv = require('cfenv');

const path = require('path')

//var React = require('react');

// create a new express server
var app = express();

// serve the client files out of ./client/build
app.use(express.static(path.join(__dirname, 'client', 'build')));
//app.use(express.static(__dirname + '/public'));

// get the app environment from Cloud Foundry
var appEnv = cfenv.getAppEnv();

// start server on the specified port and binding host
app.listen(appEnv.port, '0.0.0.0', function() {
  // print a message when the server starts listening
  console.log("server starting on " + appEnv.url);
});
