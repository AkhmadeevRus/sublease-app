package search

import (
	"github.com/AkhmadeevRus/sublease-app/pkg/property"
	"github.com/jmoiron/sqlx"
)

type ISearchRepository interface {
	GetAllProperties() ([]property.Property, error)
	SearchByFilters(filter PropertyFilter) ([]property.Property, error)
}

type SearchRepository struct {
	db *sqlx.DB
}

func NewSearchRepository(db *sqlx.DB) *SearchRepository {
	return &SearchRepository{db: db}
}

func (r *SearchRepository) GetAllProperties() ([]property.Property, error) {
	return []property.Property{}, nil
}

func (r *SearchRepository) SearchByFilters(filter PropertyFilter) ([]property.Property, error) {
	return []property.Property{}, nil
}
