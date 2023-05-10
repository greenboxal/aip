import React from 'react'

import {ArrayField, DatagridConfigurable, DateField, ReferenceArrayField, SelectColumnsButton, Show, SimpleShowLayout, TextField, TopToolbar,} from "react-admin"

import {ListLive} from '@react-admin/ra-realtime';
import {MarkdownField} from "@react-admin/ra-markdown";

export const ChannelListActions = () => (<TopToolbar>
    <SelectColumnsButton/>
</TopToolbar>)

export const ChannelList = () => (
    <ListLive actions={<ChannelListActions/>}>
        <DatagridConfigurable rowClick="show" size="small" preferenceKey="channels.datagrid">
            <TextField source="id" label="ID"/>

            <ReferenceArrayField source="subscribers" label="Members" reference="Endpoint" />

            <DateField source="metadata.created_at" label="Created At"/>
        </DatagridConfigurable>
    </ListLive>
)

export const ChannelShow = () => (
    <Show>
        <SimpleShowLayout>
            <TextField source="id" label="ID"/>

            <ReferenceArrayField source="subscribers" label="Members" reference="Endpoint" />

            <DateField source="metadata.created_at" label="Created At"/>
        </SimpleShowLayout>
    </Show>
);
