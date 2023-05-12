package config

import "go.uber.org/fx"

var Module = fx.Module(
	"config",

	fx.Provide(NewConfigManager),
	fx.Provide(NewResourceManager),
	fx.Provide(NewNetworkManager),

	RegisterConfig[ResourceConfig]("res"),
	RegisterConfig[NetworkConfig]("net"),
)

func RegisterConfig[T any](key string) fx.Option {
	return fx.Provide(func(cm *ConfigManager) (def *T, _ error) {
		var v T

		vp := &v

		if setDefaults, ok := interface{}(vp).(HasSetDefaults); ok {
			setDefaults.SetDefaults()
		}

		if err := cm.Unmarshal(key, vp); err != nil {
			return def, err
		}

		return &v, nil
	})
}

type HasSetDefaults interface {
	SetDefaults()
}
