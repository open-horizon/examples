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

// Get list of top $limit nouns for all edge nodes
const GLOBALNOUNS_LIST = gql`
query globalnouns($limit: Int!) {
    globalnouns(limit: $limit) {
        noun
        sentiment
        numberofmentions
        timeupdated
    }
}
`

// Get list of top $limit nouns for a single edge node
const EDGE_NODE_NOUNS_LIST = gql`
query nodenouns($edgenode: String!, $limit: Int!) {
    nodenouns(edgenode: $edgenode, limit: $limit) {
        noun
        sentiment
        numberofmentions
        timeupdated
    }
}
`

// Table layout for sentiments
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

export const GlobalSentiments = (props => {
    return (
        <Query
            query={GLOBALNOUNS_LIST}
            variables={{limit: NOUN_LIMIT}}
            pollInterval={1000}
        >
        {({loading, error, data}) => {
            if (loading) return "Loading..."
            if (error) return `Error! ${error.message}`

            let globalNouns = []

            if (data && data.globalnouns) {
                globalNouns = data.globalnouns.map(o => {
                    return Object.assign({}, o, {
                        id: o.noun,
                        timeupdated: moment(o.timeupdated).toString(),
                    })
                })
            }

            return (
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
            )

        }}
        </Query>
    )
})

export const EdgeNodeSentiments = (props => {

    return (
        <Query 
            query={EDGE_NODE_NOUNS_LIST} 
            variables={{limit: 20, edgenode: props.nodeId}}
            pollInterval={1000}
        >
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
})
