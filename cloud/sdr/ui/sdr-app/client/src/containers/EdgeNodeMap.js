import 'mapbox-gl/dist/mapbox-gl.css'
import './EdgeNodeMap.css'

import React, {Component} from 'react'
import ReactMapGL, {
  Marker,
  Popup,
  NavigationControl,
} from 'react-map-gl'
import {geolocated} from 'react-geolocated'
import ReactSVG from 'react-svg'

import MapMarker from '../components/MapMarker'
import MapMarkerPopup from '../components/MapMarkerPopup'

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

  _renderCityMarker = (city, index) => {
    return (
      <Marker key={`marker-${index}`}
        longitude={city.longitude}
        latitude={city.latitude}
      >
        <MapMarker size={20} onClick={() => this.setState({popupInfo: city})} />
      </Marker>
    )
  }

  _renderPopup() {
    const {popupInfo} = this.state

    return popupInfo && (
      <Popup tipSize={5}
        anchor="top"
        longitude={popupInfo.longitude}
        latitude={popupInfo.latitude}
        onClose={() => this.setState({popupInfo: null})} >
        <MapMarkerPopup info={popupInfo} />
      </Popup>
    );
  }

  render() {

    const cities = [{
      latitude: 41.1264849,
      longitude: -73.7140195,
    }]

    return (
      <div className="bx--row">
        <div className="bx--col-xs-12">
          <div>
            <ReactMapGL
              mapStyle='mapbox://styles/mapbox/dark-v9'
              {...this.state.viewport}
              onViewportChange={(viewport) => this.setState({viewport})}
              mapboxApiAccessToken={MAPBOX_TOKEN}
              className="edge-node-map"
            >
              {cities.map(this._renderCityMarker)}

              {this._renderPopup()}
            </ReactMapGL>
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