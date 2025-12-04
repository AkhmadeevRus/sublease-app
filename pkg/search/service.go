package search

import (
	"github.com/AkhmadeevRus/sublease-app/pkg/property"
	"github.com/google/uuid"
)

type ISearchService interface {
	GetAllProperties() ([]property.Property, error)
	GetPropertyById(id uuid.UUID) (property.Property, error)
	SearchByFilters(filter PropertyFilter) ([]property.Property, error)
}

type SearchService struct {
	repo ISearchRepository
}

func NewSearchService(repo ISearchRepository) *SearchService {
	return &SearchService{repo: repo}
}

func (s *SearchService) GetAllProperties() ([]property.Property, error) {
	return s.repo.GetAllProperties()
}

func (s *SearchService) GetPropertyById(id uuid.UUID) (property.Property, error) {
	return s.repo.GetPropertyById(id)
}

func (s *SearchService) SearchByFilters(filter PropertyFilter) ([]property.Property, error) {
	if err := filter.Validate(); err != nil {
		return nil, err
	}
	return s.repo.SearchByFilters(filter)
}
