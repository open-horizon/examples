import React, {Component} from 'react'
import {
  CloudHeader,
} from 'carbon-addons-cloud-react'
import {
  Tile,
  DropdownV2,
} from 'carbon-components-react'

import './Header.css'

class Header extends Component {

  state = {
    renderUserCb: () => {},
  }

  async componentDidMount() {
    const res = await fetch('/token', {
      method: 'GET',
      credentials: 'include',
      headers: {
        'Accept': 'application/json; charset=utf-8',
      },
    })
    const json = await res.json()
    if (json.tokens && json.tokens.identityTokenPayload && json.tokens.identityTokenPayload.email) {

      this.setState({
        renderUserCb: () => {
          return <ul className="list">
            <li>{json.tokens.identityTokenPayload.email}</li>
            <li><a href="/logout" className="bx--link">Log Out</a></li>
          </ul>
        }
      })
    } else {
      console.error('Error with fetching login token')
      window.location.href = '/login'
    }
  }

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
        renderUser={this.state.renderUserCb}
      />
      </div>
    )
  }
}

export default Header