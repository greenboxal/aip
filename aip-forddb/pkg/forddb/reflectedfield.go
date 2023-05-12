package forddb

import (
	"reflect"
)

type ReflectedField struct {
	name   string
	parent BasicType
	typ    BasicType
	field  reflect.StructField
}

func NewReflectedField(name string, parent BasicType, typ BasicType, field reflect.StructField) BasicField {
	return &ReflectedField{
		name:   name,
		parent: parent,
		typ:    typ,
		field:  field,
	}
}

func (f *ReflectedField) Name() string {
	return f.name
}

func (f *ReflectedField) Parent() BasicType {
	return f.parent
}

func (f *ReflectedField) BasicType() BasicType {
	return f.typ
}

func (f *ReflectedField) IsOptional() bool {
	return false
}

func (f *ReflectedField) GetValue(receiver any) any {
	v := reflect.ValueOf(receiver)

	for _, i := range f.field.Index {
		v = v.Field(i)
	}

	return v.Interface()
}

func (f *ReflectedField) SetValue(receiver, value any) {
	v := reflect.ValueOf(receiver)

	for _, i := range f.field.Index {
		v = v.Field(i)
	}

	v.Set(reflect.ValueOf(value))
}
