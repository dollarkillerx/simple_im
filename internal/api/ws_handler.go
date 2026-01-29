package api

import (
	"net/http"
	"strings"

	"simple_im/internal/ws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

func (a *ApiServer) WebSocket(ctx *gin.Context) {
	// Get token from query parameter
	token := ctx.Query("token")
	if token == "" {
		// Try to get from Authorization header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}
	}

	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	claims, err := a.jwtManager.ParseToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to upgrade websocket connection")
		return
	}

	client := ws.NewClient(a.hub, conn, claims.UserID, claims.Username)
	a.hub.Register(client)

	go client.WritePump()
	go client.ReadPump()

	log.Info().Int64("user_id", claims.UserID).Str("username", claims.Username).Msg("websocket client connected")
}
