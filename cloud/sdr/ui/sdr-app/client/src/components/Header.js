/* global localStorage */

import React, {Component} from 'react'
import {
  CloudHeader,
} from 'carbon-addons-cloud-react'
import {
  Modal
} from 'carbon-components-react'

import './Header.css'

class Header extends Component {

  state = {
    renderUserCb: () => {},
    modalOpen: false,
  }

  /**
   * Use an async version of componentDidMount so that
   * the header can properly render the user's email and log out button.
   * This also checks to make sure the user is logged in, otherwise they
   * will be redirected to the login route.
   */
  async componentDidMount() {
    this.isReturningUser()
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
            <li><a href="/change_password" className="bx--link">Change Password</a></li>
            <br />
            <li><a href="/logout" className="bx--link">Log Out</a></li>
          </ul>
        }
      })
    } else {
      console.error('Error with fetching login token')
      window.location.href = '/login'
    }
  }

  openDocs() {
    window.open('https://console.test.cloud.ibm.com/docs/services/edge-fabric/poc/sdr.html')
  }

  aboutModal() {
    return <Modal
      modalHeading="About This App"
      open={this.state.modalOpen}
      primaryButtonText="Ok"
      onPrimarySubmit={this.closeModal.bind(this)}
      secondaryButtonText="Docs"
      onSecondarySubmit={this.openDocs.bind(this)}
      onRequestClose={this.closeModal.bind(this)}
      onRequestSubmit={this.closeModal.bind(this)}
    >
      <p>
        This is the web interface for the Edge-Fabric Software Defined Radio (SDR) example program.  
        Note that this is just an example program, not a production service.
        <br />
        Documentation for using the SDR application is <a href="https://console.test.cloud.ibm.com/docs/services/edge-fabric/poc/sdr.html">here</a>
        &nbsp;and a developer walkthrough of the main software components is <a href="https://console.test.cloud.ibm.com/docs/services/edge-fabric/dev/sdr.html">here</a>.
        <br />
        The open source code for the cloud side of this example is <a href="https://github.com/open-horizon/examples/tree/master/cloud/sdr">here</a>.
        <br />
        The open source code for the Edge Node side (which runs where the SDR sensors are located) is
        &nbsp;<a href="https://github.com/open-horizon/examples/tree/master/edge/services/sdr">here</a> (the low level service that manages the hardware),
        &nbsp;and <a href="https://github.com/open-horizon/examples/tree/master/edge/msghub/sdr2msghub">here</a> (the higher level service that sends data to the cloud).
      </p>
    </Modal>
  }

  isReturningUser() {
    if (localStorage.getItem('sdrReturningUser-v0.0.6') === null) {
      localStorage.setItem('sdrReturningUser-v0.0.6', true)
      this.openModal()
      return false
    }
    this.closeModal()
    return true
  }

  closeModal() {
    this.setState({ modalOpen: false })
  }

  openModal() {
    this.setState({ modalOpen: true })
  }

  render() {
    return (
      <div>
        {this.aboutModal()}
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