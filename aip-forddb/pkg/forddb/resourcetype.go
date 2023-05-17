package forddb

import "github.com/greenboxal/aip/aip-forddb/pkg/typesystem"

type FilterableField struct {
	Field     typesystem.Field
	Operators []string
}

type BasicResourceType interface {
	BasicType

	ResourceName() ResourceTypeName

	IDType() BasicType
	ResourceType() BasicType

	FilterableFields() []FilterableField

	CreateID(name string) BasicResourceID
}

type ResourceType[ID ResourceID[T], T Resource[ID]] interface {
	Type[T]

	BasicResourceType
}

func DefineResourceType[ID ResourceID[T], T Resource[ID]](name string) ResourceType[ID, T] {
	t := newResourceType[ID, T](name)

	TypeSystem().Register(t)

	return t
}
