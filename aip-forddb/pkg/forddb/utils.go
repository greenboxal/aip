package forddb

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"

	"github.com/hashicorp/go-multierror"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/pelletier/go-toml/v2"
	"github.com/samber/lo"
	"sigs.k8s.io/yaml"

	"github.com/greenboxal/aip/aip-forddb/pkg/typesystem"
)

func Put[T BasicResource](ctx context.Context, db Database, res T) (def T, _ error) {
	resource, err := db.Put(ctx, res)

	if err != nil {
		return def, err
	}

	if resource == nil {
		return def, ErrNotFound
	}

	return resource.(T), nil
}

func List[T BasicResource](
	ctx context.Context,
	db Database,
	typ BasicResourceType,
	opts ...ListOption,
) ([]T, error) {
	resources, err := db.List(ctx, typ.GetResourceID(), opts...)

	if err != nil {
		return nil, err
	}

	return lo.Map(resources, func(resource BasicResource, _ int) T {
		return resource.(T)
	}), nil
}

func Get[T BasicResource](ctx context.Context, db Database, id ResourceID[T]) (def T, _ error) {
	typ := id.BasicResourceType()
	resource, err := db.Get(ctx, typ.GetResourceID(), id)

	if err != nil {
		return def, err
	}

	if resource == nil {
		return def, ErrNotFound
	}

	return resource.(T), nil
}

func TypeOf(val any) BasicType {
	return TypeSystem().LookupByType(reflect.TypeOf(val))
}

func LookupTypeByName(name string) BasicResourceType {
	return TypeSystem().LookupBySingularName(name)
}

func LookupTypeByNamePlural(name string) BasicResourceType {
	return TypeSystem().LookupByPluralName(name)
}

func ImportPath(db Database, path string) error {
	var merr error

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && isSupportedFile(path) {
			var raw RawResource

			data, err := os.ReadFile(path)

			if err != nil {
				merr = multierror.Append(merr, err)
				return nil
			}

			switch filepath.Ext(path) {
			case ".toml":
				if err := toml.Unmarshal(data, &raw); err != nil {
					merr = multierror.Append(merr, err)
					return nil
				}

				data, err = json.Marshal(raw)

				if err != nil {
					merr = multierror.Append(merr, err)
					return nil
				}

			case ".yaml", ".yml", ".json":
				if err := yaml.Unmarshal(data, &raw); err != nil {
					merr = multierror.Append(merr, err)
					return nil
				}

				data, err = json.Marshal(raw)

				if err != nil {
					merr = multierror.Append(merr, err)
					return nil
				}
			default:
				return nil
			}

			resource, err := ipld.Decode(data, dagjson.Decode)

			if err != nil {
				merr = multierror.Append(merr, err)
				return nil
			}

			_, err = db.Put(context.TODO(), typesystem.Unwrap(resource).(BasicResource))

			if err != nil {
				merr = multierror.Append(merr, err)
				return nil
			}
		}

		return nil
	})

	if err != nil {
		merr = multierror.Append(merr, err)
	}

	return merr
}

func isSupportedFile(path string) bool {
	switch filepath.Ext(path) {
	case ".toml", ".yaml", ".yml", ".json":
		return true
	}

	return false
}
