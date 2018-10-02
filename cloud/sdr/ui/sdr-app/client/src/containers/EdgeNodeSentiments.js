import React, {Component} from 'react'
import {graphql} from 'react-apollo'
import {gql} from 'apollo-boost'

import {EdgeSentiments} from '../Sentiment'

const EDGE_NODE_NOUNS_LIST = gql`
query nodenouns($edgenode: String, $limit: Int) {
  nodenouns(edgenode: $edgenode, limit: $limit) {
        noun
        sentiment
        numberofmentions
        timeupdated
    }
}
`

class EdgeNodeSentiments extends Component {
  state = {

  }

  componentDidMount() {
    console.log('props', this.props)
  }

  render() {
    return (
      <EdgeSentiments />
    )
  }
}

export default graphql(EDGE_NODE_NOUNS_LIST)(EdgeNodeSentiments)