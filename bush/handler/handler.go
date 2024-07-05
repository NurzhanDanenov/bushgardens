package handler

import (
	"bush/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	userHandler := router.Group("/user")
	{
		userHandler.POST("/register", h.Register)
		userHandler.POST("/login", h.Login)
	}

	apiHandler := router.Group("/api", h.userIdentity)
	{
		imageHandler := apiHandler.Group("/images")
		{
			imageHandler.POST("/uploading", h.createAva)
			imageHandler.GET("/downloading", h.getAllImages)
			//imageHandler.GET("/:id", h.getImageById)
			//imageHandler.DELETE("/:id", h.deleteImage)
			//imageHandler.GET("/metrics", gin.WrapH(promhttp.Handler()))
		}
	}
	return router
}
