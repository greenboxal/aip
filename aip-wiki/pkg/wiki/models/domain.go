package models

import (
	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

type DomainID struct {
	forddb.StringResourceID[*Domain]
}

type Domain struct {
	forddb.ResourceBase[DomainID, *Domain] `json:"metadata"`

	Spec DomainSpec `json:"spec"`
}

type DomainSpec struct {
	Host string `json:"host"`
}
