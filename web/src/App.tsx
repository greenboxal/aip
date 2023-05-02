import React from 'react';
import { styled, createTheme, ThemeProvider } from '@mui/material/styles';
import { ApolloClient, InMemoryCache, ApolloProvider, gql } from '@apollo/client';
import {AppBar, Container, IconButton, Toolbar, Typography, Box} from "@mui/material";
import { ResourcesPage } from './pages/resources';

const mdTheme = createTheme();

const client = new ApolloClient({
    uri: 'http://127.0.0.1:30100/v1/graphql',
    cache: new InMemoryCache(),
});


function App() {
    return (
        <ApolloProvider client={client}>
            <ThemeProvider theme={mdTheme}>
                <Box sx={{ display: 'flex' }}>
                    <AppBar position="static">
                        <Toolbar variant="dense">
                            <Typography variant="h6" color="inherit" component="div">
                                Photos
                            </Typography>
                        </Toolbar>
                    </AppBar>
                    <Container>
                        <ResourcesPage />
                    </Container>
                </Box>
            </ThemeProvider>
        </ApolloProvider>
    );
}

export default App;
