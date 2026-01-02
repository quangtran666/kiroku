package models

import (
	"errors"
	"time"
)

// Priority levels for todos
const (
	PriorityNone   = 0
	PriorityLow    = 1
	PriorityMedium = 2
	PriorityHigh   = 3
)

// ErrEmptyTitle is returned when note title is empty
var ErrEmptyTitle = errors.New("title cannot be empty")

// Note represents a note or todo item
type Note struct {
	ID         int64      `json:"id"`
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	FolderID   *int64     `json:"folder_id,omitempty"`
	TemplateID *int64     `json:"template_id,omitempty"`
	IsTodo     bool       `json:"is_todo"`
	IsDone     bool       `json:"is_done"`
	Priority   int        `json:"priority"`
	DueDate    *time.Time `json:"due_date,omitempty"`
	Tags       string     `json:"tags"`
	Starred    bool       `json:"starred"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// Validate validates the note fields
func (n *Note) Validate() error {
	if n.Title == "" {
		return ErrEmptyTitle
	}
	if n.Priority < PriorityNone || n.Priority > PriorityHigh {
		n.Priority = PriorityNone
	}
	return nil
}

// ToggleDone toggles the done status of a todo
func (n *Note) ToggleDone() {
	if n.IsTodo {
		n.IsDone = !n.IsDone
	}
}

// ToggleStar toggles the starred status
func (n *Note) ToggleStar() {
	n.Starred = !n.Starred
}

// SetPriority sets the priority level
func (n *Note) SetPriority(p int) {
	if p >= PriorityNone && p <= PriorityHigh {
		n.Priority = p
	}
}

// PriorityString returns a human-readable priority string
func (n *Note) PriorityString() string {
	switch n.Priority {
	case PriorityHigh:
		return "High"
	case PriorityMedium:
		return "Medium"
	case PriorityLow:
		return "Low"
	default:
		return "None"
	}
}

// PriorityIcon returns an icon for the priority
func (n *Note) PriorityIcon() string {
	switch n.Priority {
	case PriorityHigh:
		return "🔴"
	case PriorityMedium:
		return "🟡"
	case PriorityLow:
		return "🟢"
	default:
		return ""
	}
}

// StatusIcon returns an icon for todo status
func (n *Note) StatusIcon() string {
	if !n.IsTodo {
		return "📝"
	}
	if n.IsDone {
		return "☑"
	}
	return "☐"
}

// ListOptions contains options for listing notes
type ListOptions struct {
	FolderID  *int64
	IsTodo    *bool
	IsDone    *bool
	Starred   *bool
	Priority  *int
	OrderBy   string
	OrderDesc bool
	Limit     int
	Offset    int
}

// SearchResult represents a search result with highlight info
type SearchResult struct {
	Note    Note
	Snippet string
	Rank    float64
}
