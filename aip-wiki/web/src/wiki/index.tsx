import React from "react";
import {Route} from "react-router";
import {Routes} from "react-router-dom";
import {createTheme, ThemeProvider} from "@mui/material";

import WikiLayout from "./layouts/WikiLayout";
import Home from "./pages/Home";
import Article from "./pages/Article";

const theme = createTheme({
    typography: {
        fontSize: 14,

        h1: { fontSize: '1.8em', },
        h2: { fontSize: '1.5em', },
        h3: { fontSize: '1.2em', },
    },
});

const Wiki = () => (
    <ThemeProvider theme={theme}>
        <Routes>
            <Route path="/" element={<WikiLayout />}>
                <Route index element={<Home />} />

                <Route path="wiki/:slug" element={<Article />} />
            </Route>
        </Routes>
    </ThemeProvider>
)

export default Wiki