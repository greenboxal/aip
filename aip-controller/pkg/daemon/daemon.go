package daemon

import (
	"github.com/greenboxal/aip/aip-controller/pkg/collective/comms"
	"github.com/greenboxal/aip/aip-controller/pkg/collective/supervisor"
)

type Daemon struct {
	routing *comms.Routing
	manager *supervisor.Manager
}

func NewDaemon(routing *comms.Routing, manager *supervisor.Manager) *Daemon {
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
		"-i", name,
	}

	args = append(args, extraArgs...)

	_, err = d.manager.Spawn(
		supervisor.WithID(name),
		supervisor.WithArgs(args...),
		supervisor.WithPort(port),
	)

	if err != nil {
		return err
	}

	return nil
}
