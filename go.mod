module github.com/greenboxal/aip

go 1.20

require (
	github.com/slack-go/slack v0.12.2
	github.com/stretchr/testify v1.8.2
	go.uber.org/fx v1.19.2
	go.uber.org/multierr v1.6.0
	go.uber.org/zap v1.24.0
	golang.org/x/exp v0.0.0-20230425010034-47ecfdc1ba53
	golang.org/x/sync v0.1.0
	golang.org/x/time v0.3.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.7 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/dig v1.16.1 // indirect
	golang.org/x/sys v0.1.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace go.uber.org/fx => github.com/uber-go/fx v1.19.2

replace go.uber.org/dig => github.com/uber-go/dig v1.16.1
