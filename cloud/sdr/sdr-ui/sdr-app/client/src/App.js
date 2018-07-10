import React, { Component } from 'react';
import logo from './logo.svg';
import './App.css';
import Sentiments from './Sentiment';


class App extends Component {
  render() {
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
        <Sentiments />
      </div>
    );
  }
}

export default App;
