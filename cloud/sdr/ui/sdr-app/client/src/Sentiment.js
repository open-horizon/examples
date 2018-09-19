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

const Sentiment = graphql(GLOBALNOUNS_LIST)(props => { return (
    <table className="Sentiment-table">
        <thead>
            <tr>
                <th className="Sentiment-cell">Keyword</th>
                <th>Sentiment</th>
                <th>Number of Mentions</th>
                <th>Last Updated</th>
            </tr>
        </thead>
        <tbody>
        {props.data.loading ? '' : props.data.globalnouns.map((row) =>
            <tr key={row.noun}>
                <td>{row.noun}</td> <td>{row.sentiment}</td> <td>{row.numberofmentions}</td> <td>{row.timeupdated}</td>
            </tr>
        )}
        </tbody>
    </table>
)}
);

/* leaving this here as reference for now...
class Sentiment extends Component {
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