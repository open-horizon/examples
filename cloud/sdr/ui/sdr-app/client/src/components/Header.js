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
    const res = await fetch('/token')
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
      console.log('res result from fetching tokens: ', res)
      console.log('json result from fetching tokens: ', json)
      // window.location.href = '/login'
    }
  }

  async renderUser() {
    const res = await fetch('/token')
    const json = await res.json()
    if (json.tokens && json.tokens.identityTokenPayload && json.tokens.identityTokenPayload.email) {
      return <ul className="list">
        <li>{json.tokens.identityTokenPayload.email}</li>
        <li><a href="/logout" className="bx--link">Log Out</a></li>
      </ul>
    } else {
      console.log('could not get', json)
    }

    // fetch('/token')
    //     .then((res) => {
    //       return res.json()
    //     })
    //     .then((json) => {
    //       console.log('json', json)
    //       if (json.tokens && json.tokens.identityTokenPayload && json.tokens.identityTokenPayload.email) {
    //         renderItems = <ul className="list">
    //           <li>{json.tokens.identityTokenPayload.email}</li>
    //           <li><a href="/logout" className="bx--link">Log Out</a></li>
    //         </ul>
    //       }
    //     })
    //     .catch((err) => {
    //       window.location.href = '/login'
    //     })
    // if (typeof renderItems === "undefined") {
    //   console.log('renderitems is undefined')
    //   // window.location.href = '/login'
    // } else {
    //   return renderItems
    // }
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