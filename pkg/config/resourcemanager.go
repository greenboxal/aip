package config

import (
	"os"
	"path"
)

type ResourceManager struct {
	cm *ConfigManager

	dataDirectory    string
	privateDirectory string
}

func NewResourceManager(cm *ConfigManager) *ResourceManager {
	return &ResourceManager{
		cm: cm,

		// TODO: Read from ConfigManager
		dataDirectory:    "./data",
		privateDirectory: "./private",
	}
}

func (rm *ResourceManager) GetDataDirectory(subPath string) string {
	p := path.Join(rm.dataDirectory, subPath)

	if err := os.MkdirAll(p, 0750); err != nil {
		panic(err)
	}

	return p
}

func (rm *ResourceManager) GetPrivateDirectory(subPath string) string {
	p := path.Join(rm.privateDirectory, subPath)

	if err := os.MkdirAll(p, 0700); err != nil {
		panic(err)
	}

	return p
}
