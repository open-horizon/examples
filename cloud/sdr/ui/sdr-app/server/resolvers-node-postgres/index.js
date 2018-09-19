/*eslint-env node*/
const psql = require('./adapter').psql;
// console.log('in resolver-node-postgres/index.js');

// should match type Query in schema.js, one function per endpoint
exports.resolvers = {
    Query: {
        // globalnouns() {  // <- also works
        globalnouns: () => {
            // Graphql will automatically wait for this promise to be fulfilled:
            // https://graphql.org/learn/execution/#asynchronous-resolvers - "During execution, GraphQL will wait for Promises, Futures, and Tasks to complete before continuing and will do so with optimal concurrency."
            console.log('running globalnouns resolver');
            return psql.query('select noun, sentiment, numberofmentions, timeupdated from globalnouns order by timeupdated desc limit 20').then((res) => res.rows);
            // return [{ noun: 'wedding', sentiment: 'positive', numberofmentions: 10 }];
        },
        noun: (obj, args) => {
            console.log('running noun resolver');
            return psql.query(`select noun, sentiment, numberofmentions, timeupdated from nouns where noun = '${args.noun}'`).then((res) => res.rows[0]);
        }
    }
};
