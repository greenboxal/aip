package forddb

import (
	"bytes"
	"encoding/json"
	"io"
)

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
