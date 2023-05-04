package forddb

import (
	"reflect"

	"github.com/greenboxal/aip/aip-controller/pkg/utils"
)

var basicResourceType = reflect.TypeOf((*BasicResource)(nil)).Elem()
var basicResourceIdType = reflect.TypeOf((*BasicResourceID)(nil)).Elem()
var basicResourcePointerType = reflect.TypeOf((*BasicResourcePointer)(nil)).Elem()
var basicResourceSlotType = reflect.TypeOf((*BasicResourceSlot)(nil)).Elem()
var resourceMetadataType = reflect.TypeOf((*ResourceMetadata[ResourceTypeID, BasicResourceType])(nil)).Elem()
var resourceMetadataTypeName = utils.GetParsedTypeName(resourceMetadataType)

func IsBasicResource(t reflect.Type) bool {
	t = DerefPointer(t)

	if utils.GetParsedTypeName(t).FullName() == resourceMetadataTypeName.FullName() {
		return false
	}

	if t.Implements(basicResourceType) {
		return true
	}

	if reflect.PtrTo(t).Implements(basicResourceType) {
		return true
	}

	return false
}

func IsBasicResourcePointer(t reflect.Type) bool {
	t = DerefPointer(t)

	if t.Implements(basicResourcePointerType) {
		return true
	}

	if reflect.PtrTo(t).Implements(basicResourcePointerType) {
		return true
	}

	return false
}

func IsBasicResourceSlot(t reflect.Type) bool {
	t = DerefPointer(t)

	if t.Implements(basicResourceSlotType) {
		return true
	}

	if reflect.PtrTo(t).Implements(basicResourceSlotType) {
		return true
	}

	return false
}

func IsBasicResourceId(t reflect.Type) bool {
	t = DerefPointer(t)

	if t.Implements(basicResourceIdType) {
		return true
	}

	if reflect.PtrTo(t).Implements(basicResourceIdType) {
		return true
	}

	return false
}

func DerefType[T any]() reflect.Type {
	t := reflect.TypeOf((*T)(nil)).Elem()

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t
}

func DerefPointer(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t
}