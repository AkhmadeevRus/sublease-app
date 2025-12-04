package di

import (
	"github.com/AkhmadeevRus/sublease-app/pkg/auth"
	emailsmtp "github.com/AkhmadeevRus/sublease-app/pkg/email_smtp"
	"github.com/AkhmadeevRus/sublease-app/pkg/property"
	"github.com/AkhmadeevRus/sublease-app/pkg/search"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	AuthRepository           auth.IAuthRepository
	PropertyRepository       property.IPropertyRepository
	SearchRepository         search.ISearchRepository
	EmailSmtpRepository      emailsmtp.IEmailSmtpRepository
	EmailSmtpCacheRepository emailsmtp.IEmailCacheRepository
}

func NewRepository(db *sqlx.DB, cacheDb *redis.Client, emailCfg *emailsmtp.EmailCfg) *Repository {
	return &Repository{
		AuthRepository:           auth.NewAuthRepository(db),
		PropertyRepository:       property.NewPropertyRepository(db),
		SearchRepository:         search.NewSearchRepository(db),
		EmailSmtpRepository:      emailsmtp.NewEmailSmtpRepository(db, emailCfg),
		EmailSmtpCacheRepository: emailsmtp.NewEmailCacheRepository(cacheDb, emailCfg.CodeExp),
	}
}
