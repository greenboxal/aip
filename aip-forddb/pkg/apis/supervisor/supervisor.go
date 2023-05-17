package supervisor

import (
	"context"

	"github.com/samber/lo"

	supervisor2 "github.com/greenboxal/aip/aip-controller/pkg/collective/supervisor"
	"github.com/greenboxal/aip/aip-controller/pkg/daemon"
)

type ListChildrenRequest struct {
	Empty bool `json:"empty"`
}

type ListChildrenResponse struct {
	Children []string `json:"children"`
}

type StartChildRequest struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
}

type StartChildResponse struct {
}

type SupervisorAPI struct {
	daemon  *daemon.Daemon
	manager *supervisor2.Manager
}

func NewSupervisorApi(d *daemon.Daemon, m *supervisor2.Manager) *SupervisorAPI {
	return &SupervisorAPI{daemon: d, manager: m}
}

func (api *SupervisorAPI) StartChild(ctx context.Context, req *StartChildRequest) (*StartChildResponse, error) {
	err := api.daemon.StartSupervised(req.Name, req.Args...)

	if err != nil {
		return nil, err
	}

	return &StartChildResponse{}, nil
}

func (api *SupervisorAPI) StopChild(name string) error {
	children := api.manager.Child(name)

	if children == nil {
		return nil
	}

	return children.Close()
}

func (api *SupervisorAPI) ListChildren(ctx context.Context, req *ListChildrenRequest) (*ListChildrenResponse, error) {
	children := api.manager.Children()

	names := lo.Map(children, func(child *supervisor2.Supervisor, _index int) string {
		return child.Config.ID
	})

	return &ListChildrenResponse{
		Children: names,
	}, nil
}
