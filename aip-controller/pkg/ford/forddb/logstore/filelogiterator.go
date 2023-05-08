package logstore

import (
	"context"
	"encoding/json"

	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

type fileLogIterator struct {
	ls      *FileLogStore
	options forddb.LogIteratorOptions

	current    *forddb.LogEntryRecord
	currentLsn forddb.LSN

	err error
}

func newFileLogIterator(ls *FileLogStore, options forddb.LogIteratorOptions) *fileLogIterator {
	return &fileLogIterator{
		ls:      ls,
		options: options,
	}
}

func (l *fileLogIterator) Error() error {
	return l.err
}

func (l *fileLogIterator) Entry() *forddb.LogEntry {
	return &l.current.LogEntry
}

func (l *fileLogIterator) Record() *forddb.LogEntryRecord {
	return l.current
}

func (l *fileLogIterator) CurrentLsn() forddb.LSN {
	return l.currentLsn
}

func (l *fileLogIterator) SetLSN(ctx context.Context, lsn forddb.LSN) error {
	l.currentLsn = lsn

	return l.invalidate(ctx)
}

func (l *fileLogIterator) SeekRelative(ctx context.Context, offset int64) error {
	l.currentLsn.Clock += uint64(offset)

	return l.invalidate(ctx)
}

func (l *fileLogIterator) Next(ctx context.Context) bool {
	if err := l.SeekRelative(ctx, 1); err != nil {
		l.err = err
		return false
	}

	return true
}

func (l *fileLogIterator) Previous(ctx context.Context) bool {
	if err := l.SeekRelative(ctx, -1); err != nil {
		l.err = err
		return false
	}

	return true
}

func (l *fileLogIterator) Reset(lsn forddb.LSN) {
	l.currentLsn = lsn
	l.current = nil
}

func (l *fileLogIterator) Close() error {
	return nil
}

func (l *fileLogIterator) invalidate(ctx context.Context) error {
	var record forddb.LogEntryRecord

	if l.current != nil && l.current.LSN.Equals(l.currentLsn) {
		return nil
	}

	index := l.currentLsn.Clock
	head, err := l.ls.log.LastIndex()

	if err != nil {
		return err
	}

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

				head, err = l.ls.log.LastIndex()

				if err != nil {
					return err
				}
			}
		}
	}

	data, err := l.ls.log.Read(index)

	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &record); err != nil {
		return err
	}

	l.current = &record

	return nil
}
