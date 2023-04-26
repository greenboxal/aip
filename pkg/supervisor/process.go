package supervisor

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"golang.org/x/sync/errgroup"

	"github.com/greenboxal/aip/pkg/collective"
)

type Process struct {
	ctx    context.Context
	cancel context.CancelFunc

	incomingCh chan collective.Message
	outgoingCh chan collective.Message

	cmd *exec.Cmd

	stdoutPipe io.ReadCloser
	stdinPipe  io.WriteCloser
}

func NewProcess(ctx context.Context, program string, args ...string) (*Process, error) {
	ctx, cancel := context.WithCancel(ctx)

	cmd := exec.CommandContext(ctx, program, args...)

	cmd.Stderr = os.Stderr

	stdinPipe, err := cmd.StdinPipe()

	if err != nil {
		panic(err)
	}

	stdoutPipe, err := cmd.StdoutPipe()

	if err != nil {
		panic(err)
	}

	return &Process{
		ctx:    ctx,
		cancel: cancel,

		cmd:        cmd,
		stdinPipe:  stdinPipe,
		stdoutPipe: stdoutPipe,

		incomingCh: make(chan collective.Message),
		outgoingCh: make(chan collective.Message),
	}, nil
}

func (s *Process) Run() error {
	wg, gctx := errgroup.WithContext(s.ctx)

	wg.Go(func() error {
		return s.cmd.Run()
	})

	wg.Go(func() error {
		reader := bufio.NewReader(s.stdoutPipe)

		for {
			select {
			case <-gctx.Done():
				return gctx.Err()
			default:
			}

			line, _, err := reader.ReadLine()

			if err != nil {
				return err
			}

			msg, err := DecodeMessage(line)

			if err != nil {
				fmt.Printf("Invalid message from stdout: %s\n", err)
				continue
			}

			s.outgoingCh <- msg
		}
	})

	wg.Go(func() error {
		for {
			select {
			case <-gctx.Done():
				return gctx.Err()

			case msg := <-s.incomingCh:
				data, err := EncodeMessage(msg)

				if err != nil {
					panic(err)
				}

				_, err = s.stdinPipe.Write([]byte(string(data) + "\n"))

				if err != nil {
					return err
				}
			}
		}
	})

	return wg.Wait()
}

func (s *Process) Close() error {
	s.cancel()

	_ = s.stdinPipe.Close()
	_ = s.stdoutPipe.Close()

	return nil
}

func (s *Process) Incoming() chan<- collective.Message {
	return s.incomingCh
}

func (s *Process) Outgoing() <-chan collective.Message {
	return s.outgoingCh
}
