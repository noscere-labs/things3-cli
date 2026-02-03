package bear

import (
	"time"
)

// Note represents a Bear note with all its metadata
type Note struct {
	ID         string    `json:"note_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Tags       []string  `json:"tags"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
	IsTrashed  bool      `json:"is_trashed"`
	Pinned     bool      `json:"pinned,omitempty"`
}

// Tag represents a tag that can be applied to notes
type Tag struct {
	Name string `json:"name"`
}

// Response is the standard response envelope for all API operations
// This wraps both successful operations and errors in a consistent format
type Response struct {
	Success bool        `json:"success"`
	Error   string      `json:"error,omitempty"`
	Code    string      `json:"error_code,omitempty"`
	Details string      `json:"details,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// NoteListResponse wraps a list of notes with count
type NoteListResponse struct {
	Count int     `json:"count"`
	Notes []Note  `json:"notes"`
}

// TagListResponse wraps a list of tags with count
type TagListResponse struct {
	Count int   `json:"count"`
	Tags  []Tag `json:"tags"`
}

// CreateNoteOptions contains parameters for creating a new note
type CreateNoteOptions struct {
	Title      string
	Content    string
	Tags       []string
	FilePath   string
	Pin        bool
	NoWindow   bool
	Timestamp  bool
}

// ReadNoteOptions contains parameters for reading a note
type ReadNoteOptions struct {
	ID             string
	Title          string
	Header         string
	ExcludeTrashed bool
}

// UpdateNoteOptions contains parameters for updating a note
type UpdateNoteOptions struct {
	ID        string
	Content   string
	FilePath  string
	Mode      string // append, prepend, replace, replace_all
	Header    string
	Tags      []string
	NewLine   bool
	Timestamp bool
}

// ListNotesOptions contains parameters for listing notes
type ListNotesOptions struct {
	Tag    string
	Search string
	Filter string // all, untagged, todo, today, locked
	Token  string
}

// ArchiveNoteOptions contains parameters for archiving a note
type ArchiveNoteOptions struct {
	ID       string
	NoWindow bool
}

// TagsListOptions contains parameters for listing tags
type TagsListOptions struct {
	Token string
}

// RenameTagOptions contains parameters for renaming a tag
type RenameTagOptions struct {
	Name    string
	NewName string
}

// DeleteTagOptions contains parameters for deleting a tag
type DeleteTagOptions struct {
	Name string
}
