import React from 'react'

import {DatagridConfigurable, SelectColumnsButton, Show, SimpleShowLayout, TextField, TopToolbar,} from "react-admin"

import {ListLive} from '@react-admin/ra-realtime';

export const JobListActions = () => (<TopToolbar>
    <SelectColumnsButton/>
</TopToolbar>)

export const JobList = () => (
    <ListLive actions={<JobListActions/>}>
        <DatagridConfigurable rowClick="show" size="small" preferenceKey="images.datagrid">
            <TextField source="id" label="ID"/>
            <TextField source="status.state" label="State"/>
            <TextField source="spec.handler" label="Handler"/>
        </DatagridConfigurable>
    </ListLive>
)

export const JobShow = () => (
    <Show>
        <SimpleShowLayout>
            <TextField source="id" label="ID"/>
            <TextField source="status.state" label="State"/>

            <TextField source="spec.handler" label="Handler"/>
            <TextField source="spec.payload" label="Payload"/>

            <TextField source="status.last_error" label="Last Error"/>
            <TextField source="status.progress" label="Progress"/>
            <TextField source="status.result" label="Result"/>
        </SimpleShowLayout>
    </Show>
);
