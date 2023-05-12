import React from 'react'

import {
    ChipField,
    Datagrid,
    DatagridConfigurable,
    DateField,
    ReferenceArrayField,
    ReferenceField,
    SelectColumnsButton,
    Show,
    SimpleShowLayout,
    TextField,
    TopToolbar,
} from "react-admin"

import {ListLive, ShowLive} from '@react-admin/ra-realtime';
import {MarkdownField} from "@react-admin/ra-markdown";
import {SpanListDatagrid, SpanShow} from "./Span";

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
    <ShowLive>
        <SimpleShowLayout>
            <TextField source="id" label="ID"/>

            <ReferenceField reference="Span" source="root_span_id" label="Root Span">
                <TextField source="id" label="ID"/>
                <TextField source="name" label="Name"/>
                <DateField source="started_at" label="Started At" />
                <DateField source="completed_at" label="Completed At" />

                <ReferenceField source="trace_id" reference="Trace" label="Trace" link="show">
                    <ChipField source="metadata.id"/>
                </ReferenceField>

                <ReferenceField reference="Span" source="parent_id" label="Parent" link="show">
                    <ChipField source="metadata.id"/>
                </ReferenceField>

                <ReferenceArrayField reference="Span" source="inner_span_ids">
                    <SpanListDatagrid />
                </ReferenceArrayField>
            </ReferenceField>
        </SimpleShowLayout>
    </ShowLive>
);
