package forddb

type BasicResourceType interface {
	BasicType

	ResourceName() ResourceTypeName

	IDType() BasicType
	ResourceType() BasicType

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
