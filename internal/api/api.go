package api

import (
	"fmt"

	"simple_im/internal/conf"
	"simple_im/internal/middleware"
	"simple_im/internal/storage"
	"simple_im/internal/ws"
	"simple_im/pkg/common/jwt"

	"github.com/gin-gonic/gin"
)

type ApiServer struct {
	storage    *storage.Storage
	hub        *ws.Hub
	conf       conf.Config
	app        *gin.Engine
	rpcHandler *RpcHandler
	jwtManager *jwt.JWTManager
}

func NewApiServer(storage *storage.Storage, hub *ws.Hub, config conf.Config) *ApiServer {
	jwtManager := jwt.NewJWTManager(config.JWTConfiguration.Secret, config.JWTConfiguration.Expire)

	return &ApiServer{
		storage:    storage,
		hub:        hub,
		conf:       config,
		jwtManager: jwtManager,
	}
}

func (a *ApiServer) Run() error {
	if a.conf.ServiceConfiguration.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	a.app = gin.New()
	a.app.Use(middleware.HttpRecover())
	a.app.Use(gin.Logger())
	a.app.Use(middleware.Cors())

	a.rpcHandler = NewRpcHandler(a.storage, a.hub, a.jwtManager)
	a.registerRpcMethods()
	a.Router()

	addr := fmt.Sprintf("0.0.0.0:%s", a.conf.ServiceConfiguration.Port)
	fmt.Printf("Server starting on %s\n", addr)
	return a.app.Run(addr)
}

func (a *ApiServer) Router() {
	a.app.GET("/health", a.HealthCheck)
	a.app.POST("/api/rpc", a.Rpc)
	a.app.GET("/ws", a.WebSocket)

	// File upload/download
	a.app.POST("/api/upload", middleware.JWTAuth(a.jwtManager), a.Upload)
	a.app.Static("/files", a.conf.UploadConfiguration.SavePath)
}

func (a *ApiServer) HealthCheck(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"status": "ok",
	})
}

func (a *ApiServer) registerRpcMethods() {
	// Basic methods
	a.rpcHandler.RegisterMethod(&PingMethod{})

	// User methods
	a.rpcHandler.RegisterMethod(NewUserRegisterMethod(a.storage, a.jwtManager))
	a.rpcHandler.RegisterMethod(NewUserLoginMethod(a.storage, a.jwtManager))
	a.rpcHandler.RegisterMethod(NewUserInfoMethod(a.storage))

	// Friend methods
	a.rpcHandler.RegisterMethod(NewFriendListMethod(a.storage))
	a.rpcHandler.RegisterMethod(NewFriendAddMethod(a.storage, a.hub))
	a.rpcHandler.RegisterMethod(NewFriendAcceptMethod(a.storage, a.hub))
	a.rpcHandler.RegisterMethod(NewFriendPendingMethod(a.storage))

	// Group methods
	a.rpcHandler.RegisterMethod(NewGroupCreateMethod(a.storage))
	a.rpcHandler.RegisterMethod(NewGroupListMethod(a.storage))
	a.rpcHandler.RegisterMethod(NewGroupInfoMethod(a.storage))
	a.rpcHandler.RegisterMethod(NewGroupJoinMethod(a.storage))

	// Message methods
	a.rpcHandler.RegisterMethod(NewMessageSendMethod(a.storage, a.hub))
	a.rpcHandler.RegisterMethod(NewMessageHistoryMethod(a.storage))
}
