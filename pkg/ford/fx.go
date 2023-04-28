package ford

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"sigs.k8s.io/yaml"

	"github.com/greenboxal/aip/pkg/ford/forddb"
)

var Module = fx.Module(
	"ford",

	fx.Provide(NewDatabase),
	fx.Provide(NewManager),
	fx.Provide(NewTaskReconciler),
	fx.Provide(NewAgentReconciler),
	fx.Provide(NewPortReconciler),

	fx.Invoke(func(
		logger *zap.SugaredLogger,
		db forddb.Database,
		tr *TaskReconciler,
		pr *PortReconciler,
		ar *AgentReconciler,
	) error {
		var resources []forddb.RawResource

		err := filepath.Walk("./data", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && isSupportedFile(path) {
				var rawResource map[string]interface{}

				data, err := os.ReadFile(path)

				if err != nil {
					logger.Warn(err)
					return nil
				}

				switch filepath.Ext(path) {
				case ".toml":
					if err := toml.Unmarshal(data, &rawResource); err != nil {
						logger.Warn(err)
						return nil
					}
				case ".yaml", ".yml", ".json":
					if err := yaml.Unmarshal(data, &rawResource); err != nil {
						logger.Warn(err)
						return nil
					}
				default:
					return nil
				}

				resources = append(resources, rawResource)
			}

			return nil
		})

		if err != nil {
			return err
		}

		for _, rawResource := range resources {
			resource, err := forddb.Decode(rawResource)

			if err != nil {
				logger.Warn(err)
				return nil
			}

			_, err = forddb.CreateOrUpdate(db, resource)

			if err != nil {
				return err
			}
		}

		return nil
	}),
)

func isSupportedFile(path string) bool {
	switch filepath.Ext(path) {
	case ".toml", ".yaml", ".yml", ".json":
		return true
	}

	return false
}

func NewDatabase() forddb.Database {
	return forddb.NewInMemory()
}
