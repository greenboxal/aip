package forddbimpl

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/greenboxal/aip/pkg/ford/forddb"
)

type MemoryLogStore struct {
	m    sync.RWMutex
	cond *sync.Cond

	clock   uint64
	records []forddb.LogEntryRecord
}

func NewMemoryLogStore() *MemoryLogStore {
	mls := &MemoryLogStore{}

	mls.cond = sync.NewCond(&sync.RWMutex{})
	mls.records = make([]forddb.LogEntryRecord, 256)

	return mls
}

func (mls *MemoryLogStore) RecordCount() uint64 {
	mls.m.RLock()
	defer mls.m.RUnlock()

	return uint64(len(mls.records))
}

func (mls *MemoryLogStore) OpenStream(id forddb.LogStreamID) (forddb.LogStream, error) {
	fa, err := os.CreateTemp("", "forddb-memorylogstore-*.log")

	if err != nil {
		return nil, err
	}

	fb, err := os.CreateTemp("", "forddb-memorylogstore-*.log")

	if err != nil {
		return nil, err
	}

	return forddb.NewLogStream(mls, id, fa, fb)
}

func (mls *MemoryLogStore) Append(ctx context.Context, log forddb.LogEntry) (forddb.LogEntryRecord, error) {
	mls.m.Lock()
	defer mls.m.Unlock()

	mls.clock++

	record := forddb.LogEntryRecord{
		LSN:   forddb.MakeLSN(mls.clock, time.Now()),
		Entry: log,
	}

	mls.records = append(mls.records, record)

	mls.cond.Broadcast()

	return record, nil
}

func (mls *MemoryLogStore) Iterator(options ...forddb.LogIteratorOption) forddb.LogIterator {
	return &memoryLogStoreIterator{ls: mls, options: forddb.NewLogIteratorOptions(options...)}
}

func (mls *MemoryLogStore) Close() error {
	return nil
}
