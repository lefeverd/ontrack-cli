package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

const (
	configFileName = ".ontrack-cli-config.yaml"
)

// Root configuration
type RootConfig struct {
	// Default configuration name
	Selected string
	// List of configurations
	Configurations []Config
}

// Configuration content
type Config struct {
	// Name of the configuration
	Name string
	// URL of the remote server
	URL string
	// Username for the remote server (when using basic authentication)
	Username string
	// Password for the remote server (when using basic authentication)
	Password string
	// Token for the remote server (when using token-based authentication)
	Token string
}

// Reads the configuration
func ReadRootConfiguration() RootConfig {
	var root RootConfig
	home, _ := homedir.Dir()
	configFilePath := filepath.Join(home, configFileName)

	// If the config file does not exist, returns an empty root config
	if _, err := os.Stat(configFilePath); err != nil {
		if os.IsNotExist(err) {
			return root
		}
	}

	reader, _ := os.Open(configFilePath)
	buf, _ := ioutil.ReadAll(reader)
	yaml.Unmarshal(buf, &root)
	return root
}

// Adds a new configuration and set as default
func AddConfiguration(config Config) {
	root := ReadRootConfiguration()
	// TODO Check if the configuration name already exists
	// Default selected configuration is the added one
	// Adds the configuration to the list
	newRoot := RootConfig{
		Selected:       config.Name,
		Configurations: append(root.Configurations, config),
	}
	// Saves the root configuration back
	home, _ := homedir.Dir()
	configFilePath := filepath.Join(home, configFileName)
	buf, _ := yaml.Marshal(newRoot)
	_, _ = os.OpenFile(configFilePath, os.O_CREATE|os.O_WRONLY, 0600)
	_ = ioutil.WriteFile(configFilePath, buf, 0600)
}
