package di

import (
	"github.com/AkhmadeevRus/sublease-app/pkg/auth"
	"github.com/AkhmadeevRus/sublease-app/pkg/property"
	"github.com/AkhmadeevRus/sublease-app/pkg/search"
)

type Service struct {
	AuthService     auth.IAuthService
	PropertyService property.IPropertyService
	SearchService   search.ISearchService
}

func NewService(repos *Repository) *Service {
	return &Service{
		AuthService:     auth.NewAuthService(repos.AuthRepository),
		PropertyService: property.NewPropertyService(repos.PropertyRepository),
		SearchService:   search.NewSearchService(repos.SearchRepository),
	}
}
