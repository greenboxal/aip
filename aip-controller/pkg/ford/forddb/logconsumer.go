package forddb

import (
	"context"
	"errors"

	"github.com/jbenet/goprocess"
	goprocessctx "github.com/jbenet/goprocess/context"
)

type LogStreamHandler func(ctx context.Context, record *LogEntryRecord) error

type LogConsumer struct {
	LogStore LogStore
	StreamID LogStreamID
	Handler  LogStreamHandler
}

func (lc *LogConsumer) Run(proc goprocess.Process) {
	ctx := goprocessctx.OnClosingContext(proc)
	stream, err := lc.LogStore.OpenStream(lc.StreamID)

	if err != nil {
		panic(err)
	}

	shouldRetry := false

	for true {
		if !shouldRetry && !stream.Next(ctx) {
			if stream.Error() != nil {
				panic(stream.Error())
			}

			break
		}

		if shouldRetry {
			shouldRetry = false
		}

		record := stream.Record()

		if record.Kind == LogEntryKindInvalid {
			panic(errors.New("invalid log entry kind"))
		}

		if err := lc.Handler(ctx, record); err != nil {
			shouldRetry = true
		}

		if err := stream.SaveCheckpoint(); err != nil {
			panic(err)
		}
	}
}
