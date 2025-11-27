package auth

import "github.com/google/uuid"

type UserRole string

const (
	RoleLandlord UserRole = "landlord"
	RoleBuyer    UserRole = "buyer"
)

type User struct {
	Id               uuid.UUID `json:"-" db:"id"`
	Username         string    `json:"username" binding:"required"`
	Name             string    `json:"name" binding:"required"`
	Password         string    `json:"password" binding:"required"`
	Phone            string    `json:"phone" binding:"required"`
	Email            string    `json:"email" binding:"required" db:"email"`
	IsEmailConfirmed bool      `json:"-" db:"confirmed_email"`
	Role             UserRole  `json:"role" binding:"required" db:"role"`
}
