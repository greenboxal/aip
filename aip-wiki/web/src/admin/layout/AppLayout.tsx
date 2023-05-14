import React from "react";
import {Layout, LayoutProps} from "react-admin";
import {AppLocationContext, Breadcrumb} from "@react-admin/ra-navigation";
import {ReactQueryDevtools} from "react-query/devtools";
import AppBar from "./AppBar";
import AppMenu from "./AppMenu";

export const AppLayout: React.FC<LayoutProps> = ({children, ...rest}) => (<>
    <AppLocationContext>
        <Layout {...rest} menu={AppMenu} appBar={AppBar}>
            <Breadcrumb></Breadcrumb>
            {children}
        </Layout>
    </AppLocationContext>
    <ReactQueryDevtools/>
</>)

export default AppLayout;
