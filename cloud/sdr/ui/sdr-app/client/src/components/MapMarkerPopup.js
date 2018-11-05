import React, {PureComponent} from 'react'
import moment from 'moment'
import {
  Tag,
} from 'carbon-components-react'

export default class MapMarkerPopup extends PureComponent {

  render() {
    const {info, data} = this.props
    const displayName = `${info.city}, ${info.state}`

    let renderList = undefined

    // if there are more than one nodes at this lat/lng, render a list instead
    if (info.length > 1) {
      renderList = <div>
        <h1>Multiple Nodes Found</h1>
        <div>{info.map(o => <Tag key={o.edgenode} type="custom"><a href={`/app/edge-nodes/details?id=${o.edgenode}`}>{o.edgenode}</a><br /></Tag>)}</div>  
      </div>
    }

    return (
      <div>
        {info.length > 1 &&
          renderList
        }
        {info.length === 1 &&
          <div>
            <h1><a href={`/app/edge-nodes/details?id=${info[0].edgenode}`}>{info[0].edgenode}</a></h1>
            <p>
              Latitude: {info[0].latitude}
              <br />
              Longitiude: {info[0].longitude}
              <br />
              Last Updated: {moment(info[0].timeupdated).toString()}
              <br />
              Top Noun: {data.edgenodetopnoun.noun}
              <br />
              Number of Mentions: {data.edgenodetopnoun.numberofmentions}
              <br />
              Sentiment: {data.edgenodetopnoun.sentiment}
            </p>
          </div>
        }
        
      </div>
    )
  }
}