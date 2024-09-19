package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	LocalDir  string `json:"local_dir"`
	RemoteDir string `json:"remote_dir"`
}

var AppConfig Config

// LoadConfig reads the configuration from a JSON file.
func LoadConfig() error {
	file, err := os.Open("config/config.json")
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&AppConfig)
	if err != nil {
		return err
	}
	return nil
}

// SaveConfig writes the configuration to a JSON file.
func SaveConfig() error {
	byteValue, err := json.MarshalIndent(AppConfig, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile("config/config.json", byteValue, 0644)
	if err != nil {
		return err
	}
	return nil
}
