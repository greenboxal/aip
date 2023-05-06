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
    MenuProps, TabbedShowLayout, ArrayField, Create, SimpleForm, TextInput,
} from "react-admin"

import { MarkdownField } from "@react-admin/ra-markdown"

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

export const PageCreate = () => (
    <Create>
        <SimpleForm>
            <TextInput source="title" label="Title" />
            <TextInput source="voice" label="Voice" />
            <TextInput source="language" label="Language" />
        </SimpleForm>
    </Create>
);
