package auth

import (
	"fmt"
	"strings"

	"github.com/AkhmadeevRus/sublease-app/pkg/apperror"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *AuthHandler) UserIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		apperror.NewErrorResponse(c, apperror.NewUnauthorizedError("authorization header is required", "MISSING_AUTH_HEADER"))
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		apperror.NewErrorResponse(c, apperror.NewUnauthorizedError("invalid authorization header format", "INVALID_AUTH_FORMAT"))
		return
	}

	userId, err := h.service.ParseToken(headerParts[1])
	if err != nil {
		apperror.NewErrorResponse(c, err)
		return
	}
	c.Set(userCtx, userId)
}

func GetUserId(c *gin.Context) (uuid.UUID, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		err := apperror.NewInternalError(fmt.Errorf("user id not found in context")).WithCode("USER_ID_MISSING")
		return uuid.Nil, err
	}

	idInt, ok := id.(uuid.UUID)
	if !ok {
		err := apperror.NewInternalError(fmt.Errorf("user id has invalid type in context")).WithCode("INVALID_USER_ID_TYPE")
		return uuid.Nil, err
	}

	return idInt, nil
}
