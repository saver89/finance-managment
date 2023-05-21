package http_delivery_v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/saver89/finance-management/internal/delivery/http/v1/handlers"
	"github.com/saver89/finance-management/internal/service"
	"github.com/saver89/finance-management/pkg/logger"
)

type ServiceParams struct {
	AccountService     service.AccountService
	TransactionService service.TransactionService
}

type HttpDelivery struct {
	group    *gin.RouterGroup
	services ServiceParams
	log      logger.Logger

	accountHandlers     *handlers.AccountHandlers
	transactionHandlers *handlers.TransactionHandlers
}

func NewHttpDelivery(
	log logger.Logger,
	group *gin.RouterGroup,
	services ServiceParams,
) *HttpDelivery {
	httpDelivery := &HttpDelivery{
		group:               group,
		services:            services,
		log:                 log,
		accountHandlers:     handlers.NewAccountHandlers(log, services.AccountService),
		transactionHandlers: handlers.NewTransactionHandler(log, services.TransactionService),
	}

	httpDelivery.MapRoutes()

	return httpDelivery
}

func (hd *HttpDelivery) MapRoutes() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("transaction-type", handlers.ValidTransactionType)
	}

	hd.group.POST("/account", hd.accountHandlers.CreateAccount)
	hd.group.GET("/account/:id", hd.accountHandlers.GetAccount)
	hd.group.GET("/accounts", hd.accountHandlers.ListAccount)

	hd.group.POST("/transfer", hd.transactionHandlers.CreateTransfer)
}
