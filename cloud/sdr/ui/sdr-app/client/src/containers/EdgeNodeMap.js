import 'mapbox-gl/dist/mapbox-gl.css'
import './EdgeNodeMap.css'

import React, {Component} from 'react'
import ReactMapGL, {
  Marker,
  Popup,
  NavigationControl,
} from 'react-map-gl'
import {
  Breadcrumb,
  BreadcrumbItem,
  Tile,
  Loading,
} from 'carbon-components-react'
import { SizeMe } from 'react-sizeme'
import {geolocated} from 'react-geolocated'
import { graphql, Query } from 'react-apollo'
import { gql } from 'apollo-boost'

import MapMarker from '../components/MapMarker'
import MapMarkerPopup from '../components/MapMarkerPopup'

import {
  MAPBOX_TOKEN,
} from '../config/settings'

const EDGE_NODE_LIST = gql`
{
  edgenodes {
    edgenode
    latitude
    longitude
    timeupdated
  }
}
`

const EDGE_NODE_TOP_NOUN = gql`
query edgenodetopnoun($edgenode: String!) {
    edgenodetopnoun(edgenode: $edgenode) {
        noun
        sentiment
        numberofmentions
        timeupdated
    }
}
`

class EdgeNodeMap extends Component {
  state = {
    viewport: {
      width: 800,
      height: 800,
      latitude: (this.props.isGeolocationAvailable && this.props.isGeolocationEnabled && this.props.coords && this.props.coords.latitude) || 41.1264849,
      longitude: (this.props.isGeolocationAvailable && this.props.isGeolocationEnabled && this.props.coords && this.props.coords.longitude) || -73.7140195,
      zoom: 8,
    }
  }

  _renderCityMarker = (edgeNodes, index) => {
    return (
      <Marker key={`marker-${index}`}
        longitude={edgeNodes[0].longitude}
        latitude={edgeNodes[0].latitude}
      >
        <MapMarker size={20} onClick={() => this.setState({popupInfo: edgeNodes})} />
      </Marker>
    )
  }

  _renderPopup() {
    const {popupInfo} = this.state

    return popupInfo && (
      <Popup tipSize={5}
        anchor="top"
        longitude={popupInfo[0].longitude}
        latitude={popupInfo[0].latitude}
        onClose={() => this.setState({popupInfo: null})} >
        {popupInfo.length === 1 ?
          <Query 
            query={EDGE_NODE_TOP_NOUN} 
            variables={{edgenode: popupInfo[0].edgenode}}
            pollInterval={1000}
          >
            {({loading, error, data}) => {
              if (loading) return <Loading withOverlay={false} />
              if (error) return `Error! ${error.message}`

              return <MapMarkerPopup info={popupInfo} data={data} />
            }}
        </Query>
        : <MapMarkerPopup info={popupInfo} />
        }
      </Popup>
    );
  }

  render() {
    let edgeNodes = []
    if (this.props && this.props.data && this.props.data.edgenodes) {
      edgeNodes = this.props.data.edgenodes
    }

    // hash key is: {lat}-{lng}
    // hash val is: array of edgeNode
    let edgeNodeHash = {}

    for (let i = 0; i < edgeNodes.length; i++) {
      const checkKey = edgeNodes[i].latitude + '-' + edgeNodes[i].longitude
      if (typeof edgeNodeHash[checkKey] === 'undefined') { // create new in hash
        edgeNodeHash[checkKey] = [edgeNodes[i]]
      } else { // there's another edge node w/ same lat and lng, group them together
        edgeNodeHash[checkKey].push(edgeNodes[i])
      }
    }

    let edgeNodesDeduped = []

    for (let i = 0; i < Object.keys(edgeNodeHash).length; i++) {
      edgeNodesDeduped.push(edgeNodeHash[Object.keys(edgeNodeHash)[i]])
    }

    return (
      <div>
        <Breadcrumb noTrailingSlash={false}>
          <BreadcrumbItem href="/app/edge-nodes">
            Edge Nodes
          </BreadcrumbItem>
        </Breadcrumb>
        <br />
        <div className="bx--row">
          <div className="bx--col-xs-12">
            <ReactMapGL
              mapStyle='mapbox://styles/mapbox/dark-v9'
              {...this.state.viewport}
              onViewportChange={(viewport) => this.setState({viewport})}
              mapboxApiAccessToken={MAPBOX_TOKEN}
              className="edge-node-map"
            >
              {edgeNodesDeduped.map(this._renderCityMarker)}

              {this._renderPopup()}
            </ReactMapGL>
          </div>
        </div>
      </div>
    )
  }
}

export default graphql(EDGE_NODE_LIST, {
  pollInterval: 1000,
})(geolocated({
  positionOptions: {
    enableHighAccuracy: false,
  },
  userDecisionTimeout: 5000,
})(EdgeNodeMap))