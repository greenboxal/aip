package forddb

import (
	"context"
	"errors"
)

type LogStreamHandler func(ctx context.Context, record *LogEntryRecord) error

type LogConsumer struct {
	LogStore LogStore
	StreamID LogStreamID
	Handler  LogStreamHandler
}

func (lc *LogConsumer) Run(ctx context.Context) error {
	stream, err := lc.LogStore.OpenStream(lc.StreamID)

	if err != nil {
		return err
	}

	for true {
		if !stream.Next(ctx) {
			if stream.Error() != nil {
				return stream.Error()
			}

			break
		}

		record := stream.Record()

		if record.Kind == LogEntryKindInvalid {
			return errors.New("invalid log entry kind")
		}

		if err := lc.Handler(ctx, record); err != nil {
			return err
		}

		if err := stream.SaveCheckpoint(); err != nil {
			return err
		}
	}

	return nil
}
