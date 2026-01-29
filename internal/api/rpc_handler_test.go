package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"simple_im/pkg/common/resp"

	"github.com/gin-gonic/gin"
)

func TestRpcHandler_HandleRpcRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	handler := NewRpcHandler(env.Storage, env.Hub, env.JWTManager)
	handler.RegisterMethod(&PingMethod{})

	// Test ping method
	reqBody := resp.RpcRequest{
		JsonRPC: "2.0",
		Method:  "ping",
		Id:      "1",
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/rpc", bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.HandleRpcRequest(ctx)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response resp.RpcResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.Error != nil {
		t.Errorf("Unexpected error: %v", response.Error)
	}
	if response.Id != "1" {
		t.Errorf("Expected ID '1', got '%s'", response.Id)
	}
}

func TestRpcHandler_MethodNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	handler := NewRpcHandler(env.Storage, env.Hub, env.JWTManager)

	reqBody := resp.RpcRequest{
		JsonRPC: "2.0",
		Method:  "nonexistent.method",
		Id:      "1",
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/rpc", bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.HandleRpcRequest(ctx)

	var response resp.RpcResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.Error == nil {
		t.Error("Expected error for non-existent method")
	}
	if response.Error.Code != -32601 {
		t.Errorf("Expected error code -32601, got %d", response.Error.Code)
	}
}

func TestRpcHandler_InvalidJsonRpcVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)

	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	handler := NewRpcHandler(env.Storage, env.Hub, env.JWTManager)

	reqBody := resp.RpcRequest{
		JsonRPC: "1.0", // Invalid version
		Method:  "ping",
		Id:      "1",
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/rpc", bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.HandleRpcRequest(ctx)

	var response resp.RpcResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.Error == nil {
		t.Error("Expected error for invalid JSON-RPC version")
	}
}

func TestRpcHandler_AuthRequired(t *testing.T) {
	gin.SetMode(gin.TestMode)

	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	handler := NewRpcHandler(env.Storage, env.Hub, env.JWTManager)
	handler.RegisterMethod(NewUserInfoMethod(env.Storage))

	// Request without Authorization header
	reqBody := resp.RpcRequest{
		JsonRPC: "2.0",
		Method:  "user.info",
		Id:      "1",
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/rpc", bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.HandleRpcRequest(ctx)

	var response resp.RpcResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.Error == nil {
		t.Error("Expected error for missing authorization")
	}
}

func TestRpcHandler_AuthWithValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	// Create test user
	user, _ := env.CreateTestUser("authuser", "password")

	// Generate token
	token, _ := env.JWTManager.GenerateToken(user.ID, user.Username)

	handler := NewRpcHandler(env.Storage, env.Hub, env.JWTManager)
	handler.RegisterMethod(NewUserInfoMethod(env.Storage))

	reqBody := resp.RpcRequest{
		JsonRPC: "2.0",
		Method:  "user.info",
		Id:      "1",
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/rpc", bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Request.Header.Set("Authorization", "Bearer "+token)

	handler.HandleRpcRequest(ctx)

	var response resp.RpcResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.Error != nil {
		t.Errorf("Unexpected error: %v", response.Error)
	}
}

func TestRpcHandler_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	env, err := SetupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test env: %v", err)
	}

	handler := NewRpcHandler(env.Storage, env.Hub, env.JWTManager)
	handler.RegisterMethod(NewUserInfoMethod(env.Storage))

	reqBody := resp.RpcRequest{
		JsonRPC: "2.0",
		Method:  "user.info",
		Id:      "1",
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/rpc", bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Request.Header.Set("Authorization", "Bearer invalid.token.here")

	handler.HandleRpcRequest(ctx)

	var response resp.RpcResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.Error == nil {
		t.Error("Expected error for invalid token")
	}
}

func TestRpcHandler_RegisterMethod(t *testing.T) {
	env, _ := SetupTestEnv()
	handler := NewRpcHandler(env.Storage, env.Hub, env.JWTManager)

	handler.RegisterMethod(&PingMethod{})

	method, ok := handler.getMethod("ping")
	if !ok {
		t.Error("Method should be registered")
	}
	if method.Name() != "ping" {
		t.Errorf("Expected method name 'ping', got '%s'", method.Name())
	}
}

func TestPingMethod(t *testing.T) {
	method := &PingMethod{}

	if method.Name() != "ping" {
		t.Errorf("Expected 'ping', got '%s'", method.Name())
	}

	if method.RequireAuth() {
		t.Error("Ping should not require auth")
	}

	result, err := method.Execute(nil, nil)
	if err != nil {
		t.Fatalf("Ping failed: %v", err)
	}

	resultMap := result.(map[string]interface{})
	if resultMap["pong"] != true {
		t.Error("Expected pong: true")
	}
	if resultMap["message"] != "pong" {
		t.Error("Expected message: pong")
	}
}
