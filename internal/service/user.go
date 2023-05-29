package service

import (
	"context"
	"database/sql"

	"github.com/saver89/finance-management/internal/delivery/http/v1/handlers/request"
	db "github.com/saver89/finance-management/internal/repository/postgres/sqlc"
	"github.com/saver89/finance-management/internal/service/response"
	"github.com/saver89/finance-management/pkg/logger"
	"github.com/saver89/finance-management/pkg/password"
)

type UserService interface {
	CreateUser(ctx context.Context, req *request.CreateUserRequest) (*response.CreateUserResponse, error)
	GetUser(ctx context.Context, req *request.GetUserRequest) (*response.GetUserResponse, error)
}

type userService struct {
	store db.Store
	log   logger.Logger
}

func NewUserService(store db.Store, log logger.Logger) UserService {
	return &userService{
		store: store,
		log:   log,
	}
}

func (us *userService) CreateUser(ctx context.Context, req *request.CreateUserRequest) (*response.CreateUserResponse, error) {

	arg := db.CreateUserParams{
		OfficeID:   req.OfficeID,
		FirstName:  sql.NullString{String: req.FirstName, Valid: true},
		LastName:   sql.NullString{String: req.LastName, Valid: true},
		MiddleName: sql.NullString{String: req.MiddleName, Valid: true},
		Email:      sql.NullString{String: req.Email, Valid: true},
	}

	hashPassword, err := password.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	arg.PasswordHash = hashPassword

	user, err := us.store.CreateUser(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &response.CreateUserResponse{User: user}, nil
}

func (us *userService) GetUser(ctx context.Context, req *request.GetUserRequest) (*response.GetUserResponse, error) {
	user, err := us.store.GetUser(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &response.GetUserResponse{User: user}, nil
}
