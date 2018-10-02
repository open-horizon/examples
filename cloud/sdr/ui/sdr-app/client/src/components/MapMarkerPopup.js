import React, {PureComponent} from 'react'
import moment from 'moment'

export default class MapMarkerPopup extends PureComponent {

  render() {
    const {info} = this.props
    const displayName = `${info.city}, ${info.state}`

    let renderList = undefined

    if (info.length > 1) {
      renderList = <div>
        <h1>Multiple Nodes Found</h1>
        <div>{info.map(o => <div key={o.edgenode}><a href={`/edge-nodes/${o.edgenode}`}>{o.edgenode}&n</a><br /></div>)}</div>  
      </div>
    }

    return (
      <div>
        {info.length > 1 &&
          renderList
        }
        {info.length === 1 &&
          <div>
            <h1><a href={`/edge-nodes/${info[0].edgenode}`}>{info[0].edgenode}</a></h1>
            <p>
              Latitude: {info[0].latitude}
              <br />
              Longitiude: {info[0].longitude}
              <br />
              Last Updated: {moment(info[0].timeupdated).toString()}
            </p>
          </div>
        }
        
      </div>
    )
  }
}