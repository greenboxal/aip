package forddb

import (
	"reflect"
	"strings"

	"github.com/greenboxal/aip/aip-forddb/pkg/typesystem"
)

type TypeAnnotation struct{}

func AnnotationFromType(t reflect.Type) *TypeMetadata {
	var metadata TypeMetadata

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if !field.Anonymous || field.Type != reflect.TypeOf(TypeAnnotation{}) {
			continue
		}

		tag := field.Tag.Get("ford")
		parts := strings.Split(tag, ",")
		name := parts[0]
		//opts := parts[1:]

		if name == "-" {
			return nil
		}

		metadata.Name = name

		return &metadata
	}

	return nil
}

type TypeMetadata struct {
	ResourceBase[TypeID, *TypeMetadata]

	Name          string                   `json:"name"`
	Kind          Kind                     `json:"kind"`
	PrimitiveKind typesystem.PrimitiveKind `json:"primitive_kind"`
	Scope         ResourceScope            `json:"scope"`

	IsRuntimeOnly bool `json:"is_runtime_only"`
}
