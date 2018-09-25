import React, {Component} from 'react'
import {
  InteriorLeftNav, InteriorLeftNavItem,
} from 'carbon-addons-cloud-react'

class Nav extends Component {
  render() {
    return (
      <InteriorLeftNav>
        <InteriorLeftNavItem href="#1" label="Global Keywords">
          
        </InteriorLeftNavItem>
        <InteriorLeftNavItem href="#2" label="Edge Nodes">
          
        </InteriorLeftNavItem>
      </InteriorLeftNav>
    )
  }
}

export default Nav