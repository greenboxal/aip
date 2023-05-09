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

export const JobListActions = () => (<TopToolbar>
    <SelectColumnsButton />
</TopToolbar>)

export const JobList = () => (
    <ListLive actions={<JobListActions />}>
        <DatagridConfigurable rowClick="show" size="small" preferenceKey="images.datagrid">
            <TextField source="id" label="ID" />
            <TextField source="status.state" label="State" />
            <TextField source="spec.handler" label="Handler" />
        </DatagridConfigurable>
    </ListLive>
)

export const JobShow = () => (
    <Show>
        <SimpleShowLayout>
            <TextField source="id" label="ID" />
            <TextField source="status.state" label="State" />

            <TextField source="spec.handler" label="Handler" />
            <TextField source="spec.payload" label="Payload" />

            <TextField source="status.last_error" label="Last Error" />
            <TextField source="status.progress" label="Progress" />
            <TextField source="status.result" label="Result" />
        </SimpleShowLayout>
    </Show>
);
