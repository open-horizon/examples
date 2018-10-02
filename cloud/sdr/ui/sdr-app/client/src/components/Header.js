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
        productName="Cloud"
        logoHref="https://www.ibm.com/cloud/"
      />
      </div>
    )
  }
}

export default Header