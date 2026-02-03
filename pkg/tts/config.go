package tts

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/yourusername/bear-cli/pkg/util"
)

// NewTTSConfig creates a TTSConfig from environment variables and config files
func NewTTSConfig() (*TTSConfig, error) {
	// Load configuration with priority order
	murfCfg, err := util.GetMurfConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load MURF config: %w", err)
	}

	config := &TTSConfig{
		APIKey:    murfCfg["api_key"],
		VoiceID:   murfCfg["voice_id"],
		Format:    murfCfg["format"],
		OutputDir: murfCfg["output_dir"],
		MaxLength: 5000,
		MinLength: 10,
	}

	// Parse sample rate
	sampleRateStr := murfCfg["sample_rate"]
	if sampleRateStr != "" {
		sampleRate, err := strconv.Atoi(sampleRateStr)
		if err == nil {
			config.SampleRate = sampleRate
		}
	}

	// Parse enabled flag
	if enabledStr := murfCfg["enabled"]; enabledStr != "" {
		config.Enabled = strings.ToLower(enabledStr) != "false"
	}

	// Parse auto-play flag
	if autoPlayStr := murfCfg["auto_play"]; autoPlayStr != "" {
		config.AutoPlay = strings.ToLower(autoPlayStr) == "true"
	}

	// Expand ~ in output directory
	if strings.HasPrefix(config.OutputDir, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			config.OutputDir = filepath.Join(home, config.OutputDir[1:])
		}
	}

	// Validate config
	if !config.Enabled {
		return config, nil
	}

	if config.APIKey == "" {
		return nil, fmt.Errorf("MURF_API_KEY not configured")
	}

	return config, nil
}

// GenerateOutputPath creates a unique output path for audio file
func (c *TTSConfig) GenerateOutputPath(format string) (string, error) {
	// Ensure output directory exists
	if err := os.MkdirAll(c.OutputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate filename with timestamp
	timestamp := getCurrentTimestamp()
	format = strings.ToLower(format)
	if format == "" {
		format = strings.ToLower(c.Format)
	}
	filename := fmt.Sprintf("bear-tts-%s.%s", timestamp, format)

	return filepath.Join(c.OutputDir, filename), nil
}

// getCurrentTimestamp returns a timestamp suitable for filename
func getCurrentTimestamp() string {
	return time.Now().Format("20060102-150405")
}
