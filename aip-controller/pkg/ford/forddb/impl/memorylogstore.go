package forddbimpl

import (
	"context"
	"os"
	"sync"
	"time"

	logstore2 "github.com/greenboxal/aip/aip-controller/pkg/ford/forddb/logstore"
)

type MemoryLogStore struct {
	m    sync.RWMutex
	cond *sync.Cond

	clock   uint64
	records []logstore2.LogEntryRecord
}

func NewMemoryLogStore() *MemoryLogStore {
	mls := &MemoryLogStore{}

	mls.cond = sync.NewCond(&sync.RWMutex{})
	mls.records = make([]logstore2.LogEntryRecord, 0, 256)

	return mls
}

func (mls *MemoryLogStore) RecordCount() uint64 {
	mls.m.RLock()
	defer mls.m.RUnlock()

	return uint64(len(mls.records))
}

func (mls *MemoryLogStore) OpenStream(id logstore2.LogStreamID) (logstore2.LogStream, error) {
	fa, err := os.CreateTemp("", "forddb-memorylogstore-*.log")

	if err != nil {
		return nil, err
	}

	fb, err := os.CreateTemp("", "forddb-memorylogstore-*.log")

	if err != nil {
		return nil, err
	}

	return logstore2.NewLogStream(mls, id, fa, fb)
}

func (mls *MemoryLogStore) Append(ctx context.Context, log logstore2.LogEntry) (logstore2.LogEntryRecord, error) {
	mls.m.Lock()
	defer mls.m.Unlock()

	mls.clock = uint64(len(mls.records))

	record := logstore2.LogEntryRecord{
		LSN:      logstore2.MakeLSN(mls.clock, time.Now()),
		LogEntry: log,
	}

	mls.records = append(mls.records, record)

	mls.cond.Broadcast()

	return record, nil
}

func (mls *MemoryLogStore) Iterator(options ...logstore2.LogIteratorOption) logstore2.LogIterator {
	return &memoryLogStoreIterator{ls: mls, options: logstore2.NewLogIteratorOptions(options...)}
}

func (mls *MemoryLogStore) Close() error {
	return nil
}
