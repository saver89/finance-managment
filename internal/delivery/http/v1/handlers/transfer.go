package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saver89/finance-management/internal/delivery/http/v1/handlers/request"
	"github.com/saver89/finance-management/internal/service"
	httpErrors "github.com/saver89/finance-management/pkg"
	"github.com/saver89/finance-management/pkg/logger"
)

type TransactionHandlers struct {
	transactionService service.TransactionService
	log                logger.Logger
}

func NewTransactionHandler(log logger.Logger, transactionService service.TransactionService) *TransactionHandlers {
	return &TransactionHandlers{
		transactionService: transactionService,
		log:                log,
	}
}

func (h *TransactionHandlers) CreateTransfer(ctx *gin.Context) {
	var req request.CreateTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(httpErrors.ErrorResponse(h.log, httpErrors.NewBadRequestError(err.Error())))
		return
	}

	if err := h.ValidateTransfer(ctx, &req); err != nil {
		ctx.JSON(httpErrors.ErrorResponse(h.log, err))
		return
	}

	transaction, err := h.transactionService.CreateTransfer(ctx, &req)
	if err != nil {
		ctx.JSON(httpErrors.ErrorResponse(h.log, err))
		return
	}

	ctx.JSON(http.StatusOK, transaction)
}

func (h *TransactionHandlers) ValidateTransfer(ctx *gin.Context, req *request.CreateTransferRequest) error {
	return h.transactionService.ValidateTransfer(ctx, req)
}
