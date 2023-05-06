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

import DefaultIcon from '@mui/icons-material/ViewList'

import { ReactQueryDevtools } from "react-query/devtools"

import buildGraphQLProvider from "ra-data-graphql-simple"

export const PageList = () => (
    <List>
        <Datagrid>
            <TextField source="spec.title" label="Title" />
            <TextField source="spec.language" label="Language" />
            <TextField source="spec.voice" label="Voice" />

            <TextField source="id" label="ID" />

            <ShowButton />
        </Datagrid>
    </List>
);

export const PageShow = () => (
    <Show>
        <TabbedShowLayout>
            <TabbedShowLayout.Tab label="General">
                <TextField source="id" label="ID" />

                <DateField source="metadata.created_at" label="Created At" />
                <DateField source="metadata.updated_at" label="Updated At" />

                <TextField source="spec.title" label="Title" />
                <TextField source="spec.language" label="Language" />
                <TextField source="spec.voice" label="Voice" />

                <ArrayField source="status.images" label="Images">
                    <Datagrid bulkActionButtons={false}>
                        <TextField source="title" label="Title" />
                        <TextField source="source" label="Source" />
                    </Datagrid>
                </ArrayField>

                <ArrayField source="status.links" label="Links">
                    <Datagrid bulkActionButtons={false}>
                        <TextField source="title" label="Title" />
                        <TextField source="to" label="To" />
                    </Datagrid>
                </ArrayField>


                <MarkdownField source="status.markdown" label="Markdown Contents" />
            </TabbedShowLayout.Tab>
            <TabbedShowLayout.Tab label="HTML Preview">
                <RichTextField source="status.html" label="HTML Contents" />
            </TabbedShowLayout.Tab>
        </TabbedShowLayout>
    </Show>
);

export const ImageList = () => (
    <List>
        <Datagrid>
            <ImageField source="status.url" title="Thumbnail" label="" />

            <TextField source="id" label="ID" />

            <TextField source="spec.path" label="Path" />
            <TextField source="spec.prompt" label="Prompt" />
            <UrlField source="status.url" label="URL" />

            <ShowButton />
        </Datagrid>
    </List>
);

export const ImageShow = () => (
    <Show>
        <SimpleShowLayout>
            <TextField source="id" label="ID" />

            <DateField source="metadata.created_at" label="Created At" />
            <DateField source="metadata.updated_at" label="Updated At" />

            <TextField source="spec.path" label="Path" />
            <TextField source="spec.prompt" label="Prompt" />
            <UrlField source="status.url" label="URL" />

            <ImageField source="status.url" title="Preview" />
        </SimpleShowLayout>
    </Show>
);

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
            clientOptions: {
                uri: 'https://desciclo.ai/v1/graphql',
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
            <Resource name="Image" hasShow list={ImageList} show={ImageShow}/>
            <Resource name="Page" hasShow list={PageList} show={PageShow}/>

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
