package logstore

import (
	"context"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
)

type oldFileLogIterator struct {
	ls      *OldFileLogStore
	options forddb.LogIteratorOptions

	current        *forddb.LogEntryRecord
	currentLsn     forddb.LSN
	currentSegment *fileStoreSegment

	err error
}

func newOldFileLogIterator(ls *OldFileLogStore, options forddb.LogIteratorOptions) *oldFileLogIterator {
	return &oldFileLogIterator{
		ls:      ls,
		options: options,
	}
}

func (l *oldFileLogIterator) Error() error {
	return l.err
}

func (l *oldFileLogIterator) Entry() *forddb.LogEntry {
	return &l.current.LogEntry
}

func (l *oldFileLogIterator) Record() *forddb.LogEntryRecord {
	return l.current
}

func (l *oldFileLogIterator) CurrentLsn() forddb.LSN {
	return l.currentLsn
}

func (l *oldFileLogIterator) SetLSN(ctx context.Context, lsn forddb.LSN) error {
	if lsn.Clock == FileSegmentBaseSeekToHead {
		lsn = l.ls.currentSegment.tail
	}

	l.currentLsn = lsn

	return l.invalidate(ctx)
}

func (l *oldFileLogIterator) SeekRelative(ctx context.Context, offset int64) error {
	l.currentLsn.Clock += uint64(offset)

	return l.invalidate(ctx)
}

func (l *oldFileLogIterator) Next(ctx context.Context) bool {
	if err := l.SeekRelative(ctx, 1); err != nil {
		l.err = err
		return false
	}

	return true
}

func (l *oldFileLogIterator) Previous(ctx context.Context) bool {
	if err := l.SeekRelative(ctx, -1); err != nil {
		l.err = err
		return false
	}

	return true
}

func (l *oldFileLogIterator) Reset(lsn forddb.LSN) {
	l.currentLsn = lsn
	l.current = nil
}

func (l *oldFileLogIterator) Close() error {
	return nil
}

func (l *oldFileLogIterator) invalidate(ctx context.Context) error {
	if l.current != nil && l.current.LSN.Equals(l.currentLsn) {
		return nil
	}

	if l.currentSegment == nil {
		if l.currentSegment != nil {
			if err := l.currentSegment.Close(); err != nil {
				return err
			}
		}

		segment, err := l.ls.openSegment(l.currentLsn, true, false)

		if err != nil {
			return err
		}

		l.currentSegment = segment
	}

	index := l.currentLsn.Clock
	head := l.currentSegment.tail.Clock

	if index > head {
		index = head + 1

		if l.options.Block {
			l.ls.cond.L.Lock()
			defer l.ls.cond.L.Unlock()

			for index > head {
				select {
				case <-ctx.Done():
					return ctx.Err()

				default:
				}

				l.ls.cond.Wait()

				head = l.currentSegment.tail.Clock
			}
		}
	}

	if err := l.currentSegment.Seek(l.currentLsn); err != nil {
		return err
	}

	entry, err := l.currentSegment.Read(l.current)

	if err != nil {
		return err
	}

	l.current = entry

	return nil
}
