package handlers

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	db "github.com/saver89/finance-management/internal/repository/postgres/sqlc"
	"github.com/saver89/finance-management/internal/service"
	"github.com/saver89/finance-management/pkg/logger"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

type Server struct {
	store  db.Store
	router *gin.Engine
	log    logger.Logger
}

func NewServer(store db.Store, log logger.Logger) *Server {
	server := &Server{store: store, log: log}
	server.router = gin.Default()

	v1Group := server.router.Group("/v1")
	accountService := service.NewAccountService(store, server.log)

	accountHandlers := NewAccountHandlers(server.log, accountService)

	v1Group.POST("/account", accountHandlers.CreateAccount)
	v1Group.GET("/account/:id", accountHandlers.GetAccount)
	v1Group.GET("/accounts", accountHandlers.ListAccount)

	return server
}
