package http_delivery_v1

import (
	"github.com/gin-gonic/gin"
	"github.com/saver89/finance-management/internal/delivery/http/v1/handlers"
	"github.com/saver89/finance-management/internal/service"
)

type ServiceParams struct {
	AccountService service.AccountService
}

type HttpDelivery struct {
	group           *gin.RouterGroup
	services        ServiceParams
	accountHandlers *handlers.AccountHandlers
}

func InitHttpV1(group *gin.RouterGroup, services ServiceParams) {
	httpDelivery := &HttpDelivery{
		group:           group,
		services:        services,
		accountHandlers: handlers.NewAccountHandlers(services.AccountService),
	}

	httpDelivery.MapRoutes()
}

func (hd *HttpDelivery) MapRoutes() {
	hd.group.POST("/account", hd.accountHandlers.CreateAccount)
	hd.group.GET("/account/:id", hd.accountHandlers.GetAccount)
	hd.group.GET("/accounts", hd.accountHandlers.ListAccount)
}
