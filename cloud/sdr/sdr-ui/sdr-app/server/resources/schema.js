const {makeExecutableSchema, addMockFunctionsToSchema} = require('graphql-tools')
const resolvers = require('../resolvers')
const casual = require('casual')

const typeDefs = `
scalar Date

type Noun {
    noun: String
    sentiment: String
    numberofmentions: Int
    timeupdated: Date
}

type Query {
    nouns: [Noun]
    noun(id: ID!): Noun
  }
  `

const schema = makeExecutableSchema({typeDefs, resolvers})
module.exports = {schema}
