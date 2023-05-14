import React from "react";
import {MenuProps} from "react-admin";
import {MultiLevelMenu} from "@react-admin/ra-navigation";
import {ResourceMenuItem} from "./ResourceMenuItem";

export const AppMenu: React.FC<MenuProps> = (props: MenuProps) => (
    <MultiLevelMenu {...props}>
        <MultiLevelMenu.Item label="Content Management" name="content-management">
            <ResourceMenuItem name="Page"/>
            <ResourceMenuItem name="Image"/>
            <ResourceMenuItem name="Domain"/>
            <ResourceMenuItem name="Layout"/>
        </MultiLevelMenu.Item>

        <MultiLevelMenu.Item label="Messaging" name="messaging">
            <ResourceMenuItem name="Channel" />
            <ResourceMenuItem name="Endpoint" />
            <ResourceMenuItem name="Message" />
        </MultiLevelMenu.Item>

        <MultiLevelMenu.Item label="Production Pipeline" name="production-pipeline">
            <ResourceMenuItem name="Agent"/>
            <ResourceMenuItem name="Profile"/>
            <ResourceMenuItem name="Task"/>
            <ResourceMenuItem name="Pipeline"/>
            <ResourceMenuItem name="Team"/>
        </MultiLevelMenu.Item>

        <MultiLevelMenu.Item label="Infrastructure" name="infrastructure">
            <ResourceMenuItem name="Job"/>
            <ResourceMenuItem name="Memory"/>
        </MultiLevelMenu.Item>

        <MultiLevelMenu.Item label="Tracing" name="tracing">
            <ResourceMenuItem name="Trace"/>
            <ResourceMenuItem name="Span"/>
        </MultiLevelMenu.Item>
    </MultiLevelMenu>
)

export default AppMenu;
