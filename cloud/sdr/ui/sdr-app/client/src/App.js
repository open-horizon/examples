// Main component for the UI portion of the SDR Edge App.
// Built on react: https://reactjs.org/docs/getting-started.html

import React, { Component } from 'react';
//import logo from './logo.svg';
import './App.css';
import {GlobalSentiments, EdgeNodeSentiments} from './Sentiment';


class App extends Component {
  render() {
    // Return components to render. This is JSX, see https://reactjs.org/docs/introducing-jsx.html
    return (
      <div className="App">
        <header className="App-header">
          {/* <img src={logo} className="App-logo" alt="logo" /> */}
          <h1 className="App-title">IBM SDR Insights Edge Application</h1>
        </header>
        <h1 className="Page-title">Global Keyword Sentiments</h1>
        <GlobalSentiments />
        {/*todo: just an example, remove this */}
        <br></br>
        <h1 className="Page-title">Keyword Sentiments From a Sample Edge Node</h1>
        <EdgeNodeSentiments />
      </div>
    );
  }
}

export default App;
