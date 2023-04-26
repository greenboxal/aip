package daemon

import (
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"

	"github.com/greenboxal/aip/pkg/supervisor"
)

type Daemon struct {
	routing *Routing
	manager *supervisor.Manager
}

func NewDaemon(routing *Routing, manager *supervisor.Manager) *Daemon {
	return &Daemon{
		routing: routing,
		manager: manager,
	}
}

func (d *Daemon) StartSupervised(name string, extraArgs ...string) error {
	port, err := d.routing.AddPort(name)

	if err != nil {
		return err
	}

	if err := port.Subscribe(name); err != nil {
		return err
	}

	if err := port.Subscribe("aip-bod-room"); err != nil {
		return err
	}

	args := []string{
		"bash",
		"-c",
		`set -eu; if [ -f .env ]; then source .env; fi; exec python -m aip chat $*`,
		"--",
		"--raw",
		"-i", name,
	}

	args = append(args, extraArgs...)

	cfg := supervisor.Config{
		ID:      name,
		Program: "/usr/bin/env",
		Args:    args,
	}

	_, err = d.manager.Spawn(&cfg, port)

	if err != nil {
		return err
	}

	return nil
}

func (d *Daemon) Run() error {
	for _, str := range os.Args[1:] {
		var name, profile string

		components := strings.SplitN(str, ":", 2)

		if len(components) == 1 {
			var header profileHeader

			data, err := os.ReadFile(components[0])

			if err != nil {
				return err
			}

			if err := toml.Unmarshal(data, &header); err != nil {
				return err
			}

			name = header.Name
			profile = components[0]
		} else {
			name = components[0]
			profile = components[1]
		}

		if err := d.StartSupervised(name, "-p", profile); err != nil {
			return err
		}
	}

	return nil
}

type profileHeader struct {
	Name string `toml:"name"`
}
