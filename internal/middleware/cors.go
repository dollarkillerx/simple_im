package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Token, X-Requested-With")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}
