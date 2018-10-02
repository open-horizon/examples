import React, {Component} from 'react'
import {graphql} from 'react-apollo'
import {gql} from 'apollo-boost'

import {EdgeNodeSentiments} from '../Sentiment'

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

class EdgeNodeDetails extends Component {
  state = {

  }

  componentDidMount() {
    console.log('props', this.props)
  }

  render() {
    return (
      <EdgeNodeSentiments nodeId={this.props.location && this.props.location.pathname && this.props.location.pathname.split('/').splice(2,5).join('/')} />
    )
  }
}

export default graphql(EDGE_NODE_NOUNS_LIST)(EdgeNodeDetails)