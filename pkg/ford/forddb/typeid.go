package forddb

import (
	"encoding/json"
	"hash/fnv"
	"reflect"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
)

type ResourceTypeID string

func (id ResourceTypeID) BasicResourceType() BasicResourceType {
	return TypeSystem().LookupByIDType(reflect.TypeOf(id))
}

func (id ResourceTypeID) Name() string                            { return id.String() }
func (id ResourceTypeID) PrimitiveKind() PrimitiveKind            { return PrimitiveKindString }
func (id ResourceTypeID) AsBasicResourceID() BasicResourceID      { return id }
func (id ResourceTypeID) AsCid() cid.Cid                          { panic("unsupported") }
func (id ResourceTypeID) AsLink() ipld.Link                       { panic("unsupported") }
func (id ResourceTypeID) String() string                          { return string(id) }
func (id *ResourceTypeID) SetValueString(value string)            { *id = ResourceTypeID(value) }
func (id ResourceTypeID) Type() BasicResourceType                 { return TypeSystem().LookupByID(id) }
func (id ResourceTypeID) MarshalJSON() ([]byte, error)            { return json.Marshal(string(id)) }
func (id ResourceTypeID) MarshalText() (text []byte, err error)   { return []byte(id), nil }
func (id ResourceTypeID) MarshalBinary() (data []byte, err error) { return []byte(id), nil }
func (id *ResourceTypeID) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, (*string)(id))
}

func (id *ResourceTypeID) UnmarshalText(data []byte) error {
	*(*string)(id) = string(data)

	return nil
}

func (id *ResourceTypeID) UnmarshalBinary(data []byte) error {
	*(*string)(id) = string(data)

	return nil
}

func (id ResourceTypeID) LinkPrototype() ipld.LinkPrototype {
	return id.BasicResourceType().SchemaLinkPrototype()
}

func (id ResourceTypeID) Hash64() uint64 {
	h := fnv.New64()

	_, err := h.Write([]byte(id.Name()))

	if err != nil {
		panic(err)
	}

	return h.Sum64()
}
