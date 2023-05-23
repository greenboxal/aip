import {DataProvider, GET_LIST, GET_MANY, GET_MANY_REFERENCE} from "react-admin";
import {ApolloClient, gql, HttpLink, InMemoryCache, split} from "@apollo/client";
import {GraphQLWsLink} from "@apollo/client/link/subscriptions";
import {createClient} from "graphql-ws";
import {getMainDefinition} from "@apollo/client/utilities";
import buildGraphQLProvider, {buildQuery} from "./base";
import {CREATE} from "ra-core";
import {addSearchMethod} from "@react-admin/ra-search";
import {setContext} from "@apollo/client/link/context";
import {auth0} from "../../authProvider";
import {IntrospectionResult} from "ra-data-graphql";
import {IntrospectionField} from "graphql";

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

const API_HTTP_BASE_URL = process.env.REACT_APP_AIP_API_HTTP_BASE_URL || "http://localhost:30100/v1/graphql"
const API_WS_BASE_URL = process.env.REACT_APP_AIP_API_WS_BASE_URL || "ws://localhost:30100/v1/graphql/ws"

export default async function buildDataProvider(): Promise<{ client: ApolloClient<any>, dataProvider: DataProvider }> {
    const httpLink = new HttpLink({
        uri: API_HTTP_BASE_URL,
    })

    const wsLink = new GraphQLWsLink(createClient({
        url: API_WS_BASE_URL,
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

    const sanitizeResource = (data: any) => {
        if (!data) {
            return data
        }

        const metadata: any = data.metadata || {}

        data.id = metadata.id || data.id

        return data
    }

    const gqlBuildQuery = (introspectionResults: IntrospectionResult) => {
        const baseQuery = buildQuery(introspectionResults)

        return (raFetchMethod: string,
                resource: string,
                queryType: IntrospectionField
        ) => {
            const resourceQuery = baseQuery(raFetchMethod, resource, queryType)
            const baseParseResponse = resourceQuery.parseResponse

            resourceQuery.parseResponse = (response) => {
                const baseResponse = baseParseResponse(response)

                if (
                    raFetchMethod === GET_LIST ||
                    raFetchMethod === GET_MANY ||
                    raFetchMethod === GET_MANY_REFERENCE
                ) {
                    return {
                        data: baseResponse.data.map(sanitizeResource),
                        total: baseResponse.total,
                    };
                }

                return {
                    data: sanitizeResource(baseResponse.data)
                };
            }

            return resourceQuery
        }
    }

    const gqlProvider = await buildGraphQLProvider({
        client: client,
        buildQuery: gqlBuildQuery,

        introspection: {
            operationNames: {
                [CREATE]: (type) => {
                    switch (type.name) {
                        case "Page":
                            return "wikiPageManagerGetPage"
                    }

                    return "create" + type.name
                },
            },
        },
    })

    const dataProvider = enhanceDataProvider(client, gqlProvider)

    return {
        client,
        dataProvider,
    }
}