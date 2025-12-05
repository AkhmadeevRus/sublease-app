package emailsmtp

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type EmailCfg struct {
	OwnerEmail    string
	OwnerPassword string
	Address       string
	CodeLength    int
	CodeExp       time.Duration
}

type EmailSmtpRepository struct {
	db  *sqlx.DB
	cfg *EmailCfg
}

func NewEmailCfg(ownerEmail, ownerPassword, addr string, codeLength int, codeExp time.Duration) *EmailCfg {
	return &EmailCfg{OwnerEmail: ownerEmail,
		OwnerPassword: ownerPassword,
		Address:       addr,
		CodeLength:    codeLength,
		CodeExp:       codeExp,
	}
}

type IEmailSmtpRepository interface {
	CheckEmailConfirm(email string) (bool, error)
	ConfirmEmail(email string) error
	SendConfirmEmailMessage(email, code string) error
	SendMessage(email, messageText, title string) error
	GenerateConfirmCode() string
}

func NewEmailSmtpRepository(db *sqlx.DB, cfg *EmailCfg) *EmailSmtpRepository {
	return &EmailSmtpRepository{db: db, cfg: cfg}
}

func (r *EmailSmtpRepository) CheckEmailConfirm(email string) (bool, error) {
	var status bool
	sql, args, err := sq.Select("confirmed_email").
		From("users").
		Where(sq.Eq{"email": email}).
		ToSql()
	if err != nil {
		return false, err
	}
	err = r.db.Get(&status, sql, args...)
	if err != nil {
		return false, err
	}
	return status, err
}

func (r *EmailSmtpRepository) ConfirmEmail(email string) error {
	sql, args, err := sq.Update("users").
		Set("confirmed_email", true).
		Where(sq.Eq{"email": email}).
		ToSql()
	if err != nil {
		return err
	}
	res, err := r.db.Exec(sql, args...)
	n, _ := res.RowsAffected()
	if n == 0 && err == nil {
		return errors.New("alredy confirmed")
	}
	return err
}

func (r *EmailSmtpRepository) SendConfirmEmailMessage(email, code string) error {
	baseText := `confirm your email with this code: %s. 
			If you don't ask this code just ignore this message.`
	text := fmt.Sprintf(baseText, code)
	subject := fmt.Sprintf("Email confirm code %s", code)
	return r.SendMessage(email, text, subject)
}

func (r *EmailSmtpRepository) SendMessage(email, messageText, title string) error {
	toEmail := email
	fromEmail := r.cfg.OwnerEmail
	subjectBody := fmt.Sprintf("Subject:%s\n\n %s", title, messageText)
	status := smtp.SendMail(
		r.cfg.Address,
		smtp.PlainAuth("", fromEmail, r.cfg.OwnerPassword, "smtp.gmail.com"),
		fromEmail,
		[]string{toEmail},
		[]byte(subjectBody),
	)
	return status
}

func (r *EmailSmtpRepository) GenerateConfirmCode() string {
	code := rand.Intn(int(math.Pow10(r.cfg.CodeLength)))
	return strconv.Itoa(code)
}
