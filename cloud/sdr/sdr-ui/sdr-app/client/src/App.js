// Main component for the UI portion of the SDR Edge App.
// Built on react: https://reactjs.org/docs/getting-started.html

import React, { Component } from 'react';
import logo from './logo.svg';
import './App.css';
import Sentiment from './Sentiment';


class App extends Component {
  render() {
    // Return components to render. This is JSX, see https://reactjs.org/docs/introducing-jsx.html
    return (
      <div className="App">
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h1 className="App-title">Welcome to the IBM SDR Insights App</h1>
        </header>
        <h1>Insights</h1>
        {/* <p className="App-intro">
          To get started, edit <code>src/App.js</code> and save to reload.
        </p> */}
        <Sentiment />
      </div>
    );
  }
}

export default App;
