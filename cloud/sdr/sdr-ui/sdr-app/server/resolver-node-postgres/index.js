const psql = require('./adapter').psql;

// should match type Query in schema.js, one function per endpoint
exports.resolver = {
    Query: {
        // nouns() {
        nouns: (root, args, context) => {
            // I think graphql will automatically wait for this promise to be fulfilled:
            // https://graphql.org/learn/execution/#asynchronous-resolvers - "During execution, GraphQL will wait for Promises, Futures, and Tasks to complete before continuing and will do so with optimal concurrency."
            return psql.query('select noun, sentiment, numberofmentions from nouns').then((res) => res.rows);
            // return [{ noun: 'wedding', sentiment: 'positive', numberofmentions: 10 }];
        }
    }
};
