// Package utils handles configuration
// and API token retrieval
package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	APIToken string `json:"api_token"`
}

var configFileName = ".token.json"

func getConfigFilePath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, configFileName), nil
}

func LoadConfig() (*Config, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
		return  nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func SaveToken(token string) error {
	cfg := Config{
		APIToken: token,
	}
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0600)
}
