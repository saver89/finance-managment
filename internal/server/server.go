package server

import (
	"github.com/gin-gonic/gin"
	http_delivery "github.com/saver89/finance-management/internal/delivery/http/v1"
	db "github.com/saver89/finance-management/internal/repository/postgres/sqlc"
	"github.com/saver89/finance-management/internal/service"
	"github.com/saver89/finance-management/pkg/logger"
)

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
	transactionService := service.NewTransactionService(store, server.log)
	userService := service.NewUserService(store, server.log)

	http_delivery.NewHttpDelivery(server.log, v1Group, http_delivery.ServiceParams{
		AccountService:     accountService,
		TransactionService: transactionService,
		UserService:        userService,
	})

	return server
}

func (s *Server) Run(addr string) error {
	if err := s.runHttpServer(addr); err != nil {
		return err
	}

	return nil
}
