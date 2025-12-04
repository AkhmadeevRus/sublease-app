package property

import (
	"net/http"

	"github.com/AkhmadeevRus/sublease-app/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IPropertyHandler interface {
	CreateProperty(ctx *gin.Context)
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
		auth.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input Property
	if err := c.BindJSON(&input); err != nil {
		auth.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.CreateProperty(input, userId)
	if err != nil {
		auth.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}

type GetMyPropertiesResponse struct {
	Data []Property `json:"data"`
}

func (h *PropertyHandler) DeleteProperty(c *gin.Context) {
	userId, err := auth.GetUserId(c)
	if err != nil {
		auth.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id := c.Param("id")
	if id == "" {
		auth.NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	propertyId, err := uuid.Parse(id)
	if err != nil {
		auth.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.DeleteProperty(userId, propertyId)
	if err != nil {
		auth.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}

func (h *PropertyHandler) UpdateProperty(c *gin.Context) {
	userId, err := auth.GetUserId(c)
	if err != nil {
		auth.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id := c.Param("id")
	if id == "" {
		auth.NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	propertyId, err := uuid.Parse(id)
	if err != nil {
		auth.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input PropertyUpdate
	if err := c.BindJSON(&input); err != nil {
		auth.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.service.UpdateProperty(userId, propertyId, input)
	if err != nil {
		auth.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}
