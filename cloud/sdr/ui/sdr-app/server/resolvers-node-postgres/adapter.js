/*eslint-env node*/
const { Pool } = require('pg');
const settings = require('../config/settings');
const connStr = settings.postgresUrl; // add your psql details

// First create a pool, and then use that to create the client. Both are actually promises, that will connect on demand?
console.log(`Connecting (eventually) to ${connStr}`);
const psql = new Pool({ connectionString: connStr });
// verify the connection to the db
psql.query('SELECT NOW()')
    .then((res) => console.log('Connected to db: ' + res.rows[0].now))
    .catch((e) => setImmediate(() => { throw e; }));

psql.query('select noun, sentiment, numberofmentions from globalnouns order by timeupdated desc limit 5').then((res) => {
    console.log('globalnouns table:');
    console.log(res.rows);
}).catch((e) => console.error(e.stack));

/* To test what psql.query() returns once the promise is fulfilled...
(async function() {
    // const mydata = await psql.manyOrNone('select noun, sentiment, numberofmentions from globalnouns').then((data) => data.slice());
    const res = await psql.query('select noun, sentiment, numberofmentions from globalnouns');
    console.log('when waiting for the data:');
    console.log(res.rows);
}()); */

exports.psql = psql;
