import 'mapbox-gl/dist/mapbox-gl.css'

import React, {Component} from 'react'
import ReactMapGL from 'react-map-gl'
import {geolocated} from 'react-geolocated'

import {
  MAPBOX_TOKEN,
} from '../config/settings'

class EdgeNodeMap extends Component {
  state = {
    viewport: {
      height: 800,
      width: 800,
      latitude: (this.props.isGeolocationAvailable && this.props.isGeolocationEnabled && this.props.coords && this.props.coords.latitude) || 41.1264849,
      longitude: (this.props.isGeolocationAvailable && this.props.isGeolocationEnabled && this.props.coords && this.props.coords.longitude) || -73.7140195,
      zoom: 8,
    }
  }

  render() {
    return (
      <div className="bx--row">
        <div className="bx--col-xs-12">
          <div>
            <ReactMapGL
              {...this.state.viewport}
              onViewportChange={(viewport) => this.setState({viewport})}
              mapboxApiAccessToken={MAPBOX_TOKEN}
            />
          </div>
        </div>
      </div>
    )
  }
}

export default geolocated({
  positionOptions: {
    enableHighAccuracy: false,
  },
  userDecisionTimeout: 5000,
})(EdgeNodeMap)