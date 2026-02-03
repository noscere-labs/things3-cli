package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/bear-cli/pkg/bear"
	"github.com/yourusername/bear-cli/pkg/formatter"
	"github.com/yourusername/bear-cli/pkg/tts"
	"github.com/yourusername/bear-cli/pkg/util"
)

// createCmd creates a new note in Bear
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new note in Bear",
	Long: `Create a new note with optional title, content, tags, and file attachments.

Examples:
  bear create --title "Meeting Notes" --content "Discussed Q1 roadmap" --tags "work,important"
  bear create --title "Project Plan" --file ~/Documents/plan.pdf --tags "projects"
  bear create --content "Quick note" --pin`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get command-line flags
		title, _ := cmd.Flags().GetString("title")
		content, _ := cmd.Flags().GetString("content")
		tagsStr, _ := cmd.Flags().GetString("tags")
		filePath, _ := cmd.Flags().GetString("file")
		pin, _ := cmd.Flags().GetBool("pin")
		timestamp, _ := cmd.Flags().GetBool("timestamp")

		// Validate input
		if title == "" && content == "" && filePath == "" {
			formatter.PrintError(
				"At least one of title, content, or file must be provided",
				"INVALID_ARGUMENTS",
				"",
			)
			return nil
		}

		// Parse tags from comma-separated string
		tags := util.ParseTags(tagsStr)

		// Create Bear client
		client, err := bear.NewClient()
		if err != nil {
			formatter.PrintError(
				"Failed to initialize Bear client",
				"CLIENT_ERROR",
				err.Error(),
			)
			return nil
		}

		// Create the note
		note, err := client.CreateNote(bear.CreateNoteOptions{
			Title:     title,
			Content:   content,
			Tags:      tags,
			FilePath:  filePath,
			Pin:       pin,
			Timestamp: timestamp,
		})

		if err != nil {
			formatter.PrintError(
				fmt.Sprintf("Failed to create note: %v", err),
				"BEAR_ERROR",
				err.Error(),
			)
			return nil
		}

		// Format and print success response
		formatter.PrintSuccess(note)
		return nil
	},
}

// readCmd reads an existing note from Bear
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read a note from Bear",
	Long: `Read and display a note by ID or title.

Examples:
  bear read --id "7E4B681B-..."
  bear read --title "Meeting Notes"
  bear read --title "Meeting Notes" --header "Action Items"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get command-line flags
		id, _ := cmd.Flags().GetString("id")
		title, _ := cmd.Flags().GetString("title")
		header, _ := cmd.Flags().GetString("header")
		excludeTrashed, _ := cmd.Flags().GetBool("exclude-trashed")

		// Validate that ID or Title is provided
		if id == "" && title == "" {
			formatter.PrintError(
				"Must provide either --id or --title",
				"INVALID_ARGUMENTS",
				"",
			)
			return nil
		}

		// Create Bear client
		client, err := bear.NewClient()
		if err != nil {
			formatter.PrintError(
				"Failed to initialize Bear client",
				"CLIENT_ERROR",
				err.Error(),
			)
			return nil
		}

		// Read the note
		note, err := client.ReadNote(bear.ReadNoteOptions{
			ID:             id,
			Title:          title,
			Header:         header,
			ExcludeTrashed: excludeTrashed,
		})

		if err != nil {
			formatter.PrintError(
				fmt.Sprintf("Failed to read note: %v", err),
				"BEAR_ERROR",
				err.Error(),
			)
			return nil
		}

		// Format and print success response
		formatter.PrintSuccess(note)
		return nil
	},
}

// updateCmd modifies an existing note
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing note in Bear",
	Long: `Modify a note by appending, prepending, or replacing content.

Modes:
  append      - Add content to the end (default)
  prepend     - Add content to the beginning
  replace     - Replace content but keep title
  replace_all - Replace entire note including title

Examples:
  bear update --id "7E4B681B-..." --content "New item" --mode append
  bear update --id "7E4B681B-..." --content "Replaced content" --mode replace_all
  bear update --id "7E4B681B-..." --file document.pdf --mode append`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get command-line flags
		id, _ := cmd.Flags().GetString("id")
		content, _ := cmd.Flags().GetString("content")
		filePath, _ := cmd.Flags().GetString("file")
		mode, _ := cmd.Flags().GetString("mode")
		header, _ := cmd.Flags().GetString("header")
		tagsStr, _ := cmd.Flags().GetString("tags")
		newLine, _ := cmd.Flags().GetBool("new-line")
		timestamp, _ := cmd.Flags().GetBool("timestamp")

		// Validate that ID is provided
		if id == "" {
			formatter.PrintError(
				"Note ID (--id) is required",
				"INVALID_ARGUMENTS",
				"",
			)
			return nil
		}

		// Validate that content or file is provided
		if content == "" && filePath == "" {
			formatter.PrintError(
				"Must provide either --content or --file",
				"INVALID_ARGUMENTS",
				"",
			)
			return nil
		}

		// Parse tags from comma-separated string
		tags := util.ParseTags(tagsStr)

		// Create Bear client
		client, err := bear.NewClient()
		if err != nil {
			formatter.PrintError(
				"Failed to initialize Bear client",
				"CLIENT_ERROR",
				err.Error(),
			)
			return nil
		}

		// Update the note
		note, err := client.UpdateNote(bear.UpdateNoteOptions{
			ID:        id,
			Content:   content,
			FilePath:  filePath,
			Mode:      mode,
			Header:    header,
			Tags:      tags,
			NewLine:   newLine,
			Timestamp: timestamp,
		})

		if err != nil {
			formatter.PrintError(
				fmt.Sprintf("Failed to update note: %v", err),
				"BEAR_ERROR",
				err.Error(),
			)
			return nil
		}

		// Format and print success response
		formatter.PrintSuccess(note)
		return nil
	},
}

// listCmd lists notes with optional filtering
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List notes from Bear with optional filtering",
	Long: `List notes from Bear. Can filter by tag, search term, or status.

Filters:
  --tag TAG         - Show notes with a specific tag
  --search TERM     - Search notes by content (requires API token)
  --filter TYPE     - Filter by type: all, untagged, todo, today, locked

Examples:
  bear list --tag "work"
  bear list --search "roadmap" --token "API_TOKEN"
  bear list --filter untagged
  bear list --filter todo`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get command-line flags
		tag, _ := cmd.Flags().GetString("tag")
		search, _ := cmd.Flags().GetString("search")
		// filter, _ := cmd.Flags().GetString("filter")
		token, _ := cmd.Flags().GetString("token")

		// If search is requested, require a token
		if search != "" && token == "" {
			// Try to get token from config
			var err error
			token, err = util.GetToken()
			if token == "" || err != nil {
				formatter.PrintError(
					"API token required for search operations",
					"INVALID_ARGUMENTS",
					"Provide with --token or set with 'bear config set-token'",
				)
				return nil
			}
		}

		// Create Bear client
		client, err := bear.NewClient()
		if err != nil {
			formatter.PrintError(
				"Failed to initialize Bear client",
				"CLIENT_ERROR",
				err.Error(),
			)
			return nil
		}

		// List notes based on filter type
		var result interface{}

		if search != "" {
			// Search operation
			resp, err := client.SearchNotes(bear.ListNotesOptions{
				Search: search,
				Token:  token,
			})
			if err != nil {
				formatter.PrintError(
					fmt.Sprintf("Failed to search notes: %v", err),
					"BEAR_ERROR",
					err.Error(),
				)
				return nil
			}
			result = resp
		} else if tag != "" {
			// Tag filter operation
			// Ensure we have a token for tag operations
			if token == "" {
				var err error
				token, err = util.GetToken()
				if token == "" || err != nil {
					formatter.PrintError(
						"API token required for tag list operations",
						"INVALID_ARGUMENTS",
						"Provide with --token or set with 'bear config set-token'",
					)
					return nil
				}
			}

			resp, err := client.ListNotesByTag(bear.ListNotesOptions{
				Tag:   tag,
				Token: token,
			})
			if err != nil {
				formatter.PrintError(
					fmt.Sprintf("Failed to list notes by tag: %v", err),
					"BEAR_ERROR",
					err.Error(),
				)
				return nil
			}
			result = resp
		} else {
			// Default to all notes
			formatter.PrintError(
				"List all notes not yet implemented without filters",
				"NOT_IMPLEMENTED",
				"Use --tag or --search to filter notes",
			)
			return nil
		}

		// Format and print success response
		formatter.PrintSuccess(result)
		return nil
	},
}

// archiveCmd archives (trashes) a note
var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archive (move to trash) a note in Bear",
	Long: `Archive a note by moving it to Bear's trash.

Examples:
  bear archive --id "7E4B681B-..."
  bear archive --id "7E4B681B-..." --no-window`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get command-line flags
		id, _ := cmd.Flags().GetString("id")
		noWindow, _ := cmd.Flags().GetBool("no-window")

		// Validate that ID is provided
		if id == "" {
			formatter.PrintError(
				"Note ID (--id) is required",
				"INVALID_ARGUMENTS",
				"",
			)
			return nil
		}

		// Create Bear client
		client, err := bear.NewClient()
		if err != nil {
			formatter.PrintError(
				"Failed to initialize Bear client",
				"CLIENT_ERROR",
				err.Error(),
			)
			return nil
		}

		// Archive the note
		err = client.ArchiveNote(bear.ArchiveNoteOptions{
			ID:       id,
			NoWindow: noWindow,
		})

		if err != nil {
			formatter.PrintError(
				fmt.Sprintf("Failed to archive note: %v", err),
				"BEAR_ERROR",
				err.Error(),
			)
			return nil
		}

		// Format and print success response
		formatter.PrintSuccess(map[string]interface{}{
			"note_id": id,
			"status":  "archived",
		})
		return nil
	},
}

// speakCmd converts a Bear note to speech
var speakCmd = &cobra.Command{
	Use:   "speak",
	Short: "Convert a Bear note to speech using TTS",
	Long: `Read a note from Bear and convert its content to speech using MURF AI TTS.

The note content will be retrieved, cleaned of code blocks and markdown,
and converted to natural-sounding audio. Audio files are saved to ~/.config/bear-cli/audio/

Examples:
  bear speak --id "7E4B681B-..."
  bear speak --title "Meeting Notes"
  bear speak --id "7E4B681B-..." --voice "en-UK-emma"
  bear speak --id "7E4B681B-..." --play
  bear speak --id "7E4B681B-..." --output ~/audio/meeting.mp3`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get command-line flags
		id, _ := cmd.Flags().GetString("id")
		title, _ := cmd.Flags().GetString("title")
		voice, _ := cmd.Flags().GetString("voice")
		output, _ := cmd.Flags().GetString("output")
		play, _ := cmd.Flags().GetBool("play")
		header, _ := cmd.Flags().GetString("header")

		// Validate that ID or Title is provided
		if id == "" && title == "" {
			formatter.PrintError(
				"Must provide either --id or --title",
				"INVALID_ARGUMENTS",
				"",
			)
			return nil
		}

		// Create Bear client
		bearClient, err := bear.NewClient()
		if err != nil {
			formatter.PrintError(
				"Failed to initialize Bear client",
				"CLIENT_ERROR",
				err.Error(),
			)
			return nil
		}

		// Read the note
		note, err := bearClient.ReadNote(bear.ReadNoteOptions{
			ID:     id,
			Title:  title,
			Header: header,
		})

		if err != nil {
			formatter.PrintError(
				fmt.Sprintf("Failed to read note: %v", err),
				"BEAR_ERROR",
				err.Error(),
			)
			return nil
		}

		// Create TTS client
		ttsClient, err := tts.NewClient()
		if err != nil {
			formatter.PrintError(
				"MURF TTS not configured",
				"MURF_NOT_CONFIGURED",
				"Set API key with: bear config set-murf --api-key YOUR_KEY",
			)
			return nil
		}

		// Generate speech
		options := tts.TTSOptions{
			Text:       note.Content,
			VoiceID:    voice,
			OutputPath: output,
			AutoPlay:   play,
		}

		result, err := ttsClient.GenerateSpeech(note.Content, options)
		if err != nil {
			formatter.PrintError(
				fmt.Sprintf("Failed to generate speech: %v", err),
				"TTS_ERROR",
				err.Error(),
			)
			return nil
		}

		// Check if generation was successful
		if !result.Success {
			formatter.PrintError(
				result.Error,
				result.ErrorCode,
				"",
			)
			return nil
		}

		// Format and print success response
		formatter.PrintSuccess(map[string]interface{}{
			"note_id":        note.ID,
			"note_title":     note.Title,
			"audio_path":     result.AudioPath,
			"text_length":    result.TextLength,
			"cleaned_length": result.CleanedLength,
			"format":         result.Format,
			"voice_id":       result.VoiceID,
			"auto_played":    result.AutoPlayed,
		})
		return nil
	},
}

// tagsCmd manages Bear tags
var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "Manage Bear tags",
	Long: `Work with Bear tags: list, rename, or delete tags.

Subcommands:
  list    - List all tags (requires API token)
  rename  - Rename a tag
  delete  - Delete a tag`,
}

// tagsListCmd lists all tags
var tagsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tags in Bear",
	Long: `Retrieve all tags used in Bear (requires API token).

Example:
  bear tags list --token "API_TOKEN"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get command-line flags
		token, _ := cmd.Flags().GetString("token")

		// If no token provided, try to load from config
		if token == "" {
			var err error
			token, err = util.GetToken()
			if token == "" || err != nil {
				formatter.PrintError(
					"API token required for tags operation",
					"INVALID_ARGUMENTS",
					"Provide with --token or set with 'bear config set-token'",
				)
				return nil
			}
		}

		// Create Bear client
		client, err := bear.NewClient()
		if err != nil {
			formatter.PrintError(
				"Failed to initialize Bear client",
				"CLIENT_ERROR",
				err.Error(),
			)
			return nil
		}

		// Get all tags
		result, err := client.GetAllTags(bear.TagsListOptions{
			Token: token,
		})

		if err != nil {
			formatter.PrintError(
				fmt.Sprintf("Failed to list tags: %v", err),
				"BEAR_ERROR",
				err.Error(),
			)
			return nil
		}

		// Format and print success response
		formatter.PrintSuccess(result)
		return nil
	},
}

// tagsRenameCmd renames a tag
var tagsRenameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename a tag",
	Long: `Rename a tag across all notes.

Example:
  bear tags rename --name "old-tag" --new-name "new-tag"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get command-line flags
		name, _ := cmd.Flags().GetString("name")
		newName, _ := cmd.Flags().GetString("new-name")

		// Validate parameters
		if name == "" || newName == "" {
			formatter.PrintError(
				"Both --name and --new-name are required",
				"INVALID_ARGUMENTS",
				"",
			)
			return nil
		}

		// Create Bear client
		client, err := bear.NewClient()
		if err != nil {
			formatter.PrintError(
				"Failed to initialize Bear client",
				"CLIENT_ERROR",
				err.Error(),
			)
			return nil
		}

		// Rename the tag
		err = client.RenameTag(bear.RenameTagOptions{
			Name:    name,
			NewName: newName,
		})

		if err != nil {
			formatter.PrintError(
				fmt.Sprintf("Failed to rename tag: %v", err),
				"BEAR_ERROR",
				err.Error(),
			)
			return nil
		}

		// Format and print success response
		formatter.PrintSuccess(map[string]interface{}{
			"old_name": name,
			"new_name": newName,
		})
		return nil
	},
}

// tagsDeleteCmd deletes a tag
var tagsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a tag",
	Long: `Delete a tag from all notes (does not delete notes, only removes the tag).

Example:
  bear tags delete --name "old-tag"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get command-line flags
		name, _ := cmd.Flags().GetString("name")

		// Validate parameters
		if name == "" {
			formatter.PrintError(
				"Tag name (--name) is required",
				"INVALID_ARGUMENTS",
				"",
			)
			return nil
		}

		// Create Bear client
		client, err := bear.NewClient()
		if err != nil {
			formatter.PrintError(
				"Failed to initialize Bear client",
				"CLIENT_ERROR",
				err.Error(),
			)
			return nil
		}

		// Delete the tag
		err = client.DeleteTag(bear.DeleteTagOptions{
			Name: name,
		})

		if err != nil {
			formatter.PrintError(
				fmt.Sprintf("Failed to delete tag: %v", err),
				"BEAR_ERROR",
				err.Error(),
			)
			return nil
		}

		// Format and print success response
		formatter.PrintSuccess(map[string]interface{}{
			"deleted_tag": name,
		})
		return nil
	},
}

// configCmd manages CLI configuration
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage bear CLI configuration",
	Long: `Manage configuration including API token storage.

Subcommands:
  set-token - Store API token for persistent use
  get-token - Display stored API token (masked)
  show      - Show current configuration`,
}

// configSetTokenCmd stores an API token
var configSetTokenCmd = &cobra.Command{
	Use:   "set-token",
	Short: "Store API token",
	Long: `Store your Bear API token for use with operations that require it.

Example:
  bear config set-token --token "123456-789ABC-DEF012"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get command-line flags
		token, _ := cmd.Flags().GetString("token")

		// Validate token is provided
		if token == "" {
			formatter.PrintError(
				"Token (--token) is required",
				"INVALID_ARGUMENTS",
				"",
			)
			return nil
		}

		// Save token to config
		err := util.SetToken(token)
		if err != nil {
			formatter.PrintError(
				fmt.Sprintf("Failed to save token: %v", err),
				"CONFIG_ERROR",
				err.Error(),
			)
			return nil
		}

		// Format and print success response
		formatter.PrintSuccess(map[string]interface{}{
			"status": "token saved",
		})
		return nil
	},
}

// configGetTokenCmd retrieves the stored token
var configGetTokenCmd = &cobra.Command{
	Use:   "get-token",
	Short: "Display stored API token",
	Long: `Retrieve the stored API token (displayed in masked form for security).

Example:
  bear config get-token`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load token from config
		token, err := util.GetToken()
		if err != nil || token == "" {
			formatter.PrintError(
				"No token configured",
				"CONFIG_ERROR",
				"Set one with 'bear config set-token --token YOUR_TOKEN'",
			)
			return nil
		}

		// Mask the token for display
		maskedToken := util.MaskToken(token)

		// Format and print success response
		formatter.PrintSuccess(map[string]interface{}{
			"token": maskedToken,
		})
		return nil
	},
}

// configShowCmd displays current configuration
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long: `Display the current configuration including token status and settings.

Example:
  bear config show`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load current configuration
		config, err := util.LoadConfig()
		if err != nil {
			formatter.PrintError(
				fmt.Sprintf("Failed to load config: %v", err),
				"CONFIG_ERROR",
				err.Error(),
			)
			return nil
		}

		// Get config file path
		configPath, _ := util.ConfigPath()

		// Prepare response
		response := map[string]interface{}{
			"token_set":     config.Token != "",
			"token":         util.MaskToken(config.Token),
			"callback_port": config.CallbackPort,
			"timeout_sec":   config.CallbackTimeoutSeconds,
			"show_window":   config.ShowWindow,
			"output_format": config.OutputFormat,
			"config_path":   configPath,
			"last_updated":  config.LastUpdated,
		}

		// Format and print success response
		formatter.PrintSuccess(response)
		return nil
	},
}

// configSetMurfCmd configures MURF TTS settings
var configSetMurfCmd = &cobra.Command{
	Use:   "set-murf",
	Short: "Configure MURF TTS settings",
	Long: `Set MURF API key and TTS preferences.

Configuration is saved to ~/.config/bear-cli/config.json

Examples:
  bear config set-murf --api-key "your-api-key"
  bear config set-murf --api-key "your-api-key" --voice "en-US-sara" --auto-play`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flags
		apiKey, _ := cmd.Flags().GetString("api-key")
		voice, _ := cmd.Flags().GetString("voice")
		format, _ := cmd.Flags().GetString("format")
		sampleRate, _ := cmd.Flags().GetInt("sample-rate")
		outputDir, _ := cmd.Flags().GetString("output-dir")
		autoPlay, _ := cmd.Flags().GetBool("auto-play")

		// Validate at least one setting is provided
		if apiKey == "" && voice == "" && format == "" && sampleRate == 0 && outputDir == "" && !autoPlay {
			formatter.PrintError(
				"At least one setting must be provided",
				"INVALID_ARGUMENTS",
				"Use --api-key, --voice, --format, --sample-rate, --output-dir, or --auto-play",
			)
			return nil
		}

		// Save config
		err := util.SetMurfConfig(apiKey, voice, format, sampleRate, outputDir, autoPlay)
		if err != nil {
			formatter.PrintError(
				fmt.Sprintf("Failed to save MURF config: %v", err),
				"CONFIG_ERROR",
				err.Error(),
			)
			return nil
		}

		// Format and print success response
		response := map[string]interface{}{
			"status": "MURF configuration updated",
		}
		if apiKey != "" {
			response["api_key"] = util.MaskAPIKey(apiKey)
		}
		if voice != "" {
			response["voice_id"] = voice
		}
		if format != "" {
			response["format"] = format
		}
		if sampleRate > 0 {
			response["sample_rate"] = sampleRate
		}
		if outputDir != "" {
			response["output_dir"] = outputDir
		}
		response["auto_play"] = autoPlay

		formatter.PrintSuccess(response)
		return nil
	},
}

// configShowMurfCmd displays MURF TTS configuration
var configShowMurfCmd = &cobra.Command{
	Use:   "show-murf",
	Short: "Display MURF TTS configuration",
	Long: `Show current MURF TTS configuration from ~/.config/bear-cli/config.json

Example:
  bear config show-murf`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load current configuration
		config, err := util.LoadConfig()
		if err != nil {
			formatter.PrintError(
				fmt.Sprintf("Failed to load config: %v", err),
				"CONFIG_ERROR",
				err.Error(),
			)
			return nil
		}

		// Check if any MURF config is set
		isConfigured := config.MurfAPIKey != ""

		// Prepare response with masked API key
		response := map[string]interface{}{
			"configured":    isConfigured,
			"voice_id":      config.MurfVoiceID,
			"format":        config.MurfFormat,
			"sample_rate":   config.MurfSampleRate,
			"output_dir":    config.MurfOutputDir,
			"auto_play":     config.MurfAutoPlay,
			"enabled":       config.MurfEnabled,
		}

		if isConfigured {
			response["api_key"] = util.MaskAPIKey(config.MurfAPIKey)
		} else {
			response["api_key"] = "not configured"
		}

		formatter.PrintSuccess(response)
		return nil
	},
}

// init sets up all commands and their flags
func init() {
	// Create command flags
	createCmd.Flags().StringP("title", "t", "", "Note title")
	createCmd.Flags().StringP("content", "c", "", "Note content")
	createCmd.Flags().StringP("tags", "g", "", "Comma-separated tags (e.g., 'work,urgent')")
	createCmd.Flags().StringP("file", "f", "", "File path to attach to note")
	createCmd.Flags().BoolP("pin", "p", false, "Pin note to top of list")
	createCmd.Flags().Bool("timestamp", false, "Prepend current date/time to content")

	// Read command flags
	readCmd.Flags().StringP("id", "i", "", "Note ID")
	readCmd.Flags().StringP("title", "t", "", "Note title (for lookup)")
	readCmd.Flags().StringP("header", "e", "", "Extract specific header section")
	readCmd.Flags().Bool("exclude-trashed", false, "Skip trashed notes")

	// Update command flags
	updateCmd.Flags().StringP("id", "i", "", "Note ID (required)")
	updateCmd.Flags().StringP("content", "c", "", "Content to add/update")
	updateCmd.Flags().StringP("file", "f", "", "File path to attach")
	updateCmd.Flags().StringP("mode", "m", "append", "Update mode: append, prepend, replace, replace_all")
	updateCmd.Flags().StringP("header", "e", "", "Target specific header section")
	updateCmd.Flags().StringP("tags", "g", "", "Comma-separated tags to add/update")
	updateCmd.Flags().Bool("new-line", false, "Add content on new line (append mode only)")
	updateCmd.Flags().Bool("timestamp", false, "Prepend date/time to added content")

	// List command flags
	listCmd.Flags().StringP("tag", "t", "", "Filter by tag")
	listCmd.Flags().StringP("search", "s", "", "Search notes by term (requires token)")
	listCmd.Flags().StringP("filter", "f", "", "Filter type: all, untagged, todo, today, locked")
	listCmd.Flags().StringP("token", "k", "", "API token (or use config)")

	// Archive command flags
	archiveCmd.Flags().StringP("id", "i", "", "Note ID (required)")
	archiveCmd.Flags().Bool("no-window", false, "Don't show Bear window")

	// Speak command flags
	speakCmd.Flags().StringP("id", "i", "", "Note ID")
	speakCmd.Flags().StringP("title", "t", "", "Note title (for lookup)")
	speakCmd.Flags().StringP("voice", "v", "", "Override voice ID")
	speakCmd.Flags().StringP("output", "o", "", "Custom output path for audio file")
	speakCmd.Flags().BoolP("play", "p", false, "Auto-play audio after generation")
	speakCmd.Flags().StringP("header", "e", "", "Extract specific header section")

	// Tags list command flags
	tagsListCmd.Flags().StringP("token", "k", "", "API token (or use config)")

	// Tags rename command flags
	tagsRenameCmd.Flags().StringP("name", "n", "", "Current tag name (required)")
	tagsRenameCmd.Flags().StringP("new-name", "w", "", "New tag name (required)")

	// Tags delete command flags
	tagsDeleteCmd.Flags().StringP("name", "n", "", "Tag name to delete (required)")

	// Config set-token command flags
	configSetTokenCmd.Flags().StringP("token", "k", "", "API token (required)")

	// Config set-murf command flags
	configSetMurfCmd.Flags().StringP("api-key", "k", "", "MURF API key")
	configSetMurfCmd.Flags().StringP("voice", "v", "", "Voice ID (e.g., en-UK-mason)")
	configSetMurfCmd.Flags().StringP("format", "f", "", "Audio format (MP3, WAV, FLAC, OGG)")
	configSetMurfCmd.Flags().IntP("sample-rate", "r", 0, "Sample rate in Hz (8000, 16000, 22050, 24000, 44100, 48000)")
	configSetMurfCmd.Flags().StringP("output-dir", "d", "", "Output directory for audio files")
	configSetMurfCmd.Flags().BoolP("auto-play", "a", false, "Auto-play audio after generation")

	// Add subcommands to tags command
	tagsCmd.AddCommand(tagsListCmd)
	tagsCmd.AddCommand(tagsRenameCmd)
	tagsCmd.AddCommand(tagsDeleteCmd)

	// Add subcommands to config command
	configCmd.AddCommand(configSetTokenCmd)
	configCmd.AddCommand(configGetTokenCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetMurfCmd)
	configCmd.AddCommand(configShowMurfCmd)
}

// GetCommands returns all available commands for the root command
func GetCommands() []*cobra.Command {
	return []*cobra.Command{
		createCmd,
		readCmd,
		updateCmd,
		listCmd,
		archiveCmd,
		speakCmd,
		tagsCmd,
		configCmd,
	}
}
