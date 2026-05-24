package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	Database struct {
		File string `yaml:"File"`
	} `yaml:"database"`
	TaskManagers []ManagerConfig `yaml:"task_managers"`
}

type ManagerConfig struct {
	Name          string `yaml:"name"`
	DisplayName   string `yaml:"display_name"`
	ActivePath    string `yaml:"active_path"`
	CompletedPath string `yaml:"completed_path"`
}

func NewConfig() *Config {
	return &Config{}
}

func (cfg *Config) Load(configPath string, createIfNoExist bool) error {
	// Check config File
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		if createIfNoExist {
			err = nil
			conf := defaultConfig()
			err := conf.Save(configPath)
			if err != nil {
				return ErrCannotCreateConfig
			}
		} else {
			return ErrConfigNotFound
		}
	}
	if err != nil {
		return fmt.Errorf("cannot check config: %v", err)
	}

	// Read config File
	data, err := os.ReadFile(configPath)
	if err != nil {
		return ErrCannotReadConfig
	}

	// Pase config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return ErrCannotParseConfig
	}

	return nil
}

func (cfg *Config) Save(configPath string) error {
	// Make yaml config
	configData, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	// Check config dir
	dir := filepath.Dir(configPath)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return ErrCannotCreateDir
	}

	// Write config to File
	err = os.WriteFile(configPath, configData, 0644)
	if err != nil {
		return ErrCannotCreateConfig
	}

	return nil
}
