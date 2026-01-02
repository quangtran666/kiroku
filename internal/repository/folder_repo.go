package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/tranducquang/kiroku/internal/database"
	"github.com/tranducquang/kiroku/internal/models"
)

// FolderRepository handles folder database operations
type FolderRepository struct {
	db *database.DB
}

// NewFolderRepository creates a new folder repository
func NewFolderRepository(db *database.DB) *FolderRepository {
	return &FolderRepository{db: db}
}

// Create creates a new folder
func (r *FolderRepository) Create(ctx context.Context, folder *models.Folder) error {
	if err := folder.Validate(); err != nil {
		return err
	}

	query := `
		INSERT INTO folders (name, parent_id, icon, position, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	folder.CreatedAt = now
	folder.UpdatedAt = now

	result, err := r.db.ExecContext(ctx, query,
		folder.Name,
		folder.ParentID,
		folder.Icon,
		folder.Position,
		folder.CreatedAt,
		folder.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create folder: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id: %w", err)
	}
	folder.ID = id

	return nil
}

// GetByID retrieves a folder by ID
func (r *FolderRepository) GetByID(ctx context.Context, id int64) (*models.Folder, error) {
	query := `
		SELECT id, name, parent_id, icon, position, created_at, updated_at
		FROM folders
		WHERE id = ?
	`

	folder := &models.Folder{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&folder.ID,
		&folder.Name,
		&folder.ParentID,
		&folder.Icon,
		&folder.Position,
		&folder.CreatedAt,
		&folder.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get folder by id: %w", err)
	}

	return folder, nil
}

// Update updates an existing folder
func (r *FolderRepository) Update(ctx context.Context, folder *models.Folder) error {
	if err := folder.Validate(); err != nil {
		return err
	}

	query := `
		UPDATE folders
		SET name = ?, parent_id = ?, icon = ?, position = ?, updated_at = ?
		WHERE id = ?
	`

	folder.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		folder.Name,
		folder.ParentID,
		folder.Icon,
		folder.Position,
		folder.UpdatedAt,
		folder.ID,
	)
	if err != nil {
		return fmt.Errorf("update folder: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

// Delete deletes a folder by ID
func (r *FolderRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM folders WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete folder: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

// GetAll retrieves all folders
func (r *FolderRepository) GetAll(ctx context.Context) ([]*models.Folder, error) {
	query := `
		SELECT id, name, parent_id, icon, position, created_at, updated_at
		FROM folders
		ORDER BY position ASC, name ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get all folders: %w", err)
	}
	defer rows.Close()

	var folders []*models.Folder
	for rows.Next() {
		var folder models.Folder
		err := rows.Scan(
			&folder.ID,
			&folder.Name,
			&folder.ParentID,
			&folder.Icon,
			&folder.Position,
			&folder.CreatedAt,
			&folder.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan folder: %w", err)
		}
		folders = append(folders, &folder)
	}

	return folders, nil
}

// GetRootFolders retrieves all root folders (no parent)
func (r *FolderRepository) GetRootFolders(ctx context.Context) ([]*models.Folder, error) {
	query := `
		SELECT id, name, parent_id, icon, position, created_at, updated_at
		FROM folders
		WHERE parent_id IS NULL
		ORDER BY position ASC, name ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get root folders: %w", err)
	}
	defer rows.Close()

	var folders []*models.Folder
	for rows.Next() {
		var folder models.Folder
		err := rows.Scan(
			&folder.ID,
			&folder.Name,
			&folder.ParentID,
			&folder.Icon,
			&folder.Position,
			&folder.CreatedAt,
			&folder.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan folder: %w", err)
		}
		folders = append(folders, &folder)
	}

	return folders, nil
}

// GetChildren retrieves child folders of a parent folder
func (r *FolderRepository) GetChildren(ctx context.Context, parentID int64) ([]*models.Folder, error) {
	query := `
		SELECT id, name, parent_id, icon, position, created_at, updated_at
		FROM folders
		WHERE parent_id = ?
		ORDER BY position ASC, name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, parentID)
	if err != nil {
		return nil, fmt.Errorf("get children folders: %w", err)
	}
	defer rows.Close()

	var folders []*models.Folder
	for rows.Next() {
		var folder models.Folder
		err := rows.Scan(
			&folder.ID,
			&folder.Name,
			&folder.ParentID,
			&folder.Icon,
			&folder.Position,
			&folder.CreatedAt,
			&folder.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan folder: %w", err)
		}
		folders = append(folders, &folder)
	}

	return folders, nil
}

// CountNotes returns the number of notes in a folder
func (r *FolderRepository) CountNotes(ctx context.Context, folderID int64) (int, error) {
	query := `SELECT COUNT(*) FROM notes WHERE folder_id = ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, folderID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count notes in folder: %w", err)
	}

	return count, nil
}
