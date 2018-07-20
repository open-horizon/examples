// const GraphQLSchema = require('graphql').GraphQLSchema;
const makeExecutableSchema = require('graphql-tools').makeExecutableSchema;

const resolver = require('./resolver-node-postgres').resolver;

// define our user type
// then define a users query, which must return an array type that optionally contains 0 or more User types
const typeDefs = `
type Noun {
    noun: String!
    sentiment: String
    numberofmentions: Int
}

type Query {
    nouns: [Noun]
}
`;

exports.schema = makeExecutableSchema({
  typeDefs,
  resolver,
});
