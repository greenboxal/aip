import React, {createElement} from 'react'

import {
    Admin,
    Resource,
    ListGuesser,
    DataProvider,
    ShowGuesser,
    Layout,
    LayoutProps,
    defaultTheme,
    useResourceDefinitions,
    useGetResourceLabel,
    useCreatePath,
    MenuProps,
    TitlePortal,
    AppBar,
    LocalesMenuButton,
    ToggleThemeButton,
    RaThemeOptions,
} from "react-admin"

import buildGraphQLProvider from 'ra-data-graphql-simple';
import {addSearchMethod, Search} from "@react-admin/ra-search";

import { MultiLevelMenu, AppLocationContext, Breadcrumb} from '@react-admin/ra-navigation'
import { ReactQueryDevtools } from "react-query/devtools"
import {ImageList, ImageShow} from "./resources/Image";
import {JobList, JobShow} from "./resources/Job";
import {PageCreate, PageList, PageShow} from "./resources/Page";
import DefaultIcon from '@mui/icons-material/ViewList'

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

const darkTheme = {
    ...defaultTheme,
    palette: {
        mode: 'dark',
    },
} as RaThemeOptions

export const MyAppBar = () => (
    <AppBar>
        <TitlePortal />
        <Search />
        <ToggleThemeButton lightTheme={defaultTheme} darkTheme={darkTheme}  />
        <LocalesMenuButton />
    </AppBar>
);

const AppMenu: React.FC<MenuProps> = (props: MenuProps) => (
    <MultiLevelMenu {...props}>
        <MultiLevelMenu.Item label="Content Management" name="content-management">
            <ResourceMenuItem name="Page" />
            <ResourceMenuItem name="Image" />
        </MultiLevelMenu.Item>

        <MultiLevelMenu.Item label="Infrastructure" name="infrastructure">
            <ResourceMenuItem name="Job" />
            <ResourceMenuItem name="Memory" />
        </MultiLevelMenu.Item>
    </MultiLevelMenu>
)

const AppLayout: React.FC<LayoutProps> = ({ children, ...rest }) => (<>
    <AppLocationContext>
        <Layout {...rest} menu={AppMenu} appBar={MyAppBar}>
            <Breadcrumb></Breadcrumb>
            {children}
        </Layout>
    </AppLocationContext>
    <ReactQueryDevtools />
</>)

const Dashboard: React.FC = () => (
    <main>

    </main>
)

const App: React.FC = () => {
    const [baseDataProvider, setBaseDataProvider] = React.useState<DataProvider>(null)

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
                setBaseDataProvider(() => graphQlDataProvider)
            })
    }, [])

    if (!baseDataProvider) {
        return (<div>Loading</div>)
    }

    const dataProvider = addSearchMethod(baseDataProvider, [
        "Image",
        "Page",
    ]);

    return (
        <Admin
            layout={AppLayout}
            dataProvider={dataProvider}
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
        </Admin>
    )
}

export default App
