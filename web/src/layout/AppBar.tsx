import {AppBar, defaultTheme, LocalesMenuButton, RaThemeOptions, TitlePortal, ToggleThemeButton} from "react-admin";
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
        <TitlePortal />
        <Search />
        <ToggleThemeButton lightTheme={defaultTheme} darkTheme={darkTheme}  />
        <LocalesMenuButton />
    </AppBar>
);

export default MyAppBar;