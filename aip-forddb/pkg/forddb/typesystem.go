package forddb

import (
	"reflect"
	"strings"
	"sync"

	"github.com/greenboxal/aip/aip-forddb/pkg/typesystem"
	"github.com/greenboxal/aip/aip-sdk/pkg/utils"
)

var typeSystem = NewResourceTypeSystem()

func TypeSystem() *ResourceTypeSystem {
	return typeSystem
}

type ResourceTypeSystem struct {
	m sync.Mutex

	resourceTypes   map[TypeID]BasicResourceType
	idTypeMap       map[reflect.Type]BasicResourceType
	resourceTypeMap map[reflect.Type]BasicResourceType
	singularNameMap map[string]BasicResourceType
	pluralNameMap   map[string]BasicResourceType
	typeMap         map[reflect.Type]BasicType
}

func NewResourceTypeSystem() *ResourceTypeSystem {
	typeType := newResourceType[TypeID, BasicResourceType]("type")
	typeType.SetRuntimeOnly()

	rts := &ResourceTypeSystem{
		resourceTypes:   make(map[TypeID]BasicResourceType, 32),
		idTypeMap:       make(map[reflect.Type]BasicResourceType, 32),
		typeMap:         make(map[reflect.Type]BasicType, 32),
		singularNameMap: make(map[string]BasicResourceType, 32),
		pluralNameMap:   make(map[string]BasicResourceType, 32),
		resourceTypeMap: make(map[reflect.Type]BasicResourceType, 32),
	}

	rts.Register(typeType)

	return rts
}

func (rts *ResourceTypeSystem) Register(t BasicResourceType) {
	if rts.doRegister(t, true) {
		t.Initialize(rts)
	}
}
func (rts *ResourceTypeSystem) doRegister(t BasicType, lock bool) bool {
	if lock {
		rts.m.Lock()
		defer rts.m.Unlock()
	}

	if _, ok := rts.typeMap[t.RuntimeType()]; ok {
		return false
	}

	rts.typeMap[t.RuntimeType()] = t

	if t, ok := t.(BasicResourceType); ok {
		if existing, ok := rts.idTypeMap[t.IDType().RuntimeType()]; ok {
			if existing != t {
				panic("duplicate resource type ID type")
			}
		}

		rts.singularNameMap[strings.ToLower(t.ResourceName().Name)] = t
		rts.pluralNameMap[strings.ToLower(t.ResourceName().Plural)] = t
		rts.resourceTypes[t.GetResourceID()] = t
		rts.resourceTypeMap[t.RuntimeType()] = t
		rts.idTypeMap[t.IDType().RuntimeType()] = t
	}

	return true
}

func (rts *ResourceTypeSystem) LookupBySingularName(name string) BasicResourceType {
	name = strings.ToLower(name)

	return rts.singularNameMap[name]
}

func (rts *ResourceTypeSystem) LookupByPluralName(name string) BasicResourceType {
	name = strings.ToLower(name)

	return rts.pluralNameMap[name]
}

func (rts *ResourceTypeSystem) LookupByID(id TypeID) BasicResourceType {
	return rts.resourceTypes[id]
}

func (rts *ResourceTypeSystem) LookupByIDType(typ reflect.Type) BasicResourceType {
	typ = DerefPointer(typ)

	return rts.idTypeMap[typ]
}

func (rts *ResourceTypeSystem) LookupByResourceType(typ reflect.Type) BasicResourceType {
	typ = DerefPointer(typ)

	return rts.resourceTypeMap[typ]
}

func (rts *ResourceTypeSystem) LookupByType(rt reflect.Type) (result BasicType) {
	isNew := false

	if IsBasicResource(rt) {
		return rts.LookupByResourceType(rt)
	}

	defer func() {
		if isNew {
			result.Initialize(rts)
		}
	}()

	typ := typesystem.TypeFrom(DerefPointer(rt))

	if existing := rts.typeMap[rt]; existing != nil {
		return existing
	}

	rts.m.Lock()
	defer rts.m.Unlock()

	if existing := rts.typeMap[rt]; existing != nil {
		return existing
	}

	metadata := AnnotationFromType(rt)

	if metadata == nil {
		metadata = &TypeMetadata{}
	}

	if metadata.Name == "" {
		name := utils.GetParsedTypeName(rt).NormalizedFullNameWithArguments()

		metadata.Name = name
	}

	if metadata.Kind == KindInvalid {
		kind := KindValue

		if IsBasicResourceId(rt) {
			kind = KindId
		}

		metadata.Kind = kind
	}

	if metadata.PrimitiveKind == typesystem.PrimitiveKindInvalid {
		metadata.PrimitiveKind = typ.PrimitiveKind()
	}

	result = newBasicType(typ, *metadata)

	isNew = true

	rts.doRegister(result, false)

	return
}

func (rts *ResourceTypeSystem) ResourceTypes() []BasicResourceType {
	var result = make([]BasicResourceType, 0, len(rts.resourceTypes))

	for _, typ := range rts.resourceTypes {
		result = append(result, typ)
	}

	return result
}
