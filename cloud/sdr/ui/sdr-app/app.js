/*eslint-env node*/

/*
Node.js ibm cloud express server for the SDR app. This serves both the client UI (from its build dir)
and the data from the SDR postgresql DB.
*/

// const requiredPort = 6002;   // required because we have to hardcode the PORT env var to this in the cloud service, and the proxy setting in the client package.json

const path = require('path');
const winston = require('winston');   // for logging: https://github.com/winstonjs/winston
const bodyParser = require('body-parser');
const cors = require('cors');
var express = require('express');   // see: http://expressjs.com
var cfenv = require('cfenv');   // provides access to your Cloud Foundry environment, see: https://www.npmjs.com/package/cfenv

// PostgreSQL
//const { postgraphile } = require("postgraphile");

// Graphql server-side modules
const { graphqlExpress, graphiqlExpress } = require('apollo-server-express');
const { apolloUploadExpress } = require('apollo-upload-server');
const {schema} = require('./server/schema');
//const _ = require('lodash');

// create a new express server
var appSvr = express();

// get the app environment from Cloud Foundry
var appEnv = cfenv.getAppEnv();
// console.log(process.env);
// console.log('HOME: ' + process.env.HOME);
// console.log(appEnv);
// const envStr = JSON.stringify(process.env)
// console.log('Env: ' + envStr)

appSvr.use('*', cors());

// Configure what express should use/serve:
// serve the client files out of ./client/build (which is built by running 'npm run build' in that dir)
appSvr.use(express.static(path.join(__dirname, 'client', 'build')));
//app.use(express.static(__dirname + '/public'));

//const {postgresUrl} = require('./server/config/settings')
//app.use(postgraphile(process.env.SDR_DB_URL || postgresUrl || "postgres://localhost/"));

appSvr.use('/graphql', bodyParser.json(), apolloUploadExpress(), graphqlExpress({ schema }));
appSvr.use('/graphiql', graphiqlExpress({ endpointURL: '/graphql', }));    // only needed for developers to interactively browse the db

// Set up logging
var logFile = process.env.HOME + '/logs/sdr-app.log';
if (process.env.CF_INSTANCE_IP) {
  logFile = process.env.HOME + '/../logs/sdr-app.log';     // when running in the cloud service, HOME is set to /home/vcap/app, but the logs should go in the existing /home/vcap/logs
}
console.log('Configuring winston logging to ' + logFile);
const logger = winston.createLogger({
  level: 'info',    // and below
  format: winston.format.combine( winston.format.timestamp(), winston.format.json() ),
  transports: [
    // new winston.transports.Console(),
    new winston.transports.File({ filename: logFile })
  ]
});
logger.info('Winston logging configured.');

/* Verify the port number is correct in all of the environments we run it
const port = requiredPort
if (process.env.PORT != requiredPort || appEnv.port != requiredPort) {
  // I think we have to set the PORT env var in the cloud service to get the proxy in front of the service configured correctly
  logError(`the PORT env var (${process.env.PORT}) must be set to ${requiredPort}`)
} */

// logger.info('Env: ' + envStr)
logger.info(`Important Environment Variables: cenv: url=${appEnv.url},bind=${appEnv.bind},port=${appEnv.port}, USER: ${process.env.USER}, HOME: ${process.env.HOME}, CF_INSTANCE_IP: ${process.env.CF_INSTANCE_IP}, CF_INSTANCE_PORTS: ${process.env.CF_INSTANCE_PORTS}, VCAP_APP_PORT: ${process.env.VCAP_APP_PORT}, NODE_ENV: ${process.env.NODE_ENV}, npm_config_node_version: ${process.env.npm_config_node_version}`);

appSvr.get('*', (req, res) => {
  res.sendFile(path.join(__dirname + '/client/build/index.html'))
})

// start server on the specified port and binding host
appSvr.listen(appEnv.port || 6005, '0.0.0.0', function() {
  // print a message when the server starts listening
  const listeningStr = `SDR express server listening on ${appEnv.url} (port ${appEnv.port})`;
  console.log(listeningStr);
  logger.info(listeningStr);
});


/**
 * Write the error to both the console and log
 * @param {string} str - the msg to log
function logError(str) {
  console.error('Error: ' + str);
  logger.error('Error: ' + str);
}
 */

/**
 * Write the warning to both the console and log
 * @param {string} str - the msg to log
function logWarning(str) {
  console.warning('Warning: ' + str);
  logger.warning('Warning: ' + str);
}
 */
