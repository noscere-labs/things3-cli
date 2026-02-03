package util

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Config represents the bear-cli configuration stored in ~/.config/bear-cli/config.json
type Config struct {
	Token                  string    `json:"token"`
	CallbackPort           int       `json:"callback_port"`
	CallbackTimeoutSeconds int       `json:"callback_timeout_seconds"`
	ShowWindow             bool      `json:"show_window"`
	OutputFormat           string    `json:"output_format"`
	LastUpdated            time.Time `json:"last_updated"`
	MurfAPIKey             string    `json:"murf_api_key"`
	MurfVoiceID            string    `json:"murf_voice_id"`
	MurfFormat             string    `json:"murf_format"`
	MurfSampleRate         int       `json:"murf_sample_rate"`
	MurfOutputDir          string    `json:"murf_output_dir"`
	MurfAutoPlay           bool      `json:"murf_auto_play"`
	MurfEnabled            bool      `json:"murf_enabled"`
}

// DefaultConfig returns a Config with sensible defaults
func DefaultConfig() Config {
	return Config{
		CallbackPort:           8765,
		CallbackTimeoutSeconds: 10,
		ShowWindow:             false,
		OutputFormat:           "json",
		Token:                  "",
		LastUpdated:            time.Now(),
	}
}

// ConfigPath returns the path to the config file (~/.config/bear-cli/config.json)
func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "bear-cli", "config.json"), nil
}

// EnsureConfigDir creates the ~/.config/bear-cli directory if it doesn't exist
func EnsureConfigDir() error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Also ensure audio directory exists
	audioDir := filepath.Join(dir, "audio")
	if err := os.MkdirAll(audioDir, 0755); err != nil {
		return fmt.Errorf("failed to create audio directory: %w", err)
	}

	return nil
}

// LoadConfig reads and parses the config file, returning defaults if not found
func LoadConfig() (Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return Config{}, err
	}

	// If config doesn't exist, return defaults
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	// Read and unmarshal the config file
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
	// Update the timestamp when saving
	config.LastUpdated = time.Now()

	// Ensure the config directory exists
	if err := EnsureConfigDir(); err != nil {
		return err
	}

	path, err := ConfigPath()
	if err != nil {
		return err
	}

	// Marshal to JSON with indentation for readability
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file with restricted permissions (user-readable only)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetToken retrieves the stored API token
// Checks environment variable first, then config file
func GetToken() (string, error) {
	// Environment variable takes precedence
	if token := os.Getenv("BEAR_TOKEN"); token != "" {
		return token, nil
	}

	config, err := LoadConfig()
	if err != nil {
		return "", err
	}

	return config.Token, nil
}

// SetToken stores the API token in the config file
func SetToken(token string) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	config.Token = token
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

// LoadEnvFile loads environment variables from ~/.config/bear-cli/.env if it exists
func LoadEnvFile() (map[string]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return map[string]string{}, nil // Return empty map if can't get home dir
	}

	envPath := filepath.Join(home, ".config", "bear-cli", ".env")
	envVars := make(map[string]string)

	// If file doesn't exist, that's okay - just return empty map
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		return envVars, nil
	}

	// Read the file
	data, err := os.ReadFile(envPath)
	if err != nil {
		return map[string]string{}, fmt.Errorf("failed to read .env file: %w", err)
	}

	// Parse KEY=value format
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Remove surrounding quotes if present
			value = strings.Trim(value, "\"'")
			envVars[key] = value
		}
	}

	return envVars, nil
}

// GetMurfConfig loads MURF configuration from multiple sources with priority:
// 1. .env file (~/.config/bear-cli/.env)
// 2. Environment variables (MURF_*)
// 3. Config file (~/.config/bear-cli/config.json)
// 4. Defaults
func GetMurfConfig() (map[string]string, error) {
	config := make(map[string]string)

	// Load defaults
	config["voice_id"] = "en-UK-mason"
	config["format"] = "MP3"
	config["sample_rate"] = "24000"
	config["enabled"] = "true"
	config["auto_play"] = "false"

	home, err := os.UserHomeDir()
	if err == nil {
		config["output_dir"] = filepath.Join(home, ".config", "bear-cli", "audio")
	}

	// Load from .env file (highest priority)
	envFileVars, err := LoadEnvFile()
	if err == nil {
		for key, value := range envFileVars {
			// Map environment variable names to config keys
			switch key {
			case "MURF_API_KEY":
				config["api_key"] = value
			case "MURF_VOICE_ID":
				config["voice_id"] = value
			case "MURF_FORMAT":
				config["format"] = value
			case "MURF_SAMPLE_RATE":
				config["sample_rate"] = value
			case "MURF_OUTPUT_DIR":
				config["output_dir"] = value
			case "MURF_AUTO_PLAY":
				config["auto_play"] = value
			case "MURF_ENABLED":
				config["enabled"] = value
			}
		}
	}

	// Load from environment variables (medium priority)
	if apiKey := os.Getenv("MURF_API_KEY"); apiKey != "" {
		config["api_key"] = apiKey
	}
	if voiceID := os.Getenv("MURF_VOICE_ID"); voiceID != "" {
		config["voice_id"] = voiceID
	}
	if format := os.Getenv("MURF_FORMAT"); format != "" {
		config["format"] = format
	}
	if sampleRate := os.Getenv("MURF_SAMPLE_RATE"); sampleRate != "" {
		config["sample_rate"] = sampleRate
	}
	if outputDir := os.Getenv("MURF_OUTPUT_DIR"); outputDir != "" {
		config["output_dir"] = outputDir
	}
	if autoPlay := os.Getenv("MURF_AUTO_PLAY"); autoPlay != "" {
		config["auto_play"] = autoPlay
	}
	if enabled := os.Getenv("MURF_ENABLED"); enabled != "" {
		config["enabled"] = enabled
	}

	// Load from config file (low priority)
	cfg, err := LoadConfig()
	if err == nil {
		if cfg.MurfAPIKey != "" {
			config["api_key"] = cfg.MurfAPIKey
		}
		if cfg.MurfVoiceID != "" {
			config["voice_id"] = cfg.MurfVoiceID
		}
		if cfg.MurfFormat != "" {
			config["format"] = cfg.MurfFormat
		}
		if cfg.MurfSampleRate > 0 {
			config["sample_rate"] = fmt.Sprintf("%d", cfg.MurfSampleRate)
		}
		if cfg.MurfOutputDir != "" {
			config["output_dir"] = cfg.MurfOutputDir
		}
		if cfg.MurfEnabled {
			config["enabled"] = "true"
		}
		if cfg.MurfAutoPlay {
			config["auto_play"] = "true"
		}
	}

	return config, nil
}

// SetMurfConfig saves MURF settings to the config file
func SetMurfConfig(apiKey, voiceID, format string, sampleRate int, outputDir string, autoPlay bool) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	if apiKey != "" {
		config.MurfAPIKey = apiKey
	}
	if voiceID != "" {
		config.MurfVoiceID = voiceID
	}
	if format != "" {
		config.MurfFormat = format
	}
	if sampleRate > 0 {
		config.MurfSampleRate = sampleRate
	}
	if outputDir != "" {
		config.MurfOutputDir = outputDir
	}
	config.MurfAutoPlay = autoPlay
	config.MurfEnabled = true

	return SaveConfig(config)
}

// MaskAPIKey returns a masked version of the API key for display
func MaskAPIKey(apiKey string) string {
	if len(apiKey) <= 12 {
		return "***"
	}
	return apiKey[:6] + "***" + apiKey[len(apiKey)-6:]
}
