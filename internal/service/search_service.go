package service

import (
	"context"

	"github.com/tranducquang/kiroku/internal/models"
	"github.com/tranducquang/kiroku/internal/repository"
)

// SearchService handles search business logic
type SearchService struct {
	searchRepo *repository.SearchRepository
}

// NewSearchService creates a new search service
func NewSearchService(searchRepo *repository.SearchRepository) *SearchService {
	return &SearchService{
		searchRepo: searchRepo,
	}
}

// Search performs a full-text search
func (s *SearchService) Search(ctx context.Context, query string, opts models.ListOptions) ([]repository.SearchResult, error) {
	return s.searchRepo.Search(ctx, query, opts)
}

// SearchByTag searches notes by tag
func (s *SearchService) SearchByTag(ctx context.Context, tag string, opts models.ListOptions) ([]*models.Note, error) {
	return s.searchRepo.SearchByTag(ctx, tag, opts)
}
