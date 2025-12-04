package auth

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"time"

	emailsmtp "github.com/AkhmadeevRus/sublease-app/pkg/email_smtp"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	TokeknTTL = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId uuid.UUID `json:"user_id"`
}

type IAuthService interface {
	CreateUser(user User) error
	GenerateToken(username, password string) (string, error)
	ParseToken(accessToken string) (uuid.UUID, error)
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
		logrus.Errorf("err while send confirm email message:%s", err.Error())
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
		return "", fmt.Errorf("email not confirmed")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokeknTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (s *AuthService) ParseToken(accessToken string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signin method")
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return uuid.Nil, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(os.Getenv("PASSWORD_SALT"))))
}
