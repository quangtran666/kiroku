package service

import (
	"context"
	"fmt"

	"github.com/tranducquang/kiroku/internal/models"
	"github.com/tranducquang/kiroku/internal/repository"
)

// TemplateService handles template business logic.
// It depends on repository interfaces, not concrete types (DI principle).
type TemplateService struct {
	templateRepo repository.TemplateRepositoryInterface
}

// NewTemplateService creates a new template service with the given repository.
func NewTemplateService(templateRepo repository.TemplateRepositoryInterface) *TemplateService {
	return &TemplateService{
		templateRepo: templateRepo,
	}
}

// Create creates a new template.
func (s *TemplateService) Create(ctx context.Context, template *models.Template) error {
	if err := template.Validate(); err != nil {
		return fmt.Errorf("validate template: %w", err)
	}
	return s.templateRepo.Create(ctx, template)
}

// GetByID retrieves a template by ID.
func (s *TemplateService) GetByID(ctx context.Context, id int64) (*models.Template, error) {
	return s.templateRepo.GetByID(ctx, id)
}

// GetByName retrieves a template by name.
func (s *TemplateService) GetByName(ctx context.Context, name string) (*models.Template, error) {
	return s.templateRepo.GetByName(ctx, name)
}

// Update updates an existing template.
func (s *TemplateService) Update(ctx context.Context, template *models.Template) error {
	if err := template.Validate(); err != nil {
		return fmt.Errorf("validate template: %w", err)
	}
	return s.templateRepo.Update(ctx, template)
}

// Delete deletes a template by ID.
func (s *TemplateService) Delete(ctx context.Context, id int64) error {
	return s.templateRepo.Delete(ctx, id)
}

// List retrieves all templates.
func (s *TemplateService) List(ctx context.Context) ([]models.Template, error) {
	return s.templateRepo.List(ctx)
}

// GetDefault retrieves the default template.
func (s *TemplateService) GetDefault(ctx context.Context) (*models.Template, error) {
	return s.templateRepo.GetDefault(ctx)
}
