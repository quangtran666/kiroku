// Package commands contains tea.Cmd builders for TUI operations.
// This follows BubbleTea best practice #4: Use Command Builders.
package commands

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/tranducquang/kiroku/internal/models"
	"github.com/tranducquang/kiroku/internal/tui/constants"
	"github.com/tranducquang/kiroku/internal/tui/messages"
)

// NoteService defines the interface for note operations.
type NoteService interface {
	GetAllNotes(ctx context.Context) ([]*models.Note, error)
	GetTodos(ctx context.Context, showCompleted bool) ([]*models.Note, error)
	GetStarred(ctx context.Context) ([]*models.Note, error)
	GetByFolder(ctx context.Context, folderID int64) ([]*models.Note, error)
	Create(ctx context.Context, note *models.Note) error
	Update(ctx context.Context, note *models.Note) error
	Delete(ctx context.Context, id int64) error
	ToggleStar(ctx context.Context, id int64) error
	ToggleTodo(ctx context.Context, id int64) error
	SetPriority(ctx context.Context, id int64, priority int) error
}

// FolderService defines the interface for folder operations.
type FolderService interface {
	GetTree(ctx context.Context) ([]*models.Folder, error)
	Create(ctx context.Context, folder *models.Folder) error
	Delete(ctx context.Context, id int64) error
	ToggleStar(ctx context.Context, id int64) error
}

// TemplateService defines the interface for template operations.
type TemplateService interface {
	List(ctx context.Context) ([]models.Template, error)
}

// SearchService defines the interface for search operations.
type SearchService interface {
	Search(ctx context.Context, query string, opts models.ListOptions) ([]models.SearchResult, error)
}

// LoadDataParams contains parameters for loading initial data.
type LoadDataParams struct {
	FolderService   FolderService
	NoteService     NoteService
	TemplateService TemplateService
}

// LoadData returns a command that loads all initial data.
func LoadData(params LoadDataParams) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		folders, err := params.FolderService.GetTree(ctx)
		if err != nil {
			return messages.NewError(err, "load folders")
		}

		notes, err := params.NoteService.GetAllNotes(ctx)
		if err != nil {
			return messages.NewError(err, "load notes")
		}

		templates, err := params.TemplateService.List(ctx)
		if err != nil {
			return messages.NewError(err, "load templates")
		}

		return messages.DataLoadedMsg{
			Folders:   folders,
			Notes:     notes,
			Templates: templates,
		}
	}
}

// ReloadNotesParams contains parameters for reloading notes.
type ReloadNotesParams struct {
	NoteService   NoteService
	CurrentFilter string
	CurrentFolder *models.Folder
	ShowCompleted bool
}

// ReloadNotes returns a command that reloads notes based on the current filter.
func ReloadNotes(params ReloadNotesParams) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		var notes []*models.Note
		var err error

		switch params.CurrentFilter {
		case constants.FilterAll:
			notes, err = params.NoteService.GetAllNotes(ctx)
		case constants.FilterTodos:
			notes, err = params.NoteService.GetTodos(ctx, params.ShowCompleted)
		case constants.FilterStarred:
			notes, err = params.NoteService.GetStarred(ctx)
		default:
			if params.CurrentFolder != nil {
				notes, err = params.NoteService.GetByFolder(ctx, params.CurrentFolder.ID)
			} else {
				notes, err = params.NoteService.GetAllNotes(ctx)
			}
		}

		if err != nil {
			return messages.NewError(err, "reload notes")
		}

		return messages.DataLoadedMsg{Notes: notes}
	}
}

// SearchParams contains parameters for search.
type SearchParams struct {
	SearchService SearchService
	Query         string
}

// Search returns a command that performs a search.
func Search(params SearchParams) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		results, err := params.SearchService.Search(ctx, params.Query, models.ListOptions{})
		if err != nil {
			return messages.NewError(err, "search")
		}

		notes := make([]*models.Note, len(results))
		for i, r := range results {
			note := r.Note
			notes[i] = &note
		}

		return messages.SearchResultsMsg{
			Query: params.Query,
			Notes: notes,
		}
	}
}

// ClearStatusAfter returns a command that clears the status after a duration.
func ClearStatusAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(time.Time) tea.Msg {
		return messages.StatusClearMsg{}
	})
}

// CreateNoteParams contains parameters for creating a note.
type CreateNoteParams struct {
	NoteService   NoteService
	Title         string
	IsTodo        bool
	CurrentFolder *models.Folder
}

// CreateNote returns a command that creates a new note.
func CreateNote(params CreateNoteParams) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		var folderID *int64
		if params.CurrentFolder != nil {
			folderID = &params.CurrentFolder.ID
		}

		note := &models.Note{
			Title:    params.Title,
			FolderID: folderID,
			IsTodo:   params.IsTodo,
		}

		if err := params.NoteService.Create(ctx, note); err != nil {
			return messages.NewError(err, "create note")
		}

		return messages.NoteCreatedMsg{Note: note}
	}
}

// DeleteNote returns a command that deletes a note.
func DeleteNote(noteService NoteService, noteID int64) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		if err := noteService.Delete(ctx, noteID); err != nil {
			return messages.NewError(err, "delete note")
		}
		return messages.NoteDeletedMsg{NoteID: noteID}
	}
}

// ToggleStar returns a command that toggles a note's starred status.
func ToggleStar(noteService NoteService, noteID int64) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		if err := noteService.ToggleStar(ctx, noteID); err != nil {
			return messages.NewError(err, "toggle star")
		}
		return messages.NoteUpdatedMsg{}
	}
}

// ToggleTodo returns a command that toggles a todo's done status.
func ToggleTodo(noteService NoteService, noteID int64) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		if err := noteService.ToggleTodo(ctx, noteID); err != nil {
			return messages.NewError(err, "toggle todo")
		}
		return messages.NoteUpdatedMsg{}
	}
}

// CyclePriority returns a command that cycles a note's priority.
func CyclePriority(noteService NoteService, noteID int64, currentPriority int) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		newPriority := (currentPriority + 1) % constants.PriorityMax
		if err := noteService.SetPriority(ctx, noteID, newPriority); err != nil {
			return messages.NewError(err, "set priority")
		}
		return messages.NoteUpdatedMsg{}
	}
}

// UpdateNote returns a command that updates a note.
func UpdateNote(noteService NoteService, note *models.Note) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		if err := noteService.Update(ctx, note); err != nil {
			return messages.NewError(err, "update note")
		}
		return messages.NoteUpdatedMsg{Note: note}
	}
}

// CreateFolderParams contains parameters for creating a folder.
type CreateFolderParams struct {
	FolderService FolderService
	Name          string
	ParentID      *int64
}

// CreateFolder returns a command that creates a new folder.
func CreateFolder(params CreateFolderParams) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		folder := &models.Folder{
			Name:     params.Name,
			ParentID: params.ParentID,
		}

		if err := params.FolderService.Create(ctx, folder); err != nil {
			return messages.NewError(err, "create folder")
		}

		return ReloadFolders(params.FolderService)()
	}
}

// ReloadFolders returns a command that reloads folders.
func ReloadFolders(folderService FolderService) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		folders, err := folderService.GetTree(ctx)
		if err != nil {
			return messages.NewError(err, "reload folders")
		}
		return messages.DataLoadedMsg{Folders: folders}
	}
}

// ToggleFolderStar returns a command that toggles a folder's starred status.
func ToggleFolderStar(folderService FolderService, folderID int64) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		if err := folderService.ToggleStar(ctx, folderID); err != nil {
			return messages.NewError(err, "toggle folder star")
		}
		// Reload folders to update UI
		return ReloadFolders(folderService)()
	}
}

// DeleteFolder returns a command that deletes a folder.
func DeleteFolder(folderService FolderService, folderID int64) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		if err := folderService.Delete(ctx, folderID); err != nil {
			return messages.NewError(err, "delete folder")
		}
		// Reload folders
		return ReloadFolders(folderService)()
	}
}
