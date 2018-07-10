import React, { Component } from 'react';
import './Sentiment.css';

class Sentiment extends Component {
    state = { insights: ['Trump: negative', 'Soccer: positive'] }

    /*
    componentDidMount() {
        this.state.insights = ['Trump', 'Soccer']
    }
    */

    render() {
        const listItems = this.state.insights.map((word) =>
            <li>{word}</li>
        );
        return (
            <ul className="Sentiment-list">
            {listItems}
            </ul>
        )
    }
  }
  
export default Sentiment;