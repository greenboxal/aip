package forddb

import "reflect"

type BasicField interface {
	Name() string
	Parent() BasicType

	BasicType() BasicType

	IsOptional() bool

	GetValue(receiver any) any
	SetValue(receiver, value any)
}

type reflectedField struct {
	name   string
	parent BasicType
	typ    BasicType
	field  reflect.StructField
}

func NewReflectedField(name string, parent BasicType, typ BasicType, field reflect.StructField) BasicField {
	return &reflectedField{
		name:   name,
		parent: parent,
		typ:    typ,
		field:  field,
	}
}

func (f *reflectedField) Name() string {
	return f.name
}

func (f *reflectedField) Parent() BasicType {
	return f.parent
}

func (f *reflectedField) BasicType() BasicType {
	return f.typ
}

func (f *reflectedField) IsOptional() bool {
	return false
}

func (f *reflectedField) GetValue(receiver any) any {
	v := reflect.ValueOf(receiver)

	for _, i := range f.field.Index {
		v = v.Field(i)
	}

	return v.Interface()
}

func (f *reflectedField) SetValue(receiver, value any) {
	v := reflect.ValueOf(receiver)

	for _, i := range f.field.Index {
		v = v.Field(i)
	}

	v.Set(reflect.ValueOf(value))
}
