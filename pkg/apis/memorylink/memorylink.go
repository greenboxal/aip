package memorylink

import (
	"context"

	"go.uber.org/zap"

	"github.com/greenboxal/aip/pkg/ford"
	"github.com/greenboxal/aip/pkg/ford/forddb"
	"github.com/greenboxal/aip/pkg/indexing"
)

type MemoryLink struct {
	logger *zap.SugaredLogger
	db     forddb.Database
	index  indexing.Index
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
	MemoryId indexing.MemoryID `json:"memory_id"`
}

type OneShotGetMemoryResponse struct {
	Memory indexing.Memory `json:"memory"`
}

func (ml *MemoryLink) OneShotGetMemory(ctx context.Context, req *OneShotGetMemoryRequest) (*OneShotGetMemoryResponse, error) {
	tx, err := ml.index.OpenSession(ctx, indexing.SessionOptions{
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
	PreviousMemory *indexing.Memory    `json:"previous_memory"`
	NewMemory      indexing.MemoryData `json:"new_memory"`
}

type OneShotPutMemoryResponse struct {
	NewMemory indexing.Memory `json:"new_memory"`
}

func (ml *MemoryLink) OneShotPutMemory(ctx context.Context, req *OneShotPutMemoryRequest) (*OneShotPutMemoryResponse, error) {
	var previousMemoryId *indexing.MemoryID

	if req.PreviousMemory != nil {
		previousMemoryId = &req.PreviousMemory.ID
	}

	tx, err := ml.index.OpenSession(ctx, indexing.SessionOptions{
		ReadOnly:        true,
		InitialMemoryID: req.PreviousMemory.ID,
	})

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

	return &OneShotPutMemoryResponse{
		NewMemory: newMemory,
	}, nil
}
