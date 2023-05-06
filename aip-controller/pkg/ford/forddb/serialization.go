package forddb

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/mashingan/smapping"
)

type RawResource smapping.Mapped

func (r RawResource) GetResourceMetadata() *Metadata {
	var result Metadata

	metadata := r["metadata"]

	if metadata == nil {
		return nil
	}

	metadataMapped, ok := metadata.(map[string]interface{})

	if !ok {
		return nil
	}

	if err := smapping.FillStruct(&result, metadataMapped); err != nil {
		panic(err)
	}

	return &result
}

func (r RawResource) GetResourceBasicID() BasicResourceID {
	metadata := r["metadata"]

	if metadata == nil {
		return nil
	}

	metadataMapped, ok := metadata.(map[string]interface{})

	if !ok {
		return nil
	}

	idVal, ok := metadataMapped["id"]

	if !ok {
		return nil
	}

	return r.GetResourceTypeID().Type().CreateID(idVal.(string))
}

func (r RawResource) GetResourceTypeID() ResourceTypeID {
	metadata := r["metadata"]

	if metadata == nil {
		return ""
	}

	metadataMapped, ok := metadata.(map[string]interface{})

	if !ok {
		return ""
	}

	kindVal, ok := metadataMapped["kind"]

	if !ok {
		return ""
	}

	return NewStringID[ResourceTypeID](kindVal.(string))
}

func (r RawResource) GetResourceVersion() uint64 {
	metadata := r["metadata"]

	if metadata == nil {
		return 0
	}

	metadataMapped, ok := metadata.(map[string]interface{})

	if !ok {
		return 0
	}

	versionVal, ok := metadataMapped["version"]

	if !ok {
		return 0
	}

	if v, ok := versionVal.(uint64); ok {
		return v
	}

	if v, ok := versionVal.(float64); ok {
		return uint64(v)
	}

	return 0
}

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
