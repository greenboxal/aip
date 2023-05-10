import React from 'react'

import {ApolloClient, ApolloProvider} from '@apollo/client';

import {Admin, CustomRoutes, DataProvider, ListGuesser, Resource, ShowGuesser,} from "react-admin"

import {Route} from "react-router";

import buildDataProvider from "./data";

import AppLayout from "./layout/AppLayout";
import Dashboard from "./dashboard";
import ChatPage from "./chat";

import {ImageList, ImageShow} from "./resources/Image";
import {JobList, JobShow} from "./resources/Job";
import {PageCreate, PageList, PageShow} from "./resources/Page";
import {authProvider} from "./authProvider";
import {MessageList, MessageShow} from "./resources/Message";
import {ChannelList, ChannelShow} from "./resources/Channel";
import {EndpointList, EndpointShow} from "./resources/Endpoint";

const App: React.FC = () => {
    const [dataProviderAndClient, setDataProviderAndClient] = React.useState<{
        client: ApolloClient<any>,
        dataProvider: DataProvider,
    }>(null)

    React.useEffect(() => {
        buildDataProvider()
            .then(result => {
                setDataProviderAndClient(() => result)
            })
    }, [])

    if (!dataProviderAndClient) {
        return (<div>Loading</div>)
    }

    const {client, dataProvider} = dataProviderAndClient

    return (
        <ApolloProvider client={client}>
            <Admin
                layout={AppLayout}
                dataProvider={dataProvider}
                authProvider={authProvider}
                dashboard={Dashboard}
            >

                <Resource
                    name="Image"
                    list={ImageList}
                    show={ImageShow}
                    recordRepresentation={(record) => `${record.id} : ${record.spec.title}`}
                />

                <Resource
                    name="Page"
                    list={PageList}
                    show={PageShow}
                    create={PageCreate}
                    recordRepresentation={(record) => `${record.id} : ${record.spec.title}`}
                />

                <Resource
                    name="Job"
                    list={JobList}
                    show={JobShow}
                />

                <Resource
                    name="Memory"
                    list={ListGuesser}
                    show={ShowGuesser}
                />

                <Resource
                    name="Channel"
                    list={ChannelList}
                    show={ChannelShow}
                />

                <Resource
                    name="Message"
                    list={MessageList}
                    show={MessageShow}
                />

                <Resource
                    name="Endpoint"
                    list={EndpointShow}
                    show={EndpointList}
                />

                <CustomRoutes>
                    <Route path="/chat" element={<ChatPage/>}/>
                </CustomRoutes>
            </Admin>
        </ApolloProvider>
    )
}

export default App
