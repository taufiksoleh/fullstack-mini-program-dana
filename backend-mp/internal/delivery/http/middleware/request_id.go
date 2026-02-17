package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taufiksoleh/backend-mp/pkg/reqctx"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" {
			id = uuid.New().String()
		}
		c.Set("request_id", id)
		c.Header("X-Request-ID", id)

		c.Request = c.Request.WithContext(reqctx.WithRequestID(c.Request.Context(), id))

		c.Next()
	}
}
