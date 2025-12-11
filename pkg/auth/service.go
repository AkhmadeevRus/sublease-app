package auth

import (
	"crypto/sha1"
	"fmt"
	"os"
	"time"

	"github.com/AkhmadeevRus/sublease-app/pkg/apperror"
	emailsmtp "github.com/AkhmadeevRus/sublease-app/pkg/email_smtp"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId uuid.UUID `json:"user_id"`
}

type IAuthService interface {
	CreateUser(user User) error
	GenerateToken(username, password string) (string, error)
	ParseToken(accessToken string) (uuid.UUID, error)
	GetUserByEmail(email string) (User, error)
	UpdatePassword(email, code, newPassword string) error
}

type AuthService struct {
	repo         IAuthRepository
	emailService emailsmtp.IEmailSmtpService
}

func NewAuthService(repo IAuthRepository, emailService emailsmtp.IEmailSmtpService) *AuthService {
	return &AuthService{repo: repo, emailService: emailService}
}

func (s *AuthService) CreateUser(user User) error {
	user.Password = s.generatePasswordHash(user.Password)
	err := s.repo.CreateUser(user)
	if err != nil {
		return err
	}
	if err := s.emailService.SendConfirmEmailMessage(user.Email); err != nil {
		return err
	}
	return nil
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, s.generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	isConfirmed, err := s.emailService.CheckEmailConfirm(user.Email)
	if err != nil {
		return "", err
	}

	if !isConfirmed {
		return "", apperror.NewUnauthorizedError("email not confirmed", "EMAIL_NOT_CONFIRMED")
	}

	tokenTTL, err := time.ParseDuration(os.Getenv("TOKEN_TTL"))
	if err != nil {
		return "", apperror.NewInternalError(fmt.Errorf("filed to parse duration(tokenTTL):%w", err))
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (s *AuthService) ParseToken(accessToken string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperror.NewUnauthorizedError("invalid signing method", "INVALID_TOKEN")
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return uuid.Nil, apperror.NewUnauthorizedError(err.Error(), "INVALID_TOKEN")
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return uuid.Nil, apperror.NewUnauthorizedError("invalid token claims", "INVALID_TOKEN")
	}

	return claims.UserId, nil
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(os.Getenv("PASSWORD_SALT"))))
}

func (s *AuthService) GetUserByEmail(email string) (User, error) {
	return s.repo.GetUserByEmail(email)
}

func (s *AuthService) UpdatePassword(email, code, password string) error {
	passwordHash := s.generatePasswordHash(password)
	return s.repo.UpdatePassword(email, passwordHash)
}
