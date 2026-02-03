package tts

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// PlayAudio plays an audio file asynchronously
func PlayAudio(audioPath string) error {
	// Verify file exists
	if _, err := os.Stat(audioPath); err != nil {
		return fmt.Errorf("audio file not found: %w", err)
	}

	// Detect OS and use appropriate player
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		// macOS: use afplay
		cmd = exec.Command("afplay", audioPath)
	case "linux":
		// Linux: try multiple players in order of preference
		players := []string{"mpg123", "ffplay", "aplay"}
		for _, player := range players {
			if _, err := exec.LookPath(player); err == nil {
				cmd = exec.Command(player, audioPath)
				break
			}
		}
		if cmd == nil {
			return fmt.Errorf("no audio player found (tried: %v)", players)
		}
	default:
		return fmt.Errorf("audio playback not supported on %s", runtime.GOOS)
	}

	// Run detached so we don't wait for playback to complete
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start audio player: %w", err)
	}

	// Don't wait - let it play in background
	return nil
}
