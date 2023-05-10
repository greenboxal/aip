import React from "react";
import {ChatList, Input, MessageList, MessageType} from "react-chat-elements";
import {Container, Card, Stack} from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";
import {Button} from "react-admin";
import { styled } from '@mui/material/styles';
import clsx from 'clsx';
import {useSubscription} from "@apollo/client";

const chatMessageListRef = React.createRef()
const inputReference = React.createRef()

/*
const GET_MESSAGES = gql`
subscription realTimeEvents($endpoint: String!) {
    realTimeEvents(endpoint: $endpoint) {
        message_event {
            message {
                id
                channel
            }
        }
    }
}
`*/

const Topics = () => (
    <ChatList
        className="chat-list"
        lazyLoadingImage=""
        id="lol"
        dataSource={[
            {
                id: 1,
                avatar: 'https://avatars.githubusercontent.com/u/80540635?v=4',
                alt: 'kursat_avatar',
                title: 'Kursat',
                subtitle: "Why don't we go to the No Way Home movie this weekend ?",
                date: new Date(),
                unread: 3,
            }, {
                id: 1,
                avatar: 'https://avatars.githubusercontent.com/u/80540635?v=4',
                alt: 'kursat_avatar',
                title: 'Kursat',
                subtitle: "Why don't we go to the No Way Home movie this weekend ?",
                date: new Date(),
                unread: 3,
            }, {
                id: 1,
                avatar: 'https://avatars.githubusercontent.com/u/80540635?v=4',
                alt: 'kursat_avatar',
                title: 'Kursat',
                subtitle: "Why don't we go to the No Way Home movie this weekend ?",
                date: new Date(),
                unread: 3,
            }
        ]}/>
)


const Messages = () => {
    const { data } = useSubscription(GET_MESSAGES);

    if (!data) {
        return null;
    }

    const messageList = data

    return (
        <MessageList referance={chatMessageListRef} dataSource={[
            {
                id: 1,
                position: 'right',
                avatar: 'https://avatars.githubusercontent.com/u/80540635?v=4',
                type: 'text',
                title: "Hello",
                text: 'Lorem ipsum dolor sit amet, consectetur adipisicing elit',
                date: new Date(),
            }, {
                id: 1,
                position: 'left',
                avatar: 'https://avatars.githubusercontent.com/u/80540635?v=4',
                type: 'text',
                title: "Hello",
                text: 'Lorem ipsum dolor sit amet, consectetur adipisicing elit',
                date: new Date(),
            }] as MessageType[]} lockable={true}/>
    )
}

const MessageInput = () => (
    <Input
        referance={inputReference}
        placeholder='Type here...'
        multiline={true}
        maxHeight={200}
        rightButtons={<Button color='primary' label='Send'/>}
    />
)

const ChatPage = () => {
    return (
        <Root className={clsx('chat-page')} >
            <Grid container spacing={0}>
                <Grid xs={4}>
                    <Card>
                        <Topics/>
                    </Card>
                </Grid>
                <Grid xs={8}>
                    <Stack>
                        <Messages/>
                        <MessageInput/>
                    </Stack>
                </Grid>
            </Grid>
        </Root>
    );
};

const PREFIX = 'Chat';

export const ChatClasses = {
    main: `${PREFIX}-main`,
    noActions: `${PREFIX}-noActions`,
    card: `${PREFIX}-card`,
};

const Root = styled('div', {
    name: PREFIX,
    overridesResolver: (props, styles) => styles.root,
})(({ theme }) => ({
    [`& .${ChatClasses.main}`]: {
        display: 'flex',
    },
    [`& .${ChatClasses.noActions}`]: {
        marginTop: '1em',
    },
    [`& .${ChatClasses.card}`]: {
        flex: '1 1 auto',
    },
}));

export default ChatPage;