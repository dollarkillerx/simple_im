package server

import (
	"simple_im/internal/api"
	"simple_im/internal/conf"
	"simple_im/internal/storage"
	"simple_im/internal/ws"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Server struct {
	storage   *storage.Storage
	apiServer *api.ApiServer
	wsHub     *ws.Hub
	conf      conf.Config
}

func NewServer(db *gorm.DB, redisClient *redis.Client, config conf.Config) *Server {
	st := storage.NewStorage(redisClient, db)
	hub := ws.NewHub()
	apiServer := api.NewApiServer(st, hub, config)

	return &Server{
		storage:   st,
		apiServer: apiServer,
		wsHub:     hub,
		conf:      config,
	}
}

func (s *Server) Run() error {
	// Start WebSocket hub
	go s.wsHub.Run()

	// Start API server
	return s.apiServer.Run()
}
