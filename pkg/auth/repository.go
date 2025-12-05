package auth

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type IAuthRepository interface {
	CreateUser(user User) error
	GetUser(username, password string) (User, error)
}

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(user User) error {
	var usernameCount int
	sql, args, err := sq.Select("count(*)").
		From("users").
		Where(sq.Eq{"username": user.Username}).
		ToSql()
	if err != nil {
		return err
	}
	err = r.db.Get(&usernameCount, sql, args...)
	if err != nil {
		return err
	}
	if usernameCount > 0 {
		return fmt.Errorf("this username is already in use")
	}

	var emailCount int
	sql, args, err = sq.Select("count(*)").
		From("users").
		Where(sq.Eq{"email": user.Email}).
		ToSql()
	if err != nil {
		return err
	}
	err = r.db.Get(&emailCount, sql, args...)
	if err != nil {
		return err
	}
	if emailCount > 0 {
		return fmt.Errorf("this email is already in use")
	}

	user.Id = uuid.New()
	sql, args, err = sq.Insert("users").
		Columns("id", "name", "username", "password", "phone", "role", "email").
		Values(user.Id, user.Name, user.Username, user.Password, user.Phone, user.Role, user.Email).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(sql, args...)
	return err
}

func (r *AuthRepository) GetUser(username, password string) (User, error) {
	var user User
	sql, args, err := sq.Select("id").
		From("users").
		Where(sq.Eq{"username": username, "password": password}).
		ToSql()
	if err != nil {
		return User{}, err
	}
	err = r.db.Get(&user, sql, args...)
	return user, err
}
