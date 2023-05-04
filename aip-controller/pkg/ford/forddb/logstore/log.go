package logstore

import (
	"time"

	"github.com/ipfs/go-cid"

	forddb2 "github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

type LogEntryKind int

const (
	LogEntryKindInvalid LogEntryKind = iota
	LogEntryKindLookup
	LogEntryKindGet
	LogEntryKindSet
	LogEntryKindDelete
)

type LSN struct {
	Clock uint64    `json:"clock"`
	TS    time.Time `json:"ts"`
}

func (l LSN) IsZero() bool {
	return l.Clock == 0 && l.TS.IsZero()
}

func (l LSN) Equals(other LSN) bool {
	return l.Clock == other.Clock && l.TS == other.TS
}

func (l LSN) Less(other LSN) bool {
	return l.Clock < other.Clock
}

func (l LSN) LessEqual(other LSN) bool {
	return l.Clock <= other.Clock
}

func (l LSN) Greater(other LSN) bool {
	return l.Clock > other.Clock
}

func (l LSN) GreaterEqual(other LSN) bool {
	return l.Clock >= other.Clock
}

func (l LSN) IsBetween(a, b LSN) bool {
	return l.GreaterEqual(a) && l.Less(b)
}

func MakeLSN(clock uint64, ts time.Time) LSN {
	return LSN{
		Clock: clock,
		TS:    ts,
	}
}

type LogEntry struct {
	ID   forddb2.BasicResourceID `json:"id"`
	Type forddb2.ResourceTypeID  `json:"type"`

	Kind LogEntryKind `json:"kind"`

	Version uint64 `json:"version,omitempty"`

	CurrentCid  *cid.Cid `json:"current_cid,omitempty"`
	PreviousCid *cid.Cid `json:"previous_cid,omitempty"`

	Previous forddb2.RawResource `json:"previous,omitempty"`
	Current  forddb2.RawResource `json:"current,omitempty"`

	// Runtime only
	CachedPrevious forddb2.BasicResource `json:"-"`
	CachedCurrent  forddb2.BasicResource `json:"-"`
}

type LogEntryRecord struct {
	LSN
	LogEntry `json:"entry"`
}