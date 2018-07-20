# IBM Cloud Node.js Starter App for SDR

This application demonstrates a simple, reusable Node.js web application based on the Express framework.

## Run the App Sever Locally

[Install Node.js](https://nodejs.org/en/download/) and then:

```
npm install   # to install the app's dependencies listed in package.json
npm start`   # to start the app express svr. Will start on localhost:6001 by default
```

Access the running app in a browser at <http://localhost:6001>

## Add React Client

```
npm install -g create-react-app
create-react-app client   # see https://github.com/facebook/create-react-app
```

## Run Client Locally

```
cd client
npm start
```

## Build Production Client

```
cd client
npm run build
```

## Push Updated App to Cloud Service

[Install Bluemix CLI](https://console.bluemix.net/docs/cli/reference/bluemix_cli/get_started.html), then:

```
bx app push sdr-app
```

Then browse https://sdr-app.mybluemix.net/

## Notes

- Add graphql to get postgres data
    - Apollo intro: https://flaviocopes.com/apollo/
    - GraphQL: http://graphql.github.io/learn/queries/
    - GraphQL, PostgreSQL, and pg-promise: https://blog.cloudboost.io/postgresql-and-graphql-2da30c6cde26
    - Pg-promise: http://vitaly-t.github.io/pg-promise/index.html
    - GraphQL & PostgreSQL Quickstart: https://medium.com/@james_mensch/node-js-graphql-postgresql-quickstart-91ffa4374663
    - Postgraphile: https://www.graphile.org/postgraphile/quick-start-guide/
- Eslint disables: https://gist.github.com/cletusw/e01a85e399ab563b1236
- JSDoc: http://usejsdoc.org/about-getting-started.html
