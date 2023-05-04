package config

import "go.uber.org/fx"

var Module = fx.Module(
	"config",

	fx.Provide(NewConfigManager),
	fx.Provide(NewResourceManager),
)

func RegisterConfig[T any](key string) fx.Option {
	return fx.Provide(func(cm *ConfigManager) (def *T, _ error) {
		var v T

		if err := cm.Unmarshal(key, &v); err != nil {
			return def, err
		}

		return &v, nil
	})
}
