package di

import (
	"github.com/AkhmadeevRus/sublease-app/pkg/auth"
	emailsmtp "github.com/AkhmadeevRus/sublease-app/pkg/email_smtp"
	"github.com/AkhmadeevRus/sublease-app/pkg/property"
	"github.com/AkhmadeevRus/sublease-app/pkg/search"
)

type Service struct {
	AuthService      auth.IAuthService
	PropertyService  property.IPropertyService
	SearchService    search.ISearchService
	EmailSmtpService emailsmtp.IEmailSmtpService
}

func NewService(repos *Repository, emailSmtpService *emailsmtp.EmailSmtpService) *Service {
	return &Service{
		AuthService:      auth.NewAuthService(repos.AuthRepository, emailSmtpService),
		PropertyService:  property.NewPropertyService(repos.PropertyRepository),
		SearchService:    search.NewSearchService(repos.SearchRepository),
		EmailSmtpService: emailSmtpService,
	}
}
