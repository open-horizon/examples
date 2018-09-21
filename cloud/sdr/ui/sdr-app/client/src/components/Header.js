import React, {Component} from 'react'
import {
  CloudHeader,
} from 'carbon-addons-cloud'

import './Header.css'

class Header extends Component {
  render() {
    return (
      <div>
      <CloudHeader 
        companyName="IBM"
        productName="Cloud"
      />
      </div>
    )
  }
}

export default Header