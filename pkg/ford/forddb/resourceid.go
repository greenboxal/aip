package forddb

import (
	"encoding"
	"encoding/json"
	"reflect"

	"github.com/ipfs/go-cid"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
)

type BasicResourceID interface {
	json.Marshaler
	encoding.TextMarshaler
	encoding.BinaryMarshaler

	BasicResourceID() BasicResourceID
	String() string
	MarshalJSON() ([]byte, error)
}

type ResourceID[T BasicResource] interface {
	BasicResourceID
}

type StringResourceID[T BasicResource] string

type IStringResourceID interface {
	setValueString(value string)
}

func NewStringID[ID BasicResourceID](name string) (result ID) {
	t := reflect.TypeOf(result)
	idVal := reflect.New(t)
	idStr := idVal.Interface().(IStringResourceID)
	idStr.setValueString(name)
	return idVal.Elem().Interface().(ID)
}

func (s StringResourceID[T]) MarshalText() (text []byte, err error) {
	return []byte(s), nil
}

func (s StringResourceID[T]) MarshalBinary() (data []byte, err error) {
	return []byte(s), nil
}

func (s StringResourceID[T]) BasicResourceID() BasicResourceID {
	return s
}

func (s StringResourceID[T]) String() string {
	return string(s)
}

func (s StringResourceID[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}

func (s *StringResourceID[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, (*string)(s))
}

func (s *StringResourceID[T]) setValueString(value string) {
	*s = StringResourceID[T](value)
}

type CidResourceID[T BasicResource] cid.Cid

type ICidResourceID interface {
	setValueCid(value cid.Cid)
}

func (s CidResourceID[T]) BasicResourceID() BasicResourceID {
	return s
}

func (s CidResourceID[T]) String() string {
	return cid.Cid(s).String()
}

func (s CidResourceID[T]) MarshalJSON() ([]byte, error) {
	return cid.Cid(s).MarshalJSON()
}

func (s CidResourceID[T]) MarshalText() ([]byte, error) {
	return cid.Cid(s).MarshalText()
}

func (s CidResourceID[T]) MarshalBinary() ([]byte, error) {
	return cid.Cid(s).MarshalBinary()
}

func (s *CidResourceID[T]) UnmarshalJSON(data []byte) error {
	return (*cid.Cid)(s).UnmarshalJSON(data)
}

func (s *CidResourceID[T]) UnmarshalText(data []byte) error {
	return (*cid.Cid)(s).UnmarshalText(data)
}

func (s *CidResourceID[T]) UnmarshalBinary(data []byte) error {
	return (*cid.Cid)(s).UnmarshalBinary(data)
}

func (s *CidResourceID[T]) setValueString(value string) {
	id, err := cid.Decode(value)

	if err != nil {
		panic(err)
	}

	s.setValueCid(id)
}

func (s *CidResourceID[T]) setValueCid(value cid.Cid) {
	*s = CidResourceID[T](value)
}

type LinkResourceID[T BasicResource] cidlink.Link

type ILinkResourceID interface {
	setValueLink(value cidlink.Link)
}

func (s LinkResourceID[T]) BasicResourceID() BasicResourceID {
	return s
}

func (s LinkResourceID[T]) String() string {
	return cidlink.Link(s).String()
}

func (s LinkResourceID[T]) MarshalJSON() ([]byte, error) {
	return cidlink.Link(s).MarshalJSON()
}

func (s LinkResourceID[T]) MarshalText() ([]byte, error) {
	return cidlink.Link(s).MarshalText()
}

func (s LinkResourceID[T]) MarshalBinary() ([]byte, error) {
	return cidlink.Link(s).MarshalBinary()
}

func (s *LinkResourceID[T]) UnmarshalJSON(data []byte) error {
	return (*cidlink.Link)(s).UnmarshalJSON(data)
}

func (s *LinkResourceID[T]) UnmarshalText(data []byte) error {
	return (*cidlink.Link)(s).UnmarshalText(data)
}

func (s *LinkResourceID[T]) UnmarshalBinary(data []byte) error {
	return (*cidlink.Link)(s).UnmarshalBinary(data)
}

func (s *LinkResourceID[T]) setValueLink(value cidlink.Link) {
	*s = LinkResourceID[T](value)
}
