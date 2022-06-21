package middleware

import (
	"github.com/fkcs/gateway/internal/infrastructure/logger"
	errord "github.com/fkcs/gateway/internal/utils/error"
	"github.com/fkcs/gateway/internal/utils/types"
	"github.com/gin-gonic/gin"
	"time"
)

func CostTime(c *gin.Context) {
	logger.Logger().Infof("start")
	start := time.Now()
	c.Next()
	time.Since(start)
	logger.Logger().Infof("end")
}

func InterfaceInvalid(c *gin.Context) {
	logger.Logger().Infof("%v", c.Request.URL.Path)
	errord.BadRequest(types.InterfaceInvalid, c)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
