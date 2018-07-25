const psql = require('./adapter').psql; // our adapter from psqlAdapter.js

// should match type Query in schema.js
// one function per endpoint
exports.resolver = {
    Query: {
        // nouns() {
        nouns: (root, args, context) => {
            const nounsQuery = 'select noun, sentiment, numberofmentions from nouns';
            // I think graphql will automatically wait for this promise to be fulfilled
            // return psql.manyOrNone(nounsQuery).then( (data) => data.slice() );
            return psql.manyOrNone(nounsQuery);
        }
    }
};
