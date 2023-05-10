import React from "react";
import {DataProvider} from "react-admin";
import {ApolloClient, gql, HttpLink, InMemoryCache, split} from "@apollo/client";
import {GraphQLWsLink} from "@apollo/client/link/subscriptions";
import {createClient} from "graphql-ws";
import {getMainDefinition} from "@apollo/client/utilities";
import buildGraphQLProvider from "ra-data-graphql-simple";
import {CREATE} from "ra-core";
import {addSearchMethod} from "@react-admin/ra-search";
import {setContext} from "@apollo/client/link/context";
import {auth0} from "../authProvider";

function enhanceDataProvider(client: ApolloClient<any>, baseDataProvider: DataProvider): DataProvider {
    let subscriptions: any = {};

    const searchDataProvider = addSearchMethod(baseDataProvider, [
        "Image",
        "Page",
    ]);

    const dataProvider = {
        ...searchDataProvider,

        subscribe: async (topic: string, subscriptionCallback: any) => {
            let sub = subscriptions[topic]

            if (sub == null) {
                sub = {
                    callbacks: [],
                }

                subscriptions[topic] = sub

                const resourceType = topic.replace(/^resource\//, '')

                sub.observable = client.subscribe({
                    variables: {
                        resourceType: resourceType,
                    },

                    query: gql`
                        subscription Sub($resourceType: String!) {
                            resourceChanged(resourceType: $resourceType) {
                                type

                                payload {
                                    ids
                                }
                            }
                        }
                    `
                })

                sub.subscription = sub.observable.subscribe((data: any) => {
                    if (data.data.resourceChanged == null) {
                        return
                    }

                    dataProvider.publish(topic, data.data.resourceChanged)
                })
            }

            sub.callbacks.push(subscriptionCallback)

            return Promise.resolve({data: null});
        },

        unsubscribe: async (topic: string, subscriptionCallback: any) => {
            let sub = subscriptions[topic]

            if (sub == null) {
                return
            }

            sub.callbacks = sub.callbacks.filter(
                (it: any) => it !== subscriptionCallback
            )

            if (sub.callbacks.length === 0) {
                sub.subscription.unsubscribe()

                delete subscriptions[topic]
            }

            return Promise.resolve({data: null});
        },

        publish: (topic: string, event: any) => {
            if (!topic) {
                return Promise.reject(new Error('missing topic'));
            }

            if (!event.type) {
                return Promise.reject(new Error('missing event type'));
            }

            let sub = subscriptions[topic]

            if (sub == null) {
                return
            }

            sub.callbacks.forEach((callback: any) => callback(event));

            return Promise.resolve({data: null});
        },
    }

    return dataProvider;
}

export default function buildDataProvider(): Promise<{ client: ApolloClient<any>, dataProvider: DataProvider }> {
    const httpLink = new HttpLink({
        uri: 'http://localhost:30100/v1/graphql',
    })

    const wsLink = new GraphQLWsLink(createClient({
        url: 'ws://localhost:30100/v1/graphql/ws',
        keepAlive: 10000,
        connectionParams: {
            reconnect: true,
        }
    }))

    const authLink = setContext((_, {headers}) => {
        return auth0.getTokenSilently().then((token) => {
            // return the headers to the context so httpLink can read them
            return {
                headers: {
                    ...headers,
                    authorization: token ? `Bearer ${token}` : "",
                }
            }
        })
    });

    const splitLink = split(
        ({query}) => {
            const definition = getMainDefinition(query);
            return (
                definition.kind === 'OperationDefinition' &&
                definition.operation === 'subscription'
            );
        },
        authLink.concat(wsLink),
        authLink.concat(httpLink),
    );

    const client = new ApolloClient({
        link: splitLink,
        connectToDevTools: true,

        cache: new InMemoryCache().restore({}),
    })

    return buildGraphQLProvider({
        client: client,

        introspection: {
            operationNames: {
                [CREATE]: (type) => {
                    switch (type.name) {
                        case "Page":
                            return "wikiPageManagerGetPage"
                    }

                    return undefined
                },
            },
        },
    })
        .then(graphQlDataProvider => {
            return {client, dataProvider: enhanceDataProvider(client, graphQlDataProvider)}
        })
}