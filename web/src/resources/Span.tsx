import React from 'react'

import {ChipField, DatagridConfigurable, DateField, ReferenceArrayField, ReferenceField, SelectColumnsButton, Show, SimpleShowLayout, TextField, TopToolbar,} from "react-admin"

import {ListLive} from '@react-admin/ra-realtime';
import {MarkdownField} from "@react-admin/ra-markdown";

export const SpanListActions = () => (<TopToolbar>
    <SelectColumnsButton/>
</TopToolbar>)

export const SpanList = () => (
    <ListLive actions={<SpanListActions/>} sort={{ field: 'metadata.created_at', order: 'DESC' }}>
        <DatagridConfigurable rowClick="show" size="small" preferenceKey="spans.datagrid">
            <TextField source="id" label="ID"/>
            <TextField source="name" label="Name"/>
            <DateField source="metadata.started_at" label="Started At" />
            <DateField source="metadata.completed_at" label="Completed At" />
        </DatagridConfigurable>
    </ListLive>
)

export const SpanShow = () => (
    <Show>
        <SimpleShowLayout>
            <TextField source="id" label="ID"/>
            <TextField source="name" label="Name"/>
            <DateField source="metadata.started_at" label="Started At" />
            <DateField source="metadata.completed_at" label="Completed At" />

            <ReferenceField source="trace_id" reference="Trace" label="Trace" link="show">
                <ChipField source="metadata.id"/>
            </ReferenceField>
        </SimpleShowLayout>
    </Show>
);
