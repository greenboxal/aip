package ford

import "reflect"

type ResourceTypeID string
type ResourceID any

type BasicResourceType interface {
	ID() ResourceTypeID
	Name() string
	IDType() reflect.Type
	ResourceType() reflect.Type
}

type ResourceType[ID ResourceID, T Resource] interface {
	BasicResourceType
}

type resourceType[ID ResourceID, T Resource] struct {
	name         string
	idType       reflect.Type
	resourceType reflect.Type
}

func (r *resourceType[ID, T]) ID() ResourceTypeID {
	return ResourceTypeID(r.name)
}

func (r *resourceType[ID, T]) Name() string {
	return r.name
}

func (r *resourceType[ID, T]) IDType() reflect.Type {
	return r.idType
}

func (r *resourceType[ID, T]) ResourceType() reflect.Type {
	return r.resourceType
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
	return s.idTypeMap[typ]
}

func (s *resourceTypeSystem) LookupByResourceType(typ reflect.Type) BasicResourceType {
	return s.resourceTypeMap[typ]
}

var typeSystem = &resourceTypeSystem{
	resourceTypes: make(map[ResourceTypeID]BasicResourceType, 32),
}

func DefineResourceType[ID ResourceID, T Resource](name string) ResourceType[ID, T] {
	t := &resourceType[ID, T]{
		name:         name,
		idType:       reflect.TypeOf((*ID)(nil)).Elem(),
		resourceType: reflect.TypeOf((*T)(nil)).Elem(),
	}

	typeSystem.Register(t)

	return t
}

type Resource interface {
	GetID() ResourceID
	GetType() ResourceTypeID
}

type ResourceBase[ID ResourceID, T Resource] struct {
	ID ResourceID `json:"ID"`
}

func (r *ResourceBase[ID, T]) GetID() ResourceID {
	return r.ID
}

func (r *ResourceBase[ID, T]) GetType() ResourceTypeID {
	return typeSystem.LookupByResourceType(reflect.TypeOf(r)).ID()
}

type Database interface {
	Get(typ ResourceTypeID, id ResourceID) (Resource, error)

	GetPipeline(id PipelineID) (*Pipeline, error)
	GetAgent(id AgentID) (*Agent, error)
	GetTeam(id TeamID) (*Team, error)
	GetProfile(id PipelineID) (string, error)
	GetTask(id PipelineID) (string, error)
}
