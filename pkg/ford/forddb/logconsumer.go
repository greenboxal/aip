package forddb

import "context"

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

	for stream.Next(ctx) {
		if err := lc.Handler(ctx, stream.Record()); err != nil {
			return err
		}

		if err := stream.SaveCheckpoint(); err != nil {
			return err
		}
	}

	return nil
}
