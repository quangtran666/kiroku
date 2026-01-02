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

// TemplateRepository handles template database operations
type TemplateRepository struct {
	db *database.DB
}

// NewTemplateRepository creates a new template repository
func NewTemplateRepository(db *database.DB) *TemplateRepository {
	return &TemplateRepository{db: db}
}

// Create creates a new template
func (r *TemplateRepository) Create(ctx context.Context, template *models.Template) error {
	if err := template.Validate(); err != nil {
		return err
	}

	query := `
		INSERT INTO templates (name, content, description, is_default, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	template.CreatedAt = now
	template.UpdatedAt = now

	result, err := r.db.ExecContext(ctx, query,
		template.Name,
		template.Content,
		template.Description,
		template.IsDefault,
		template.CreatedAt,
		template.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create template: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id: %w", err)
	}
	template.ID = id

	return nil
}

// GetByID retrieves a template by ID
func (r *TemplateRepository) GetByID(ctx context.Context, id int64) (*models.Template, error) {
	query := `
		SELECT id, name, content, description, is_default, created_at, updated_at
		FROM templates
		WHERE id = ?
	`

	template := &models.Template{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&template.ID,
		&template.Name,
		&template.Content,
		&template.Description,
		&template.IsDefault,
		&template.CreatedAt,
		&template.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get template by id: %w", err)
	}

	return template, nil
}

// GetByName retrieves a template by name
func (r *TemplateRepository) GetByName(ctx context.Context, name string) (*models.Template, error) {
	query := `
		SELECT id, name, content, description, is_default, created_at, updated_at
		FROM templates
		WHERE name = ?
	`

	template := &models.Template{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&template.ID,
		&template.Name,
		&template.Content,
		&template.Description,
		&template.IsDefault,
		&template.CreatedAt,
		&template.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get template by name: %w", err)
	}

	return template, nil
}

// Update updates an existing template
func (r *TemplateRepository) Update(ctx context.Context, template *models.Template) error {
	if err := template.Validate(); err != nil {
		return err
	}

	query := `
		UPDATE templates
		SET name = ?, content = ?, description = ?, is_default = ?, updated_at = ?
		WHERE id = ?
	`

	template.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		template.Name,
		template.Content,
		template.Description,
		template.IsDefault,
		template.UpdatedAt,
		template.ID,
	)
	if err != nil {
		return fmt.Errorf("update template: %w", err)
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

// Delete deletes a template by ID
func (r *TemplateRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM templates WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete template: %w", err)
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

// List retrieves all templates
func (r *TemplateRepository) List(ctx context.Context) ([]models.Template, error) {
	query := `
		SELECT id, name, content, description, is_default, created_at, updated_at
		FROM templates
		ORDER BY is_default DESC, name ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list templates: %w", err)
	}
	defer rows.Close()

	var templates []models.Template
	for rows.Next() {
		var template models.Template
		err := rows.Scan(
			&template.ID,
			&template.Name,
			&template.Content,
			&template.Description,
			&template.IsDefault,
			&template.CreatedAt,
			&template.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan template: %w", err)
		}
		templates = append(templates, template)
	}

	return templates, nil
}

// GetDefault retrieves the default template
func (r *TemplateRepository) GetDefault(ctx context.Context) (*models.Template, error) {
	query := `
		SELECT id, name, content, description, is_default, created_at, updated_at
		FROM templates
		WHERE is_default = 1
		LIMIT 1
	`

	template := &models.Template{}
	err := r.db.QueryRowContext(ctx, query).Scan(
		&template.ID,
		&template.Name,
		&template.Content,
		&template.Description,
		&template.IsDefault,
		&template.CreatedAt,
		&template.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get default template: %w", err)
	}

	return template, nil
}
