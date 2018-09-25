import React, {Component} from 'react'
import ReactMapGL from 'react-map-gl'

class EdgeNodeMap extends Component {
  state = {
    viewport: {
      height: 400,
      latitude: 0,
      longitude: 0,
      zoom: 8,
    }
  }

  render() {
    return (
      <ReactMapGL
        {...this.state.viewport}
        onViewportChange={(viewport) => this.setState({viewport})}
      />
    )
  }
}

export default EdgeNodeMap