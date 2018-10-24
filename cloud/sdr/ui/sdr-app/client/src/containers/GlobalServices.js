import React, {Component} from 'react'

import {
  Breadcrumb,
  BreadcrumbItem,
} from 'carbon-components-react'

import {GlobalSentiments} from '../Sentiment'

import './GlobalServices.css'

class GlobalServices extends Component {
  render() {
    return (
      <div>
        <Breadcrumb noTrailingSlash={false}>
          <BreadcrumbItem href="/app/global-keywords">
            Global Keywords
          </BreadcrumbItem>
        </Breadcrumb>
        <br />
        <div className="bx--row">
          <div className="bx--col-xs-12">
            <GlobalSentiments />
          </div>
        </div>
      </div>
    )
  }
}

export default GlobalServices