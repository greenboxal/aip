package tracing

import (
	"context"
	"sync"
	"time"
)

var spanContextKey = "spanContextKey"
var spanTracerKey = "spanTracerKey"

func TracerFromContext(ctx context.Context) Tracer {
	v := ctx.Value(spanTracerKey)

	if v == nil {
		return nil
	}

	return v.(Tracer)
}

func WithTracer(ctx context.Context, tracer Tracer) context.Context {
	return context.WithValue(ctx, spanTracerKey, tracer)
}

func StartTrace(ctx context.Context, name string) (context.Context, SpanContext) {
	var parentId SpanID

	t := TracerFromContext(ctx)
	parent := getSpanContext(ctx)

	if parent != nil {
		parentId = parent.SpanID()
	}

	if t == nil {
		t = GetGlobalTracer()
	}

	tc := newTraceContext(t, NewTraceID(), parentId, name)

	ctx = context.WithValue(ctx, spanContextKey, tc.rootSpan)

	tc.rootSpan.Start()

	return ctx, tc.rootSpan
}

func StartSpan(ctx context.Context, name string) (context.Context, SpanContext) {
	parent := getSpanContext(ctx)

	if parent == nil {
		return ctx, newSpanContext(nil, SpanID{}, name)
	}

	span := newSpanContext(parent.traceCtx, parent.SpanID(), name)
	ctx = context.WithValue(ctx, spanContextKey, span)

	span.Start()

	parent.onChildStarted(span)

	return ctx, span
}

func newTraceContext(tracer Tracer, id TraceID, parentSpanId SpanID, name string) *traceContext {
	tc := &traceContext{}

	tc.tracer = tracer
	tc.trace.ID = id
	tc.rootSpan = newSpanContext(tc, parentSpanId, name)

	return tc
}

type SpanContext interface {
	TraceID() TraceID
	ParentID() SpanID
	SpanID() SpanID

	SetAttribute(key string, value string)

	End()
}

type traceContext struct {
	m        sync.Mutex
	trace    Trace
	tracer   Tracer
	rootSpan *spanContext
	finished bool
}

func (tc *traceContext) TraceID() TraceID {
	return tc.trace.ID
}

func (tc *traceContext) RootSpan() *spanContext {
	return tc.rootSpan
}

func (tc *traceContext) onSpanStarted(sc *spanContext) {
	tc.m.Lock()
	defer tc.m.Unlock()

	if tc.finished {
		panic("trace already finished")
	}

	tc.trace.SpanIds = append(tc.trace.SpanIds, sc.span.ID)

	if tc.tracer != nil {
		tc.tracer.OnSpanStarted(sc, sc.span)
	}

	if sc == tc.rootSpan {
		tc.trace.RootSpanID = tc.rootSpan.SpanID()
		tc.trace.Name = tc.rootSpan.span.Name
		tc.trace.StartedAt = tc.rootSpan.span.StartedAt
	}
}

func (tc *traceContext) onSpanFinished(sc *spanContext) {
	tc.m.Lock()
	defer tc.m.Unlock()

	if tc.finished {
		panic("trace already finished")
	}

	if tc.tracer != nil {
		tc.tracer.OnSpanFinished(sc, sc.span)
	}

	if sc == tc.rootSpan {
		tc.trace.CompletedAt = tc.rootSpan.span.CompletedAt
		tc.finished = true

		if tc.tracer != nil {
			tc.tracer.OnTraceFinished(tc, tc.trace)
		}
	}
}

func newSpanContext(traceCtx *traceContext, parentId SpanID, name string) *spanContext {
	sc := &spanContext{}

	sc.traceCtx = traceCtx
	sc.span.ID = NewSpanID()
	sc.span.ParentID = parentId
	sc.span.Name = name

	if traceCtx != nil {
		sc.span.TraceID = traceCtx.trace.ID
	}

	return sc
}

func getSpanContext(ctx context.Context) *spanContext {
	value := ctx.Value(spanContextKey)

	if value == nil {
		return nil
	}

	return value.(*spanContext)
}

type spanContext struct {
	m sync.Mutex

	traceCtx *traceContext

	span     Span
	finished bool
}

func (sc *spanContext) SetAttribute(key string, value string) {
	if sc.traceCtx == nil {
		return
	}

	sc.m.Lock()
	defer sc.m.Unlock()

	for i, v := range sc.span.Attributes {
		if v.Key == key {
			sc.span.Attributes[i].Value = value
			return
		}
	}

	sc.span.Attributes = append(sc.span.Attributes, SpanAttribute{
		Key:   key,
		Value: value,
	})
}

func (sc *spanContext) TraceID() TraceID {
	return sc.traceCtx.trace.ID
}

func (sc *spanContext) ParentID() SpanID {
	return sc.span.ParentID
}

func (sc *spanContext) SpanID() SpanID {
	return sc.span.ID
}

func (sc *spanContext) Start() {
	if sc.traceCtx == nil {
		return
	}

	sc.m.Lock()
	defer sc.m.Unlock()

	if sc.finished {
		panic("span already finished")
	}

	if !sc.span.StartedAt.IsZero() {
		panic("span already started")
	}

	sc.span.StartedAt = time.Now()
	sc.traceCtx.onSpanStarted(sc)
}

func (sc *spanContext) End() {
	if sc.traceCtx == nil {
		return
	}

	sc.m.Lock()
	defer sc.m.Unlock()

	if sc.finished {
		panic("span already finished")
	}

	if sc.span.StartedAt.IsZero() {
		panic("span not started")
	}

	sc.span.CompletedAt = time.Now()
	sc.span.Duration = sc.span.CompletedAt.Sub(sc.span.StartedAt)
	sc.finished = true

	sc.traceCtx.onSpanFinished(sc)
}

func (sc *spanContext) onChildStarted(span *spanContext) {
	sc.m.Lock()
	defer sc.m.Unlock()

	sc.span.InnerSpanIds = append(sc.span.InnerSpanIds, span.span.ID)
}
