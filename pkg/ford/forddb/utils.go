package forddb

import (
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	"github.com/pelletier/go-toml/v2"
	"sigs.k8s.io/yaml"
)

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
