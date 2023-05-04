package forddbimpl

import (
	"context"

	logstore2 "github.com/greenboxal/aip/aip-controller/pkg/ford/forddb/logstore"
)

type memoryLogStoreIterator struct {
	ls      *MemoryLogStore
	options logstore2.LogIteratorOptions

	current        *logstore2.LogEntryRecord
	currentLsn     logstore2.LSN
	currentSegment *fileStoreSegment

	err error
}

func (mls *memoryLogStoreIterator) Next(ctx context.Context) bool {
	if mls.current == nil {
		if err := mls.invalidate(ctx); err != nil {
			mls.err = err
			return false
		}

		return mls.current != nil
	}

	if err := mls.SeekRelative(ctx, 1); err != nil {
		mls.err = err
		return false
	}

	return true
}

func (mls *memoryLogStoreIterator) Previous(ctx context.Context) bool {
	if mls.current == nil {
		if err := mls.invalidate(ctx); err != nil {
			mls.err = err
			return false
		}

		return mls.current != nil
	}

	if err := mls.SeekRelative(ctx, -1); err != nil {
		mls.err = err
		return false
	}

	return true
}

func (mls *memoryLogStoreIterator) SetLSN(ctx context.Context, lsn logstore2.LSN) error {
	mls.currentLsn = lsn

	return mls.invalidate(ctx)
}

func (mls *memoryLogStoreIterator) SeekRelative(ctx context.Context, offset int64) error {
	clock := mls.currentLsn.Clock + uint64(offset)

	if clock > mls.ls.RecordCount() && !mls.options.Block {
		return nil
	}

	mls.currentLsn.Clock = clock

	return mls.invalidate(ctx)
}

func (mls *memoryLogStoreIterator) Error() error {
	return mls.err
}

func (mls *memoryLogStoreIterator) Entry() *logstore2.LogEntry {
	if mls.current == nil {
		return nil
	}

	return &mls.current.LogEntry
}

func (mls *memoryLogStoreIterator) Record() *logstore2.LogEntryRecord {
	return mls.current
}

func (mls *memoryLogStoreIterator) CurrentLsn() logstore2.LSN {
	return mls.currentLsn
}

func (mls *memoryLogStoreIterator) Reset(lsn logstore2.LSN) {
	mls.current = nil
	mls.currentLsn = lsn
}

func (mls *memoryLogStoreIterator) Close() error {
	return nil
}

func (mls *memoryLogStoreIterator) invalidate(ctx context.Context) error {
	if mls.current != nil && mls.current.LSN.Equals(mls.currentLsn) {
		return nil
	}

	head := mls.ls.RecordCount()

	if mls.currentLsn.Clock == 0 {
		mls.currentLsn.Clock = 1
	} else if mls.currentLsn.Clock == 0xFFFFFFFFFFFFFFFF {
		mls.currentLsn.Clock = head
	}

	index := mls.currentLsn.Clock - 1

	mls.current = nil

	if index >= head {
		if mls.options.Block {
			mls.ls.cond.L.Lock()
			defer mls.ls.cond.L.Unlock()

			for index >= head {
				index = head

				select {
				case <-ctx.Done():
					return ctx.Err()

				default:
				}

				mls.ls.cond.Wait()

				head = mls.ls.RecordCount()
			}
		}
	}

	if index >= 0 && index < head {
		current := mls.ls.records[index]
		mls.current = &current
	}

	return nil
}
