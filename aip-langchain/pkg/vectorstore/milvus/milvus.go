package milvus

import (
	"context"
	"os"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"go.uber.org/fx"
)

type Config struct {
	Endpoint string `json:"endpoint"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *Config) SetDefaults() {
	if c.Endpoint == "" {
		c.Endpoint = os.Getenv("MILVUS_ENDPOINT")
	}

	if c.Username == "" {
		c.Username = os.Getenv("MILVUS_USERNAME")
	}

	if c.Password == "" {
		c.Password = os.Getenv("MILVUS_PASSWORD")
	}
}

type Milvus struct {
	client client.Client
}

func NewMilvus(
	lc fx.Lifecycle,
	config *Config,
) (*Milvus, error) {
	m := &Milvus{}

	c, err := client.NewDefaultGrpcClientWithURI(
		context.Background(),
		config.Endpoint,
		config.Username,
		config.Password,
	)

	if err != nil {
		return nil, err
	}

	m.client = c

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return m.Close()
		},
	})

	return m, nil
}

func (m *Milvus) Client() client.Client {
	return m.client
}

func (m *Milvus) Close() error {
	return m.client.Close()
}
