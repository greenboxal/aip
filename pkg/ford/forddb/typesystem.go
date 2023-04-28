package forddb

import "reflect"

var typeSystem = &resourceTypeSystem{
	resourceTypes:   make(map[ResourceTypeID]BasicResourceType, 32),
	idTypeMap:       make(map[reflect.Type]BasicResourceType, 32),
	resourceTypeMap: make(map[reflect.Type]BasicResourceType, 32),
}

var typeType = &resourceType[ResourceTypeID, BasicResourceType]{
	resourceType: reflect.TypeOf((*BasicResourceType)(nil)).Elem(),
	idType:       reflect.TypeOf((*ResourceTypeID)(nil)).Elem(),
}

func init() {
	typeType.isRuntimeOnly = true
	typeType.ResourceMetadata.Name = "type"
	typeType.ResourceMetadata.ID = typeType.MakeId(typeType.ResourceMetadata.Name).(ResourceTypeID)

	typeSystem.Register(typeType)
}

type resourceTypeSystem struct {
	resourceTypes   map[ResourceTypeID]BasicResourceType
	idTypeMap       map[reflect.Type]BasicResourceType
	resourceTypeMap map[reflect.Type]BasicResourceType
}

func (s *resourceTypeSystem) Register(t BasicResourceType) {
	s.resourceTypes[t.ID()] = t
	s.idTypeMap[t.IDType()] = t
	s.resourceTypeMap[t.ResourceType()] = t
}

func (s *resourceTypeSystem) LookupByID(id ResourceTypeID) BasicResourceType {
	return s.resourceTypes[id]
}

func (s *resourceTypeSystem) LookupByIDType(typ reflect.Type) BasicResourceType {
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	return s.idTypeMap[typ]
}

func (s *resourceTypeSystem) LookupByResourceType(typ reflect.Type) BasicResourceType {
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	return s.resourceTypeMap[typ]
}
