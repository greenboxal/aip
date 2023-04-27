package agent

type Manifest struct {
	Name string `json:"name"`

	// Communication
	PortName           string   `json:"port_name"`
	SubscribedChannels []string `json:"subscribed_channels"`

	// AGI Configuration
	Memory MemoryDefinition `json:"memory"`
	Tools  []Tool           `json:"tools"`

	// Supervision
	ExtraArguments []string `json:"extra_arguments"`
}

type Tool struct {
	Name string `json:"name"`
}

type MemoryDefinition struct {
}
