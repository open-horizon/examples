import React, {Component} from 'react'
import {
  CloudHeader,
} from 'carbon-addons-cloud-react'

import './Header.css'

class Header extends Component {


  render() {
    return (
      <div>
      <CloudHeader 
        companyName="IBM"
        productName="SDR Sentiment Viewer"
        logoHref="https://www.ibm.com/cloud/"
        links={[
          { href: 'https://console.bluemix.net/catalog/', linkText: 'Catalog' },
          { href: 'https://console.stage1.bluemix.net/docs/services/edge-fabric/', linkText: 'Docs' },
        ]}
        className="cloud-header-fixed"
      />
      </div>
    )
  }
}

export default Header