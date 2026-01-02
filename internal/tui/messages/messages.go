// Package messages contains all custom tea.Msg types for TUI communication.
// This follows the BubbleTea best practice of separating message types.
package messages

import (
	"fmt"

	"github.com/tranducquang/kiroku/internal/models"
)

// DataLoadedMsg indicates that data has been loaded from the database.
type DataLoadedMsg struct {
	Folders   []*models.Folder
	Notes     []*models.Note
	Templates []models.Template
}

// ErrorMsg wraps an error as a message with context.
type ErrorMsg struct {
	Err     error
	Context string
}

func (e ErrorMsg) Error() string {
	if e.Context == "" {
		return e.Err.Error()
	}
	return fmt.Sprintf("%s: %v", e.Context, e.Err)
}

// NewError creates a new ErrorMsg with context.
func NewError(err error, context string) ErrorMsg {
	return ErrorMsg{Err: err, Context: context}
}

// StatusClearMsg indicates that the status message should be cleared.
type StatusClearMsg struct{}

// EditorFinishedMsg indicates that the external editor process has completed.
type EditorFinishedMsg struct {
	TempFile string
	NoteID   int64
}

// SearchResultsMsg contains search results.
type SearchResultsMsg struct {
	Query string
	Notes []*models.Note
}

// NoteCreatedMsg indicates a note was created successfully.
type NoteCreatedMsg struct {
	Note *models.Note
}

// NoteDeletedMsg indicates a note was deleted.
type NoteDeletedMsg struct {
	NoteID int64
}

// NoteUpdatedMsg indicates a note was updated.
type NoteUpdatedMsg struct {
	Note *models.Note
}
