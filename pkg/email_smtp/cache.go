package emailsmtp

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type IEmailCacheRepository interface {
	GetConfirmCode(Email string) (string, time.Duration, error)
	SaveConfirmCode(email, code string) error
	DeleteConfirmCode(email string) error
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
	return err
}

func (c *EmailCacheRepository) GetConfirmCode(email string) (string, time.Duration, error) {
	var ttl time.Duration
	ctx := context.Background()
	code, err := c.db.Get(ctx, email).Result()
	if err != nil {
		return "", ttl, err
	}
	ttl, err = c.db.TTL(ctx, email).Result()
	return code, ttl, err
}

func (c *EmailCacheRepository) DeleteConfirmCode(email string) error {
	key := fmt.Sprintf("email:confirm:%s", email)
	ctx := context.Background()
	return c.db.Del(ctx, key).Err()
}
