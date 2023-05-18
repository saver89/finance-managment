package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saver89/finance-management/internal/delivery/http/v1/handlers/request"
	"github.com/saver89/finance-management/internal/service"
	httpErrors "github.com/saver89/finance-management/pkg"
	"github.com/saver89/finance-management/pkg/logger"
)

type AccountHandlers struct {
	accountService service.AccountService
	log            logger.Logger
}

func NewAccountHandlers(log logger.Logger, accountService service.AccountService) *AccountHandlers {
	return &AccountHandlers{
		accountService: accountService,
		log:            log,
	}
}

func (a *AccountHandlers) CreateAccount(ctx *gin.Context) {
	var req request.CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		a.log.Errorf("error while binding json: %+v", err)
		ctx.JSON(httpErrors.ErrorResponse(httpErrors.NewBadRequestError(err.Error())))
		return
	}

	account, err := a.accountService.CreateAccount(ctx, &req)
	if err != nil {
		a.log.Errorf("error while creating account: %+v", err)
		ctx.JSON(httpErrors.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (a *AccountHandlers) GetAccount(ctx *gin.Context) {
	var req request.GetAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		a.log.Errorf("error while binding uri: %+v", err)
		ctx.JSON(httpErrors.ErrorResponse(httpErrors.NewBadRequestError(err.Error())))
		return
	}

	account, err := a.accountService.GetAccount(ctx, &req)
	if err != nil {
		a.log.Errorf("error while getting account: %+v", err)
		ctx.JSON(httpErrors.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (a *AccountHandlers) ListAccount(ctx *gin.Context) {
	var req request.ListAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		a.log.Errorf("error while binding query: %+v", err)
		ctx.JSON(httpErrors.ErrorResponse(httpErrors.NewBadRequestError(err.Error())))
		return
	}

	accounts, err := a.accountService.ListAccount(ctx, &req)
	if err != nil {
		a.log.Errorf("error while listing accounts: %+v", err)
		ctx.JSON(httpErrors.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
