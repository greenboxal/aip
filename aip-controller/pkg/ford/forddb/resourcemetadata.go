package forddb

import (
	"time"
)

type ResourceMetadata interface {
	GetResourceMetadata() *Metadata
	GetResourceBasicID() BasicResourceID
	GetResourceTypeID() TypeID
	GetResourceVersion() uint64
}

type Metadata struct {
	Kind      TypeID        `json:"kind"`
	Scope     ResourceScope `json:"scope"`
	Namespace string        `json:"namespace"`
	Version   uint64        `json:"version"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func (r *Metadata) GetResourceMetadata() *Metadata { return r }
func (r *Metadata) GetResourceVersion() uint64     { return r.Version }
func (r *Metadata) GetResourceTypeID() TypeID {
	return r.Kind
}
