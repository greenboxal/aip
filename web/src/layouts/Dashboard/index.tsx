import React from 'react';

import {
    AppBar,
    Container,
    IconButton,
    Toolbar,
    Typography,
    Box,
    Drawer,
    List,
    ListItem,
    ListItemButton,
    ListItemIcon,
    ListItemText, Divider
} from "@mui/material";

type Props = {
    children?: React.ReactNode
};


const DashboardLayout: React.FC<Props> = ({ children }) => (
    <Box sx={{ display: 'flex' }}>
        {/*<AppBar position="static">
            <Toolbar variant="dense">
                <Typography variant="h6" color="inherit" component="div">
                    AIP
                </Typography>
                <IconButton edge="start" color="inherit" aria-label="menu" sx={{ mr: 2 }}>
                    Resources
                </IconButton>
            </Toolbar>
        </AppBar>
        <Drawer variant="permanent">
            <Toolbar />
            <Box sx={{ overflow: 'auto' }}>
                <List>
                    {['Inbox', 'Starred', 'Send email', 'Drafts'].map((text, index) => (
                        <ListItem key={text} disablePadding>
                            <ListItemButton>
                                <ListItemIcon>
                                    Hello
                                </ListItemIcon>
                                <ListItemText primary={text} />
                            </ListItemButton>
                        </ListItem>
                    ))}
                </List>
                <Divider />
                <List>
                    {['All mail', 'Trash', 'Spam'].map((text, index) => (
                        <ListItem key={text} disablePadding>
                            <ListItemButton>
                                <ListItemIcon>
                                    Hi
                                </ListItemIcon>
                                <ListItemText primary={text} />
                            </ListItemButton>
                        </ListItem>
                    ))}
                </List>
            </Box>
        </Drawer>*/}
        <main>
            <Container>
                { children }
            </Container>
        </main>
    </Box>
)

export default DashboardLayout