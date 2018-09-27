// Displays the global word sentiment list

// import React, { Component } from 'react';
import React from 'react';
import { graphql } from 'react-apollo';
import { gql } from 'apollo-boost';

import {
    DataTable,
} from 'carbon-components-react';

import moment from 'moment'

import './Sentiment.css';

const {
    TableContainer,
    Table,
    TableHead,
    TableRow,
    TableBody,
    TableCell,
    TableHeader,
  } = DataTable

const NOUN_LIMIT = 20;
const TEMP_EDGE_NODE = 'ibm/isaac_x86_desktop';     //todo: just to test nodenouns table, remove eventually
const TEMP_EDGE_NODE_LIMIT = 5;     //todo: just to test nodenouns table, remove eventually

const GLOBALNOUNS_LIST = gql`
{
    globalnouns(limit: ${NOUN_LIMIT}) {
        noun
        sentiment
        numberofmentions
        timeupdated
    }
}
`;

const EDGE_NODE_NOUNS_LIST = gql`
{
    nodenouns(edgenode: "${TEMP_EDGE_NODE}", limit: ${TEMP_EDGE_NODE_LIMIT}) {
        noun
        sentiment
        numberofmentions
        timeupdated
    }
}
`;

const globalSentimentHeaders = [
    {
        key: 'noun',
        header: 'Keyword',
    }, {
        key: 'sentiment',
        header: 'Sentiment',
    }, {
        key: 'numberofmentions',
        header: 'Number of Mentions',
    }, {
        key: 'timeupdated',
        header: 'Last Updated',
    },
]

export const GlobalSentiments = graphql(GLOBALNOUNS_LIST)(props => { 

    let globalNouns = []

    if (props && props.data && props.data.globalnouns) {
        globalNouns = props.data.globalnouns.map(o => {
            return Object.assign({}, o, {
                id: o.noun,
                timeupdated: moment(o.timeupdated).toString(),
            })
        })
    }
    
    return (
        <div>
            <DataTable
                headers={globalSentimentHeaders}
                rows={globalNouns}
                render={({rows, headers, getHeaderProps}) => (
                    <TableContainer title="Global Keyword Sentiments">
                        <Table>
                            <TableHead>
                                <TableRow>
                                    {headers.map(header => (
                                        <TableHeader {...getHeaderProps({header})}>
                                            {header.header}
                                        </TableHeader>
                                    ))}
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {rows.map(row => (
                                    <TableRow key={row.id}>
                                        {row.cells.map(cell => (
                                            <TableCell key={cell.id}>
                                                {cell.value}
                                            </TableCell>
                                        ))}
                                    </TableRow>
                                ))}
                            </TableBody>
                        </Table>
                    </TableContainer>
                )}
            />
        </div>
    )
});

// client.query({ query: GLOBALNOUNS_LIST }).then(console.log);

// export const GlobalSentiments = graphql(GLOBALNOUNS_LIST)(props => { return (
//     <div>
//     <p className="page-description">The top {NOUN_LIMIT} keywords mentioned on all of the edge nodes:</p>
//     <table className="Sentiment-table">
//         <thead>
//             <tr>
//                 <th className="Sentiment-cell">Keyword</th>
//                 <th>Sentiment</th>
//                 <th>Number of Mentions</th>
//                 <th>Last Updated</th>
//             </tr>
//         </thead>
//         <tbody>
//         {props.data.loading ? '' : props.data.globalnouns.map((row) =>
//             <tr key={row.noun}>
//                 <td>{row.noun}</td> <td>{row.sentiment}</td> <td>{row.numberofmentions}</td> <td>{row.timeupdated}</td>
//             </tr>
//         )}
//         </tbody>
//     </table>
//     </div>
// )}
// );

export const EdgeNodeSentiments = graphql(EDGE_NODE_NOUNS_LIST)(props => { return (
    <div>
    <p>The top {TEMP_EDGE_NODE_LIMIT} keywords mentioned on edge node <strong>{TEMP_EDGE_NODE}</strong>:</p>
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
        {props.data.loading ? '' : props.data.nodenouns.map((row) =>
            <tr key={row.noun}>
                <td>{row.noun}</td> <td>{row.sentiment}</td> <td>{row.numberofmentions}</td> <td>{row.timeupdated}</td>
            </tr>
        )}
        </tbody>
    </table>
    </div>
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
  
//export default Sentiment;