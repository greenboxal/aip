package forddb

import (
	"os"
	"path/filepath"
	"reflect"

	"github.com/hashicorp/go-multierror"
	"github.com/ipld/go-ipld-prime/schema"
	"github.com/pelletier/go-toml/v2"
	"sigs.k8s.io/yaml"
)

func TypeSystem() *ResourceTypeSystem {
	return typeSystem
}

func IsBasicResource(t reflect.Type) bool {
	t = derefPointer(t)

	if t.Implements(basicResourceType) {
		return true
	}

	if reflect.PtrTo(t).Implements(basicResourceType) {
		return true
	}

	return false
}

func IsBasicResourcePointer(t reflect.Type) bool {
	t = derefPointer(t)

	if t.Implements(basicResourcePointerType) {
		return true
	}

	if reflect.PtrTo(t).Implements(basicResourcePointerType) {
		return true
	}

	return false
}

func IsBasicResourceSlot(t reflect.Type) bool {
	t = derefPointer(t)

	if t.Implements(basicResourceSlotType) {
		return true
	}

	if reflect.PtrTo(t).Implements(basicResourceSlotType) {
		return true
	}

	return false
}

func IsBasicResourceId(t reflect.Type) bool {
	t = derefPointer(t)

	if t.Implements(basicResourceIdType) {
		return true
	}

	if reflect.PtrTo(t).Implements(basicResourceIdType) {
		return true
	}

	return false
}

func ImportPath(db Database, path string) error {
	var merr error

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && isSupportedFile(path) {
			var rawResource RawResource

			data, err := os.ReadFile(path)

			if err != nil {
				merr = multierror.Append(merr, err)
				return nil
			}

			switch filepath.Ext(path) {
			case ".toml":
				if err := toml.Unmarshal(data, &rawResource); err != nil {
					merr = multierror.Append(merr, err)
					return nil
				}
			case ".yaml", ".yml", ".json":
				if err := yaml.Unmarshal(data, &rawResource); err != nil {
					merr = multierror.Append(merr, err)
					return nil
				}
			default:
				return nil
			}

			resource, err := Decode(rawResource)

			if err != nil {
				merr = multierror.Append(merr, err)
				return nil
			}

			_, err = db.Put(resource)

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

func normalizedName(str string) string {
	return NormalizeTypeNameRegex.ReplaceAllString(str, "")
}

func normalizedTypeName(typ reflect.Type) schema.TypeName {
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	name := typ.Name()

	if name == "" {
		name = typ.Kind().String()
	}

	return NormalizeTypeNameRegex.ReplaceAllString(name, "")
}
