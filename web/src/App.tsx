import React, {createElement} from 'react'

import {
    Admin,
    Resource,
    ListGuesser,
    DataProvider,
    List,
    Datagrid,
    TextField,
    DateField,
    BooleanField,
    ShowGuesser,
    Layout,
    LayoutProps,
    Menu,
    SimpleShowLayout,
    Show,
    RichTextField,
    ShowButton,
    UrlField,
    ImageField,
    defaultTheme,
    useResourceDefinitions,
    useGetResourceLabel,
    useCreatePath,
    MenuItemLink,
    MenuProps, TabbedShowLayout, ArrayField,
} from "react-admin"

import { MarkdownField } from "@react-admin/ra-markdown"

import { MultiLevelMenu, AppLocationContext, Breadcrumb} from '@react-admin/ra-navigation'
import { ReactQueryDevtools } from "react-query/devtools"
import {ImageList, ImageShow} from "./resources/Image";
import {PageCreate, PageList, PageShow} from "./resources/Page";
import DefaultIcon from '@mui/icons-material/ViewList'

import buildGraphQLProvider, { buildQuery } from 'ra-data-graphql-simple';
import { CREATE } from 'ra-core';

export const ResourceMenuItem = ({ name }: { name: string }) => {
    const resources = useResourceDefinitions();
    const getResourceLabel = useGetResourceLabel();
    const createPath = useCreatePath();
    if (!resources || !resources[name]) return null;

    return (
        <MultiLevelMenu.Item
            name={name}
            to={createPath({
                resource: name,
                type: 'list',
            })}
            label={<>{getResourceLabel(name, 2)}</>}
            icon={
                resources[name].icon ? (
                    createElement(resources[name].icon)
                ) : (
                    <DefaultIcon />
                )
            }
        />
    );
};


const AppMenu: React.FC<MenuProps> = (props: MenuProps) => (
    <MultiLevelMenu {...props}>
        <ResourceMenuItem name="Page" />
        <ResourceMenuItem name="Image" />
    </MultiLevelMenu>
)

const AppLayout: React.FC<LayoutProps> = ({ children, ...rest }) => (<>
    <AppLocationContext>
        <Layout {...rest} menu={AppMenu}>
            <Breadcrumb></Breadcrumb>
            {children}
        </Layout>
    </AppLocationContext>
    <ReactQueryDevtools />
</>)

const App: React.FC = () => {
    const [dataProvider, setDataProvider] = React.useState<DataProvider>(null)

    React.useEffect(() => {
        buildGraphQLProvider({
            introspection: {
                operationNames: {
                    [CREATE]: (type) => {
                        switch (type.name) {
                            case "Page": return "wikiPageManagerGetPage"
                        }

                        return undefined
                    },
                },
            },

            clientOptions: {
                uri: 'http://localhost:30100/v1/graphql',
                connectToDevTools: true,
            },
        })
            .then(graphQlDataProvider => {
                setDataProvider(() => graphQlDataProvider)
            })
    }, [])

    if (!dataProvider) {
        return (<div>Loading</div>)
    }

    const theme: any = {
        ...defaultTheme,
        palette: {
            //mode: 'dark',
        },
    };

    return (
        <Admin layout={AppLayout} theme={theme} dataProvider={dataProvider}>
            <Resource name="Image" list={ImageList} show={ImageShow} />
            <Resource name="Page" list={PageList} show={PageShow} create={PageCreate} />

            <Resource name="Agent" list={ListGuesser}/>
            <Resource name="Profile" list={ListGuesser}/>

            <Resource name="Pipeline" list={ListGuesser}/>
            <Resource name="Task" list={ListGuesser}/>

            <Resource name="Memory" list={ListGuesser}/>
            <Resource name="MemoryData" list={ListGuesser}/>
        </Admin>
    )
}

export default App
