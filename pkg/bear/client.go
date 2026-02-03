package bear

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/yourusername/bear-cli/pkg/util"
)

// Client handles all communication with Bear via x-callback-url scheme
// It builds URLs, executes them via macOS URL handlers, and captures responses
type Client struct {
	Token          string          // API token for advanced operations
	ShowWindow     bool            // Whether to show Bear window when executing commands
	CallbackPort   int             // Port for callback server
	CallbackServer *CallbackServer // The callback server instance
	timeout        time.Duration   // Timeout for waiting for responses
}

// NewClient creates a new Bear client with default settings
func NewClient() (*Client, error) {
	// Try to load token from config/environment
	token, err := util.GetToken()
	if err != nil {
		token = "" // It's OK if token is not configured yet
	}

	config, err := util.LoadConfig()
	if err != nil {
		config = util.DefaultConfig()
	}

	return &Client{
		Token:        token,
		ShowWindow:   config.ShowWindow,
		CallbackPort: config.CallbackPort,
		timeout:      time.Duration(config.CallbackTimeoutSeconds) * time.Second,
	}, nil
}

// buildBearURL constructs an x-callback-url for Bear
// action: The Bear action (e.g., "create", "open-note")
// params: Map of parameters to include in the URL
func (c *Client) buildBearURL(action string, params map[string]string) string {
	// Create base URL with the action
	baseURL := fmt.Sprintf("bear://x-callback-url/%s", action)

	// Add callback URL so Bear knows where to send the response
	callbackURL := fmt.Sprintf("http://localhost:%d/callback", c.CallbackPort)
	params["x-success"] = callbackURL
	params["x-error"] = callbackURL

	// Add token if we have one
	if c.Token != "" {
		params["token"] = c.Token
	}

	// Add x-window parameter to control window visibility
	if !c.ShowWindow {
		params["x-window"] = "false"
	}

	// Build query string
	queryStr := util.EncodeParams(params)

	return baseURL + "?" + queryStr
}

// executeURL opens a Bear URL and waits for the response
// This function:
// 1. Starts the callback server
// 2. Opens the URL via macOS `open` command
// 3. Waits for Bear to call back
// 4. Stops the server and returns the response
func (c *Client) executeURL(bearURL string) (map[string]string, error) {
	// Create and start callback server
	c.CallbackServer = NewCallbackServer(c.CallbackPort)
	if err := c.CallbackServer.Start(); err != nil {
		return nil, fmt.Errorf("failed to start callback server: %w", err)
	}
	defer c.CallbackServer.Stop()

	// Execute the Bear URL using macOS open command
	// This tells macOS to open the URL with the Bear app
	cmd := exec.Command("open", bearURL)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to execute URL: %w", err)
	}

	// Wait for Bear to respond via our callback server
	response, err := c.CallbackServer.WaitForResponse(c.timeout)
	if err != nil {
		return nil, err
	}

	// Check if Bear returned an error
	if errMsg, exists := response["x-error"]; exists {
		return response, fmt.Errorf("bear error: %s", errMsg)
	}

	return response, nil
}

// CreateNote creates a new note in Bear
func (c *Client) CreateNote(opts CreateNoteOptions) (*Note, error) {
	params := make(map[string]string)

	// Add title if provided
	if opts.Title != "" {
		params["title"] = opts.Title
	}

	// Add content if provided
	if opts.Content != "" {
		content := opts.Content
		if opts.Timestamp {
			content = util.GetTimestamp() + "\n" + content
		}
		params["text"] = content
	}

	// Add tags if provided
	if len(opts.Tags) > 0 {
		params["tags"] = util.JoinTags(opts.Tags)
	}

	// Handle file attachment if provided
	if opts.FilePath != "" {
		expandedPath, err := util.ExpandHomePath(opts.FilePath)
		if err != nil {
			return nil, fmt.Errorf("invalid file path: %w", err)
		}

		fileContent, err := util.ReadFileAsBase64(expandedPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}

		fileName := filepath.Base(expandedPath)
		params["file"] = fileContent
		params["filename"] = fileName
	}

	// Add optional flags
	if opts.Pin {
		params["pin"] = "yes"
	}

	// Execute the create URL
	response, err := c.executeURL(c.buildBearURL("create", params))
	if err != nil {
		return nil, err
	}

	// Parse the response into a Note
	note := &Note{
		ID:    response["identifier"],
		Title: response["title"],
	}

	// Parse timestamps if available
	if createdStr := response["creationDate"]; createdStr != "" {
		if parsed, err := time.Parse(time.RFC3339, createdStr); err == nil {
			note.CreatedAt = parsed
		}
	}

	return note, nil
}

// ReadNote reads an existing note from Bear by ID or title
func (c *Client) ReadNote(opts ReadNoteOptions) (*Note, error) {
	if opts.ID == "" && opts.Title == "" {
		return nil, fmt.Errorf("must provide either ID or Title")
	}

	params := make(map[string]string)

	// Prefer ID if both provided
	if opts.ID != "" {
		params["id"] = opts.ID
	} else {
		params["title"] = opts.Title
	}

	// Add optional header parameter to extract specific section
	if opts.Header != "" {
		params["header"] = opts.Header
	}

	// Add exclude-trashed flag if requested
	if opts.ExcludeTrashed {
		params["exclude-trashed"] = "yes"
	}

	// Execute the open-note URL
	response, err := c.executeURL(c.buildBearURL("open-note", params))
	if err != nil {
		return nil, err
	}

	// Parse response into Note object
	note := &Note{
		ID:      response["identifier"],
		Title:   response["title"],
		Content: response["note"],
	}

	// Parse tags if present
	if tagsStr := response["tags"]; tagsStr != "" {
		note.Tags = util.ParseTags(tagsStr)
	}

	// Parse timestamps
	if createdStr := response["creationDate"]; createdStr != "" {
		if parsed, err := time.Parse(time.RFC3339, createdStr); err == nil {
			note.CreatedAt = parsed
		}
	}
	if modifiedStr := response["modificationDate"]; modifiedStr != "" {
		if parsed, err := time.Parse(time.RFC3339, modifiedStr); err == nil {
			note.ModifiedAt = parsed
		}
	}

	// Parse trashed status
	if trashed := response["is_trashed"]; trashed == "yes" {
		note.IsTrashed = true
	}

	return note, nil
}

// UpdateNote modifies an existing note
// Mode can be: append, prepend, replace (content only), or replace_all (including title)
func (c *Client) UpdateNote(opts UpdateNoteOptions) (*Note, error) {
	if opts.ID == "" {
		return nil, fmt.Errorf("note ID is required for update")
	}

	if opts.Mode == "" {
		opts.Mode = "append" // Default to append
	}

	// Validate mode
	validModes := []string{"append", "prepend", "replace", "replace_all"}
	modeValid := false
	for _, valid := range validModes {
		if opts.Mode == valid {
			modeValid = true
			break
		}
	}
	if !modeValid {
		return nil, fmt.Errorf("invalid mode: %s (must be append, prepend, replace, or replace_all)", opts.Mode)
	}

	params := make(map[string]string)
	params["id"] = opts.ID

	// Handle text content updates
	if opts.Content != "" {
		content := opts.Content
		if opts.Timestamp {
			content = util.GetTimestamp() + "\n" + content
		}

		// Different parameters for different modes
		if opts.Mode == "replace_all" {
			params["title"] = opts.Content // For replace_all, this sets the whole note
			params["text"] = opts.Content
		} else {
			params["text"] = content
			params["mode"] = opts.Mode
		}
	}

	// Handle file updates
	if opts.FilePath != "" {
		expandedPath, err := util.ExpandHomePath(opts.FilePath)
		if err != nil {
			return nil, fmt.Errorf("invalid file path: %w", err)
		}

		fileContent, err := util.ReadFileAsBase64(expandedPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}

		fileName := filepath.Base(expandedPath)
		params["file"] = fileContent
		params["filename"] = fileName
	}

	// Handle header-specific updates
	if opts.Header != "" {
		params["header"] = opts.Header
	}

	// Add new line if requested (for append mode)
	if opts.NewLine && opts.Mode == "append" {
		params["new_line"] = "yes"
	}

	// Determine action based on content type
	action := "add-text"
	if opts.FilePath != "" {
		action = "add-file"
	}

	// Execute the update URL
	_, err := c.executeURL(c.buildBearURL(action, params))
	if err != nil {
		return nil, err
	}

	// Read back the updated note to return current state
	return c.ReadNote(ReadNoteOptions{ID: opts.ID})
}

// ArchiveNote archives a note (moves to trash)
func (c *Client) ArchiveNote(opts ArchiveNoteOptions) error {
	if opts.ID == "" {
		return fmt.Errorf("note ID is required for archive")
	}

	params := make(map[string]string)
	params["id"] = opts.ID

	if !opts.NoWindow {
		params["show-window"] = "yes"
	}

	_, err := c.executeURL(c.buildBearURL("archive", params))
	return err
}

// SearchNotes searches for notes matching a query
// Requires API token to be set
func (c *Client) SearchNotes(opts ListNotesOptions) (*NoteListResponse, error) {
	if opts.Token == "" && c.Token == "" {
		return nil, fmt.Errorf("API token required for search operations")
	}

	params := make(map[string]string)
	if opts.Token != "" {
		params["token"] = opts.Token
	}

	params["term"] = opts.Search

	response, err := c.executeURL(c.buildBearURL("search", params))
	if err != nil {
		return nil, err
	}

	// Parse the array response
	// Bear returns notes as a semi-colon separated list
	var notes []Note
	if notesStr := response["results"]; notesStr != "" {
		// This would need to be parsed based on Bear's actual response format
		// For now, we'll return the raw response
		_ = notesStr
	}

	return &NoteListResponse{
		Notes: notes,
		Count: len(notes),
	}, nil
}

// ListNotesByTag retrieves notes with a specific tag
// Requires API token to be set
func (c *Client) ListNotesByTag(opts ListNotesOptions) (*NoteListResponse, error) {
	if opts.Token == "" && c.Token == "" {
		return nil, fmt.Errorf("API token required for list operations")
	}

	params := make(map[string]string)
	if opts.Token != "" {
		params["token"] = opts.Token
	}
	params["name"] = opts.Tag

	response, err := c.executeURL(c.buildBearURL("open-tag", params))
	if err != nil {
		return nil, err
	}

	// Parse notes from response - Bear returns JSON array
	var notes []Note
	if notesStr := response["notes"]; notesStr != "" {
		// Bear returns a JSON array of note objects
		var rawNotes []struct {
			Identifier       string `json:"identifier"`
			Title            string `json:"title"`
			Tags             string `json:"tags"`
			CreationDate     string `json:"creationDate"`
			ModificationDate string `json:"modificationDate"`
			Pin              string `json:"pin"`
		}

		if err := json.Unmarshal([]byte(notesStr), &rawNotes); err != nil {
			return nil, fmt.Errorf("failed to parse notes response: %w", err)
		}

		for _, rn := range rawNotes {
			note := Note{
				ID:     rn.Identifier,
				Title:  rn.Title,
				Pinned: rn.Pin == "yes",
			}

			// Parse creation date
			if rn.CreationDate != "" {
				if parsed, err := time.Parse(time.RFC3339, rn.CreationDate); err == nil {
					note.CreatedAt = parsed
				}
			}

			// Parse modification date
			if rn.ModificationDate != "" {
				if parsed, err := time.Parse(time.RFC3339, rn.ModificationDate); err == nil {
					note.ModifiedAt = parsed
				}
			}

			// Parse tags (Bear returns them as a JSON string like "[\"go\"]")
			if rn.Tags != "" {
				var tagList []string
				if err := json.Unmarshal([]byte(rn.Tags), &tagList); err == nil {
					note.Tags = tagList
				}
			}

			notes = append(notes, note)
		}
	}

	return &NoteListResponse{
		Notes: notes,
		Count: len(notes),
	}, nil
}

// GetAllTags retrieves all tags from Bear
// Requires API token to be set
func (c *Client) GetAllTags(opts TagsListOptions) (*TagListResponse, error) {
	if opts.Token == "" && c.Token == "" {
		return nil, fmt.Errorf("API token required for tags operation")
	}

	params := make(map[string]string)
	if opts.Token != "" {
		params["token"] = opts.Token
	}

	response, err := c.executeURL(c.buildBearURL("tags", params))
	if err != nil {
		return nil, err
	}

	// Parse tags from response
	// Bear returns tags as a semi-colon separated list
	var tags []Tag
	if tagsStr := response["tags"]; tagsStr != "" {
		tagNames := strings.Split(tagsStr, ";")
		for _, name := range tagNames {
			trimmed := strings.TrimSpace(name)
			if trimmed != "" {
				tags = append(tags, Tag{Name: trimmed})
			}
		}
	}

	return &TagListResponse{
		Tags:  tags,
		Count: len(tags),
	}, nil
}

// RenameTag renames a tag across all notes
func (c *Client) RenameTag(opts RenameTagOptions) error {
	params := make(map[string]string)
	params["name"] = opts.Name
	params["new_name"] = opts.NewName

	_, err := c.executeURL(c.buildBearURL("rename-tag", params))
	return err
}

// DeleteTag deletes a tag from all notes
func (c *Client) DeleteTag(opts DeleteTagOptions) error {
	params := make(map[string]string)
	params["name"] = opts.Name

	_, err := c.executeURL(c.buildBearURL("delete-tag", params))
	return err
}

// parseJSONResponse is a helper that parses JSON responses from Bear
// Used internally for commands that return JSON data
func parseJSONResponse(rawJSON string, v interface{}) error {
	return json.Unmarshal([]byte(rawJSON), v)
}
