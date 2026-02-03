package tts

// TTSConfig represents MURF TTS configuration
type TTSConfig struct {
	APIKey      string
	VoiceID     string
	Format      string // MP3, WAV, FLAC, OGG
	SampleRate  int
	OutputDir   string
	AutoPlay    bool
	Enabled     bool
	MaxLength   int
	MinLength   int
}

// TTSOptions represents options for a single TTS generation request
type TTSOptions struct {
	Text       string // The text to convert to speech
	VoiceID    string // Optional override for voice ID
	OutputPath string // Optional custom output path
	AutoPlay   bool   // Whether to auto-play the audio
}

// TTSResult represents the result of a TTS generation
type TTSResult struct {
	Success       bool   // Whether generation was successful
	AudioPath     string // Path to generated audio file
	TextLength    int    // Original text length in characters
	CleanedLength int    // Cleaned text length in characters
	Format        string // Audio format used
	VoiceID       string // Voice ID used
	AutoPlayed    bool   // Whether audio was auto-played
	Error         string // Error message if generation failed
	ErrorCode     string // Error code for machine parsing
}
