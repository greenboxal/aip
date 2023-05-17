package forddb

import (
	"encoding/json"
	"hash/fnv"
	"reflect"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"

	"github.com/greenboxal/aip/aip-forddb/pkg/typesystem"
)

type TypeID string

func (id TypeID) BasicResourceType() BasicResourceType {
	return TypeSystem().LookupByIDType(reflect.TypeOf(id))
}

func (id TypeID) Name() string                            { return id.String() }
func (id TypeID) PrimitiveKind() PrimitiveKind            { return typesystem.PrimitiveKindString }
func (id TypeID) AsBasicResourceID() BasicResourceID      { return id }
func (id TypeID) AsCid() cid.Cid                          { panic("unsupported") }
func (id TypeID) AsLink() ipld.Link                       { panic("unsupported") }
func (id TypeID) String() string                          { return string(id) }
func (id *TypeID) SetValueString(value string)            { *id = TypeID(value) }
func (id TypeID) Type() BasicResourceType                 { return TypeSystem().LookupByID(id) }
func (id TypeID) MarshalJSON() ([]byte, error)            { return json.Marshal(string(id)) }
func (id TypeID) MarshalText() (text []byte, err error)   { return []byte(id), nil }
func (id TypeID) MarshalBinary() (data []byte, err error) { return []byte(id), nil }
func (id *TypeID) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, (*string)(id))
}

func (id *TypeID) UnmarshalText(data []byte) error {
	*(*string)(id) = string(data)

	return nil
}

func (id *TypeID) UnmarshalBinary(data []byte) error {
	*(*string)(id) = string(data)

	return nil
}

func (id TypeID) Hash64() uint64 {
	h := fnv.New64()

	_, err := h.Write([]byte(id.Name()))

	if err != nil {
		panic(err)
	}

	return h.Sum64()
}
