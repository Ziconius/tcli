package main

import (
	"encoding/json"
	"log/slog"
	"os"
	"time"

	"tcli/src/utils"

	"gopkg.in/yaml.v2"
)

// This needs updated.
type Auth struct {
	TenantName string `yaml:"tenant_name"`
	TenantURL  string
	APIKey     string `yaml:"api_key"`
}

func (ac *Auth) LoadConfig() error {
	// UserHomeDir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Failed to find user home directory", "error", err)

		return err
	}
	configFileDir := homeDir + "/.tcli/auth"
	fb, err := utils.FileContents(configFileDir)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(fb, ac)
	if err != nil {
		return err
	}
	// Ensures all connections are directed to tines.com.
	ac.TenantURL = "https://" + ac.TenantName + ".tines.com/"

	return nil
}

// This is the data which will be stored on disk, this should limit sensitive data, and should be detailed.
type StoredConfig struct {
	Verbose     bool      // Define verbose logging
	LastUpdated time.Time // We will use this to pull outdated stories.
	Commands    []StoryConfig
}

// This is the data which is used to create an interactable story, or a tcli 'command'.
type StoryConfig struct {
	StoryID     int
	Description string

	CommandName string // This is the name used in the command i.e. `tcli get-users`, or `tcli analyse -ip=192.168.0.1`.
	URL         string
	Input       struct{}
	Request     []string
	// Raw output for now.
	// Output string
}

func (sc *StoredConfig) LoadConfig() error {
	// UserHomeDir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Failed to find user home directory", "error", err)

		return err
	}
	configFileDir := homeDir + "/.tcli/commandConfig"
	fb, err := utils.FileContents(configFileDir)
	if err != nil {
		return err
	}
	err = json.Unmarshal(fb, sc)
	if err != nil {
		return err
	}

	return nil
}

func (sc *StoredConfig) WriteConfig() error {
	jsonBytes, err := json.Marshal(sc)
	if err != nil {
		return err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Failed to find user home directory", "error", err)

		return err
	}
	// TODO: Temp rewrite for local testing
	configFileDir := homeDir + "/.tcli/commandConfig"
	err = utils.WriteBytes(configFileDir, jsonBytes)
	if err != nil {
		return err
	}

	return nil
}

// TODO:
// ValidateState returns a boolean representing a valid & unexpired cache.
// An error is only returned in the instance that an indivudal check has errored.
func (sc *StoredConfig) ValidateState() (bool, error) {
	return true, nil
}

// TODO:
func (sc *StoredConfig) ValidConfigCache() bool {
	if len(sc.Commands) < 1 {
		slog.Warn("No commands found in current cache.")

		return false
	}

	if !sc.LastUpdated.After(time.Now().Add(-72 * time.Hour)) {
		slog.Warn("The config cache is expired, & requires download.")

		return false
	}

	return true // TODO: Hard coding.
}
