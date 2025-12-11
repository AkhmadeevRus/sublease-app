package apperror

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Error   string `json:"error"`
	Type    string `json:"type"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

func NewErrorResponse(c *gin.Context, err error) {

	var appErr *AppError
	if errors.As(err, &appErr) {
		status := GetHttpStatusByErrorType(appErr.Type)

		c.AbortWithStatusJSON(status, errorResponse{
			Error:   appErr.Error(),
			Type:    string(appErr.Type),
			Code:    appErr.Code,
			Message: appErr.Message,
		})
		return
	}

	appErr = NewInternalError(err)

	c.AbortWithStatusJSON(GetHttpStatusByErrorType(appErr.Type), errorResponse{
		Error:   appErr.Error(),
		Type:    string(appErr.Type),
		Message: appErr.Message,
	})
}
