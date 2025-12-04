package search

import (
	"net/http"

	"github.com/AkhmadeevRus/sublease-app/pkg/auth"
	"github.com/AkhmadeevRus/sublease-app/pkg/property"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ISearchHandler interface {
	GetAllProperties(c *gin.Context)
	GetPropertyById(c *gin.Context)
	SearchByFilters(c *gin.Context)
}

type SearchHandler struct {
	service ISearchService
}

func NewSearchHandler(service ISearchService) *SearchHandler {
	return &SearchHandler{service: service}
}

func (h *SearchHandler) GetAllProperties(c *gin.Context) {
	myProperties, err := h.service.GetAllProperties()
	if err != nil {
		auth.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, property.GetMyPropertiesResponse{
		Data: myProperties,
	})
}

func (h *SearchHandler) GetPropertyById(c *gin.Context) {
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
	property, err := h.service.GetPropertyById(propertyId)
	if err != nil {
		auth.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, property)
}

func (h *SearchHandler) SearchByFilters(c *gin.Context) {
	var filter PropertyFilter

	if err := c.BindJSON(&filter); err != nil {
		auth.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	properties, err := h.service.SearchByFilters(filter)
	if err != nil {
		auth.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, property.GetMyPropertiesResponse{
		Data: properties,
	})
}
