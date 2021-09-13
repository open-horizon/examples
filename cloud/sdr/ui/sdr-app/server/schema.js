/*eslint-env node*/
// const GraphQLSchema = require('graphql').GraphQLSchema;

/* Define the graphql queries supported. Examples of queries to the svr:
    query { nouns {noun sentiment numberofmentions timeupdated } }
    query { noun(noun: "trump") { noun sentiment numberofmentions timeupdated } }
*/
const typeDefs = `
scalar Date

type Noun {
    noun: String
    sentiment: String
    numberofmentions: Int
    timeupdated: Date
}

type EdgeNode {
    edgenode: String
    latitude: Float
    longitude: Float
    timeupdated: Date
}

type Query {
    globalnouns(limit: Int): [Noun]!
    nodenouns(edgenode: String!, limit: Int): [Noun]
    edgenodetopnoun(edgenode: String!): Noun
    edgenodes(limit: Int): [EdgeNode]
}
`

exports.typeDefs = typeDefs;
