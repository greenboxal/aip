package forddb

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

func EncodeIpld(v reflect.Value, typ BasicType) (ipld.Node, error) {
	if typ == nil {
		typ = TypeSystem().LookupByType(v.Type())
	}

	fmt.Printf("Encoding %s: %s %s %s\n", typ.Name(), typ.Kind(), typ.PrimitiveKind(), v.Type().String())

	switch typ.PrimitiveKind() {
	case PrimitiveKindStruct:
		fields := typ.Fields()

		return fluent.BuildMap(basicnode.Prototype.Any, int64(len(fields)), func(la fluent.MapAssembler) {
			for _, field := range fields {
				val := field.GetValue(v)
				vn, err := EncodeIpld(val, field.BasicType())

				if err != nil {
					panic(err)
				}

				la.AssembleKey().AssignString(field.Name())
				la.AssembleValue().AssignNode(vn)
			}
		})

	case PrimitiveKindString:
		return fluent.Build(basicnode.Prototype.String, func(fna fluent.NodeAssembler) {
			fna.AssignString(v.String())
		})

	case PrimitiveKindInt:
		return fluent.Build(basicnode.Prototype.Int, func(fna fluent.NodeAssembler) {
			fna.AssignInt(v.Int())
		})

	case PrimitiveKindUnsignedInt:
		return fluent.Build(basicnode.Prototype.Int, func(fna fluent.NodeAssembler) {
			fna.AssignInt(int64(v.Uint()))
		})

	case PrimitiveKindFloat:
		return fluent.Build(basicnode.Prototype.Float, func(fna fluent.NodeAssembler) {
			fna.AssignFloat(v.Float())
		})

	case PrimitiveKindBoolean:
		return fluent.Build(basicnode.Prototype.Bool, func(fna fluent.NodeAssembler) {
			fna.AssignBool(v.Bool())
		})

	case PrimitiveKindBytes:
		return fluent.Build(basicnode.Prototype.Bytes, func(fna fluent.NodeAssembler) {
			fna.AssignBytes(v.Bytes())
		})

	case PrimitiveKindList:
		return fluent.BuildList(basicnode.Prototype.Any, int64(v.Len()), func(la fluent.ListAssembler) {
			for i := 0; i < v.Len(); i++ {
				n, err := EncodeIpld(v.Index(i), nil)

				if err != nil {
					panic(err)
				}

				la.AssembleValue().AssignNode(n)
			}
		})

	case PrimitiveKindMap:
		keys := v.MapKeys()

		return fluent.BuildMap(basicnode.Prototype.Any, int64(v.Len()), func(la fluent.MapAssembler) {
			for _, k := range keys {
				kn, err := EncodeIpld(k, nil)

				if err != nil {
					panic(err)
				}

				val := v.MapIndex(k)
				vn, err := EncodeIpld(val, nil)

				if err != nil {
					panic(err)
				}

				la.AssembleKey().AssignNode(kn)
				la.AssembleValue().AssignNode(vn)
			}

		})
	}

	return nil, errors.New("unsupported kind")
}

func EncodeWithIpld(value any) (RawResource, error) {
	var result RawResource

	n, err := EncodeIpld(reflect.ValueOf(value), nil)

	if err != nil {
		return nil, err
	}

	data, err := ipld.Encode(n, dagjson.Encode)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}
