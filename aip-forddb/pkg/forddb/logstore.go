package forddb

import (
	"context"
)

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
