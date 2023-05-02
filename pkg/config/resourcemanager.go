package config

import (
	"os"
	"path"
)

type ResourceManager struct {
	cm *ConfigManager

	dataDirectory    string
	privateDirectory string
	procDirectory    string
}

func NewResourceManager(cm *ConfigManager) *ResourceManager {
	return &ResourceManager{
		cm: cm,

		// TODO: Read from ConfigManager
		dataDirectory:    "./data",
		privateDirectory: "./private",
		procDirectory:    "./data/tmp/proc",
	}
}

func (rm *ResourceManager) GetProcDirectory(subPaths ...string) string {
	subPath := path.Join(subPaths...)
	p := path.Join(rm.procDirectory, subPath)

	if err := os.MkdirAll(p, 0750); err != nil {
		panic(err)
	}

	return p
}

func (rm *ResourceManager) GetDataDirectory(subPaths ...string) string {
	subPath := path.Join(subPaths...)
	p := path.Join(rm.dataDirectory, subPath)

	if err := os.MkdirAll(p, 0750); err != nil {
		panic(err)
	}

	return p
}

func (rm *ResourceManager) GetPrivateDirectory(subPaths ...string) string {
	subPath := path.Join(subPaths...)
	p := path.Join(rm.privateDirectory, subPath)

	if err := os.MkdirAll(p, 0700); err != nil {
		panic(err)
	}

	return p
}
