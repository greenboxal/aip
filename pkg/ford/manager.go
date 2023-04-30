package ford

import (
	"go.uber.org/zap"

	"github.com/greenboxal/aip/pkg/indexing"
	"github.com/greenboxal/aip/pkg/indexing/impl"
)

type Manager struct {
	logger *zap.SugaredLogger
	index  *impl.Index
}

func NewManager(
	logger *zap.SugaredLogger,
	storage indexing.Storage,
) (*Manager, error) {
	m := &Manager{
		logger: logger.Named("ford"),
	}

	m.index = impl.NewIndex(storage, indexing.IndexConfiguration{})

	return m, nil
}

func (m *Manager) Index() indexing.Index {
	return m.index
}
