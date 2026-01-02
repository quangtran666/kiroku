package service

import (
	"context"
	"fmt"

	"github.com/tranducquang/kiroku/internal/models"
	"github.com/tranducquang/kiroku/internal/repository"
)

// FolderService handles folder business logic
type FolderService struct {
	folderRepo *repository.FolderRepository
	noteRepo   *repository.NoteRepository
}

// NewFolderService creates a new folder service
func NewFolderService(
	folderRepo *repository.FolderRepository,
	noteRepo *repository.NoteRepository,
) *FolderService {
	return &FolderService{
		folderRepo: folderRepo,
		noteRepo:   noteRepo,
	}
}

// Create creates a new folder
func (s *FolderService) Create(ctx context.Context, folder *models.Folder) error {
	if err := folder.Validate(); err != nil {
		return fmt.Errorf("validate folder: %w", err)
	}
	return s.folderRepo.Create(ctx, folder)
}

// GetByID retrieves a folder by ID
func (s *FolderService) GetByID(ctx context.Context, id int64) (*models.Folder, error) {
	return s.folderRepo.GetByID(ctx, id)
}

// Update updates an existing folder
func (s *FolderService) Update(ctx context.Context, folder *models.Folder) error {
	if err := folder.Validate(); err != nil {
		return fmt.Errorf("validate folder: %w", err)
	}
	return s.folderRepo.Update(ctx, folder)
}

// Delete deletes a folder by ID
func (s *FolderService) Delete(ctx context.Context, id int64) error {
	return s.folderRepo.Delete(ctx, id)
}

// GetAll retrieves all folders
func (s *FolderService) GetAll(ctx context.Context) ([]*models.Folder, error) {
	return s.folderRepo.GetAll(ctx)
}

// GetTree retrieves the folder tree structure
func (s *FolderService) GetTree(ctx context.Context) ([]*models.Folder, error) {
	// Get all folders
	folders, err := s.folderRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all folders: %w", err)
	}

	// Build folder map
	folderMap := make(map[int64]*models.Folder)
	for _, folder := range folders {
		folderMap[folder.ID] = folder

		// Get note count for each folder
		count, err := s.folderRepo.CountNotes(ctx, folder.ID)
		if err != nil {
			return nil, fmt.Errorf("count notes: %w", err)
		}
		folder.NoteCount = count
	}

	// Build tree structure
	var rootFolders []*models.Folder
	for _, folder := range folders {
		if folder.ParentID == nil {
			rootFolders = append(rootFolders, folder)
		} else {
			parent, ok := folderMap[*folder.ParentID]
			if ok {
				parent.Children = append(parent.Children, folder)
			}
		}
	}

	return rootFolders, nil
}

// GetChildren retrieves child folders
func (s *FolderService) GetChildren(ctx context.Context, parentID int64) ([]*models.Folder, error) {
	return s.folderRepo.GetChildren(ctx, parentID)
}
