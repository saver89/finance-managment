package server

import (
	"github.com/gin-gonic/gin"
	http_delivery "github.com/saver89/finance-management/internal/delivery/http/v1"
	db "github.com/saver89/finance-management/internal/repository/postgres/sqlc"
	"github.com/saver89/finance-management/internal/service"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	server.router = gin.Default()

	v1Group := server.router.Group("/v1")

	accountService := service.NewAccountService(store)

	http_delivery.InitHttpV1(v1Group, http_delivery.ServiceParams{
		AccountService: accountService,
	})

	return server
}

func (s *Server) Run(addr string) error {
	if err := s.runHttpServer(addr); err != nil {
		return err
	}

	return nil
}
