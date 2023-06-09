package logstore

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
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
	mls.records = make([]forddb.LogEntryRecord, 0, 256)

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

	checkPoint, err := forddb.NewFileCheckpoint(id, fa, fb)

	if err != nil {
		return nil, err
	}

	return forddb.NewLogStream(mls, id, checkPoint)
}

func (mls *MemoryLogStore) Append(ctx context.Context, log forddb.LogEntry) (forddb.LogEntryRecord, error) {
	log = forddb.Clone(log)

	mls.m.Lock()
	defer mls.m.Unlock()

	mls.clock = uint64(len(mls.records))

	record := forddb.LogEntryRecord{
		LSN:      forddb.MakeLSN(mls.clock, time.Now()),
		LogEntry: log,
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
