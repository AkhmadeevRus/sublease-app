package search

import (
	"github.com/gin-gonic/gin"
)

type ISearchHandler interface {
	GetAllProperties(c *gin.Context)
	SearchByFilters(c *gin.Context)
}

type SearchHandler struct {
	service ISearchService
}

func NewSearchHandler(service ISearchService) *SearchHandler {
	return &SearchHandler{service: service}
}

func (h *SearchHandler) GetAllProperties(c *gin.Context) {
}

func (h *SearchHandler) SearchByFilters(c *gin.Context) {
}
