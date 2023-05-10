import React from "react";
import {MenuProps} from "react-admin";
import {MultiLevelMenu} from "@react-admin/ra-navigation";
import {ResourceMenuItem} from "./ResourceMenuItem";

export const AppMenu: React.FC<MenuProps> = (props: MenuProps) => (
    <MultiLevelMenu {...props}>
        <MultiLevelMenu.Item label="Content Management" name="content-management">
            <ResourceMenuItem name="Page"/>
            <ResourceMenuItem name="Image"/>
        </MultiLevelMenu.Item>

        <MultiLevelMenu.Item label="Messaging" name="messaging">
            <ResourceMenuItem name="Channel"/>
            <ResourceMenuItem name="Message"/>
        </MultiLevelMenu.Item>

        <MultiLevelMenu.Item label="Infrastructure" name="infrastructure">
            <ResourceMenuItem name="Job"/>
            <ResourceMenuItem name="Memory"/>
        </MultiLevelMenu.Item>

    </MultiLevelMenu>
)

export default AppMenu;
