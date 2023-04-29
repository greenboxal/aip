package documents

type Kind string

const (
	KindUnknown Kind = "application/octet-stream"
	KindText    Kind = "text/plain"
)

type Link struct {
	Kind Kind `json:"kind"`
}

type BasicTypeID interface {
	Name() string
	MimeType() string

	BasicTypeID() BasicTypeID
}

type TypeID[T BasicType] interface {
	BasicTypeID

	Type() T
}

type BasicType interface {
	BasicTypeID() BasicTypeID
}

type Type[ID BasicTypeID] interface {
	BasicType

	TypeID() ID
}

type Metadata struct {
}

type BasicDocument interface {
	Type() BasicType

	Title() string
	Text() string

	Metadata() Metadata
	Links() []Link
}
