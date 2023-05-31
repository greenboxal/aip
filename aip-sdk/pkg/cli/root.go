package cli

import (
	"github.com/spf13/cobra"
)

type Manager struct {
	Root *cobra.Command
}

func New() *Manager {
	root := &cobra.Command{
		Use: "aipctl",
	}

	return &Manager{
		Root: root,
	}
}

func (m *Manager) Run() error {
	return m.Root.Execute()
}
