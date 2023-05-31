package tracing

import (
	"context"

	"github.com/jbenet/goprocess"
	goprocessctx "github.com/jbenet/goprocess/context"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
)

var globalTracer Tracer

func SetGlobalTracer(tracer Tracer) {
	globalTracer = tracer
}

func GetGlobalTracer() Tracer {
	return globalTracer
}

type Tracer interface {
	OnSpanStarted(sc *spanContext, span Span)
	OnSpanFinished(sc *spanContext, span Span)
	OnTraceFinished(tc *traceContext, trace Trace)
}

type tracer struct {
	logger *zap.SugaredLogger
	db     forddb.Database

	spanCh  chan Span
	traceCh chan Trace

	worker goprocess.Process
}

func NewTracer(
	logger *zap.SugaredLogger,
	lc fx.Lifecycle,
	db forddb.Database,
) Tracer {
	t := &tracer{
		db:     db,
		logger: logger,

		spanCh:  make(chan Span, 128),
		traceCh: make(chan Trace, 128),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return t.Start()
		},

		OnStop: func(ctx context.Context) error {
			return t.Shutdown(ctx)
		},
	})

	SetGlobalTracer(t)

	return t
}

func (t *tracer) OnSpanStarted(sc *spanContext, span Span) {
	t.spanCh <- span
}

func (t *tracer) OnSpanFinished(sc *spanContext, span Span) {
	t.spanCh <- span
}

func (t *tracer) OnTraceFinished(tc *traceContext, trace Trace) {
	t.traceCh <- trace
}

func (t *tracer) Start() error {
	t.worker = goprocess.Go(func(proc goprocess.Process) {
		proc.Go(func(proc goprocess.Process) {
			ctx := goprocessctx.OnClosingContext(proc)

			for span := range t.spanCh {
				_, err := t.db.Put(ctx, &span, forddb.WithOnConflict(forddb.OnConflictReplace))

				if err != nil {
					t.logger.Errorw("failed to store span", "error", err)
				}
			}
		})

		proc.Go(func(proc goprocess.Process) {
			ctx := goprocessctx.OnClosingContext(proc)

			for trace := range t.traceCh {
				_, err := t.db.Put(ctx, &trace, forddb.WithOnConflict(forddb.OnConflictReplace))

				if err != nil {
					t.logger.Errorw("failed to store trace", "error", err)
				}
			}
		})

		if err := proc.CloseAfterChildren(); err != nil {
			panic(err)
		}
	})

	return nil
}

func (t *tracer) Shutdown(ctx context.Context) error {
	err := t.worker.Close()

	if err != nil {
		return err
	}

	return nil
}
