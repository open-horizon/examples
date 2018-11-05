import React, {Component} from 'react'
import {graphql} from 'react-apollo'
import {gql} from 'apollo-boost'
import {
  Breadcrumb,
  BreadcrumbItem,
} from 'carbon-components-react'
import qs from 'query-string'

import {EdgeNodeSentiments} from '../Sentiment'

// Fetch list of edge node nouns
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

  render() {
    return (
      <div>
        <Breadcrumb noTrailingSlash={false}>
          <BreadcrumbItem href="/app/edge-nodes">
            Edge Nodes
          </BreadcrumbItem>
          <BreadcrumbItem href="#">
            {window.location && window.location.search && qs.parse(window.location.search).id}
          </BreadcrumbItem>
        </Breadcrumb>
        <br />
        <EdgeNodeSentiments nodeId={window.location && window.location.search && qs.parse(window.location.search).id} />
      </div>
    )
  }
}

export default graphql(EDGE_NODE_NOUNS_LIST)(EdgeNodeDetails)