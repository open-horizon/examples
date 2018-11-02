/*eslint-env node*/
const psql = require('./adapter').psql
// console.log('in resolver-node-postgres/index.js');

// should match type Query in schema.js, one function per endpoint
exports.resolvers = {
    Query: {
        // get list of nouns for all edge nodes
        globalnouns: (obj, args) => {
            // Graphql will automatically wait for this promise to be fulfilled:
            // https://graphql.org/learn/execution/#asynchronous-resolvers - "During execution, GraphQL will wait for Promises, Futures, and Tasks to complete before continuing and will do so with optimal concurrency."
            return psql.query(`select noun, sentiment, numberofmentions, timeupdated from globalnouns order by numberofmentions desc, timeupdated desc limit ${args.limit}`)
                    .then((res) => res.rows)
            // return [{ noun: 'wedding', sentiment: 'positive', numberofmentions: 10 }];
        },
        // get list of nouns for a single edge node
        nodenouns: (obj, args) => {
            return psql.query(`select noun, sentiment, numberofmentions, timeupdated from nodenouns where edgenode = '${args.edgenode}' order by numberofmentions desc, timeupdated desc limit ${args.limit}`)
                    .then((res) => res.rows)
        },
        // get list of edge nodes
        edgenodes: (obj, args) => {
            return psql.query('select edgenode, latitude, longitude, timeupdated from edgenodes')
                    .then((res) => res.rows)
        },
        // get single top noun for a single edge node
        edgenodetopnoun: (obj, args) => {
            return psql.query(`select noun, sentiment, numberofmentions, timeupdated from nodenouns where edgenode = '${args.edgenode}' order by numberofmentions desc, timeupdated desc limit 1`)
                    .then((res) => res.rows[0])
        },
    },
}
