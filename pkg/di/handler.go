package di

import (
	"github.com/AkhmadeevRus/sublease-app/pkg/auth"
	"github.com/AkhmadeevRus/sublease-app/pkg/property"
	"github.com/AkhmadeevRus/sublease-app/pkg/search"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	AuthHandler     auth.IAuthHandler
	PropertyHandler property.IPropertyHandler
	SearchHandler   search.ISearchHandler
}

func NewHandler(services *Service) *Handler {
	return &Handler{
		AuthHandler:     auth.NewAuthHandler(services.AuthService, services.EmailSmtpService),
		PropertyHandler: property.NewPropertyHandler(services.PropertyService),
		SearchHandler:   search.NewSearchHandler(services.SearchService),
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.AuthHandler.SignUp)
		auth.POST("/sign-in", h.AuthHandler.SignIn)
		auth.POST("/confirm-email", h.AuthHandler.ConfirmEmail)
		auth.POST("/resend-confirm-email", h.AuthHandler.ResendConfirmEmail)
	}
	api := router.Group("/api", h.AuthHandler.UserIdentity)
	{
		property := api.Group("/property")
		{
			property.POST("/", h.PropertyHandler.CreateProperty)
			property.PUT("/:id", h.PropertyHandler.UpdateProperty)
			property.DELETE("/:id", h.PropertyHandler.DeleteProperty)
		}
		search := api.Group("/search")
		{
			search.GET("/", h.SearchHandler.GetAllProperties)
			search.GET("/:id", h.SearchHandler.GetPropertyById)
			search.POST("/", h.SearchHandler.SearchByFilters)
		}
	}
	return router
}
