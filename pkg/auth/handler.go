package auth

import (
	"net/http"

	emailsmtp "github.com/AkhmadeevRus/sublease-app/pkg/email_smtp"
	"github.com/gin-gonic/gin"
)

type IAuthHandler interface {
	SignUp(c *gin.Context)
	SignIn(c *gin.Context)
	UserIdentity(c *gin.Context)
	ConfirmEmail(c *gin.Context)
	ResendConfirmEmail(c *gin.Context)
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
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.service.CreateUser(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
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
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.service.GenerateToken(input.Username, input.Password)
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
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
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.emailService.ConfirmEmail(input.Email, input.Code)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
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
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.emailService.SendConfirmEmailMessage(input.Email)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}
