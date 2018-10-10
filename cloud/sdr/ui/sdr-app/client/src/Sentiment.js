// Displays the global word sentiment list

// import React, { Component } from 'react';
import React from 'react';
import {
    Query,
    graphql,
} from 'react-apollo';
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
query nodenouns($edgenode: String!, $limit: Int!) {
    nodenouns(edgenode: $edgenode, limit: $limit) {
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
                render={({ rows, headers, getHeaderProps }) => (
                    <TableContainer title="Global Keyword Sentiments">
                        <Table>
                            <TableHead>
                                <TableRow>
                                    {headers.map(header => (
                                        <TableHeader {...getHeaderProps({ header })}>
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

export const EdgeNodeSentiments = (props => {

    return (
        <Query query={EDGE_NODE_NOUNS_LIST} variables={{limit: 20, edgenode: props.nodeId}}>
        {({loading, error, data}) => {
            if (loading) return "Loading..."
            if (error) return `Error! ${error.message}`

            let nodenouns = []

            if (data && data.nodenouns) {
                nodenouns = data.nodenouns.map(o => {
                    return Object.assign({}, o, {
                        id: o.noun,
                        timeupdated: moment(o.timeupdated).toString(),
                    })
                })
            }

            return (
                <DataTable
                    headers={globalSentimentHeaders}
                    rows={nodenouns}
                    render={({ rows, headers, getHeaderProps }) => (
                        <TableContainer title={`The top ${NOUN_LIMIT} nouns on node: ${props && props.nodeId}`}>
                            <Table>
                                <TableHead>
                                    <TableRow>
                                        {headers.map(header => (
                                            <TableHeader {...getHeaderProps({ header })}>
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
            )
        }}
        </Query>
    )
});
