package search

import (
	"github.com/AkhmadeevRus/sublease-app/pkg/property"
)

type ISearchService interface {
	GetAllProperties() ([]property.Property, error)
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

func (s *SearchService) SearchByFilters(filter PropertyFilter) ([]property.Property, error) {
	return s.repo.SearchByFilters(filter)
}
