package models

import (
	"errors"
	"time"
)

// ErrEmptyFolderName is returned when folder name is empty
var ErrEmptyFolderName = errors.New("folder name cannot be empty")

// Folder represents a folder for organizing notes
type Folder struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	ParentID  *int64    `json:"parent_id,omitempty"`
	Icon      string    `json:"icon"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Runtime fields (not stored in DB)
	NoteCount int       `json:"note_count,omitempty"`
	Children  []*Folder `json:"children,omitempty"`
	Expanded  bool      `json:"-"`
	Level     int       `json:"-"`
}

// Validate validates the folder fields
func (f *Folder) Validate() error {
	if f.Name == "" {
		return ErrEmptyFolderName
	}
	if f.Icon == "" {
		f.Icon = "📁"
	}
	return nil
}

// IsRoot returns true if folder has no parent
func (f *Folder) IsRoot() bool {
	return f.ParentID == nil
}

// HasChildren returns true if folder has child folders
func (f *Folder) HasChildren() bool {
	return len(f.Children) > 0
}

// Toggle toggles the expanded state
func (f *Folder) Toggle() {
	f.Expanded = !f.Expanded
}

// FolderListOptions contains options for listing folders
type FolderListOptions struct {
	ParentID     *int64
	IncludeCount bool
	OrderBy      string
	OrderDesc    bool
}
