package forddb

import (
	"context"
)

type LogIterator interface {
	Next(ctx context.Context) bool
	Previous(ctx context.Context) bool

	SetLSN(ctx context.Context, lsn LSN) error
	SeekRelative(ctx context.Context, offset int64) error

	Error() error
	Entry() *LogEntry
	Record() *LogEntryRecord
	CurrentLsn() LSN

	Reset(lsn LSN)

	Close() error
}
