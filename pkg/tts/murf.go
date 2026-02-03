package tts

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// FindMurfScript locates the MURF CLI script using priority order
func FindMurfScript() (string, error) {
	// Priority order:
	// 1. MURF_CLI_SCRIPT environment variable
	if scriptPath := os.Getenv("MURF_CLI_SCRIPT"); scriptPath != "" {
		if _, err := os.Stat(scriptPath); err == nil {
			return scriptPath, nil
		}
	}

	// 2. ~/.config/bear-cli/murf-cli.js (global install)
	home, err := os.UserHomeDir()
	if err == nil {
		globalPath := filepath.Join(home, ".config", "bear-cli", "murf-cli.js")
		if _, err := os.Stat(globalPath); err == nil {
			return globalPath, nil
		}
	}

	// 3. ./murf/murf-cli.js (development, relative to current dir)
	if _, err := os.Stat("./murf/murf-cli.js"); err == nil {
		return "./murf/murf-cli.js", nil
	}

	return "", fmt.Errorf("MURF CLI script not found - install murf/murf-cli.js to ~/.config/bear-cli/murf-cli.js")
}

// MurfRequest represents the request to the MURF API
type MurfRequest struct {
	VoiceID    string `json:"voiceId"`
	Text       string `json:"text"`
	Format     string `json:"format"`
	SampleRate int    `json:"sampleRate"`
	Speed      int    `json:"speed"`
	Pitch      int    `json:"pitch"`
}

// ExecuteMurfScript calls the MURF TTS CLI script and returns the audio file path
func ExecuteMurfScript(text string, config *TTSConfig, outputPath string) (string, error) {
	// Find the script
	scriptPath, err := FindMurfScript()
	if err != nil {
		return "", fmt.Errorf("TTS_SCRIPT_NOT_FOUND: %w", err)
	}

	// Create MURF request
	request := MurfRequest{
		VoiceID:    config.VoiceID,
		Text:       text,
		Format:     config.Format,
		SampleRate: config.SampleRate,
	}

	// Marshal request to JSON
	requestData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Prepare environment variables for the script
	env := os.Environ()
	env = append(env, fmt.Sprintf("MURF_API_KEY=%s", config.APIKey))
	env = append(env, fmt.Sprintf("MURF_VOICE_ID=%s", config.VoiceID))
	env = append(env, fmt.Sprintf("MURF_FORMAT=%s", config.Format))
	env = append(env, fmt.Sprintf("MURF_SAMPLE_RATE=%d", config.SampleRate))
	env = append(env, fmt.Sprintf("MURF_OUTPUT_DIR=%s", config.OutputDir))

	// Execute Node.js script with the request JSON on stdin
	cmd := exec.Command("node", scriptPath)
	cmd.Env = env
	cmd.Stdin = strings.NewReader(string(requestData))

	// Capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		errOutput := strings.TrimSpace(string(output))
		return "", fmt.Errorf("TTS_GENERATION_FAILED: %s", errOutput)
	}

	// Parse output to get audio file path
	// The murf-cli.js script outputs the file path on success
	audioPath := strings.TrimSpace(string(output))

	// Filter out log messages - look for actual file path
	lines := strings.Split(audioPath, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip log lines and error messages
		if strings.HasPrefix(line, "[") || strings.HasPrefix(line, "ERROR") {
			continue
		}
		// Check if it looks like a file path
		if strings.HasPrefix(line, "/") || strings.HasPrefix(line, "~") {
			// Expand ~ if present
			if strings.HasPrefix(line, "~") {
				home, _ := os.UserHomeDir()
				line = filepath.Join(home, line[1:])
			}
			// Verify file exists
			if _, err := os.Stat(line); err == nil {
				return line, nil
			}
		}
	}

	return "", fmt.Errorf("TTS_GENERATION_FAILED: no audio file generated")
}
