import React from 'react'

import {
    ArrayField, Datagrid,
    DatagridConfigurable,
    DateField,
    ReferenceArrayField,
    ReferenceManyField,
    ReferenceOneField,
    SelectColumnsButton,
    Show,
    SimpleShowLayout,
    TextField,
    TopToolbar,
} from "react-admin"

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

            <DateField source="metadata.created_at" label="Created At"/>
            <DateField source="metadata.updated_at" label="Updated At"/>

            <ReferenceArrayField source="subscribers" label="Members" reference="Endpoint" />

            <ReferenceManyField reference="Message" target="channel" label="Messages">
                <Datagrid rowClick="show" size="small" bulkActionButtons={false} sort={{ field: 'metadata.created_at', order: 'ASC' }}>
                    <TextField source="id" label="ID"/>
                    <DateField source="metadata.created_at" label="Created At"/>
                    <TextField source="from" label="From"/>
                    <MarkdownField source="text" label="Text"/>

                </Datagrid>
            </ReferenceManyField>
        </SimpleShowLayout>
    </Show>
);
