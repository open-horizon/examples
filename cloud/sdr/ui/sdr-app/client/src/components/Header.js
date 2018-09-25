import React, {Component} from 'react'
import {
  CloudHeader,
} from 'carbon-addons-cloud-react'

import Nav from './Nav'

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
      <Nav />
      </div>
    )
  }
}

export default Header