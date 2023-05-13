import React from 'react'

import {DatagridConfigurable, DateField, SelectColumnsButton, Show, SimpleShowLayout, TextField, TopToolbar,} from "react-admin"

import {ListLive} from '@react-admin/ra-realtime';
import {MarkdownField} from "@react-admin/ra-markdown";

export const MessageListActions = () => (<TopToolbar>
    <SelectColumnsButton/>
</TopToolbar>)

export const MessageList = () => (
    <ListLive actions={<MessageListActions/>} sort={{ field: 'metadata.created_at', order: 'DESC' }}>
        <DatagridConfigurable rowClick="show" size="small" preferenceKey="messages.datagrid">
            <TextField source="id" label="ID"/>
            <TextField source="channel" label="Channel"/>
            <TextField source="from" label="From"/>
            <MarkdownField source="text" label="Text"/>

            <DateField source="metadata.created_at" label="Created At"/>
        </DatagridConfigurable>
    </ListLive>
)

export const MessageShow = () => (
    <Show>
        <SimpleShowLayout>
            <TextField source="id" label="ID"/>
            <TextField source="channel" label="Channel"/>
            <TextField source="from" label="From"/>

            <MarkdownField source="text" label="Text"/>

            <DateField source="metadata.created_at" label="Created At"/>
        </SimpleShowLayout>
    </Show>
);
