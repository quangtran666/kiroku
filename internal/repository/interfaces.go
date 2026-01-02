// Package repository contains interfaces and implementations for data access.
// This follows Go best practice: Accept interfaces, return structs.
package repository

import (
	"context"

	"github.com/tranducquang/kiroku/internal/models"
)

// NoteRepositoryInterface defines the contract for note data access.
// This enables dependency injection and easier testing.
type NoteRepositoryInterface interface {
	Create(ctx context.Context, note *models.Note) error
	GetByID(ctx context.Context, id int64) (*models.Note, error)
	Update(ctx context.Context, note *models.Note) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, opts models.ListOptions) ([]*models.Note, error)
	Count(ctx context.Context, opts models.ListOptions) (int, error)
	GetByFolder(ctx context.Context, folderID int64) ([]*models.Note, error)
	GetTodos(ctx context.Context, done *bool) ([]*models.Note, error)
	GetStarred(ctx context.Context) ([]*models.Note, error)
	GetRecent(ctx context.Context, limit int) ([]*models.Note, error)
}

// FolderRepositoryInterface defines the contract for folder data access.
type FolderRepositoryInterface interface {
	Create(ctx context.Context, folder *models.Folder) error
	GetByID(ctx context.Context, id int64) (*models.Folder, error)
	Update(ctx context.Context, folder *models.Folder) error
	Delete(ctx context.Context, id int64) error
	GetAll(ctx context.Context) ([]*models.Folder, error)
	GetChildren(ctx context.Context, parentID int64) ([]*models.Folder, error)
	CountNotes(ctx context.Context, folderID int64) (int, error)
	GetStarred(ctx context.Context) ([]*models.Folder, error)
}

// TemplateRepositoryInterface defines the contract for template data access.
type TemplateRepositoryInterface interface {
	Create(ctx context.Context, template *models.Template) error
	GetByID(ctx context.Context, id int64) (*models.Template, error)
	Update(ctx context.Context, template *models.Template) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]models.Template, error)
	GetByName(ctx context.Context, name string) (*models.Template, error)
	GetDefault(ctx context.Context) (*models.Template, error)
}

// SearchRepositoryInterface defines the contract for search data access.
type SearchRepositoryInterface interface {
	Search(ctx context.Context, query string, opts models.ListOptions) ([]SearchResult, error)
	SearchByTag(ctx context.Context, tag string, opts models.ListOptions) ([]*models.Note, error)
}

// Compile-time interface compliance checks
var (
	_ NoteRepositoryInterface     = (*NoteRepository)(nil)
	_ FolderRepositoryInterface   = (*FolderRepository)(nil)
	_ TemplateRepositoryInterface = (*TemplateRepository)(nil)
	_ SearchRepositoryInterface   = (*SearchRepository)(nil)
)
