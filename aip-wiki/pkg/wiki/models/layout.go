package models

import (
	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
)

type LayoutID struct {
	forddb.StringResourceID[*Layout]
}

type Layout struct {
	forddb.ResourceBase[LayoutID, *Layout] `json:"metadata"`

	Spec LayoutSpec `json:"spec"`
}

type LayoutSpec struct {
	Host   string   `json:"host"`
	Layout LayoutID `json:"layout"`
}
