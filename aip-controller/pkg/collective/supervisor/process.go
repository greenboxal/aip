package supervisor

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/greenboxal/aip/aip-controller/pkg/collective/msn"
)

type Process struct {
	logger *zap.SugaredLogger

	program string
	args    []string

	cmd          *exec.Cmd
	stdinWriter  io.WriteCloser
	stdoutReader io.ReadCloser

	procStartedCh chan struct{}
	outgoingCh    chan msn.Message

	handler func(m msn.Message)
}

func NewProcess(
	logger *zap.SugaredLogger,
	handler func(m msn.Message),
	program string,
	args ...string,
) (*Process, error) {
	return &Process{
		logger:  logger,
		handler: handler,
		program: program,
		args:    args,

		procStartedCh: make(chan struct{}),
		outgoingCh:    make(chan msn.Message, 128),
	}, nil
}

func (p *Process) Send(msg msn.Message) error {
	p.outgoingCh <- msg

	return nil
}

func (p *Process) Run(ctx context.Context) error {
	var err error

	stdinReader, stdinWriter, err := os.Pipe()

	if err != nil {
		return err
	}

	stdoutReader, stdoutWriter, err := os.Pipe()

	if err != nil {
		return err
	}

	p.stdinWriter = stdinWriter
	p.stdoutReader = stdoutReader

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg, gctx := errgroup.WithContext(ctx)

	var env []string

	env = append(env, os.Environ()...)
	env = append(env, "AIP_IPC_BASE_FD=3")

	p.cmd = exec.CommandContext(gctx, p.program, p.args...)

	p.cmd.Env = env
	p.cmd.Stdin = nil
	p.cmd.Stdout = os.Stdout
	p.cmd.Stderr = os.Stdout
	p.cmd.ExtraFiles = []*os.File{stdinReader, stdoutWriter, stdinWriter, stdoutReader}

	if err := p.cmd.Start(); err != nil {
		return err
	}

	wg.Go(func() error {
		for {
			select {
			case <-gctx.Done():
				return gctx.Err()

			case msg := <-p.outgoingCh:
				data, err := json.Marshal(msg)

				if err != nil {
					return err
				}

				_, err = p.stdinWriter.Write([]byte(string(data) + "\n"))

				if err != nil {
					return err
				}
			}
		}
	})

	wg.Go(func() error {
		var msg msn.Message

		reader := bufio.NewReader(p.stdoutReader)

		for {
			select {
			case <-gctx.Done():
				return gctx.Err()
			default:
			}

			line, _, err := reader.ReadLine()

			if err != nil {
				if err == io.EOF {
					return nil
				}

				return err
			}

			if err := json.Unmarshal(line, &msg); err != nil {
				return err
			}

			p.handler(msg)
		}
	})

	wg.Go(func() error {
		close(p.procStartedCh)

		return p.cmd.Wait()
	})

	return wg.Wait()
}

func (p *Process) Close() error {
	return nil
}

func (p *Process) Process() *os.Process {
	_, _ = <-p.procStartedCh

	return p.cmd.Process
}

func (p *Process) Pid() int {
	_, _ = <-p.procStartedCh

	return p.cmd.Process.Pid
}
