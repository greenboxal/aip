package forddb

import (
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
	cloned := resource.GetResourceTypeID().Type().CreateInstance()

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

	rawResource["kind"] = resource.GetResourceTypeID().Name()

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
