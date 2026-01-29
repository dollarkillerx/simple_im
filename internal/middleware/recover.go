package middleware

import (
	"net/http"
	"runtime/debug"

	"simple_im/pkg/common/resp"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func HttpRecover() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				stackTrace := string(debug.Stack())
				log.Error().
					Interface("error", r).
					Str("path", ctx.Request.URL.Path).
					Str("stack", stackTrace).
					Msg("panic recovered")

				resp.Return(ctx, http.StatusInternalServerError, "", nil, nil)
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}
