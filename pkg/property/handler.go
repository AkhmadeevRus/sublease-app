package property

import (
	"net/http"

	"github.com/AkhmadeevRus/sublease-app/pkg/auth"
	"github.com/gin-gonic/gin"
)

type IPropertyHandler interface {
	CreateProperty(ctx *gin.Context)
	GetMyProperties(ctx *gin.Context)
	DeleteProperty(ctx *gin.Context)
	UpdateProperty(ctx *gin.Context)
}

type PropertyHandler struct {
	service IPropertyService
}

func NewPropertyHandler(service IPropertyService) *PropertyHandler {
	return &PropertyHandler{service: service}
}

func (h *PropertyHandler) CreateProperty(c *gin.Context) {
	userId, err := auth.GetUserId(c)
	if err != nil {
		return
	}

	var input Property
	if err := c.BindJSON(&input); err != nil {
		auth.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.CreateProperty(input, userId)
	if err != nil {
		auth.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *PropertyHandler) GetMyProperties(ctx *gin.Context) {
}

func (h *PropertyHandler) DeleteProperty(ctx *gin.Context) {
}

func (h *PropertyHandler) UpdateProperty(ctx *gin.Context) {
}
