package auth

import (
	"fmt"

	"github.com/AkhmadeevRus/sublease-app/pkg/apperror"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type IAuthRepository interface {
	CreateUser(user User) error
	GetUser(username, password string) (User, error)
	GetUserByEmail(email string) (User, error)
	UpdatePassword(email, password string) error
}

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(user User) error {
	var usernameCount int
	sql, args, err := sq.Select("COUNT(username)").
		From("users").
		Where(sq.Eq{"username": user.Username}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return apperror.NewInternalError(fmt.Errorf("err in build sql query"))
	}
	err = r.db.Get(&usernameCount, sql, args...)
	if err != nil {
		return apperror.NewInternalError(fmt.Errorf("err in build sql query"))
	}
	if usernameCount > 0 {
		return apperror.NewConflictError("this username is alredy in use", "USERNAME_EXISTS")
	}

	var emailCount int
	sql, args, err = sq.Select("count(*)").
		From("users").
		Where(sq.Eq{"email": user.Email}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return apperror.NewInternalError(fmt.Errorf("err in build sql query"))
	}
	err = r.db.Get(&emailCount, sql, args...)
	if err != nil {
		return apperror.NewInternalError(fmt.Errorf("err in build sql query"))
	}
	if emailCount > 0 {
		return apperror.NewConflictError("this email is alredy in use", "EMAIL_EXISTS")
	}

	user.Id = uuid.New()
	sql, args, err = sq.Insert("users").
		Columns("id", "name", "username", "password", "phone", "role", "email").
		Values(user.Id, user.Name, user.Username, user.Password, user.Phone, user.Role, user.Email).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return apperror.NewInternalError(fmt.Errorf("err in build sql query"))
	}
	_, err = r.db.Exec(sql, args...)
	if err != nil {
		return apperror.NewInternalError(fmt.Errorf("err to create user:%s", err.Error()))
	}
	return nil
}

func (r *AuthRepository) GetUser(username, password string) (User, error) {
	var user User
	sql, args, err := sq.Select("*").
		From("users").
		Where(sq.Eq{"username": username, "password": password}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return User{}, apperror.NewInternalError(fmt.Errorf("err in build sql query"))
	}
	err = r.db.Get(&user, sql, args...)
	if err != nil {
		return User{}, apperror.NewNotFoundError("user not found", "USER_NOT_FOUND")
	}
	return user, nil
}

func (r *AuthRepository) GetUserByEmail(email string) (User, error) {
	var user User
	sql, args, err := sq.Select("*").
		From("users").
		Where(sq.Eq{"email": email}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return User{}, apperror.NewInternalError(fmt.Errorf("err in build sql query"))
	}

	err = r.db.Get(&user, sql, args...)
	if err != nil {
		return User{}, apperror.NewNotFoundError("user not found", "USER_NOT_FOUND")
	}
	return user, nil
}

func (r *AuthRepository) UpdatePassword(email, password string) error {
	sql, args, err := sq.Update("users").
		Set("password", password).
		Where(sq.Eq{"email": email}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return apperror.NewInternalError(fmt.Errorf("err in build sql query"))
	}
	result, err := r.db.Exec(sql, args...)
	if err != nil {
		return apperror.NewInternalError(fmt.Errorf("err to update password user:%s", err.Error()))
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return apperror.NewNotFoundError("user not found", "USER_NOT_FOUND")
	}
	return nil
}
