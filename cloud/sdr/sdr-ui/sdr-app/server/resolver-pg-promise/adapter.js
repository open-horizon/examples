const pgPromise = require('pg-promise');
const settings = require('../config/settings');

// const connStr = 'postgresql://user@host:port/database'; // add your psql details
const connStr = settings.postgresUrl; // add your psql details

const pgp = pgPromise({}); // empty pgPromise instance
console.log(`connecting to ${connStr}`);
const psql = pgp(connStr); // get connection to your db instance

//psql.manyOrNone('select noun, sentiment, numberofmentions from nouns', [true]).then((data) => {
psql.manyOrNone('select noun, sentiment, numberofmentions from nouns').then((data) => {
    console.log('nouns table:');
    console.log(data.slice());
});

// To test what psql.manyOrNone() returns once the promise is fulfilled...
(async function() {
    // const mydata = await psql.manyOrNone('select noun, sentiment, numberofmentions from nouns').then((data) => data.slice());
    const mydata = await psql.manyOrNone('select noun, sentiment, numberofmentions from nouns');
    console.log('mydata:');
    console.log(mydata);
}());

exports.psql = psql;
