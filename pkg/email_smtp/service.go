package emailsmtp

import (
	"fmt"
	"os"
	"time"

	"github.com/AkhmadeevRus/sublease-app/pkg/apperror"
)

type IEmailSmtpService interface {
	CheckEmailConfirm(email string) (bool, error)
	ConfirmEmail(email, code string) error
	SendConfirmEmailMessage(email string) error
	SendPasswordResetEmail(email string) error
	SendMessage(email, messageText, title string) error
	GenerateConfirmCode() string
	GetPasswordResetCode(email string) (string error)
}

type EmailSmtpService struct {
	repo  IEmailSmtpRepository
	cache IEmailCacheRepository
}

func NewEmailSmtpService(repo IEmailSmtpRepository, cache IEmailCacheRepository) *EmailSmtpService {
	return &EmailSmtpService{repo: repo, cache: cache}
}

func (s *EmailSmtpService) CheckEmailConfirm(email string) (bool, error) {
	return s.repo.CheckEmailConfirm(email)
}

func (s *EmailSmtpService) ConfirmEmail(email, code string) error {
	trueCode, _, err := s.cache.GetConfirmCode(email)
	if err != nil {
		return err
	}
	if trueCode != code {
		return apperror.NewBadRequestError("bad confirm email code", "INVALID_CODE")
	}
	return s.repo.ConfirmEmail(email)
}

func (s *EmailSmtpService) SendConfirmEmailMessage(email string) error {
	minTtl, _ := time.ParseDuration(os.Getenv("MIN_TTL"))
	maxTtl, _ := time.ParseDuration(os.Getenv("MAX_TTL"))
	_, ttl, err := s.cache.GetConfirmCode(email)
	if err == nil {
		return err
	} else if minTtl < ttl {
		return apperror.NewBadRequestError(fmt.Sprintf("code has already been sent %s ago", maxTtl-ttl), "CODE_TOO_SOON")
	}
	code := s.GenerateConfirmCode()
	err = s.cache.SaveConfirmCode(email, code)
	if err != nil {
		return err
	}
	err = s.repo.SendConfirmEmailMessage(email, code)
	if err != nil {
		return err
	}
	return nil
}
func (s *EmailSmtpService) SendMessage(email, messageText, title string) error {
	return s.repo.SendMessage(email, messageText, title)
}
func (s *EmailSmtpService) GenerateConfirmCode() string {
	return s.repo.GenerateConfirmCode()
}

func (s *EmailSmtpService) GetPasswordResetCode(email string) (string, error) {
	code, _, err := s.cache.GetPasswordResetCode(email)
	if err != nil {
		return "", err
	}
	return code, nil
}

func (s *EmailSmtpService) SendPasswordResetEmail(email string) error {
	minTtl, _ := time.ParseDuration(os.Getenv("MIN_TTL"))
	maxTtl, _ := time.ParseDuration(os.Getenv("MAX_TTL"))
	_, ttl, err := s.cache.GetPasswordResetCode(email)
	if err != nil {
		return err
	} else if minTtl < ttl {
		return apperror.NewBadRequestError(fmt.Sprintf("code has already been sent %s ago", maxTtl-ttl), "CODE_TOO_SOON")
	}
	code := s.GenerateConfirmCode()
	err = s.cache.SavePasswordResetCode(email, code)
	if err != nil {
		return err
	}
	err = s.repo.SendPasswordResetEmailMessage(email, code)
	if err != nil {
		return err
	}
	return nil
}
