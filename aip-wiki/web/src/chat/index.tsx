import React from "react";
import {ChatList, IChatItemProps, Input, MessageList, MessageType} from "react-chat-elements";
import {Button, InfiniteList, InfiniteListBase, ListBase, SimpleList, TextField, useRecordContext} from "react-admin";
import { styled } from '@mui/material/styles';
import clsx from 'clsx';

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

const Messages = () => {
    return (
        <InfiniteList resource="Message" filter={{
            q: 'resource.channel == "1"',
        }} sort={{ field: 'metadata.created_at', order: 'DESC' }}>
            <SimpleList
                primaryText={record => record.from}
                secondaryText={record => record.text}
            />
        </InfiniteList>
    )
}

const Topics = () => (
    <InfiniteList resource="Channel" aside={<Messages />}>
        <SimpleList
            primaryText={record => record.id}
            secondaryText={record => record.members}
        />
    </InfiniteList>
)

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
            <Topics />
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