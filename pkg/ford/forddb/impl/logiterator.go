package forddbimpl

import (
	"context"
	"errors"

	"github.com/greenboxal/aip/pkg/ford/forddb"
)

type logIterator struct {
	ls *FileLogStore

	current        *forddb.LogEntryRecord
	currentLsn     forddb.LSN
	currentSegment *fileStoreSegment

	err error
}

func newLogIterator(ls *FileLogStore) *logIterator {
	return &logIterator{
		ls: ls,
	}
}

func (l *logIterator) Error() error {
	return l.err
}

func (l *logIterator) Entry() *forddb.LogEntry {
	return &l.current.Entry
}

func (l *logIterator) Record() *forddb.LogEntryRecord {
	return l.current
}

func (l *logIterator) CurrentLsn() forddb.LSN {
	return l.currentLsn
}

func (l *logIterator) SetLSN(ctx context.Context, lsn forddb.LSN) error {
	l.currentLsn = lsn

	return l.invalidate(ctx)
}

func (l *logIterator) SeekRelative(ctx context.Context, offset int64) error {
	return errors.New("not implemented")
}

func (l *logIterator) Next(ctx context.Context) bool {
	if l.current == nil {
		if err := l.invalidate(ctx); err != nil {
			l.err = err
			return false
		}

		return l.current != nil
	}

	if err := l.SeekRelative(ctx, 1); err != nil {
		l.err = err
		return false
	}

	return true
}

func (l *logIterator) Previous(ctx context.Context) bool {
	if l.current == nil {
		if err := l.invalidate(ctx); err != nil {
			l.err = err
			return false
		}

		return l.current != nil
	}

	if err := l.SeekRelative(ctx, -1); err != nil {
		l.err = err
		return false
	}

	return true
}

func (l *logIterator) Reset(lsn forddb.LSN) {
	l.currentLsn = lsn
	l.current = nil
}

func (l *logIterator) Close() error {
	return nil
}

func (l *logIterator) invalidate(ctx context.Context) error {
	if l.current != nil && l.current.LSN.Equals(l.currentLsn) {
		return nil
	}

	if l.currentSegment == nil || !l.currentLsn.IsBetween(l.currentSegment.head, l.currentSegment.tail) {
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
