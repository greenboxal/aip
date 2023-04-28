package supervisor

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"os/exec"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/greenboxal/aip/pkg/collective"
)

type Process struct {
	logger *zap.SugaredLogger

	program string
	args    []string

	stdinWriter  *os.File
	stdinReader  *os.File
	stdoutWriter *os.File
	stdoutReader *os.File

	handler func(m collective.Message)
}

func NewProcess(
	logger *zap.SugaredLogger,
	handler func(m collective.Message),
	program string,
	args ...string,
) (*Process, error) {
	stdinWriter, stdinReader, err := os.Pipe()

	if err != nil {
		return nil, err
	}

	stdoutWriter, stdoutReader, err := os.Pipe()

	if err != nil {
		return nil, err
	}

	return &Process{
		logger:  logger,
		handler: handler,
		program: program,
		args:    args,

		stdinWriter:  stdinWriter,
		stdinReader:  stdinReader,
		stdoutWriter: stdoutWriter,
		stdoutReader: stdoutReader,
	}, nil
}

func (p *Process) SendAndReceive(msg collective.Message) error {
	data, err := json.Marshal(msg)

	if err != nil {
		return err
	}

	_, err = p.stdinWriter.Write([]byte(string(data) + "\n"))

	if err != nil {
		return err
	}

	return nil
}

func (p *Process) Run(ctx context.Context) error {
	var err error

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg, gctx := errgroup.WithContext(ctx)

	wg.Go(func() error {
		var env []string

		env = append(env, os.Environ()...)
		env = append(env, "AIP_TRANSPORT=ipc")
		env = append(env, "AIP_IPC_BASE_FD=3")

		cmd := exec.CommandContext(gctx, p.program, p.args...)

		cmd.Env = env
		cmd.Stdin = nil
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.ExtraFiles = []*os.File{p.stdinReader, p.stdoutWriter}

		return cmd.Run()
	})

	wg.Go(func() error {
		reader := bufio.NewReader(p.stdoutReader)

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

			var msg collective.Message

			if err := json.Unmarshal(line, &msg); err != nil {
				return err
			}

			p.handler(msg)
		}
	})

	return err
}

func (p *Process) Close() error {
	p.stdoutWriter.Close()
	p.stdoutReader.Close()
	p.stdinReader.Close()
	p.stdinWriter.Close()

	return nil
}
