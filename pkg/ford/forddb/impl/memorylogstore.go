package forddbimpl

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/greenboxal/aip/pkg/ford/forddb/logstore"
)

type MemoryLogStore struct {
	m    sync.RWMutex
	cond *sync.Cond

	clock   uint64
	records []logstore.LogEntryRecord
}

func NewMemoryLogStore() *MemoryLogStore {
	mls := &MemoryLogStore{}

	mls.cond = sync.NewCond(&sync.RWMutex{})
	mls.records = make([]logstore.LogEntryRecord, 0, 256)

	return mls
}

func (mls *MemoryLogStore) RecordCount() uint64 {
	mls.m.RLock()
	defer mls.m.RUnlock()

	return uint64(len(mls.records))
}

func (mls *MemoryLogStore) OpenStream(id logstore.LogStreamID) (logstore.LogStream, error) {
	fa, err := os.CreateTemp("", "forddb-memorylogstore-*.log")

	if err != nil {
		return nil, err
	}

	fb, err := os.CreateTemp("", "forddb-memorylogstore-*.log")

	if err != nil {
		return nil, err
	}

	return logstore.NewLogStream(mls, id, fa, fb)
}

func (mls *MemoryLogStore) Append(ctx context.Context, log logstore.LogEntry) (logstore.LogEntryRecord, error) {
	mls.m.Lock()
	defer mls.m.Unlock()

	mls.clock = uint64(len(mls.records))

	record := logstore.LogEntryRecord{
		LSN:      logstore.MakeLSN(mls.clock, time.Now()),
		LogEntry: log,
	}

	mls.records = append(mls.records, record)

	mls.cond.Broadcast()

	return record, nil
}

func (mls *MemoryLogStore) Iterator(options ...logstore.LogIteratorOption) logstore.LogIterator {
	return &memoryLogStoreIterator{ls: mls, options: logstore.NewLogIteratorOptions(options...)}
}

func (mls *MemoryLogStore) Close() error {
	return nil
}
