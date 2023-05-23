import React from 'react'

import {Create, DatagridConfigurable, DateField, List, SelectColumnsButton, Show, SimpleForm, SimpleShowLayout, TextField, TextInput, TopToolbar, UrlField,} from "react-admin"

export const DomainListActions = () => (<TopToolbar>
    <SelectColumnsButton/>
</TopToolbar>)

export const DomainList = () => (
    <List actions={<DomainListActions/>}>
        <DatagridConfigurable rowClick="show" size="small" preferenceKey="domains.datagrid">
            <TextField source="id" label="ID"/>
        </DatagridConfigurable>
    </List>
)

export const DomainShow = () => (
    <Show>
        <SimpleShowLayout>
            <TextField source="id" label="ID"/>

            <DateField source="metadata.created_at" label="Created At"/>
            <DateField source="metadata.updated_at" label="Updated At"/>
        </SimpleShowLayout>
    </Show>
);

export const DomainCreate = () => (
    <Create>
        <SimpleForm>
            <TextInput name="metadata.id" source="metadata.id" label="Domain" />
        </SimpleForm>
    </Create>
);
