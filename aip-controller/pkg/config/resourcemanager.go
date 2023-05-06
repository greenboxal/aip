package config

import (
	"os"
	"path"
)

type ResourceConfig struct {
	DataDirectory    string `json:"data_dir"`
	PrivateDirectory string `json:"private_dir"`
	ProcDirectory    string `json:"proc_dir"`
}

func (rc *ResourceConfig) SetDefaults() {
	rc.DataDirectory = "/tmp/aip/data"
	rc.PrivateDirectory = "/tmp/aip/private"
	rc.ProcDirectory = "/tmp/aip/data/tmp/proc"
}

type ResourceManager struct {
	config *ResourceConfig
}

func NewResourceManager(config *ResourceConfig) *ResourceManager {
	return &ResourceManager{
		config: config,
	}
}

func (rm *ResourceManager) GetProcDirectory(subPaths ...string) string {
	subPath := path.Join(subPaths...)
	p := path.Join(rm.config.ProcDirectory, subPath)

	if err := os.MkdirAll(p, 0750); err != nil {
		panic(err)
	}

	return p
}

func (rm *ResourceManager) GetDataDirectory(subPaths ...string) string {
	subPath := path.Join(subPaths...)
	p := path.Join(rm.config.DataDirectory, subPath)

	if err := os.MkdirAll(p, 0750); err != nil {
		panic(err)
	}

	return p
}

func (rm *ResourceManager) GetPrivateDirectory(subPaths ...string) string {
	subPath := path.Join(subPaths...)
	p := path.Join(rm.config.PrivateDirectory, subPath)

	if err := os.MkdirAll(p, 0700); err != nil {
		panic(err)
	}

	return p
}
