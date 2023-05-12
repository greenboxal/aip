import React from 'react'

import {Datagrid, DatagridConfigurable, DateField, ReferenceArrayField, SelectColumnsButton, Show, SimpleShowLayout, TextField, TopToolbar,} from "react-admin"

import {ListLive} from '@react-admin/ra-realtime';
import {MarkdownField} from "@react-admin/ra-markdown";
import {SpanShow} from "./Span";

export const TraceListActions = () => (<TopToolbar>
    <SelectColumnsButton/>
</TopToolbar>)

export const TraceList = () => (
    <ListLive actions={<TraceListActions/>} sort={{ field: 'metadata.created_at', order: 'DESC' }}>
        <DatagridConfigurable rowClick="show" size="small" preferenceKey="traces.datagrid">
            <TextField source="id" label="ID"/>
        </DatagridConfigurable>
    </ListLive>
)

export const TraceShow = () => (
    <Show>
        <SimpleShowLayout>
            <TextField source="id" label="ID"/>

            <ReferenceArrayField reference="Span" source="spans" label="Spans">
                <Datagrid rowClick="show" size="small">
                    <TextField source="id" label="ID"/>
                    <TextField source="name" label="Name"/>
                    <DateField source="metadata.started_at" label="Started At" />
                    <DateField source="metadata.completed_at" label="Completed At" />
                </Datagrid>
            </ReferenceArrayField>
        </SimpleShowLayout>
    </Show>
);
