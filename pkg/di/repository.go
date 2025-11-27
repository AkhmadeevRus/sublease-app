package di

import (
	"github.com/AkhmadeevRus/sublease-app/pkg/auth"
	"github.com/AkhmadeevRus/sublease-app/pkg/property"
	"github.com/AkhmadeevRus/sublease-app/pkg/search"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	AuthRepository     auth.IAuthRepository
	PropertyRepository property.IPropertyRepository
	SearchRepository   search.ISearchRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		AuthRepository:     auth.NewAuthRepository(db),
		PropertyRepository: property.NewPropertyRepository(db),
		SearchRepository:   search.NewSearchRepository(db),
	}
}
