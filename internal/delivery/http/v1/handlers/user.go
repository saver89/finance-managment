package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saver89/finance-management/internal/delivery/http/v1/handlers/request"
	"github.com/saver89/finance-management/internal/service"
	httpErrors "github.com/saver89/finance-management/pkg"
	"github.com/saver89/finance-management/pkg/logger"
)

type UserHandlers struct {
	userService service.UserService
	log         logger.Logger
}

func NewUserHandlers(log logger.Logger, userService service.UserService) *UserHandlers {
	return &UserHandlers{
		userService: userService,
		log:         log,
	}
}

func (u *UserHandlers) CreateUser(ctx *gin.Context) {
	var req request.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(httpErrors.ErrorResponse(u.log, httpErrors.NewBadRequestError(err.Error)))
		return
	}

	user, err := u.userService.CreateUser(ctx, &req)
	if err != nil {
		u.log.Errorf("error while creating user: %+v", err)
		ctx.JSON(httpErrors.ErrorResponse(u.log, err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (u *UserHandlers) GetUser(ctx *gin.Context) {
	var req request.GetUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(httpErrors.ErrorResponse(u.log, httpErrors.NewBadRequestError(err)))
		return
	}

	user, err := u.userService.GetUser(ctx, &req)
	if err != nil {
		ctx.JSON(httpErrors.ErrorResponse(u.log, err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}
