package memorylink

import (
	"context"

	"go.uber.org/zap"

	collective2 "github.com/greenboxal/aip/aip-controller/pkg/collective"
	"github.com/greenboxal/aip/aip-controller/pkg/ford"
	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
	indexing2 "github.com/greenboxal/aip/aip-controller/pkg/indexing2"
)

type MemoryLink struct {
	logger *zap.SugaredLogger
	db     forddb.Database
	index  indexing2.Index
}

func NewMemoryLink(
	logger *zap.SugaredLogger,
	db forddb.Database,
	ford *ford.Manager,
) *MemoryLink {
	return &MemoryLink{
		logger: logger.Named("memorylink"),
		db:     db,
		index:  ford.Index(),
	}
}

type OneShotGetMemoryRequest struct {
	MemoryId collective2.MemoryID `json:"memory_id"`
}

type OneShotGetMemoryResponse struct {
	Memory collective2.Memory `json:"memory"`
}

func (ml *MemoryLink) OneShotGetMemory(ctx context.Context, req *OneShotGetMemoryRequest) (*OneShotGetMemoryResponse, error) {
	tx, err := ml.index.OpenSession(ctx, indexing2.SessionOptions{
		ReadOnly:        true,
		InitialMemoryID: req.MemoryId,
	})

	if err != nil {
		return nil, err
	}

	defer tx.Discard()

	if err := tx.SeekTo(req.MemoryId); err != nil {
		return nil, err
	}

	mem := tx.GetCurrentMemory()

	return &OneShotGetMemoryResponse{
		Memory: mem,
	}, nil
}

type OneShotPutMemoryRequest struct {
	PreviousMemory *collective2.Memory    `json:"previous_memory"`
	NewMemory      collective2.MemoryData `json:"new_memory"`
}

type OneShotPutMemoryResponse struct {
	NewMemory collective2.Memory `json:"new_memory"`
}

func (ml *MemoryLink) OneShotPutMemory(ctx context.Context, req *OneShotPutMemoryRequest) (*OneShotPutMemoryResponse, error) {
	var previousMemoryId *collective2.MemoryID

	if req.PreviousMemory != nil {
		previousMemoryId = &req.PreviousMemory.ID
	}

	opts := indexing2.SessionOptions{
		ReadOnly: false,
	}

	if previousMemoryId != nil {
		opts.InitialMemoryID = *previousMemoryId
	}

	tx, err := ml.index.OpenSession(ctx, opts)

	if err != nil {
		return nil, err
	}

	defer tx.Discard()

	if previousMemoryId != nil {
		if err := tx.SeekTo(req.PreviousMemory.ID); err != nil {
			return nil, err
		}
	}

	newMemory, err := tx.Push(req.NewMemory)

	if err != nil {
		return nil, err
	}

	if err := tx.Merge(); err != nil {
		return nil, err
	}

	return &OneShotPutMemoryResponse{
		NewMemory: newMemory,
	}, nil
}
