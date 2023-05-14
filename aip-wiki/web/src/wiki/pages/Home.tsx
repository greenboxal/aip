import {Box, Container, IconButton, InputAdornment, Stack, TextField} from "@mui/material";
import SearchIcon from "@mui/icons-material/Search"

const SearchInputBox = () => (
    <TextField fullWidth
        id="search-query"
        label="" variant="outlined"
        InputProps={{
            endAdornment: <InputAdornment position="end">
                <IconButton>
                    <SearchIcon />
                </IconButton>
            </InputAdornment>,
        }}
    />
)

const Logo = () => (
    <Box sx={{
        display: 'flex',
        justifyContent: 'center',
    }}>
        <h1>Home</h1>
    </Box>
)

const Home = () => (
    <Container sx={{
        mt: 8,
        display: 'flex',
        justifyContent: 'center',
    }}>
        <Stack sx={{
            width: 500,
            maxWidth: '100%',
        }}>
            <Logo />
            <SearchInputBox />
        </Stack>
    </Container>
)

export default Home