package tracing

import (
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
	forddb.ResourceBase[TraceID, *Trace] `json:"metadata"`

	SpanIds    []SpanID `json:"span_ids"`
	RootSpanID SpanID   `json:"root_span_id"`
}

type Span struct {
	forddb.ResourceBase[SpanID, *Span] `json:"metadata"`

	TraceID  TraceID `json:"trace_id"`
	ParentID SpanID  `json:"parent_id"`

	StartedAt   time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at"`

	Name string `json:"name"`

	InnerSpanIds []SpanID `json:"inner_span_ids"`
}

func NewTraceID() TraceID {
	return forddb.NewStringID[TraceID](uuid.New().String())
}

func NewSpanID() SpanID {
	return forddb.NewStringID[SpanID](uuid.New().String())
}

func init() {
	forddb.DefineResourceType[TraceID, *Trace]("trace")
	forddb.DefineResourceType[SpanID, *Span]("span")
}