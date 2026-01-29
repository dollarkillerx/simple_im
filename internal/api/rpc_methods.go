package api

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

// Rpc is the main JSON-RPC endpoint handler
func (a *ApiServer) Rpc(ctx *gin.Context) {
	a.rpcHandler.HandleRpcRequest(ctx)
}

// PingMethod - health check method
type PingMethod struct{}

func (m *PingMethod) Name() string { return "ping" }

func (m *PingMethod) Execute(ctx context.Context, params json.RawMessage) (interface{}, error) {
	return map[string]interface{}{
		"pong":    true,
		"time":    time.Now().Unix(),
		"message": "pong",
	}, nil
}

func (m *PingMethod) RequireAuth() bool { return false }
