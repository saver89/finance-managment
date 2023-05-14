package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saver89/finance-management/internal/delivery/http/v1/handlers/request"
	"github.com/saver89/finance-management/internal/service"
	httpErrors "github.com/saver89/finance-management/pkg"
)

type AccountHandlers struct {
	accountService service.AccountService
}

func NewAccountHandlers(accountService service.AccountService) *AccountHandlers {
	return &AccountHandlers{
		accountService: accountService,
	}
}

func (a *AccountHandlers) CreateAccount(ctx *gin.Context) {
	var req request.CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpErrors.ErrorResponse(err))
		return
	}

	account, err := a.accountService.CreateAccount(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpErrors.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}
