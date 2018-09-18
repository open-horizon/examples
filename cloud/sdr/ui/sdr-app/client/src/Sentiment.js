// Displays the global word sentiment list

// import React, { Component } from 'react';
import React from 'react';
import { graphql } from 'react-apollo';
import { gql } from 'apollo-boost';
import './Sentiment.css';

const GLOBALNOUNS_LIST = gql`
{
    globalnouns {
        noun
        sentiment
        numberofmentions
        timeupdated
    }
}
`;
// client.query({ query: GLOBALNOUNS_LIST }).then(console.log);

const Sentiment = graphql(GLOBALNOUNS_LIST)(props =>
    <ul className="Sentiment-list">
        {props.data.loading ? '' : props.data.globalnouns.map((row) =>
            <li key={row.noun}>
                <strong>{row.noun}:</strong> Sentiment: {row.sentiment}, Number Of Mentions: {row.numberofmentions}, Last Updated: {row.timeupdated}
            </li>
        )}
    </ul>
);

/* class Sentiment extends Component {
    constructor(props) {
        super(props);
        this.state = { insights: ['Trump: negative', 'WorldCup: positive'] };
    }
    componentDidMount() { this.state.insights = ['Trump', 'Soccer'] }

    render() {
        const listItems = this.state.insights.map((word) =>
            <li key={word}>{word}</li>
        );
        return (
            <ul className="Sentiment-list">
            {listItems}
            </ul>
        );
    }
} */
  
export default Sentiment;