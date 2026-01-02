package service

import (
	"context"
	"fmt"

	"github.com/tranducquang/kiroku/internal/models"
	"github.com/tranducquang/kiroku/internal/repository"
)

// FolderService handles folder business logic.
// It depends on repository interfaces, not concrete types (DI principle).
type FolderService struct {
	folderRepo repository.FolderRepositoryInterface
	noteRepo   repository.NoteRepositoryInterface
}

// NewFolderService creates a new folder service with the given repositories.
func NewFolderService(
	folderRepo repository.FolderRepositoryInterface,
	noteRepo repository.NoteRepositoryInterface,
) *FolderService {
	return &FolderService{
		folderRepo: folderRepo,
		noteRepo:   noteRepo,
	}
}

// Create creates a new folder.
func (s *FolderService) Create(ctx context.Context, folder *models.Folder) error {
	if err := folder.Validate(); err != nil {
		return fmt.Errorf("validate folder: %w", err)
	}
	return s.folderRepo.Create(ctx, folder)
}

// GetByID retrieves a folder by ID.
func (s *FolderService) GetByID(ctx context.Context, id int64) (*models.Folder, error) {
	return s.folderRepo.GetByID(ctx, id)
}

// Update updates an existing folder.
func (s *FolderService) Update(ctx context.Context, folder *models.Folder) error {
	if err := folder.Validate(); err != nil {
		return fmt.Errorf("validate folder: %w", err)
	}
	return s.folderRepo.Update(ctx, folder)
}

// Delete deletes a folder by ID.
func (s *FolderService) Delete(ctx context.Context, id int64) error {
	return s.folderRepo.Delete(ctx, id)
}

// GetAll retrieves all folders.
func (s *FolderService) GetAll(ctx context.Context) ([]*models.Folder, error) {
	return s.folderRepo.GetAll(ctx)
}

// GetTree retrieves the folder tree structure with note counts.
func (s *FolderService) GetTree(ctx context.Context) ([]*models.Folder, error) {
	folders, err := s.folderRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all folders: %w", err)
	}

	folderMap := make(map[int64]*models.Folder)
	for _, folder := range folders {
		folderMap[folder.ID] = folder

		count, err := s.folderRepo.CountNotes(ctx, folder.ID)
		if err != nil {
			return nil, fmt.Errorf("count notes: %w", err)
		}
		folder.NoteCount = count
	}

	var rootFolders []*models.Folder
	for _, folder := range folders {
		if folder.ParentID == nil {
			rootFolders = append(rootFolders, folder)
			continue
		}

		parent, ok := folderMap[*folder.ParentID]
		if ok {
			parent.Children = append(parent.Children, folder)
		}
	}

	return rootFolders, nil
}

// GetChildren retrieves child folders of a parent folder.
func (s *FolderService) GetChildren(ctx context.Context, parentID int64) ([]*models.Folder, error) {
	return s.folderRepo.GetChildren(ctx, parentID)
}

// ToggleStar toggles the starred status of a folder.
func (s *FolderService) ToggleStar(ctx context.Context, id int64) error {
	folder, err := s.folderRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get folder: %w", err)
	}

	folder.Starred = !folder.Starred
	return s.folderRepo.Update(ctx, folder)
}

// GetStarred retrieves all starred folders.
func (s *FolderService) GetStarred(ctx context.Context) ([]*models.Folder, error) {
	return s.folderRepo.GetStarred(ctx)
}
