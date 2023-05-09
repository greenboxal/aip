import React, {createElement} from "react";
import {MenuProps, useCreatePath, useGetResourceLabel, useResourceDefinitions} from "react-admin";
import {MultiLevelMenu} from "@react-admin/ra-navigation";
import DefaultIcon from "@mui/icons-material/ViewList";

export const ResourceMenuItem = ({ name }: { name: string }) => {
    const resources = useResourceDefinitions();
    const getResourceLabel = useGetResourceLabel();
    const createPath = useCreatePath();
    if (!resources || !resources[name]) return null;

    return (
        <MultiLevelMenu.Item
            name={name}
            to={createPath({
                resource: name,
                type: 'list',
            })}
            label={<>{getResourceLabel(name, 2)}</>}
            icon={
                resources[name].icon ? (
                    createElement(resources[name].icon)
                ) : (
                    <DefaultIcon />
                )
            }
        />
    );
};

export const AppMenu: React.FC<MenuProps> = (props: MenuProps) => (
    <MultiLevelMenu {...props}>
        <MultiLevelMenu.Item label="Content Management" name="content-management">
            <ResourceMenuItem name="Page" />
            <ResourceMenuItem name="Image" />
        </MultiLevelMenu.Item>

        <MultiLevelMenu.Item label="Infrastructure" name="infrastructure">
            <ResourceMenuItem name="Job" />
            <ResourceMenuItem name="Memory" />
        </MultiLevelMenu.Item>
    </MultiLevelMenu>
)

export default AppMenu;
