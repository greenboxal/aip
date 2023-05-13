import {useCreatePath, useGetResourceLabel, useResourceDefinitions} from "react-admin";
import {MultiLevelMenu} from "@react-admin/ra-navigation";
import React, {createElement} from "react";
import DefaultIcon from "@mui/icons-material/ViewList";

export const ResourceMenuItem = ({name}: { name: string }) => {
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
                    <DefaultIcon/>
                )
            }
        />
    );
};
