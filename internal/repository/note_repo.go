package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tranducquang/kiroku/internal/database"
	"github.com/tranducquang/kiroku/internal/models"
)

// Common errors
var (
	ErrNotFound = errors.New("not found")
)

// NoteRepository handles note database operations
type NoteRepository struct {
	db *database.DB
}

// NewNoteRepository creates a new note repository
func NewNoteRepository(db *database.DB) *NoteRepository {
	return &NoteRepository{db: db}
}

// Create creates a new note
func (r *NoteRepository) Create(ctx context.Context, note *models.Note) error {
	if err := note.Validate(); err != nil {
		return err
	}

	query := `
		INSERT INTO notes (title, content, folder_id, template_id, is_todo, is_done, priority, due_date, tags, starred, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	note.CreatedAt = now
	note.UpdatedAt = now

	result, err := r.db.ExecContext(ctx, query,
		note.Title,
		note.Content,
		note.FolderID,
		note.TemplateID,
		note.IsTodo,
		note.IsDone,
		note.Priority,
		note.DueDate,
		note.Tags,
		note.Starred,
		note.CreatedAt,
		note.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create note: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id: %w", err)
	}
	note.ID = id

	return nil
}

// GetByID retrieves a note by ID
func (r *NoteRepository) GetByID(ctx context.Context, id int64) (*models.Note, error) {
	query := `
		SELECT id, title, content, folder_id, template_id, is_todo, is_done, priority, due_date, tags, starred, created_at, updated_at
		FROM notes
		WHERE id = ?
	`

	note := &models.Note{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&note.ID,
		&note.Title,
		&note.Content,
		&note.FolderID,
		&note.TemplateID,
		&note.IsTodo,
		&note.IsDone,
		&note.Priority,
		&note.DueDate,
		&note.Tags,
		&note.Starred,
		&note.CreatedAt,
		&note.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get note by id: %w", err)
	}

	return note, nil
}

// Update updates an existing note
func (r *NoteRepository) Update(ctx context.Context, note *models.Note) error {
	if err := note.Validate(); err != nil {
		return err
	}

	query := `
		UPDATE notes
		SET title = ?, content = ?, folder_id = ?, template_id = ?, is_todo = ?, is_done = ?, priority = ?, due_date = ?, tags = ?, starred = ?, updated_at = ?
		WHERE id = ?
	`

	note.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		note.Title,
		note.Content,
		note.FolderID,
		note.TemplateID,
		note.IsTodo,
		note.IsDone,
		note.Priority,
		note.DueDate,
		note.Tags,
		note.Starred,
		note.UpdatedAt,
		note.ID,
	)
	if err != nil {
		return fmt.Errorf("update note: %w", err)
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

// Delete deletes a note by ID
func (r *NoteRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM notes WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete note: %w", err)
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

// List retrieves notes based on options
func (r *NoteRepository) List(ctx context.Context, opts models.ListOptions) ([]*models.Note, error) {
	var conditions []string
	var args []interface{}

	if opts.FolderID != nil {
		conditions = append(conditions, "folder_id = ?")
		args = append(args, *opts.FolderID)
	}
	if opts.IsTodo != nil {
		conditions = append(conditions, "is_todo = ?")
		args = append(args, *opts.IsTodo)
	}
	if opts.IsDone != nil {
		conditions = append(conditions, "is_done = ?")
		args = append(args, *opts.IsDone)
	}
	if opts.Starred != nil {
		conditions = append(conditions, "starred = ?")
		args = append(args, *opts.Starred)
	}
	if opts.Priority != nil {
		conditions = append(conditions, "priority = ?")
		args = append(args, *opts.Priority)
	}

	query := `
		SELECT id, title, content, folder_id, template_id, is_todo, is_done, priority, due_date, tags, starred, created_at, updated_at
		FROM notes
	`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Order by
	orderBy := opts.OrderBy
	if orderBy == "" {
		orderBy = "created_at"
	}
	orderDir := "ASC"
	if opts.OrderDesc {
		orderDir = "DESC"
	}
	query += fmt.Sprintf(" ORDER BY %s %s", orderBy, orderDir)

	// Limit and offset
	if opts.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", opts.Limit)
	}
	if opts.Offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", opts.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list notes: %w", err)
	}
	defer rows.Close()

	var notes []*models.Note
	for rows.Next() {
		var note models.Note
		err := rows.Scan(
			&note.ID,
			&note.Title,
			&note.Content,
			&note.FolderID,
			&note.TemplateID,
			&note.IsTodo,
			&note.IsDone,
			&note.Priority,
			&note.DueDate,
			&note.Tags,
			&note.Starred,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan note: %w", err)
		}
		notes = append(notes, &note)
	}

	return notes, nil
}

// Count returns the total number of notes matching the options
func (r *NoteRepository) Count(ctx context.Context, opts models.ListOptions) (int, error) {
	var conditions []string
	var args []interface{}

	if opts.FolderID != nil {
		conditions = append(conditions, "folder_id = ?")
		args = append(args, *opts.FolderID)
	}
	if opts.IsTodo != nil {
		conditions = append(conditions, "is_todo = ?")
		args = append(args, *opts.IsTodo)
	}
	if opts.IsDone != nil {
		conditions = append(conditions, "is_done = ?")
		args = append(args, *opts.IsDone)
	}
	if opts.Starred != nil {
		conditions = append(conditions, "starred = ?")
		args = append(args, *opts.Starred)
	}

	query := "SELECT COUNT(*) FROM notes"
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count notes: %w", err)
	}

	return count, nil
}

// GetByFolder retrieves notes by folder ID
func (r *NoteRepository) GetByFolder(ctx context.Context, folderID int64) ([]*models.Note, error) {
	return r.List(ctx, models.ListOptions{
		FolderID:  &folderID,
		OrderBy:   "created_at",
		OrderDesc: true,
	})
}

// GetTodos retrieves all todos
func (r *NoteRepository) GetTodos(ctx context.Context, done *bool) ([]*models.Note, error) {
	isTodo := true
	return r.List(ctx, models.ListOptions{
		IsTodo:    &isTodo,
		IsDone:    done,
		OrderBy:   "priority",
		OrderDesc: true,
	})
}

// GetStarred retrieves all starred notes
func (r *NoteRepository) GetStarred(ctx context.Context) ([]*models.Note, error) {
	starred := true
	return r.List(ctx, models.ListOptions{
		Starred:   &starred,
		OrderBy:   "updated_at",
		OrderDesc: true,
	})
}

// GetRecent retrieves recent notes
func (r *NoteRepository) GetRecent(ctx context.Context, limit int) ([]*models.Note, error) {
	return r.List(ctx, models.ListOptions{
		OrderBy:   "updated_at",
		OrderDesc: true,
		Limit:     limit,
	})
}
