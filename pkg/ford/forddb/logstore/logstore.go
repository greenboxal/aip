package logstore

import "context"

type LogStreamID string

type LogStore interface {
	OpenStream(id LogStreamID) (LogStream, error)

	Append(ctx context.Context, log LogEntry) (LogEntryRecord, error)

	Iterator(options ...LogIteratorOption) LogIterator

	Close() error
}

type LogIteratorOptions struct {
	Block bool
}

type LogIteratorOption func(opts *LogIteratorOptions)

func NewLogIteratorOptions(options ...LogIteratorOption) LogIteratorOptions {
	opts := LogIteratorOptions{}

	for _, opt := range options {
		opt(&opts)
	}

	return opts
}

func WithBlockingIterator() LogIteratorOption {
	return func(opts *LogIteratorOptions) {
		opts.Block = true
	}
}

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

type LogStream interface {
	LogIterator

	StreamID() LogStreamID

	SaveCheckpoint() error
}
