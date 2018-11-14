# IBM Cloud Node.js Starter App for SDR

This application demonstrates a simple, reusable Node.js web application based on the Express framework.

---

## Prereqs

### Mapbox Token
- https://www.mapbox.com/help/how-access-tokens-work/

### IBM App ID
This app requires you have an instance of IBM App ID inside the Bluemix Console
- https://console.bluemix.net/catalog/services/app-id

Once the App ID instance is deployed, under `Identity Providers` -> `Manage` under the `Identity Providers` tab only the following should be enabled:

- Cloud Directory

Under `Identity Providers` -> `Manage` under the `Authentication Settings` tab, enter the redirect URL to your app and sign-in token expiration periods.

Under `Identity Providers` -> `Cloud Directory` under the `Settings` tab, the following settings are applied:

- Allow users to sign-up and sign-in using: `Email and password`
- Allow users to sign-up to your app: `No`
- Allow users to manage their account from your app: `Yes`


In order for users to sign into your app, you'll have to manually create their accounts (their passwords can be changed in the app by the users once they've logged in) inside `Users` -> `Add User` button.

You will also have to create a connection from the IBM App ID instance to the SDR App instance under `Connections` -> `Create Connection` button. You can give it any Access Role and Auto Generate a Service ID.

To get service credentials for some of the configs used in this app, you will have to create a new one under `Service Credentials` -> `New Credential` with any Role, and connect the auto-generated Service ID from the previous step.

## Configs

The app should have the following configurations set:

```
sdr-app/localdev-config.json
The following credentials can be retrieved from the IBM App ID service inside the Bluemix Dashboard. 


{
  "clientId": "{clientId from app id}",
  "oauthServerUrl": "{oauthServerUrl from app id}",
  "profilesUrl": "{profilesUrl from app id}",
  "secret": "{secret of service from app id}",
  "tenantId": "{tenantId from app id}"
}
```

```
sdr-app/server/config/settings.js

exports.postgresUrl = ''
```

```
sdr-app/client/src/config/settings.js

exports.MAPBOX_TOKEN = ''
```

---

## Run dev Version of the App Sever Locally

[Install Node.js](https://nodejs.org/en/download/) and then:

```
npm install   # to install the app's dependencies listed in package.json
npm start`   # to start the app express svr. Will start on http://localhost:6001 by default
```

Access the GraphQL web interface in a browser at <http://localhost:6001/graphiql> to try a query on the server like:

```
query {
    globalnouns {
      noun
      sentiment
      numberofmentions
      timeupdated
    }
}
```


## Initial Creation of React Client

For reference, the react client code was added with:

```
npm install -g create-react-app
create-react-app client   # see https://github.com/facebook/create-react-app
```

## Run dev Version of the Client Locally

```
npm install   # to install the app's dependencies listed in package.json
cd client
npm start   # to start the react front end. Will start on http://localhost:3000 by default
```

The client will send graphql queries to the dev server at http://localhost:6001/graphql because of the proxy statement in package.json (which is only honored in dev, not production)

## Run the prod Version of Server and Client Locally

Run your app in the same mode it will run when you push it to the cloud (the server serving the front-end from the client/build pack):

```
npm start
```

Then browse to http://localhost:6001/ .

## Build and Push the Updated App to the IBM Cloud Service

[Install Bluemix CLI](https://console.bluemix.net/docs/cli/reference/bluemix_cli/get_started.html), then:

```
cd client
npm run build   # build the production version of the client/front end
cd ..
ic login   # use --sso if an ibm employee
ic target -o <cforg> -s <space>
ic app push sdr-app
```

Then browse the application front-end at https://sdr-app.mybluemix.net/

---

## Notes/References

- Graphql:
    - Apollo intro: https://flaviocopes.com/apollo/
    - Apollo docs: https://www.apollographql.com/docs/
    - GraphQL: http://graphql.github.io/learn/queries/
    - GraphQL queries tutorial: https://building.buildkite.com/tutorial-getting-started-with-graphql-queries-and-mutations-11211dfe5d64
    - GraphQL IDE: https://github.com/andev-software/graphql-ide
- Interface to postgresql:
    - Node-postgres: https://node-postgres.com/
    - (not used) Pg-promise: http://vitaly-t.github.io/pg-promise/index.html
    - (not used) Postgraphile: https://www.graphile.org/postgraphile/quick-start-guide/
- Eslint:
    - Disables: https://gist.github.com/cletusw/e01a85e399ab563b1236
    - Recommended config for react: https://github.com/yannickcr/eslint-plugin-react#recommended
    - In-file env list: https://eslint.org/docs/user-guide/configuring
- JSDoc: http://usejsdoc.org/about-getting-started.html
- HTML reference: https://html.spec.whatwg.org/multipage/ and https://www.w3.org/TR/html5/syntax.html
- CSS reference: https://cssreference.io/ and https://www.w3schools.com/cssref/
- CSS Modules (for local scope css - not using): https://github.com/css-modules/css-modules and https://github.com/gajus/react-css-modules and https://github.com/webpack-contrib/css-loader
