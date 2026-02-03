package tts

import (
	"fmt"
)

// Client manages TTS operations
type Client struct {
	config *TTSConfig
}

// NewClient creates a new TTS client
func NewClient() (*Client, error) {
	config, err := NewTTSConfig()
	if err != nil {
		return nil, err
	}

	if !config.Enabled {
		return nil, fmt.Errorf("MURF TTS is not enabled")
	}

	if config.APIKey == "" {
		return nil, fmt.Errorf("MURF_NOT_CONFIGURED: MURF API key not set. Configure with: bear config set-murf --api-key YOUR_KEY")
	}

	return &Client{
		config: config,
	}, nil
}

// GenerateSpeech converts text to speech and returns the result
func (c *Client) GenerateSpeech(text string, options TTSOptions) (*TTSResult, error) {
	result := &TTSResult{
		TextLength: len(text),
	}

	// Clean text for TTS
	cleanedText := CleanTextForTTS(text)
	result.CleanedLength = len(cleanedText)

	// Validate text length
	if !ValidateTextLength(cleanedText, c.config.MinLength, c.config.MaxLength) {
		if len(cleanedText) < c.config.MinLength {
			result.Error = fmt.Sprintf("Text too short (minimum %d characters)", c.config.MinLength)
			result.ErrorCode = "INVALID_TEXT_LENGTH"
			return result, nil
		}
		if len(cleanedText) > c.config.MaxLength {
			cleanedText = TruncateText(cleanedText, c.config.MaxLength)
			result.CleanedLength = len(cleanedText)
		}
	}

	// Determine voice ID (option overrides config)
	voiceID := c.config.VoiceID
	if options.VoiceID != "" {
		voiceID = options.VoiceID
	}
	result.VoiceID = voiceID

	// Determine output path
	outputPath := options.OutputPath
	if outputPath == "" {
		var err error
		format := c.config.Format
		outputPath, err = c.config.GenerateOutputPath(format)
		if err != nil {
			result.Error = fmt.Sprintf("Failed to generate output path: %v", err)
			result.ErrorCode = "OUTPUT_PATH_ERROR"
			return result, nil
		}
	}

	result.Format = c.config.Format
	result.AudioPath = outputPath

	// Execute MURF script
	audioPath, err := ExecuteMurfScript(cleanedText, c.config, outputPath)
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		// Parse error code from error message
		if fmt.Sprintf("%v", err) != "" {
			result.ErrorCode = "TTS_GENERATION_FAILED"
		}
		return result, nil
	}

	result.AudioPath = audioPath
	result.Success = true

	// Optionally auto-play
	if options.AutoPlay || c.config.AutoPlay {
		if err := PlayAudio(audioPath); err == nil {
			result.AutoPlayed = true
		}
	}

	return result, nil
}
