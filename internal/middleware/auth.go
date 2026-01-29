package middleware

import (
	"fmt"
	"strings"

	"simple_im/pkg/common/jwt"
	"simple_im/pkg/common/resp"

	"github.com/gin-gonic/gin"
)

const (
	ContextKeyUserID   = "user_id"
	ContextKeyUsername = "username"
)

func JWTAuth(jwtManager *jwt.JWTManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			resp.Return(ctx, 200, "", nil, fmt.Errorf("authorization header required"))
			ctx.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			resp.Return(ctx, 200, "", nil, fmt.Errorf("invalid authorization header format"))
			ctx.Abort()
			return
		}

		claims, err := jwtManager.ParseToken(parts[1])
		if err != nil {
			resp.Return(ctx, 200, "", nil, fmt.Errorf("invalid token: %v", err))
			ctx.Abort()
			return
		}

		ctx.Set(ContextKeyUserID, claims.UserID)
		ctx.Set(ContextKeyUsername, claims.Username)
		ctx.Next()
	}
}

func GetUserID(ctx *gin.Context) int64 {
	if v, exists := ctx.Get(ContextKeyUserID); exists {
		if userID, ok := v.(int64); ok {
			return userID
		}
	}
	return 0
}

func GetUsername(ctx *gin.Context) string {
	if v, exists := ctx.Get(ContextKeyUsername); exists {
		if username, ok := v.(string); ok {
			return username
		}
	}
	return ""
}
