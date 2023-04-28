package forddb

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type RawResource = map[string]interface{}

var basicResourceIdType = reflect.TypeOf((*BasicResourceID)(nil)).Elem()

func Encode(resource BasicResource) (RawResource, error) {
	rawResource := RawResource{}

	data, err := json.Marshal(resource)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &rawResource); err != nil {
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

	resource := typ.New()

	data, err := json.Marshal(rawResource)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, resource); err != nil {
		return nil, err
	}

	return resource, nil
}
