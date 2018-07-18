const {
    GraphQLString,
    GraphQLList,
    GraphQLObjectType,
    GraphQLScalarType,
  } = require('graphql')
const { GraphQLUpload } = require('apollo-upload-server')
const { Kind } = require('graphql/language')
const _ = require('lodash')

/* const {
    User,
    Post,
    Score,
    Comment,
    COS,
  } = require('../connectors') */
  
//todo: change this to be represented with postgres Date value (which is sort of a string)
const dateScalarType = new GraphQLScalarType({
    name: 'Date',
    description: 'Date represented with Int',
    serialize: (value) => {
        return new Date(value) // value from client
    },
    parseValue: (value) => {
        return value.getTime() // value to client
    },
    parseLiteral: (ast) => {
        if (ast.kind === Kind.INT) {
        return parseInt(ast.value, 10) // ast value always in string
        }
        return null
    },
})

 const resolvers = {
    Query: {
      nouns: () => {
        return Noun.find({}, (err, docs) => {
          return docs
        })
      },
      noun: (obj, args) => {
        return Noun.findOne({_id: args.id}, (err, docs) => {
          if (err) return {error: `No nouns found with id: ${args.id}`}
          return docs
        })
      },
    },
/*     Mutation: {
        addPost: (obj, args) => {
          const newPost = new Post({
            title: args.title,
            abstract: args.abstract,
            author: args.authorObjectId,
            scores: [],
            isCompleted: false,
            isArchived: false,
            creationDate: Date.now(),
            modifiedDate: Date.now(),
            isLocked: false,
          })
          newPost.save((err) => {
            // console.log('err', err)
            if (err) return {
              errors: [
                { message: 'Unable to save new post' }
              ]
            }
            console.log('success', newPost)
            return newPost
          })
          return newPost
        }
    } */
}

module.exports = {
    Date: dateScalarType,
    // Upload: GraphQLUpload,
    Query: resolvers.Query,
    // Mutation: resolvers.Mutation,
}
