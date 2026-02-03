package util

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Config represents the things3-cli configuration stored in ~/.config/things3-cli/config.json
type Config struct {
	AuthToken             string    `json:"auth_token"`
	CallbackPort          int       `json:"callback_port"`
	CallbackTimeoutSeconds int      `json:"callback_timeout_seconds"`
	OutputFormat          string    `json:"output_format"`
	LastUpdated           time.Time `json:"last_updated"`
}

// DefaultConfig returns a Config with sensible defaults
func DefaultConfig() Config {
	return Config{
		CallbackPort:          8765,
		CallbackTimeoutSeconds: 10,
		OutputFormat:          "json",
		AuthToken:             "",
		LastUpdated:           time.Now(),
	}
}

// ConfigPath returns the path to the config file (~/.config/things3-cli/config.json)
func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "things3-cli", "config.json"), nil
}

// EnsureConfigDir creates the ~/.config/things3-cli directory if it doesn't exist
func EnsureConfigDir() error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	return nil
}

// LoadConfig reads and parses the config file, returning defaults if not found
func LoadConfig() (Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return Config{}, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}

// SaveConfig writes the config to the config file
func SaveConfig(config Config) error {
	config.LastUpdated = time.Now()

	if err := EnsureConfigDir(); err != nil {
		return err
	}

	path, err := ConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetAuthToken retrieves the stored Things auth token.
// Checks environment variable first, then config file.
func GetAuthToken() (string, error) {
	if token := os.Getenv("THINGS_AUTH_TOKEN"); token != "" {
		return token, nil
	}

	config, err := LoadConfig()
	if err != nil {
		return "", err
	}

	return config.AuthToken, nil
}

// SetAuthToken stores the Things auth token in the config file
func SetAuthToken(token string) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	config.AuthToken = token
	return SaveConfig(config)
}

// MaskToken returns a masked version of the token for display
// Shows first 6 chars and last 6 chars, with *** in between
func MaskToken(token string) string {
	if len(token) <= 12 {
		return "***"
	}
	return token[:6] + "***" + token[len(token)-6:]
}
