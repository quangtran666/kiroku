package repository

import (
	"context"
	"fmt"

	"github.com/tranducquang/kiroku/internal/database"
	"github.com/tranducquang/kiroku/internal/models"
)

// SearchRepository handles full-text search operations
type SearchRepository struct {
	db *database.DB
}

// NewSearchRepository creates a new search repository
func NewSearchRepository(db *database.DB) *SearchRepository {
	return &SearchRepository{db: db}
}

// SearchResult represents a search result
type SearchResult struct {
	Note    models.Note
	Snippet string
	Rank    float64
}

// Search performs a full-text search on notes
func (r *SearchRepository) Search(ctx context.Context, query string, opts models.ListOptions) ([]SearchResult, error) {
	sqlQuery := `
		SELECT 
			n.id, n.title, n.content, n.folder_id, n.template_id, 
			n.is_todo, n.is_done, n.priority, n.due_date, n.tags, n.starred,
			n.created_at, n.updated_at,
			snippet(notes_fts, 0, '<mark>', '</mark>', '...', 32) as snippet,
			rank
		FROM notes_fts
		JOIN notes n ON notes_fts.rowid = n.id
		WHERE notes_fts MATCH ?
		ORDER BY rank
	`

	if opts.Limit > 0 {
		sqlQuery += fmt.Sprintf(" LIMIT %d", opts.Limit)
	}
	if opts.Offset > 0 {
		sqlQuery += fmt.Sprintf(" OFFSET %d", opts.Offset)
	}

	rows, err := r.db.QueryContext(ctx, sqlQuery, query)
	if err != nil {
		return nil, fmt.Errorf("search notes: %w", err)
	}
	defer rows.Close()

	var results []SearchResult
	for rows.Next() {
		var result SearchResult
		err := rows.Scan(
			&result.Note.ID,
			&result.Note.Title,
			&result.Note.Content,
			&result.Note.FolderID,
			&result.Note.TemplateID,
			&result.Note.IsTodo,
			&result.Note.IsDone,
			&result.Note.Priority,
			&result.Note.DueDate,
			&result.Note.Tags,
			&result.Note.Starred,
			&result.Note.CreatedAt,
			&result.Note.UpdatedAt,
			&result.Snippet,
			&result.Rank,
		)
		if err != nil {
			return nil, fmt.Errorf("scan search result: %w", err)
		}
		results = append(results, result)
	}

	return results, nil
}

// SearchByTag searches notes by tag
func (r *SearchRepository) SearchByTag(ctx context.Context, tag string, opts models.ListOptions) ([]*models.Note, error) {
	query := `
		SELECT id, title, content, folder_id, template_id, is_todo, is_done, priority, due_date, tags, starred, created_at, updated_at
		FROM notes
		WHERE tags LIKE ?
		ORDER BY updated_at DESC
	`

	if opts.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", opts.Limit)
	}
	if opts.Offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", opts.Offset)
	}

	tagPattern := "%" + tag + "%"
	rows, err := r.db.QueryContext(ctx, query, tagPattern)
	if err != nil {
		return nil, fmt.Errorf("search by tag: %w", err)
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
