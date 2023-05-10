import {AppBar, defaultTheme, RaThemeOptions, TitlePortal, ToggleThemeButton} from "react-admin";
import {Search} from "@react-admin/ra-search";
import React from "react";

const darkTheme = {
    ...defaultTheme,
    palette: {
        mode: 'dark',
    },
} as RaThemeOptions

export const MyAppBar = () => (
    <AppBar>
        <TitlePortal key="title"/>
        <Search key="search"/>
        <ToggleThemeButton key="themes" lightTheme={defaultTheme} darkTheme={darkTheme}/>
    </AppBar>
);

export default MyAppBar;