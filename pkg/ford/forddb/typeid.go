package forddb

import (
	"encoding/json"
	"hash/fnv"
	"reflect"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
)

type ResourceTypeID string

func (s ResourceTypeID) AsCid() cid.Cid {
	panic("unsupported")
}

func (s ResourceTypeID) AsLink() ipld.Link {
	panic("unsupported")
}

func (s ResourceTypeID) AsBasicResourceID() BasicResourceID {
	return s
}

func (i ResourceTypeID) BasicResourceType() BasicResourceType {
	return typeSystem.LookupByIDType(reflect.TypeOf(i))
}

func (i ResourceTypeID) Type() BasicResourceType {
	return typeSystem.LookupByID(i)
}

func (i ResourceTypeID) Name() string {
	return i.String()
}

func (s ResourceTypeID) String() string {
	return string(s)
}

func (s ResourceTypeID) LinkPrototype() ipld.LinkPrototype {
	return s.BasicResourceType().SchemaLinkPrototype()
}

func (s ResourceTypeID) Hash64() uint64 {
	h := fnv.New64()

	_, err := h.Write([]byte(s.Name()))

	if err != nil {
		panic(err)
	}

	return h.Sum64()
}

func (s ResourceTypeID) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}

func (s ResourceTypeID) MarshalText() (text []byte, err error) {
	return []byte(s), nil
}

func (s ResourceTypeID) MarshalBinary() (data []byte, err error) {
	return []byte(s), nil
}

func (s *ResourceTypeID) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, (*string)(s))
}

func (s *ResourceTypeID) UnmarshalText(data []byte) error {
	*(*string)(s) = string(data)

	return nil
}

func (s *ResourceTypeID) UnmarshalBinary(data []byte) error {
	*(*string)(s) = string(data)

	return nil
}

func (s *ResourceTypeID) setValueString(value string) {
	*s = ResourceTypeID(value)
}
