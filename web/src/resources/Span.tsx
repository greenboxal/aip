import React from 'react'

import {
    ArrayField,
    ChipField,
    Datagrid,
    DatagridConfigurable,
    DateField, NumberField,
    ReferenceArrayField,
    ReferenceField,
    SelectColumnsButton,
    Show,
    SimpleShowLayout,
    TextField,
    TopToolbar, WrapperField,
} from "react-admin"

import {ListLive, ShowLive} from '@react-admin/ra-realtime';
import {MarkdownField} from "@react-admin/ra-markdown";

export const SpanListActions = () => (<TopToolbar>
    <SelectColumnsButton/>
</TopToolbar>)

export const SpanListDatagrid: React.FC<{
    includeTraceId?: boolean,
}> = (props) => (
    <Datagrid rowClick="show" size="small">
        <TextField source="metadata.id" label="ID"/>

        { props.includeTraceId ? (<ReferenceField source="trace_id" reference="Trace" label="Trace ID" link="show" />) : null }

        <TextField source="name" label="Name"/>

        <DateField source="started_at" label="Started At" showDate={true} showTime={true} />
        <DateField source="completed_at" label="Completed At" showDate={true} showTime={true} />
        <NumberField source="duration" label="Duration" />
    </Datagrid>
)

export const SpanList = () => (
    <ListLive actions={<SpanListActions/>} sort={{ field: 'metadata.updated_at', order: 'DESC' }}>
        <SpanListDatagrid includeTraceId={true} />
    </ListLive>
)

export const SpanShow = () => (
    <ShowLive>
        <SimpleShowLayout>
            <TextField source="id" label="ID"/>
            <TextField source="name" label="Name"/>
            <DateField source="started_at" label="Started At" showDate={true} showTime={true} />
            <DateField source="completed_at" label="Completed At" showDate={true} showTime={true} />
            <NumberField source="duration" label="Duration" />

            <ReferenceField source="trace_id" reference="Trace" label="Trace" link="show">
                <ChipField source="metadata.id"/>
            </ReferenceField>

            <ReferenceField reference="Span" source="parent_id" label="Parent" link="show">
                <ChipField source="metadata.id"/>
            </ReferenceField>

            <ArrayField source="attributes" label="Attributes">
                <Datagrid bulkActionButtons={false}>
                    <TextField source="key" label="Key"/>
                    <TextField source="value" label="Value"/>
                </Datagrid>
            </ArrayField>

            <ReferenceArrayField reference="Span" source="inner_span_ids">
                <SpanListDatagrid />
            </ReferenceArrayField>

        </SimpleShowLayout>
    </ShowLive>
);
