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

import { ListLive } from '@react-admin/ra-realtime';

import { MarkdownField } from "@react-admin/ra-markdown"

import { MultiLevelMenu, AppLocationContext, Breadcrumb} from '@react-admin/ra-navigation'

import DefaultIcon from '@mui/icons-material/ViewList'

import { ReactQueryDevtools } from "react-query/devtools"
import {Card} from "@mui/material";

export const MessageListActions = () => (<TopToolbar>
    <SelectColumnsButton />
</TopToolbar>)

export const MessageList = () => (
    <ListLive actions={<MessageListActions />}>
        <DatagridConfigurable rowClick="show" size="small" preferenceKey="messages.datagrid">
            <TextField source="id" label="ID" />
            <TextField source="channel" label="Channel" />
            <TextField source="from" label="From" />
            <TextField source="text" label="Text" />

            <DateField source="metadata.created_at" label="Created At" />
        </DatagridConfigurable>
    </ListLive>
)

export const MessageShow = () => (
    <Show>
        <SimpleShowLayout>
            <TextField source="id" label="ID" />
            <TextField source="channel" label="Channel" />
            <TextField source="from" label="From" />
            <TextField source="text" label="Text" />

            <DateField source="metadata.created_at" label="Created At" />
        </SimpleShowLayout>
    </Show>
);
