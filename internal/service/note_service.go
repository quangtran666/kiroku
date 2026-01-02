package service

import (
	"context"
	"fmt"

	"github.com/tranducquang/kiroku/internal/models"
	"github.com/tranducquang/kiroku/internal/repository"
)

// NoteService handles note business logic.
// It depends on repository interfaces, not concrete types (DI principle).
type NoteService struct {
	noteRepo     repository.NoteRepositoryInterface
	templateRepo repository.TemplateRepositoryInterface
	folderRepo   repository.FolderRepositoryInterface
}

// NewNoteService creates a new note service with the given repositories.
func NewNoteService(
	noteRepo repository.NoteRepositoryInterface,
	templateRepo repository.TemplateRepositoryInterface,
	folderRepo repository.FolderRepositoryInterface,
) *NoteService {
	return &NoteService{
		noteRepo:     noteRepo,
		templateRepo: templateRepo,
		folderRepo:   folderRepo,
	}
}

// Create creates a new note.
func (s *NoteService) Create(ctx context.Context, note *models.Note) error {
	if err := note.Validate(); err != nil {
		return fmt.Errorf("validate note: %w", err)
	}

	if note.TemplateID == nil {
		return s.noteRepo.Create(ctx, note)
	}

	template, err := s.templateRepo.GetByID(ctx, *note.TemplateID)
	if err != nil {
		return fmt.Errorf("get template: %w", err)
	}

	if note.Content == "" {
		note.Content = template.Content
	}

	return s.noteRepo.Create(ctx, note)
}

// GetByID retrieves a note by ID.
func (s *NoteService) GetByID(ctx context.Context, id int64) (*models.Note, error) {
	return s.noteRepo.GetByID(ctx, id)
}

// Update updates an existing note.
func (s *NoteService) Update(ctx context.Context, note *models.Note) error {
	if err := note.Validate(); err != nil {
		return fmt.Errorf("validate note: %w", err)
	}
	return s.noteRepo.Update(ctx, note)
}

// Delete deletes a note by ID.
func (s *NoteService) Delete(ctx context.Context, id int64) error {
	return s.noteRepo.Delete(ctx, id)
}

// GetAllNotes retrieves all notes ordered by updated_at descending.
func (s *NoteService) GetAllNotes(ctx context.Context) ([]*models.Note, error) {
	return s.noteRepo.List(ctx, models.ListOptions{
		OrderBy:   "updated_at",
		OrderDesc: true,
	})
}

// GetByFolder retrieves notes in a specific folder.
func (s *NoteService) GetByFolder(ctx context.Context, folderID int64) ([]*models.Note, error) {
	return s.noteRepo.GetByFolder(ctx, folderID)
}

// GetTodos retrieves all todos, optionally including completed ones.
func (s *NoteService) GetTodos(ctx context.Context, showCompleted bool) ([]*models.Note, error) {
	var done *bool
	if !showCompleted {
		f := false
		done = &f
	}
	return s.noteRepo.GetTodos(ctx, done)
}

// GetStarred retrieves all starred notes.
func (s *NoteService) GetStarred(ctx context.Context) ([]*models.Note, error) {
	return s.noteRepo.GetStarred(ctx)
}

// GetRecent retrieves the most recently updated notes.
func (s *NoteService) GetRecent(ctx context.Context, limit int) ([]*models.Note, error) {
	return s.noteRepo.GetRecent(ctx, limit)
}

// ToggleStar toggles the starred status of a note.
func (s *NoteService) ToggleStar(ctx context.Context, id int64) error {
	note, err := s.noteRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get note: %w", err)
	}

	note.Starred = !note.Starred
	return s.noteRepo.Update(ctx, note)
}

// ToggleTodo toggles the done status of a todo.
func (s *NoteService) ToggleTodo(ctx context.Context, id int64) error {
	note, err := s.noteRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get note: %w", err)
	}

	if !note.IsTodo {
		return fmt.Errorf("note is not a todo")
	}

	note.IsDone = !note.IsDone
	return s.noteRepo.Update(ctx, note)
}

// SetPriority sets the priority of a note.
func (s *NoteService) SetPriority(ctx context.Context, id int64, priority int) error {
	note, err := s.noteRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get note: %w", err)
	}

	note.Priority = priority
	return s.noteRepo.Update(ctx, note)
}

// MoveToFolder moves a note to a different folder.
func (s *NoteService) MoveToFolder(ctx context.Context, noteID, folderID int64) error {
	note, err := s.noteRepo.GetByID(ctx, noteID)
	if err != nil {
		return fmt.Errorf("get note: %w", err)
	}

	if _, err := s.folderRepo.GetByID(ctx, folderID); err != nil {
		return fmt.Errorf("get folder: %w", err)
	}

	note.FolderID = &folderID
	return s.noteRepo.Update(ctx, note)
}

// Count returns the total number of notes matching the options.
func (s *NoteService) Count(ctx context.Context, opts models.ListOptions) (int, error) {
	return s.noteRepo.Count(ctx, opts)
}
