package config

import (
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type ConfigManager struct {
	*koanf.Koanf
}

type ConfigSource struct {
	Provider koanf.Provider
	Parser   koanf.Parser
}

func NewConfigManager() (*ConfigManager, error) {
	cm := &ConfigManager{
		Koanf: koanf.New("."),
	}

	sources := []ConfigSource{
		// {Provider: posflag.FlagVal(nil, nil)},

		{Provider: file.Provider("aip.yaml"), Parser: yaml.Parser()},
		{Provider: file.Provider("aip.toml"), Parser: toml.Parser()},
		{Provider: file.Provider("aip.json"), Parser: json.Parser()},
		{Provider: file.Provider(".env"), Parser: dotenv.Parser()},

		{Provider: env.Provider("AIP_", "_", func(s string) string {
			return strings.ToLower(strings.TrimPrefix(s, "AIP_"))
		})},
	}

	for _, src := range sources {
		if err := cm.Load(src.Provider, src.Parser); err != nil {
			if !os.IsNotExist(err) {
				return nil, err
			}
		}
	}

	return cm, nil
}
