package emailsmtp

import (
	"context"
	"fmt"
	"time"

	"github.com/AkhmadeevRus/sublease-app/pkg/apperror"
	"github.com/redis/go-redis/v9"
)

type IEmailCacheRepository interface {
	GetConfirmCode(Email string) (string, time.Duration, error)
	GetPasswordResetCode(Email string) (string, time.Duration, error)
	SaveConfirmCode(email, code string) error
	SavePasswordResetCode(email, code string) error
	DeleteConfirmCode(email string) error
	DeletePasswordResetCode(email string) error
}

type EmailCacheRepository struct {
	db      *redis.Client
	CodeExp time.Duration
}

func NewEmailCacheRepository(db *redis.Client, codeExp time.Duration) *EmailCacheRepository {
	return &EmailCacheRepository{db: db, CodeExp: codeExp}
}

func (c *EmailCacheRepository) SaveConfirmCode(email, code string) error {
	ctx := context.Background()
	err := c.db.Set(ctx, email, code, c.CodeExp).Err()
	if err != nil {
		return apperror.NewInternalError(fmt.Errorf("err in set db(save confirm code)"))
	}
	return err
}

func (c *EmailCacheRepository) GetConfirmCode(email string) (string, time.Duration, error) {
	var ttl time.Duration
	ctx := context.Background()
	code, err := c.db.Get(ctx, email).Result()
	if err != nil {
		return "", ttl, apperror.NewInternalError(fmt.Errorf("err in get db confirm code"))
	}
	ttl, err = c.db.TTL(ctx, email).Result()
	return code, ttl, err
}

func (c *EmailCacheRepository) DeleteConfirmCode(email string) error {
	key := fmt.Sprintf("email:confirm:%s", email)
	ctx := context.Background()
	return c.db.Del(ctx, key).Err()
}

func (c *EmailCacheRepository) SavePasswordResetCode(email, code string) error {
	ctx := context.Background()
	key := fmt.Sprintf("reset:%s", email)
	err := c.db.Set(ctx, key, code, 24*time.Hour).Err()
	if err != nil {
		return apperror.NewInternalError(fmt.Errorf("err in set db(save resest password code)"))
	}
	return nil
}

func (c *EmailCacheRepository) GetPasswordResetCode(email string) (string, time.Duration, error) {
	var ttl time.Duration
	ctx := context.Background()
	key := fmt.Sprintf("reset:%s", email)
	code, err := c.db.Get(ctx, key).Result()
	if err != nil {
		return "", ttl, apperror.NewInternalError(fmt.Errorf("err in get db confirm reset password code"))
	}
	ttl, err = c.db.TTL(ctx, key).Result()
	return code, ttl, err
}

func (c *EmailCacheRepository) DeletePasswordResetCode(email string) error {
	key := fmt.Sprintf("password:reset:%s", email)
	ctx := context.Background()
	return c.db.Del(ctx, key).Err()
}
