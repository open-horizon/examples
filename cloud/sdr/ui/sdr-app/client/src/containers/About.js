import React, { Component } from 'react'

import {
  Breadcrumb,
  BreadcrumbItem
} from 'carbon-components-react'

export default class About extends Component {
  render() {
    return (
      <div>
        <Breadcrumb noTrailingSlash={false}>
          <BreadcrumbItem href="/app/about">
            About
          </BreadcrumbItem>
        </Breadcrumb>
        <br />
        <div className="bx--row">
          <div className="bx--offset-xs-1 bx--col-xs-10">
            <h1>About This App</h1>
            <br />
            <div className="bx--row">
              <div className="bx--ofset-xs-1 bx--col-xs-10">
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
              </div>
            </div>
          </div>
        </div>
      </div>
    )
  }
}
