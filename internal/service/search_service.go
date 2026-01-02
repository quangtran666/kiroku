package service

import (
	"context"

	"github.com/tranducquang/kiroku/internal/models"
	"github.com/tranducquang/kiroku/internal/repository"
)

// SearchService handles search business logic.
// It depends on repository interfaces, not concrete types (DI principle).
type SearchService struct {
	searchRepo repository.SearchRepositoryInterface
}

// NewSearchService creates a new search service with the given repository.
func NewSearchService(searchRepo repository.SearchRepositoryInterface) *SearchService {
	return &SearchService{
		searchRepo: searchRepo,
	}
}

// Search performs a full-text search on notes.
func (s *SearchService) Search(ctx context.Context, query string, opts models.ListOptions) ([]models.SearchResult, error) {
	results, err := s.searchRepo.Search(ctx, query, opts)
	if err != nil {
		return nil, err
	}

	// Convert repository.SearchResult to models.SearchResult
	modelResults := make([]models.SearchResult, len(results))
	for i, r := range results {
		modelResults[i] = models.SearchResult{
			Note:    r.Note,
			Snippet: r.Snippet,
			Rank:    r.Rank,
		}
	}

	return modelResults, nil
}

// SearchByTag searches notes by tag.
func (s *SearchService) SearchByTag(ctx context.Context, tag string, opts models.ListOptions) ([]*models.Note, error) {
	return s.searchRepo.SearchByTag(ctx, tag, opts)
}
