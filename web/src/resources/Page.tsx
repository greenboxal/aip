import React from 'react'

import {
    ArrayField,
    ChipField,
    Create,
    CreateButton,
    Datagrid,
    DatagridConfigurable,
    DateField,
    List,
    ReferenceField,
    RichTextField,
    SelectColumnsButton,
    Show,
    SimpleForm,
    TabbedShowLayout,
    TextField,
    TextInput,
    TopToolbar,
    WithRecord,
} from "react-admin"

import {CreateInDialogButton} from "@react-admin/ra-form-layout";

import {MarkdownField} from "@react-admin/ra-markdown"

export const PageListActions = () => (<TopToolbar>
    <SelectColumnsButton/>
    <CreateButton/>
</TopToolbar>)

export const PageList = () => (
    <List actions={<PageListActions/>}>
        <DatagridConfigurable rowClick="show" size="small" preferenceKey="pages.datagrid">
            <TextField source="spec.title" label="Title"/>
            <TextField source="spec.language" label="Language"/>
            <TextField source="spec.voice" label="Voice"/>

            <TextField source="id" label="ID"/>
        </DatagridConfigurable>
    </List>
)

export const PageShow = () => (
    <Show>
        <TabbedShowLayout>
            <TabbedShowLayout.Tab label="General">
                <WithRecord render={record => (
                    <CreateInDialogButton record={{
                        base_page_id: record.id,
                        title: record.spec.title,
                        voice: record.spec.voice,
                        language: record.spec.language,
                    }} redirect="show" label="Request Edit">
                        <SimpleForm defaultValues={{
                            title: "",
                            voice: "",
                            language: "",
                        }}>
                            <TextInput source="title" label="Title"/>
                            <TextInput source="voice" label="Voice"/>
                            <TextInput source="language" label="Language"/>
                        </SimpleForm>
                    </CreateInDialogButton>
                )}/>

                <TextField source="id" label="ID"/>
                <TextField source="spec.base_page_id" label="Base Page ID"/>

                <ReferenceField source="spec.base_page_id" reference="Page" label="Base Page">
                    <ChipField source="spec.title"/>
                </ReferenceField>

                <DateField source="metadata.created_at" label="Created At"/>
                <DateField source="metadata.updated_at" label="Updated At"/>

                <TextField source="spec.title" label="Title"/>
                <TextField source="spec.language" label="Language"/>
                <TextField source="spec.voice" label="Voice"/>

                <ArrayField source="status.images" label="Images">
                    <Datagrid bulkActionButtons={false}>
                        <TextField source="title" label="Title"/>
                        <TextField source="source" label="Source"/>
                    </Datagrid>
                </ArrayField>

                <ArrayField source="status.links" label="Links">
                    <Datagrid bulkActionButtons={false}>
                        <TextField source="title" label="Title"/>
                        <TextField source="to" label="To"/>
                    </Datagrid>
                </ArrayField>


                <MarkdownField source="status.markdown" label="Markdown Contents"/>
            </TabbedShowLayout.Tab>

            <TabbedShowLayout.Tab label="HTML Preview">
                <RichTextField source="status.html" label="HTML Contents"/>
            </TabbedShowLayout.Tab>
        </TabbedShowLayout>
    </Show>
);

export const PageCreate = () => (
    <Create>
        <SimpleForm>
            <TextInput source="title" label="Title"/>
            <TextInput source="voice" label="Voice"/>
            <TextInput source="language" label="Language"/>
        </SimpleForm>
    </Create>
);
