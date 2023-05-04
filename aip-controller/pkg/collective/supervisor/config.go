package supervisor

import (
	"errors"

	"github.com/greenboxal/aip/aip-controller/pkg/collective/comms/transports"
)

type Config struct {
	ID string

	Args []string

	Port transports.Port
}

func (c *Config) Validate() error {
	if c.ID == "" {
		return errors.New("id is required")
	}

	if c.Port == nil {
		return errors.New("port is required")
	}

	return nil
}

type SupervisorOption func(cfg *Config)

func WithID(id string) SupervisorOption {
	return func(cfg *Config) {
		cfg.ID = id
	}
}

func WithArgs(args ...string) SupervisorOption {
	return func(cfg *Config) {
		cfg.Args = args
	}
}

func WithPort(port transports.Port) SupervisorOption {
	return func(cfg *Config) {
		cfg.Port = port
	}
}

func NewConfig(opts ...SupervisorOption) *Config {
	cfg := &Config{}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}
