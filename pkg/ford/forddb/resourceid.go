package forddb

import (
	"encoding"
	"encoding/json"
	"hash/fnv"
	"reflect"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/multiformats/go-multihash"
)

type BasicResourceID interface {
	json.Marshaler
	encoding.TextMarshaler
	encoding.BinaryMarshaler

	BasicResourceType() BasicResourceType
	PrimitiveKind() PrimitiveKind

	AsBasicResourceID() BasicResourceID
	AsCid() cid.Cid
	AsLink() ipld.Link

	LinkPrototype() ipld.LinkPrototype

	String() string

	Hash64() uint64
}

type ResourceID[T BasicResource] interface {
	BasicResourceID
}

type IStringResourceID interface {
	BasicResourceID

	SetValueString(value string)
}

func NewStringID[ID BasicResourceID](name string) (result ID) {
	t := reflect.TypeOf(result)
	idVal := reflect.New(t)
	idStr := idVal.Interface().(IStringResourceID)
	idStr.SetValueString(name)
	return idVal.Elem().Interface().(ID)
}

type StringResourceID[T BasicResource] string

func (s StringResourceID[T]) BasicResourceType() BasicResourceType {
	return TypeSystem().LookupByResourceType(reflect.TypeOf((*T)(nil)).Elem())
}

func (s StringResourceID[T]) PrimitiveKind() PrimitiveKind       { return PrimitiveKindString }
func (s StringResourceID[T]) String() string                     { return string(s) }
func (s StringResourceID[T]) AsBasicResourceID() BasicResourceID { return s }
func (s StringResourceID[T]) AsCid() cid.Cid {
	link := s.AsLink()

	return link.(cidlink.Link).Cid
}

func (s StringResourceID[T]) AsLink() ipld.Link {
	h, err := multihash.Sum([]byte(s), multihash.SHA2_256, -1)

	if err != nil {
		panic(err)
	}

	return s.BasicResourceType().SchemaLinkPrototype().BuildLink(h)
}

func (s StringResourceID[T]) MarshalText() (text []byte, err error)   { return []byte(s), nil }
func (s StringResourceID[T]) MarshalBinary() (data []byte, err error) { return []byte(s), nil }
func (s StringResourceID[T]) MarshalJSON() ([]byte, error)            { return json.Marshal(string(s)) }

func (s *StringResourceID[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, (*string)(s))
}

func (s *StringResourceID[T]) UnmarshalText(data []byte) error {
	*(*string)(s) = string(data)

	return nil
}

func (s *StringResourceID[T]) UnmarshalBinary(data []byte) error {
	*(*string)(s) = string(data)

	return nil
}

func (s StringResourceID[T]) Hash64() uint64 {
	h := fnv.New64()

	_, err := h.Write([]byte(s.String()))

	if err != nil {
		panic(err)
	}

	return h.Sum64()
}

func (s *StringResourceID[T]) SetValueString(value string) { *s = StringResourceID[T](value) }

func (s StringResourceID[T]) LinkPrototype() ipld.LinkPrototype {
	return s.BasicResourceType().SchemaLinkPrototype()
}

type ICidResourceID interface {
	BasicResourceID

	SetValueCid(value cid.Cid)
}

type CidResourceID[T BasicResource] cid.Cid

func (s CidResourceID[T]) BasicResourceType() BasicResourceType {
	return TypeSystem().LookupByResourceType(reflect.TypeOf((*T)(nil)).Elem())
}

func (s CidResourceID[T]) PrimitiveKind() PrimitiveKind       { return PrimitiveKindBytes }
func (s CidResourceID[T]) AsBasicResourceID() BasicResourceID { return s }
func (s CidResourceID[T]) AsCid() cid.Cid                     { return cid.Cid(s) }
func (s CidResourceID[T]) AsLink() ipld.Link                  { return cidlink.Link{Cid: s.AsCid()} }
func (s CidResourceID[T]) String() string                     { return cid.Cid(s).String() }

func (s CidResourceID[T]) MarshalJSON() ([]byte, error)   { return cid.Cid(s).MarshalJSON() }
func (s CidResourceID[T]) MarshalText() ([]byte, error)   { return cid.Cid(s).MarshalText() }
func (s CidResourceID[T]) MarshalBinary() ([]byte, error) { return cid.Cid(s).MarshalBinary() }

func (s *CidResourceID[T]) UnmarshalJSON(data []byte) error { return (*cid.Cid)(s).UnmarshalJSON(data) }
func (s *CidResourceID[T]) UnmarshalText(data []byte) error { return (*cid.Cid)(s).UnmarshalText(data) }
func (s *CidResourceID[T]) UnmarshalBinary(data []byte) error {
	return (*cid.Cid)(s).UnmarshalBinary(data)
}

func (s CidResourceID[T]) LinkPrototype() ipld.LinkPrototype {
	return s.BasicResourceType().SchemaLinkPrototype()
}

func (s *CidResourceID[T]) SetValueString(value string) {
	id, err := cid.Decode(value)

	if err != nil {
		panic(err)
	}

	s.SetValueCid(id)
}

func (s *CidResourceID[T]) SetValueCid(value cid.Cid) { *s = CidResourceID[T](value) }

func (s CidResourceID[T]) Hash64() uint64 {
	h := fnv.New64()

	_, err := h.Write(cid.Cid(s).Bytes())

	if err != nil {
		panic(err)
	}

	return h.Sum64()
}

type ILinkResourceID interface {
	BasicResourceID

	SetValueLink(value cidlink.Link)
}

type LinkResourceID[T BasicResource] cidlink.Link

func (s LinkResourceID[T]) BasicResourceType() BasicResourceType {
	return TypeSystem().LookupByResourceType(reflect.TypeOf((*T)(nil)).Elem())
}

func (s LinkResourceID[T]) PrimitiveKind() PrimitiveKind { return PrimitiveKindString }

func (s LinkResourceID[T]) AsCid() cid.Cid {
	link := s.AsLink()

	if l, ok := link.(cidlink.Link); ok {
		return l.Cid
	}

	panic("unsupported")
}

func (s LinkResourceID[T]) AsLink() ipld.Link                  { return cidlink.Link(s) }
func (s LinkResourceID[T]) AsBasicResourceID() BasicResourceID { return s }
func (s LinkResourceID[T]) String() string                     { return cidlink.Link(s).String() }
func (s LinkResourceID[T]) MarshalJSON() ([]byte, error)       { return cidlink.Link(s).MarshalJSON() }
func (s LinkResourceID[T]) MarshalText() ([]byte, error)       { return cidlink.Link(s).MarshalText() }
func (s LinkResourceID[T]) MarshalBinary() ([]byte, error)     { return cidlink.Link(s).MarshalBinary() }

func (s *LinkResourceID[T]) UnmarshalJSON(data []byte) error {
	return (*cidlink.Link)(s).UnmarshalJSON(data)
}

func (s *LinkResourceID[T]) UnmarshalText(data []byte) error {
	return (*cidlink.Link)(s).UnmarshalText(data)
}

func (s *LinkResourceID[T]) UnmarshalBinary(data []byte) error {
	return (*cidlink.Link)(s).UnmarshalBinary(data)
}

func (s *LinkResourceID[T]) SetValueLink(value cidlink.Link) { *s = LinkResourceID[T](value) }

func (s LinkResourceID[T]) LinkPrototype() ipld.LinkPrototype {
	return s.BasicResourceType().SchemaLinkPrototype()
}

func (s LinkResourceID[T]) Hash64() uint64 {
	h := fnv.New64()

	_, err := h.Write(cidlink.Link(s).Bytes())

	if err != nil {
		panic(err)
	}

	return h.Sum64()
}
