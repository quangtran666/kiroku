// Package service contains interfaces for service layer.
// This follows Go best practice: Accept interfaces, return structs.
package service

import (
	"context"
	"os/exec"

	"github.com/tranducquang/kiroku/internal/models"
)

// NoteServiceInterface defines the contract for note business logic.
type NoteServiceInterface interface {
	Create(ctx context.Context, note *models.Note) error
	GetByID(ctx context.Context, id int64) (*models.Note, error)
	Update(ctx context.Context, note *models.Note) error
	Delete(ctx context.Context, id int64) error
	GetAllNotes(ctx context.Context) ([]*models.Note, error)
	GetByFolder(ctx context.Context, folderID int64) ([]*models.Note, error)
	GetTodos(ctx context.Context, showCompleted bool) ([]*models.Note, error)
	GetStarred(ctx context.Context) ([]*models.Note, error)
	GetRecent(ctx context.Context, limit int) ([]*models.Note, error)
	ToggleStar(ctx context.Context, id int64) error
	ToggleTodo(ctx context.Context, id int64) error
	SetPriority(ctx context.Context, id int64, priority int) error
	MoveToFolder(ctx context.Context, noteID, folderID int64) error
	Count(ctx context.Context, opts models.ListOptions) (int, error)
}

// FolderServiceInterface defines the contract for folder business logic.
type FolderServiceInterface interface {
	Create(ctx context.Context, folder *models.Folder) error
	GetByID(ctx context.Context, id int64) (*models.Folder, error)
	Update(ctx context.Context, folder *models.Folder) error
	Delete(ctx context.Context, id int64) error
	GetAll(ctx context.Context) ([]*models.Folder, error)
	GetTree(ctx context.Context) ([]*models.Folder, error)
	GetChildren(ctx context.Context, parentID int64) ([]*models.Folder, error)
}

// TemplateServiceInterface defines the contract for template business logic.
type TemplateServiceInterface interface {
	Create(ctx context.Context, template *models.Template) error
	GetByID(ctx context.Context, id int64) (*models.Template, error)
	GetByName(ctx context.Context, name string) (*models.Template, error)
	Update(ctx context.Context, template *models.Template) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]models.Template, error)
	GetDefault(ctx context.Context) (*models.Template, error)
}

// SearchServiceInterface defines the contract for search business logic.
type SearchServiceInterface interface {
	Search(ctx context.Context, query string, opts models.ListOptions) ([]models.SearchResult, error)
	SearchByTag(ctx context.Context, tag string, opts models.ListOptions) ([]*models.Note, error)
}

// EditorServiceInterface defines the contract for editor operations.
type EditorServiceInterface interface {
	EditNote(title, content string) (newTitle, newContent string, err error)
	PrepareEdit(title, content string) (tmpFilePath string, cmd *exec.Cmd, err error)
	ReadEditedContent(tmpFilePath, originalTitle string) (title, content string, err error)
	CreateNote(templateContent string) (title, content string, err error)
}

// Compile-time interface compliance checks
var (
	_ NoteServiceInterface     = (*NoteService)(nil)
	_ FolderServiceInterface   = (*FolderService)(nil)
	_ TemplateServiceInterface = (*TemplateService)(nil)
	_ SearchServiceInterface   = (*SearchService)(nil)
	_ EditorServiceInterface   = (*EditorService)(nil)
)
