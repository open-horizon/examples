/*eslint-env node*/

/*
Node.js ibm cloud express server for the SDR app. This serves both the client UI (from its build dir)
and the data from the SDR postgresql DB.
*/

const path = require('path');
const winston = require('winston');   // for logging: https://github.com/winstonjs/winston
const bodyParser = require('body-parser');
const cors = require('cors');
const session = require('express-session');
const passport = require('passport');
const appID = require('ibmcloud-appid');
const express = require('express')

// Graphql server-side modules
const { graphqlExpress, graphiqlExpress } = require('apollo-server-express');
const { apolloUploadExpress } = require('apollo-upload-server');
const {schema} = require('./server/schema');

const helmet = require('helmet')
const express_enforces_ssl = require('express-enforces-ssl')
const cfEnv = require('cfenv');
const cookieParser = require('cookie-parser')
const flash = require('connect-flash')
const nconf = require('nconf')

const WebAppStrategy = appID.WebAppStrategy
const userProfileManager = appID.UserProfileManager
const UnauthorizedException = appID.UnauthorizedException

const APPID_LOGIN_URL = '/ibm/bluemix/appid/login'
const APPID_CALLBACK_URL = '/ibm/bluemix/appid/callback';

const RETURNING_USER_HINT = "An identified user returned to the app with the same identity. The app accesses his identified profile and the previous selections that he made.";
const NEW_USER_HINT = "An identified user logged in for the first time. Now when he logs in with the same credentials from any device or web client, the app will show his same profile and selections.";

const port = process.env.PORT || 6006

// create a new express server
var appSvr = express();

const isLocal = cfEnv.getAppEnv().isLocal;
const config = getLocalConfig()
console.log('config configured', config)
configureSecurity()

appSvr.use(flash())

// get the app environment from Cloud Foundry
var appEnv = cfEnv.getAppEnv();

appSvr.use('*', cors());

appSvr.use('/graphql', bodyParser.json(), apolloUploadExpress(), graphqlExpress({ schema }));
appSvr.use('/graphiql', graphiqlExpress({ endpointURL: '/graphql', }));    // only needed for developers to interactively browse the db

// IBM App ID
appSvr.use(session({
  secret: require('./server/config/settings').appIDSecret,
  resave: true,
  saveUninitialized: true,
  proxy: true,
  cookie: {
    httpOnly: true,
    secure: !isLocal,
  },
}));
appSvr.use(passport.initialize());
appSvr.use(passport.session());

let webAppStrategy = new WebAppStrategy(config)

passport.use(webAppStrategy)

userProfileManager.init(config)

// Configure passportjs with user serialization/deserialization. This is required
// for authenticated session persistence accross HTTP requests. See passportjs docs
// for additional information http://passportjs.org/docs
passport.serializeUser((user, cb) => {
  cb(null, user)
})

passport.deserializeUser((obj, cb) => {
  cb(null, obj)
})

appSvr.get(
  APPID_CALLBACK_URL, 
  passport.authenticate(WebAppStrategy.STRATEGY_NAME,
  { failureRedirect: '/error', failureFlash: true, allowAnonymousLogin: false }))
appSvr.get(
  APPID_LOGIN_URL, 
  passport.authenticate(WebAppStrategy.STRATEGY_NAME, 
    { forceLogin: true, }))

function storeRefreshTokenInCookie(req, res, next) {
  console.log('storing rf token in cookie')
	if (req.session[WebAppStrategy.AUTH_CONTEXT] && req.session[WebAppStrategy.AUTH_CONTEXT].refreshToken) {
    const refreshToken = req.session[WebAppStrategy.AUTH_CONTEXT].refreshToken;
    console.log('storing cookie to res')
    /* An example of storing user's refresh-token in a cookie with expiration of a month */
    res.cookie('refreshToken', refreshToken, {maxAge: 1000 * 60 * 60 * 24 * 30 /* 30 days */});
	}
	next();
}

function isLoggedIn(req) {
	return req.session[WebAppStrategy.AUTH_CONTEXT];
}

appSvr.get('/app*', (req, res, next) => {
	if (isLoggedIn(req)) {
    console.log('/app: user is logged in')
		return next()
	}

  console.log('/app: user is not logged in')
	webAppStrategy.refreshTokens(req, req.cookies.refreshToken).finally(function() {
		next();
	});
}, passport.authenticate(WebAppStrategy.STRATEGY_NAME), 
  storeRefreshTokenInCookie,
  (req, res, next) => {
    res.sendFile(path.join(__dirname + '/client/build/index.html'))
  });

appSvr.get('/login', passport.authenticate(WebAppStrategy.STRATEGY_NAME, {successRedirect: '/app', forceLogin: true}))

appSvr.get('/logout', (req, res, next) => {
  WebAppStrategy.logout(req)

  res.clearCookie('refreshToken')
  res.redirect('/login')
})

appSvr.get('/token', (req, res) => {
  res.json({
    tokens: req.session[WebAppStrategy.AUTH_CONTEXT],
  })
})

appSvr.get('/userInformation', (req, res) => {

})

appSvr.get('/error', (req, res) => {
  let errArr = req.flash('error')
  res.redirect('/error')
})

appSvr.get('/change_password', passport.authenticate(WebAppStrategy.STRATEGY_NAME, {
  successRedirect: '/app',
  show: WebAppStrategy.CHANGE_PASSWORD,
}))

appSvr.get('/change_details', passport.authenticate(WebAppStrategy.STRATEGY_NAME, {
  successRedirect: '/app',
  show: WebAppStrategy.CHANGE_DETAILS,
}))

appSvr.get('/', passport.authenticate(WebAppStrategy.STRATEGY_NAME, {
  successRedirect: '/app',
}))
// (req, res, next) => {
//   if (!isLoggedIn(req)) {
//     webAppStrategy.refreshTokens(req, req.cookies.refreshToken)
//         .then(() => {
//           console.log('refreshing tokens')
//           res.redirect('/app')
//         })
//         .catch((err) => {
//           next()
//         })
//   } else {
//     res.redirect('/app')
//   }
// }, (req, res, next) => {
//   res.sendFile(path.join(__dirname + '/client/build/index.html'))
// })

// Configure what express should use/serve:
// serve the client files out of ./client/build (which is built by running 'npm run build' in that dir)
appSvr.use(express.static(path.join(__dirname, 'client', 'build')));

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

// appSvr.get('*', (req, res) => {
//   res.sendFile(path.join(__dirname + '/client/build/index.html'))
// })

appSvr.use((err, req, res, next) => {
  if (err instanceof UnauthorizedException) {
    WebAppStrategy.logout(req)
    res.redirect('/')
  } else {
    next(err)
  }
})

// start server on the specified port and binding host
appSvr.listen(appEnv.port || port, '0.0.0.0', function() {
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


function getLocalConfig() {
	if (!isLocal) {
		return {};
	}
	let config = {};
	const localConfig = nconf.env().file(`${__dirname}/localdev-config.json`).get();
	const requiredParams = ['clientId', 'secret', 'tenantId', 'oauthServerUrl', 'profilesUrl'];
	requiredParams.forEach(function (requiredParam) {
		if (!localConfig[requiredParam]) {
			console.error('When running locally, make sure to create a file *localdev-config.json* in the root directory. See config.template.json for an example of a configuration file.');
			console.error(`Required parameter is missing: ${requiredParam}`);
			process.exit(1);
		}
		config[requiredParam] = localConfig[requiredParam];
	});
  config['redirectUri'] = `http://localhost:${port}${APPID_CALLBACK_URL}`;
	return config;
}

function configureSecurity() {
	appSvr.use(helmet());
	appSvr.use(cookieParser());
	appSvr.use(helmet.noCache());
	appSvr.enable("trust proxy");
	if (!isLocal) {
		appSvr.use(express_enforces_ssl());
	}
}