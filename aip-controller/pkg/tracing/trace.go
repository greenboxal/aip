package tracing

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
)

type TraceID struct {
	forddb.StringResourceID[*Trace]
}

type SpanID struct {
	forddb.StringResourceID[*Span]
}

type Trace struct {
	forddb.ResourceBase[TraceID, *Trace]

	Spans      []SpanID `json:"spans"`
	RootSpanID SpanID   `json:"root_span_id"`
}

type Span struct {
	forddb.ResourceBase[SpanID, *Span]

	TraceID  TraceID `json:"trace_id"`
	ParentID SpanID  `json:"parent_id"`

	StartedAt   time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at"`
}

type SpanContext interface {
	TraceID() TraceID
	ParentID() SpanID
	SpanID() SpanID

	End()
}

const SpanContextKey = "trace_span_context"

func FromContext(ctx context.Context) (SpanContext, bool) {
	value := ctx.Value(SpanContextKey)

	if value == nil {
		return nil, false
	}

	return value.(SpanContext), true
}

func Start(ctx context.Context, name string) (context.Context, SpanContext) {
	return nil, nil
}

func NewTraceID() TraceID {
	return forddb.NewStringID[TraceID](uuid.New().String())
}

func NewSpanID() SpanID {
	return forddb.NewStringID[SpanID](uuid.New().String())
}

func NewSpanContext(traceCtx *traceContext, parentId SpanID, name string) SpanContext {
	sc := &spanContext{
		traceCtx:  traceCtx,
		spanId:    NewSpanID(),
		parentId:  parentId,
		name:      name,
		startTime: time.Now(),
	}

	traceCtx.addSpan(sc)

	return sc
}

type traceContext struct {
	m sync.Mutex

	traceId   TraceID
	spans     []*spanContext
	rootSpan  *spanContext
	startTime time.Time
	endTime   time.Time
}

func (s *traceContext) TraceID() TraceID {
	return s.traceId
}

func (s *traceContext) RootSpan() *spanContext {
	return s.rootSpan
}

func (s *traceContext) addSpan(sc *spanContext) {
	s.m.Lock()
	defer s.m.Unlock()

	s.spans = append(s.spans, sc)
}

type spanContext struct {
	traceCtx  *traceContext
	parentID  SpanID
	parentId  SpanID
	spanId    SpanID
	name      string
	startTime time.Time
	endTime   time.Time
}

func (s *spanContext) TraceID() TraceID {
	return s.traceCtx.traceId
}

func (s *spanContext) ParentID() SpanID {
	return s.parentID
}

func (s *spanContext) SpanID() SpanID {
	return s.spanId
}

func (s *spanContext) End() {
	s.endTime = time.Now()
}
