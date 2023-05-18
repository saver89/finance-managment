package http_delivery_v1

import (
	"github.com/gin-gonic/gin"
	"github.com/saver89/finance-management/internal/delivery/http/v1/handlers"
	"github.com/saver89/finance-management/internal/service"
	"github.com/saver89/finance-management/pkg/logger"
)

type ServiceParams struct {
	AccountService service.AccountService
}

type HttpDelivery struct {
	group           *gin.RouterGroup
	services        ServiceParams
	accountHandlers *handlers.AccountHandlers
	log             logger.Logger
}

func NewHttpDelivery(log logger.Logger, group *gin.RouterGroup, services ServiceParams) {
	httpDelivery := &HttpDelivery{
		group:           group,
		services:        services,
		accountHandlers: handlers.NewAccountHandlers(log, services.AccountService),
		log:             log,
	}

	httpDelivery.MapRoutes()
}

func (hd *HttpDelivery) MapRoutes() {
	hd.group.POST("/account", hd.accountHandlers.CreateAccount)
	hd.group.GET("/account/:id", hd.accountHandlers.GetAccount)
	hd.group.GET("/accounts", hd.accountHandlers.ListAccount)
}
