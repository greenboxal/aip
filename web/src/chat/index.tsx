import React from "react";
import {ChatList, Input, MessageList, MessageType} from "react-chat-elements";
import {Card} from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";
import {Button} from "react-admin";

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

const Messages = () => (
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
        <Grid container spacing={4}>
            <Grid xs={4}>
                <Card>
                    <Topics/>
                </Card>
            </Grid>
            <Grid xs={8}>
                <Card>
                    <Messages/>
                    <MessageInput/>
                </Card>
            </Grid>
        </Grid>
    );
};

export default ChatPage;