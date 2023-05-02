package forddb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/mashingan/smapping"
)

type RawResource = smapping.Mapped

func Clone[T BasicResource](resource T) T {
	return CloneResource(resource).(T)
}

func CloneResource(resource BasicResource) BasicResource {
	rawResource := smapping.MapTags(resource, "json")
	cloned := resource.GetType().Type().CreateInstance()

	if err := smapping.FillStruct(cloned, rawResource); err != nil {
		panic(err)
	}

	return cloned.(BasicResource)
}

func Encode(resource BasicResource) (RawResource, error) {
	var rawResource RawResource

	encoded, err := json.Marshal(resource)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(encoded, &rawResource); err != nil {
		return nil, err
	}

	rawResource["kind"] = resource.GetType().Name()

	return rawResource, nil
}

func Decode(rawResource RawResource) (BasicResource, error) {
	kind := rawResource["kind"].(string)
	typ := LookupTypeByName(kind)

	if typ == nil {
		return nil, fmt.Errorf("unknown resource type: %s", kind)
	}

	resource := typ.CreateInstance()

	encoded, err := json.Marshal(rawResource)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(encoded, resource); err != nil {
		return nil, err
	}

	return resource.(BasicResource), nil
}

type CodecEncoder interface {
	Encode(resource RawResource) ([]byte, error)
	EncodeTo(writer io.Writer, resource RawResource) error
}

type CodecEncodeFunc func(resource RawResource) ([]byte, error)

func (f CodecEncodeFunc) EncodeTo(writer io.Writer, resource RawResource) error {
	data, err := f.Encode(resource)

	if err != nil {
		return err
	}

	_, err = writer.Write(data)

	return err
}

func (f CodecEncodeFunc) Encode(resource RawResource) ([]byte, error) {
	return f(resource)
}

type CodecEncodeToFunc func(writer io.Writer, resource RawResource) error

func (f CodecEncodeToFunc) Encode(resource RawResource) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)

	if err := f.EncodeTo(buffer, resource); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (f CodecEncodeToFunc) EncodeTo(writer io.Writer, resource RawResource) error {
	return f(writer, resource)
}

type CodecDecoder interface {
	Decode(data []byte) (RawResource, error)
	DecodeFrom(reader io.Reader) (RawResource, error)
}

type CodecDecodeFunc func(data []byte) (RawResource, error)

func (f CodecDecodeFunc) Decode(data []byte) (RawResource, error) {
	return f(data)
}

func (f CodecDecodeFunc) DecodeFrom(reader io.Reader) (RawResource, error) {
	data, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	return f.Decode(data)
}

type CodecDecodeFromFunc func(reader io.Reader) (RawResource, error)

func (f CodecDecodeFromFunc) Decode(data []byte) (RawResource, error) {
	return f.DecodeFrom(bytes.NewReader(data))
}

func (f CodecDecodeFromFunc) DecodeFrom(reader io.Reader) (RawResource, error) {
	return f(reader)
}

type Codec struct {
	CodecEncoder
	CodecDecoder
}

var Json = Codec{
	CodecEncoder: CodecEncodeFunc(func(resource RawResource) ([]byte, error) {
		return json.Marshal(resource)
	}),

	CodecDecoder: CodecDecodeFunc(func(data []byte) (RawResource, error) {
		var resource RawResource

		if err := json.Unmarshal(data, &resource); err != nil {
			return nil, err
		}

		return resource, nil
	}),
}

func Serialize(codec Codec, resource BasicResource) ([]byte, error) {
	rawResource, err := Encode(resource)

	if err != nil {
		return nil, err
	}

	return codec.Encode(rawResource)
}

func SerializeTo(writer io.Writer, codec Codec, resource BasicResource) error {
	rawResource, err := Encode(resource)

	if err != nil {
		return err
	}

	return codec.EncodeTo(writer, rawResource)
}

func DeserializeFrom(reader io.Reader, codec Codec) (BasicResource, error) {
	rawResource, err := codec.DecodeFrom(reader)

	if err != nil {
		return nil, err
	}

	return Decode(rawResource)
}

func Deserialize(data []byte, codec Codec) (BasicResource, error) {
	rawResource, err := codec.Decode(data)

	if err != nil {
		return nil, err
	}

	return Decode(rawResource)
}
