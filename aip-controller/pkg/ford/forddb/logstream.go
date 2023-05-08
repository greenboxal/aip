package forddb

import (
	"context"
)

type LogStreamID string

type LogStream interface {
	LogIterator

	StreamID() LogStreamID

	SaveCheckpoint() error
}

type logStream struct {
	LogIterator

	id         LogStreamID
	checkpoint LogCheckpoint
}

func NewLogStream(ls LogStore, id LogStreamID, checkpoint LogCheckpoint) (LogStream, error) {
	l := &logStream{
		LogIterator: ls.Iterator(WithBlockingIterator()),

		id:         id,
		checkpoint: checkpoint,
	}

	if checkpoint != nil {
		lsn, err := checkpoint.Load(context.Background(), id)

		if err != nil {
			return nil, err
		}

		l.Reset(lsn)
	}

	return l, nil
}

func (l *logStream) StreamID() LogStreamID {
	return l.id
}

func (l *logStream) SaveCheckpoint() error {
	if l.checkpoint == nil {
		return nil
	}

	return l.checkpoint.Save(context.Background(), l.id, l.CurrentLsn())
}

func (l *logStream) Close() error {
	if l.checkpoint != nil {
		if err := l.checkpoint.Close(); err != nil {
			return err
		}
	}

	return l.LogIterator.Close()
}
