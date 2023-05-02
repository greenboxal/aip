package forddb

type BasicField interface {
	Name() string
	Parent() BasicType

	BasicType() BasicType

	IsOptional() bool

	GetValue(receiver any) any
	SetValue(receiver, value any)
}
