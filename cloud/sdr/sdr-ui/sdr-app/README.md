# IBM Cloud Node.js Starter App for SDR

This application demonstrates a simple, reusable Node.js web application based on the Express framework.

## Run the App Sever Locally

[Install Node.js](https://nodejs.org/en/download/) and then:

```
npm install   # to install the app's dependencies listed in package.json
npm start`   # to start the app
```

Access the running app in a browser at <http://localhost:6001>

## Add React Client

```
npm install -g create-react-app
create-react-app client   # see https://github.com/facebook/create-react-app
```

## Run Client Locally

```
npm start
```

## Build Production Client

```
npm run build
```

## Push Updated App to Cloud Service

Install Bluemix CLI, then:

```
bx app push sdr-app
```

Then browse https://sdr-app.mybluemix.net/

## Notes

- Add redux to manage app state across react components?
- Add graphql to get postgres data
