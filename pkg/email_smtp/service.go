package emailsmtp

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type IEmailSmtpService interface {
	CheckEmailConfirm(email string) (bool, error)
	ConfirmEmail(email, code string) error
	SendConfirmEmailMessage(email string) error
	SendMessage(email, messageText, title string) error
	GenerateConfirmCode() string
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
		return errors.New("bad confirm email code")
	}
	return s.repo.ConfirmEmail(email)
}

func (s *EmailSmtpService) SendConfirmEmailMessage(email string) error {
	minTtl, _ := time.ParseDuration(os.Getenv("MIN_TTL"))
	maxTtl, _ := time.ParseDuration(os.Getenv("MAX_TTL"))
	_, ttl, err := s.cache.GetConfirmCode(email)
	if err != nil {
		return errors.New("err in GetConfirmCode")
	} else if minTtl < ttl {
		return fmt.Errorf("code has already been sent %s ago", maxTtl-ttl)
	}
	code := s.GenerateConfirmCode()
	err = s.cache.SaveConfirmCode(email, code)
	if err != nil {
		return errors.New("err while save confirm code")
	}
	go func() {
		err = s.repo.SendConfirmEmailMessage(email, code)
		if err != nil {
			logrus.Errorf("error while sending confirm email message: %s", err.Error())
		}
	}()
	if err != nil && err.Error() == "redis: nil" {
		return nil
	}

	return err
}

func (s *EmailSmtpService) SendMessage(email, messageText, title string) error {
	return s.repo.SendMessage(email, messageText, title)
}

func (s *EmailSmtpService) GenerateConfirmCode() string {
	return s.repo.GenerateConfirmCode()
}
