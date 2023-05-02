package typesystem

import "github.com/greenboxal/aip/pkg/ford/forddb"

type BasicField interface {
	Name() string
	Parent() forddb.BasicType

	BasicType() forddb.BasicType

	IsOptional() bool

	GetValue(receiver any) any
	SetValue(receiver, value any)
}
