package api

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"simple_im/internal/storage"
	"simple_im/internal/ws"
	"simple_im/pkg/common/jwt"
	"simple_im/pkg/common/resp"

	"github.com/gin-gonic/gin"
)

type RpcMethod interface {
	Name() string
	Execute(ctx context.Context, params json.RawMessage) (interface{}, error)
	RequireAuth() bool
}

type RpcHandler struct {
	methods    map[string]RpcMethod
	mu         sync.RWMutex
	storage    *storage.Storage
	hub        *ws.Hub
	jwtManager *jwt.JWTManager
}

func NewRpcHandler(storage *storage.Storage, hub *ws.Hub, jwtManager *jwt.JWTManager) *RpcHandler {
	return &RpcHandler{
		methods:    make(map[string]RpcMethod),
		storage:    storage,
		hub:        hub,
		jwtManager: jwtManager,
	}
}

func (h *RpcHandler) RegisterMethod(method RpcMethod) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.methods[method.Name()] = method
}

func (h *RpcHandler) getMethod(name string) (RpcMethod, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	method, ok := h.methods[name]
	return method, ok
}

func (h *RpcHandler) HandleRpcRequest(ctx *gin.Context) {
	var req resp.RpcRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.InvalidRequest(ctx, "", "Invalid JSON")
		return
	}

	if req.JsonRPC != "2.0" {
		resp.InvalidRequest(ctx, req.Id, "Invalid JSON-RPC version")
		return
	}

	method, ok := h.getMethod(req.Method)
	if !ok {
		resp.MethodNotFound(ctx, req.Id, req.Method)
		return
	}

	// Create context with request info
	rpcCtx := context.Background()

	// Check authentication if required
	if method.RequireAuth() {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			resp.ErrorReturn(ctx, req.Id, fmt.Errorf("authorization required"))
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			resp.ErrorReturn(ctx, req.Id, fmt.Errorf("invalid authorization format"))
			return
		}

		claims, err := h.jwtManager.ParseToken(parts[1])
		if err != nil {
			resp.ErrorReturn(ctx, req.Id, fmt.Errorf("invalid token: %v", err))
			return
		}

		rpcCtx = context.WithValue(rpcCtx, "user_id", claims.UserID)
		rpcCtx = context.WithValue(rpcCtx, "username", claims.Username)
	}

	result, err := method.Execute(rpcCtx, req.Params)
	if err != nil {
		resp.ErrorReturn(ctx, req.Id, err)
		return
	}

	resp.SuccessReturn(ctx, req.Id, result)
}
