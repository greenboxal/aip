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

func (f *ReflectedField) GetValue(receiver reflect.Value) reflect.Value {
	v := receiver

	for v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	for _, i := range f.field.Index {
		v = v.Field(i)
	}

	return v
}

func (f *ReflectedField) SetValue(receiver, value reflect.Value) {
	v := receiver

	for _, i := range f.field.Index {
		for v.Kind() == reflect.Pointer {
			v = v.Elem()
		}

		v = v.Field(i)
	}

	v.Set(value)
}
