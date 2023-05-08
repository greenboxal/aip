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
    MenuProps, TabbedShowLayout, ArrayField, WrapperField, SimpleForm, SingleFieldList, SelectColumnsButton, TopToolbar, DatagridConfigurable,
} from "react-admin"

import { MarkdownField } from "@react-admin/ra-markdown"

import { MultiLevelMenu, AppLocationContext, Breadcrumb} from '@react-admin/ra-navigation'

import DefaultIcon from '@mui/icons-material/ViewList'

import { ReactQueryDevtools } from "react-query/devtools"
import {Card} from "@mui/material";

export const ImageListActions = () => (<TopToolbar>
    <SelectColumnsButton />
</TopToolbar>)

export const ImageList = () => (
    <List actions={<ImageListActions />}>
        <DatagridConfigurable rowClick="show" size="small" preferenceKey="images.datagrid">
            <ImageField source="status.url" title="Thumbnail" label="" />
            <TextField source="id" label="ID" />
            <TextField source="status.prompt" label="Prompt" />

            <TextField source="spec.path" label="Path" />
        </DatagridConfigurable>
    </List>
)

export const ImageShow = () => (
    <Show>
        <SimpleShowLayout>
            <TextField source="id" label="ID" />

            <DateField source="metadata.created_at" label="Created At" />
            <DateField source="metadata.updated_at" label="Updated At" />

            <TextField source="spec.path" label="Path" />
            <TextField source="status.prompt" label="Prompt" />
            <UrlField source="status.url" label="URL" />

            <ImageField source="status.url" title="Preview" />
        </SimpleShowLayout>
    </Show>
);
