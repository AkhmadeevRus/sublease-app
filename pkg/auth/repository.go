package auth

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type IAuthRepository interface {
	CreateUser(user User) (uuid.UUID, error)
	GetUser(username, password string) (User, error)
}

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(user User) (uuid.UUID, error) {
	var id uuid.UUID
	query := `INSERT INTO users (name, username, password, phone, role, email) values ($1, $2, $3, $4, $5, $6) RETURNING id`
	row := r.db.QueryRow(query, user.Name, user.Username, user.Password, user.Phone, user.Role, user.Email)
	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *AuthRepository) GetUser(username, password string) (User, error) {
	var user User
	query := "SELECT id FROM users WHERE username=$1 AND password=$2"
	err := r.db.Get(&user, query, username, password)
	return user, err
}
