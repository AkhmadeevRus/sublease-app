package auth

import (
	"fmt"

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
	query := "SELECT COUNT(*) FROM users WHERE username = $1"
	err := r.db.Get(&usernameCount, query, user.Username)
	if err != nil {
		return err
	}
	if usernameCount > 0 {
		return fmt.Errorf("this username is already in use")
	}

	var emailCount int
	query = "SELECT COUNT(*) FROM users WHERE email = $1"
	err = r.db.Get(&emailCount, query, user.Email)
	if err != nil {
		return err
	}
	if emailCount > 0 {
		return fmt.Errorf("this email is already in use")
	}

	user.Id = uuid.New()
	query = `INSERT INTO users (id, name, username, password, phone, role, email) values ($1, $2, $3, $4, $5, $6, $7)`
	_, err = r.db.Exec(query, user.Id, user.Name, user.Username, user.Password, user.Phone, user.Role, user.Email)
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthRepository) GetUser(username, password string) (User, error) {
	var user User
	query := "SELECT id FROM users WHERE username=$1 AND password=$2"
	err := r.db.Get(&user, query, username, password)
	return user, err
}
