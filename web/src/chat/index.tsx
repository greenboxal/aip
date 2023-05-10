import React from "react";
import {ChatList, MessageList, MessageType} from "react-chat-elements";
import {Card, Grid} from "@mui/material";

const chatMessageListRef = React.createRef()

const ChatPage = () => {
    return (
        <Card>
            <Grid container spacing={4}>
                <Grid item xs={4}>
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
                        ]} />
                </Grid>
                <Grid item xs={8}>
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
                        }] as MessageType[]} lockable={true} />
                </Grid>
            </Grid>
        </Card>
    );
};

export default ChatPage;