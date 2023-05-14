import React from 'react'

import {ArrayField, DatagridConfigurable, DateField, ReferenceArrayField, SelectColumnsButton, Show, SimpleShowLayout, TextField, TopToolbar,} from "react-admin"

import {ListLive} from '@react-admin/ra-realtime';
import {MarkdownField} from "@react-admin/ra-markdown";

export const EndpointListActions = () => (<TopToolbar>
    <SelectColumnsButton/>
</TopToolbar>)

export const EndpointList = () => (
    <ListLive actions={<EndpointListActions/>}>
        <DatagridConfigurable rowClick="show" size="small" preferenceKey="endpoints.datagrid">
            <TextField source="id" label="ID"/>

            <ReferenceArrayField source="subscriptions" label="Channels" reference="Channel" />

            <DateField source="metadata.created_at" label="Created At"/>
        </DatagridConfigurable>
    </ListLive>
)

export const EndpointShow = () => (
    <Show>
        <SimpleShowLayout>
            <TextField source="id" label="ID"/>

            <ReferenceArrayField source="subscriptions" label="Channels" reference="Channel" />

            <DateField source="metadata.created_at" label="Created At"/>
        </SimpleShowLayout>
    </Show>
);
