package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"

	"github.com/slack-go/slack"
)

const ChannelName = "#aip-bod-room"

func main() {
	errCh := make(chan error, 1)
	signalCh := make(chan os.Signal, 1)
	waitCh := make(chan struct{})

	signal.Notify(signalCh, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	nickname := os.Args[1]
	otherArgs := os.Args[2:]

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

	botToken := os.Getenv("SLACK_BOT_USER_TOKEN")

	api := slack.New(botToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
	)

	client := api.NewRTM()

	defer client.Disconnect()

	go func() {
		_, _, err := client.ConnectRTMContext(ctx)

		if err != nil {
			errCh <- err
			return
		}

		for ev := range client.IncomingEvents {
			switch evt := ev.Data.(type) {
			case *slack.ConnectingEvent:
				fmt.Println("Connecting to Slack with Socket Mode...")
			case *slack.ConnectionErrorEvent:
				fmt.Println("Connection failed. Retrying later...")
			case *slack.ConnectedEvent:
				fmt.Println("Connected to Slack with Socket Mode.")

				if waitCh != nil {
					close(waitCh)
					waitCh = nil
				}

			case *slack.MessageEvent:
				sender := evt.Username

				if sender == "" {
					sender = evt.User
				}

				if sender == nickname {
					continue
				}

				text := evt.Text

				line := fmt.Sprintf("%s: %s", sender, text)

				fmt.Printf("%s\n", line)

				line, err = encodeMessage(line)

				if err != nil {
					panic(err)
				}

				_, err := stdinPipe.Write([]byte(line + "\n"))

				if err != nil {
					panic(err)
				}
			}
		}
	}()

	go func() {
		reader := bufio.NewReader(stdoutPipe)

		_, _ = <-waitCh

		fmt.Printf("Joining channel...\n")

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

			line, err = decodeMessage(line)

			if err != nil {
				fmt.Printf("Invalid message from stdout: %s\n", err)
				continue
			}

			fmt.Printf("%s\n", line)

			select {
			case <-ctx.Done():
				return
			default:
			}

			_, _, err = client.PostMessage(
				ChannelName,
				slack.MsgOptionText(string(line), true),
				slack.MsgOptionUsername(nickname),
			)

			if err != nil {
				errCh <- err
				return
			}
		}
	}()

	go func() {
		client.ManageConnection()
	}()

	go func() {
		_, _ = <-waitCh

		errCh <- cmd.Run()
	}()

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
	Input string `json:"input"`
}

type OutgoingMessage struct {
	Output string `json:"output"`
}

func decodeMessage(line []byte) ([]byte, error) {
	var msg OutgoingMessage

	if err := json.Unmarshal(line, &msg); err != nil {
		return nil, err
	}

	return []byte(msg.Output), nil
}

func encodeMessage(input string) (string, error) {
	var msg IncomingMessage

	msg.Input = input

	data, err := json.Marshal(msg)

	if err != nil {
		return "", err
	}

	return string(data), nil
}
