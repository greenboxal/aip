import React from 'react'

import {DatagridConfigurable, DateField, ImageField, List, SelectColumnsButton, Show, SimpleShowLayout, TextField, TopToolbar, UrlField,} from "react-admin"

export const ImageListActions = () => (<TopToolbar>
    <SelectColumnsButton/>
</TopToolbar>)

export const ImageList = () => (
    <List actions={<ImageListActions/>}>
        <DatagridConfigurable rowClick="show" size="small" preferenceKey="images.datagrid">
            <ImageField source="status.url" title="Thumbnail" label=""/>
            <TextField source="id" label="ID"/>
            <TextField source="status.prompt" label="Prompt"/>

            <TextField source="spec.path" label="Path"/>
        </DatagridConfigurable>
    </List>
)

export const ImageShow = () => (
    <Show>
        <SimpleShowLayout>
            <TextField source="id" label="ID"/>

            <DateField source="metadata.created_at" label="Created At"/>
            <DateField source="metadata.updated_at" label="Updated At"/>

            <TextField source="spec.path" label="Path"/>
            <TextField source="status.prompt" label="Prompt"/>
            <UrlField source="status.url" label="URL"/>

            <ImageField source="status.url" title="Preview"/>
        </SimpleShowLayout>
    </Show>
);
