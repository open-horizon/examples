// Main component for the UI portion of the SDR Edge App.
// Built on react: https://reactjs.org/docs/getting-started.html

import React, { Component } from 'react';

import {
  BrowserRouter,
  Route,
  Link,
  Redirect,
} from 'react-router-dom';

import {
  Breadcrumb,
  BreadcrumbItem,
} from 'carbon-components-react'

import GlobalServices from './containers/GlobalServices';
import EdgeNodeMap from './containers/EdgeNodeMap';
import About from './containers/About';

import Header from './components/Header';
import Nav from './components/Nav';

import './App.css';
import EdgeNodeDetails from './containers/EdgeNodeDetails';

class App extends Component {
  render() {
    // Return components to render. This is JSX, see https://reactjs.org/docs/introducing-jsx.html
    return (
      <BrowserRouter>
        <div>
          <Header />
            <div className="bx--grid">
              <div className="bx--row">
                <div className="bx--offset-xs-2 bx--col-xs-12">
                  <div className="app-content">
                    <Route exact path="/" render ={() => <Redirect to="/app/global-keywords" />} />
                    <Route exact path="/app" render={() => <Redirect to="/app/global-keywords" />} />
                    <Route path="/app/global-keywords" component={GlobalServices} />
                    <Route exact path="/app/edge-nodes" component={EdgeNodeMap} />
                    <Route path="/app/edge-nodes/details" component={EdgeNodeDetails} />
                    <Route path="/app/about" component={About} />
                  </div>
                </div>
              </div>
            </div>
          <Nav />
        </div>
      </BrowserRouter>
    );
  }
}

export default App;
