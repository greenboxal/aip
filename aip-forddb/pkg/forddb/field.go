package forddb

import "reflect"

type BasicField interface {
	Name() string
	Parent() BasicType

	BasicType() BasicType

	IsOptional() bool

	GetValue(receiver reflect.Value) reflect.Value
	SetValue(receiver, value reflect.Value)
}
