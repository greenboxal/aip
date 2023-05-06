import React from 'react'
import {Grid} from "@mui/material";

const Item: React.FC<{children: any}> = ({ children }) => (
    <p>{children}</p>
)

const Resources: React.FC = () => (
    <Grid container spacing={2}>
        <Grid item xs={8}>
            <Item>xs=8</Item>
        </Grid>
        <Grid item xs={4}>
            <Item>xs=4</Item>
        </Grid>
        <Grid item xs={4}>
            <Item>xs=4</Item>
        </Grid>
        <Grid item xs={8}>
            <Item>xs=8</Item>
        </Grid>
    </Grid>
)

export default Resources