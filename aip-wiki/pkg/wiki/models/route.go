package models

import (
	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
)

type RouteBindingID struct {
	forddb.StringResourceID[*RouteBinding] `ipld:",inline"`
}

type RouteBinding struct {
	forddb.ResourceBase[RouteBindingID, *RouteBinding] `json:"metadata"`

	DomainID DomainID `json:"domain_id"`
	PageID   PageID   `json:"page_id"`
}
