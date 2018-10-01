import React, {PureComponent} from 'react'

export default class MapMarkerPopup extends PureComponent {

  render() {
    const {info} = this.props
    const displayName = `${info.city}, ${info.state}`

    return (
      <div>
        test
      </div>
    )
  }
}