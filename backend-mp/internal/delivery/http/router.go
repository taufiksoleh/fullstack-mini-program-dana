package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/taufiksoleh/backend-mp/internal/delivery/http/handler"
	"github.com/taufiksoleh/backend-mp/internal/delivery/http/middleware"
)

func NewRouter(propertyHandler *handler.PropertyHandler) *gin.Engine {
	r := gin.New()
	r.SetTrustedProxies(nil)
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	{
		properties := v1.Group("/properties")
		properties.GET("", propertyHandler.GetAll)
		properties.GET("/doc/:documentId", propertyHandler.GetByDocumentID)
		properties.GET("/:id", propertyHandler.GetByID)
		properties.POST("", propertyHandler.Create)
		properties.PUT("/:id", propertyHandler.Update)
		properties.DELETE("/:id", propertyHandler.Delete)
	}

	return r
}
