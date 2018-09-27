import React, {Component} from 'react'
import {
  InteriorLeftNav, InteriorLeftNavItem,
} from 'carbon-addons-cloud-react'
import {
  Link,
} from 'react-router-dom'

import './Nav.css'

class Nav extends Component {
  render() {
    return (
      <InteriorLeftNav className="nav-top-reset">
        <InteriorLeftNavItem>
          <Link to="/global-keywords">Global Keywords</Link>
        </InteriorLeftNavItem>
        <InteriorLeftNavItem>
          <Link to="/edge-nodes">Edge Nodes</Link>
        </InteriorLeftNavItem>
      </InteriorLeftNav>
    )
  }
}

export default Nav