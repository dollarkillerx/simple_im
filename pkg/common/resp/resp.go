package resp

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	JsonRPCVersion = "2.0"
	ErrorCode      = -32000
)

type RpcRequest struct {
	JsonRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	Id      string          `json:"id"`
}

type RpcResponse struct {
	JsonRPC string          `json:"jsonrpc"`
	Id      string          `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *RpcError       `json:"error,omitempty"`
}

type RpcError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessReturn(ctx *gin.Context, id string, result interface{}) {
	Return(ctx, http.StatusOK, id, result, nil)
}

func ErrorReturn(ctx *gin.Context, id string, err error) {
	Return(ctx, http.StatusOK, id, nil, err)
}

func Return(ctx *gin.Context, httpCode int, id string, result interface{}, err error) {
	resp := RpcResponse{
		JsonRPC: JsonRPCVersion,
		Id:      id,
	}

	if err != nil {
		resp.Error = &RpcError{
			Code:    ErrorCode,
			Message: err.Error(),
		}
	} else if result != nil {
		resultBytes, _ := json.Marshal(result)
		resp.Result = resultBytes
	}

	ctx.JSON(httpCode, resp)
}

func InvalidRequest(ctx *gin.Context, id string, message string) {
	resp := RpcResponse{
		JsonRPC: JsonRPCVersion,
		Id:      id,
		Error: &RpcError{
			Code:    -32600,
			Message: message,
		},
	}
	ctx.JSON(http.StatusOK, resp)
}

func MethodNotFound(ctx *gin.Context, id string, method string) {
	resp := RpcResponse{
		JsonRPC: JsonRPCVersion,
		Id:      id,
		Error: &RpcError{
			Code:    -32601,
			Message: "Method not found: " + method,
		},
	}
	ctx.JSON(http.StatusOK, resp)
}
