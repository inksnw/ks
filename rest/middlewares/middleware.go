package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Recv() func(c *gin.Context) {
	return func(c *gin.Context) {
		defer func() {
			if e := recover(); e != nil {
				c.AbortWithStatusJSON(400, e)
				return
			}
		}()
		c.Next()
	}
}

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//直接取前端的地址允许跨域
		allow := c.Request.Header.Get("Origin")
		c.Writer.Header().Set("Access-Control-Allow-Origin", allow)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Request.Header.Del("Origin")
		c.Next()
	}
}
