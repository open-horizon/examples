// Displays the global word sentiment list

// import React, { Component } from 'react';
import React from 'react';
import { graphql } from 'react-apollo';
import { gql } from 'apollo-boost';
import './Sentiment.css';

const NOUNS_LIST = gql`
{
    nouns {
        noun
        sentiment
        numberofmentions
        timeupdated
    }
}
`;
// client.query({ query: NOUNS_LIST }).then(console.log);

const Sentiment = graphql(NOUNS_LIST)(props =>
    <ul>
        {props.data.loading ? '' : props.data.nouns.map((row) =>
            <li key={row.noun}>
                <strong>{row.noun}:</strong> Sentiment: {row.sentiment}, Number Of Mentions: {row.numberofmentions}, Updated: {row.timeupdated}
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