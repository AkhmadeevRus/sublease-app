package auth

import (
	"net/http"

	"github.com/AkhmadeevRus/sublease-app/pkg/apperror"
	emailsmtp "github.com/AkhmadeevRus/sublease-app/pkg/email_smtp"
	"github.com/gin-gonic/gin"
)

type IAuthHandler interface {
	SignUp(c *gin.Context)
	SignIn(c *gin.Context)
	UserIdentity(c *gin.Context)
	ConfirmEmail(c *gin.Context)
	ResendConfirmEmail(c *gin.Context)
	RequestPasswordReset(c *gin.Context)
	UpdatePassword(c *gin.Context)
}

type AuthHandler struct {
	service      IAuthService
	emailService emailsmtp.IEmailSmtpService
}

func NewAuthHandler(service IAuthService, emailService emailsmtp.IEmailSmtpService) *AuthHandler {
	return &AuthHandler{service: service, emailService: emailService}
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var input User

	if err := c.BindJSON(&input); err != nil {
		apperror.NewErrorResponse(c, apperror.NewBadRequestError(err.Error(), "INVALID_INPUT"))
		return
	}

	err := h.service.CreateUser(input)
	if err != nil {
		apperror.NewErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		apperror.NewErrorResponse(c, apperror.NewBadRequestError(err.Error(), "INVALID_INPUT"))
		return
	}

	token, err := h.service.GenerateToken(input.Username, input.Password)
	if err != nil {
		apperror.NewErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

type confirmCodeInput struct {
	Email string `json:"email" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

func (h *AuthHandler) ConfirmEmail(c *gin.Context) {
	var input confirmCodeInput
	if err := c.BindJSON(&input); err != nil {
		apperror.NewErrorResponse(c, apperror.NewBadRequestError(err.Error(), "INVALID_INPUT"))
		return
	}

	err := h.emailService.ConfirmEmail(input.Email, input.Code)
	if err != nil {
		apperror.NewErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}

type ResendConfirmEmailInput struct {
	Email string `json:"email" binding:"required"`
}

func (h *AuthHandler) ResendConfirmEmail(c *gin.Context) {
	var input ResendConfirmEmailInput
	if err := c.BindJSON(&input); err != nil {
		apperror.NewErrorResponse(c, apperror.NewBadRequestError(err.Error(), "INVALID_INPUT"))
		return
	}

	err := h.emailService.SendConfirmEmailMessage(input.Email)
	if err != nil {
		apperror.NewErrorResponse(c, err)
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}

type requestPasswordResetInput struct {
	Email string `json:"email" binding:"required"`
}

func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	var input requestPasswordResetInput
	if err := c.BindJSON(&input); err != nil {
		apperror.NewErrorResponse(c, apperror.NewBadRequestError(err.Error(), "INVALID_INPUT"))
		return
	}

	err := h.emailService.SendPasswordResetEmail(input.Email)
	if err != nil {
		apperror.NewErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}

type ResetPasswordInput struct {
	Email       string `json:"email" binding:"required"`
	Code        string `json:"code" binding:"required"`
	NewPassword string `json:"password" binding:"required"`
}

func (h *AuthHandler) UpdatePassword(c *gin.Context) {
	var input ResetPasswordInput
	if err := c.BindJSON(&input); err != nil {
		apperror.NewErrorResponse(c, apperror.NewBadRequestError(err.Error(), "INVALID_INPUT"))
		return
	}

	trueCode, err := h.emailService.GetPasswordResetCode(input.Email)
	if err != nil {
		apperror.NewErrorResponse(c, err)
	} else if trueCode != input.Code {

	}

	err := h.service.UpdatePassword(input.Email, input.Code, input.NewPassword)
	if err != nil {
		apperror.NewErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}
