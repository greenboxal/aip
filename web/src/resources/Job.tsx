import React from 'react'

import {DatagridConfigurable, DateField, FilterButton, SelectColumnsButton, Show, SimpleShowLayout, TextField, TopToolbar,} from "react-admin"

import {ListLive} from '@react-admin/ra-realtime';
import {StackedFilters, textFilter} from "@react-admin/ra-form-layout";

const jobListFilters = {
    "metadata.id": textFilter({ operators: ['eq', 'q'] }),
    "status.state": textFilter({ operators: ['eq', 'q'] }),
}

export const JobListActions = () => (<TopToolbar>
    <SelectColumnsButton/>
</TopToolbar>)

export const JobList = () => (
    <ListLive actions={<JobListActions/>}  filters={<StackedFilters config={jobListFilters} />} >
        <DatagridConfigurable rowClick="show" size="small" preferenceKey="images.datagrid">
            <TextField source="id" label="ID"/>
            <TextField source="status.state" label="State"/>
            <TextField source="spec.handler" label="Handler"/>

            <DateField source="metadata.created_at" label="Created At"/>
            <DateField source="metadata.updated_at" label="Updated At"/>
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
