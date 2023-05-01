package forddbimpl

import (
	"context"

	"github.com/greenboxal/aip/pkg/ford/forddb"
)

type memoryLogStoreIterator struct {
	ls      *MemoryLogStore
	options forddb.LogIteratorOptions

	current        *forddb.LogEntryRecord
	currentLsn     forddb.LSN
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

func (mls *memoryLogStoreIterator) SetLSN(ctx context.Context, lsn forddb.LSN) error {
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

func (mls *memoryLogStoreIterator) Entry() *forddb.LogEntry {
	if mls.current == nil {
		return nil
	}

	return &mls.current.Entry
}

func (mls *memoryLogStoreIterator) Record() *forddb.LogEntryRecord {
	return mls.current
}

func (mls *memoryLogStoreIterator) CurrentLsn() forddb.LSN {
	return mls.currentLsn
}

func (mls *memoryLogStoreIterator) Reset(lsn forddb.LSN) {
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

	if mls.options.Block {
		mls.ls.cond.L.Lock()
		defer mls.ls.cond.L.Unlock()
	}

	index := mls.currentLsn.Clock - 1
	head := mls.ls.RecordCount()

	mls.current = nil

	if index >= head {
		if mls.options.Block {
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
		mls.current = &mls.ls.records[index]
		mls.currentLsn = mls.current.LSN
	}

	return nil
}
