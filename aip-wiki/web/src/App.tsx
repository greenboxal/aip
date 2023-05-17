import React from "react";

import {BrowserRouter, Routes} from "react-router-dom";
import {Route} from "react-router";

import WikiAdmin from "./admin";
import Wiki from "./wiki";
import {ApolloClient, ApolloProvider} from "@apollo/client";
import {DataProvider} from "react-admin";
import buildDataProvider from "./admin/data";

const App = () => {
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

    return <ApolloProvider client={client}>
        <BrowserRouter>
            <Routes>
                <Route path="/admin/*" element={<WikiAdmin baseName="/admin" client={client} dataProvider={dataProvider} />} />
                <Route path="/*" element={<Wiki />} />
            </Routes>
        </BrowserRouter>
    </ApolloProvider>
}

export default App