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

    // highlight the correct route
    let activeRoute = "global-keywords"
    let locationHref = window.location.href
    if (locationHref.includes("edge-nodes")) {
      activeRoute = "edge-nodes"
    } else if (locationHref.includes("about")) {
      activeRoute = "about"
    }
  
    return (
      <InteriorLeftNav className="nav-top-reset">
        <InteriorLeftNavItem className={activeRoute === 'global-keywords' && 'active-route'}>
          <Link to="/app/global-keywords">Global Keywords</Link>
        </InteriorLeftNavItem>
        <InteriorLeftNavItem className={activeRoute === 'edge-nodes' && 'active-route'}>
          <Link to="/app/edge-nodes">Edge Nodes</Link>
        </InteriorLeftNavItem>
        <br />
        <InteriorLeftNavItem className={activeRoute === 'about' && 'active-route'}>
          <Link to="/app/about">About</Link>
        </InteriorLeftNavItem>
      </InteriorLeftNav>
    )
  }
}

export default Nav