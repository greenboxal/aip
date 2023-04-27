package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"

	"github.com/greenboxal/aip/tools/ircproxy/irc"
)

const ChannelName = "#aip-bod-room"

func main() {
	nickname := os.Args[1]
	otherArgs := os.Args[2:]

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	signalCh := make(chan os.Signal, 1)
	waitCh := make(chan struct{})

	cmdArgs := []string{"-m", "aip", "chat"}
	cmdArgs = append(cmdArgs, "--raw")
	cmdArgs = append(cmdArgs, "--ai-identity", nickname)
	cmdArgs = append(cmdArgs, otherArgs...)
	cmd := exec.CommandContext(ctx, "python", cmdArgs...)

	cmd.Stderr = os.Stderr

	stdinPipe, err := cmd.StdinPipe()

	if err != nil {
		panic(err)
	}

	defer stdinPipe.Close()

	stdoutPipe, err := cmd.StdoutPipe()

	if err != nil {
		panic(err)
	}

	defer stdoutPipe.Close()

	stream, err := net.Dial("tcp", "irc.freenode.net:6667")

	if err != nil {
		panic(err)
	}

	messageBuffersMutex := sync.Mutex{}
	messageBuffers := map[string][]string{}

	reader := bufio.NewReader(stdoutPipe)

	client := irc.NewClient(stream, irc.ClientConfig{
		Nick: nickname,
		User: nickname,
		Pass: "",
		Name: nickname,
		Handler: irc.HandlerFunc(func(client *irc.Client, message *irc.Message) {
			if message.Command == "PRIVMSG" {
				messageBuffersMutex.Lock()
				defer messageBuffersMutex.Unlock()

				sender := message.Name

				if sender == "" {
					sender = message.User
				}

				msg := message.Params[1]

				if msg == "### END OF MESSAGE ###" || !strings.HasPrefix(sender, "aip-testbot-") {
					messageBuffer := messageBuffers[sender]
					messageBuffers[sender] = nil
					newMsg := strings.TrimSpace(strings.Join(messageBuffer, "\n"))

					if newMsg == "" {
						newMsg = msg
					}

					line := fmt.Sprintf("%s: %s", sender, newMsg)

					fmt.Printf("%s\n", line)

					line, err = encodeMessage(message.Params[0], line)

					if err != nil {
						panic(err)
					}

					_, err := stdinPipe.Write([]byte(line + "\n"))

					if err != nil {
						panic(err)
					}
				} else {
					messageBuffers[sender] = append(messageBuffers[sender], msg)
				}
			} else if message.Command == "MODE" {
				if message.Params[0] == nickname {
					if waitCh != nil {
						close(waitCh)
						waitCh = nil
					}
				}
			}
		}),
	})

	go func() {
		errCh <- client.RunContext(ctx)
	}()

	go func() {
		_, _ = <-waitCh

		errCh <- cmd.Run()
	}()

	go func() {
		_, _ = <-waitCh

		fmt.Printf("Joining channel...\n")

		err := client.WriteMessage(&irc.Message{
			Command: "JOIN",
			Params:  []string{ChannelName},
		})

		if err != nil {
			errCh <- err
			return
		}

		fmt.Printf("Entering main loop\n")

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			line, _, err := reader.ReadLine()

			if err != nil {
				errCh <- err
				return
			}

			msg, err := decodeMessage(line)

			if err != nil {
				fmt.Printf("Invalid message from stdout: %s\n", err)
				continue
			}

			fmt.Printf("%s\n", line)

			for _, s := range strings.Split(string(msg.Output), "\n") {
				select {
				case <-ctx.Done():
					return
				default:
				}

				err = client.WriteMessage(&irc.Message{
					Command: "PRIVMSG",
					Params: []string{
						msg.Tag,
						s,
					},
				})

				if err != nil {
					errCh <- err
					return
				}
			}

			err = client.WriteMessage(&irc.Message{
				Command: "PRIVMSG",
				Params: []string{
					msg.Tag,
					"### END OF MESSAGE ###",
				},
			})

			if err != nil {
				errCh <- err
				return
			}

		}
	}()

	signal.Notify(signalCh, os.Interrupt)

	for {
		select {
		case <-ctx.Done():
			return
		case <-signalCh:
			cancel()
		case err := <-errCh:
			if err != nil {
				panic(err)
			}
		}
	}
}

type IncomingMessage struct {
	Tag   string `json:"tag"`
	Input string `json:"input"`
}

type OutgoingMessage struct {
	Tag    string `json:"tag"`
	Output string `json:"output"`
}

func decodeMessage(line []byte) (OutgoingMessage, error) {
	var msg OutgoingMessage

	if err := json.Unmarshal(line, &msg); err != nil {
		return OutgoingMessage{}, err
	}

	msg.Output = strings.Replace(msg.Output, "\n", " #LB# ", -1)

	return msg, nil
}

func encodeMessage(tag, input string) (string, error) {
	var msg IncomingMessage

	msg.Tag = tag
	msg.Input = input
	msg.Input = strings.Replace(msg.Input, " #LB# ", "\n", -1)

	data, err := json.Marshal(msg)

	if err != nil {
		return "", err
	}

	return string(data), nil
}