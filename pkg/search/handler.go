package search

import (
	"net/http"

	"github.com/AkhmadeevRus/sublease-app/pkg/apperror"
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
		apperror.NewErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, property.GetMyPropertiesResponse{
		Data: myProperties,
	})
}

func (h *SearchHandler) GetPropertyById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		apperror.NewErrorResponse(c, apperror.NewBadRequestError("invalid id param", "INVALID_PARAM"))
		return
	}
	propertyId, err := uuid.Parse(id)
	if err != nil {
		apperror.NewErrorResponse(c, apperror.NewInternalError(err))
		return
	}
	property, err := h.service.GetPropertyById(propertyId)
	if err != nil {
		apperror.NewErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, property)
}

func (h *SearchHandler) SearchByFilters(c *gin.Context) {
	var filter PropertyFilter

	if err := c.BindJSON(&filter); err != nil {
		apperror.NewErrorResponse(c, apperror.NewBadRequestError(err.Error(), "INVALID_INPUT"))
		return
	}

	properties, err := h.service.SearchByFilters(filter)
	if err != nil {
		apperror.NewErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, property.GetMyPropertiesResponse{
		Data: properties,
	})
}
