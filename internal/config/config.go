package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds all user settings saved locally
type Config struct {
	PasswordHash  string `json:"password_hash,omitempty"`
	PasswordSet   bool   `json:"password_set"`
	DefaultPort   int    `json:"default_port"`
	DefaultFolder string `json:"default_folder,omitempty"`
	AllowUpload   bool   `json:"allow_upload"`
	AllowDelete   bool   `json:"allow_delete"`
	ShowHidden    bool   `json:"show_hidden"`
	MaxUploadMB   int    `json:"max_upload_mb"`
}

func DefaultConfig() *Config {
	return &Config{
		DefaultPort: 8080,
		AllowUpload: true,
		AllowDelete: false,
		ShowHidden:  false,
		MaxUploadMB: 500,
	}
}

// configPath returns path like ~/.sharedrop/config.json
func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".sharedrop")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, err
	}

	cfg := DefaultConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) Save() error {
	path, err := configPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".sharedrop"), nil
}
