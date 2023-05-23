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
import {authProvider} from "../authProvider";
import {MessageList, MessageShow} from "./resources/Message";
import {ChannelList, ChannelShow} from "./resources/Channel";
import {EndpointList, EndpointShow} from "./resources/Endpoint";

import ArticleIcon from '@mui/icons-material/Article';
import ImageIcon from '@mui/icons-material/Image';
import DnsIcon from '@mui/icons-material/Dns';
import SpaceDashboardIcon from '@mui/icons-material/SpaceDashboard';
import WorkIcon from '@mui/icons-material/Work';
import MemoryIcon from '@mui/icons-material/Memory';
import TopicIcon from '@mui/icons-material/Topic';
import MessageIcon from '@mui/icons-material/Message';
import PersonIcon from '@mui/icons-material/Person';
import SupportAgentIcon from '@mui/icons-material/SupportAgent';
import AccountBoxIcon from '@mui/icons-material/AccountBox';
import AssignmentIcon from '@mui/icons-material/Assignment';
import GroupIcon from '@mui/icons-material/Group';
import GroupWorkIcon from '@mui/icons-material/GroupWork';
import RouteIcon from '@mui/icons-material/Route';
import {TraceList, TraceShow} from "./resources/Trace";
import {SpanList, SpanShow} from "./resources/Span";
import {DomainCreate, DomainList, DomainShow} from "./resources/Domain";

type WikiAdminProps = {
    client: ApolloClient<any>,
    dataProvider: DataProvider,
    baseName: string,
};

const WikiAdmin: React.FC<WikiAdminProps> = ({ client, dataProvider, baseName }) => {
    return (
        <Admin
            basename={baseName}
            dataProvider={dataProvider}
            authProvider={authProvider}
            layout={AppLayout}
            dashboard={Dashboard}
        >
            <Resource
                name="Image"
                icon={ImageIcon}
                list={ImageList}
                show={ImageShow}
                recordRepresentation={(record) => `${record.id} : ${record.spec.title}`}
            />

            <Resource
                name="Page"
                icon={ArticleIcon}
                list={PageList}
                show={PageShow}
                create={PageCreate}
                recordRepresentation={(record) => `${record.id} : ${record.spec.title}`}
            />

            <Resource name="Domain" icon={DnsIcon} list={DomainList} show={DomainShow} create={DomainCreate} />
            <Resource name="RouteBinding" icon={RouteIcon} list={ListGuesser} show={ShowGuesser} />
            <Resource name="Layout" icon={SpaceDashboardIcon} list={ListGuesser} show={ShowGuesser} />

            <Resource name="Job" icon={WorkIcon} list={JobList} show={JobShow} />
            <Resource name="Memory" icon={MemoryIcon} list={ListGuesser} show={ShowGuesser} />

            <Resource name="Channel" icon={TopicIcon} list={ChannelList} show={ChannelShow} />
            <Resource name="Message" icon={MessageIcon} list={MessageList} show={MessageShow} />
            <Resource name="Endpoint" icon={PersonIcon} list={EndpointShow} show={EndpointList} />

            <Resource name="Agent" icon={SupportAgentIcon} list={ListGuesser} show={ShowGuesser} />
            <Resource name="Profile" icon={AccountBoxIcon} list={ListGuesser} show={ShowGuesser} />
            <Resource name="Task" icon={AssignmentIcon} list={ListGuesser} show={ShowGuesser} />
            <Resource name="Team" icon={GroupIcon} list={ListGuesser} show={ShowGuesser} />
            <Resource name="Pipeline" icon={GroupWorkIcon} list={ListGuesser} show={ShowGuesser} />

            <Resource name="Trace" list={TraceList} show={TraceShow} />
            <Resource name="Span" list={SpanList} show={SpanShow} />

            <CustomRoutes>
                <Route path="/chat" element={<ChatPage/>}/>
            </CustomRoutes>
        </Admin>
    )
}

export default WikiAdmin
