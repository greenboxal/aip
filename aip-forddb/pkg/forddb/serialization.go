package forddb

import (
	"reflect"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagjson"

	"github.com/greenboxal/aip/aip-forddb/pkg/typesystem"
)

func Clone[T any](resource T) T {
	return CloneResource(resource).(T)
}

func CloneResource(resource any) any {
	raw := typesystem.Wrap(resource)
	data, err := ipld.Encode(raw, dagjson.Encode)

	if err != nil {
		panic(err)
	}

	node, err := ipld.DecodeUsingPrototype(data, dagjson.Decode, typesystem.TypeOf(resource).IpldPrototype())

	if err != nil {
		panic(err)
	}

	return typesystem.Unwrap(node)
}

func Convert[T any](raw RawResource) (def T, _ error) {
	data, err := ipld.Encode(raw, dagjson.Encode)

	if err != nil {
		return def, err
	}

	node, err := ipld.DecodeUsingPrototype(data, dagjson.Decode, typesystem.TypeOf(def).IpldPrototype())

	if err != nil {
		return def, err
	}

	return typesystem.Unwrap(node).(T), nil
}

func ConvertNode[T any](node ipld.Node) (def T, _ error) {
	data, err := ipld.Encode(node, dagjson.Encode)

	if err != nil {
		return def, err
	}

	result, err := ipld.DecodeUsingPrototype(data, dagjson.Decode, typesystem.TypeOf(def).IpldPrototype())

	if err != nil {
		return def, err
	}

	u := typesystem.Unwrap(result)
	v := reflect.ValueOf(u)
	t := reflect.TypeOf((*T)(nil)).Elem()

	if t.Kind() == reflect.Ptr && v.Kind() != reflect.Ptr {
		v = v.Addr()
	} else if t.Kind() != reflect.Ptr && v.Kind() == reflect.Ptr {
		for v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
	}

	return v.Interface().(T), nil
}

func DecodeAs[T any](rawResource RawResource) (def T, _ error) {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	r := typesystem.Unwrap(rawResource)

	if r == nil {
		return def, nil
	}

	t, ok := r.(T)

	if !ok {
		v := reflect.ValueOf(r)

		if v.CanConvert(typ) {
			v = v.Convert(typ)
		} else if v.Kind() == reflect.Interface {
			v = v.Elem().Addr()
		} else if v.Kind() != reflect.Pointer {
			if v.CanAddr() {
				v = v.Addr()
			} else {
				v2 := reflect.New(v.Type())
				v2.Elem().Set(v)
				v = v2.Convert(typ)
			}
		}

		t, ok = v.Interface().(T)

		if !ok {
			panic("invalid type")
		}
	}

	return t, nil
}
