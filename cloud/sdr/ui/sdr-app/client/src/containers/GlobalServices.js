import React, {Component} from 'react'

import Sentiment from '../Sentiment'

import './GlobalServices.css'

class GlobalServices extends Component {
  render() {
    return (
      <div>
      <div className="bx--row">
        <div className="bx--col-xs-12">
          <h1 className="page-title">Global Keyword Sentiments</h1>
        </div>
      </div>
      <div className="bx--row">
        <div className="bx--col-xs-12">
          <p className="page-description">The top 20 keywords mentioned on all of the edge nodes.</p>
        </div>
      </div>
      <br />
      <div className="bx--row">
        <div className="bx--col-xs-12">
          <Sentiment />
        </div>
      </div>
      </div>
    )
  }
}

export default GlobalServices