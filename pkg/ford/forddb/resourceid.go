package forddb

import (
	"encoding/json"

	"github.com/ipfs/go-cid"
)

type BasicResourceID interface {
	BasicResourceID() BasicResourceID
	String() string
	MarshalJSON() ([]byte, error)
}

type ResourceID[T BasicResource] interface {
	BasicResourceID
}

type StringResourceID[T BasicResource] string

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

func (s *StringResourceID[T]) setValue(value string) {
	*s = StringResourceID[T](value)
}

type CidResourceID[T BasicResource] cid.Cid

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

func (s *CidResourceID[T]) setValue(value string) {
	id, err := cid.Decode(value)

	if err != nil {
		panic(err)
	}

	*s = CidResourceID[T](id)
}
