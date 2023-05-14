package service

import (
	"context"

	"github.com/saver89/finance-management/internal/delivery/http/v1/handlers/request"
	db "github.com/saver89/finance-management/internal/repository/postgres/sqlc"
	"github.com/saver89/finance-management/internal/service/response"
)

type AccountService interface {
	CreateAccount(ctx context.Context, req *request.CreateAccountRequest) (*response.CreateAccountResponse, error)
}

type accountService struct {
	store *db.Store
}

func NewAccountService(store *db.Store) AccountService {
	return &accountService{
		store: store,
	}
}

func (a *accountService) CreateAccount(ctx context.Context, req *request.CreateAccountRequest) (*response.CreateAccountResponse, error) {

	arg := db.CreateAccountParams{
		Name:       req.Name,
		OfficeID:   req.OfficeID,
		CurrencyID: req.CurrencyID,
		CreatedBy:  req.CreatedBy,
	}
	account, err := a.store.CreateAccount(ctx, arg)
	if err != nil {
		return nil, err
	}

	response := response.CreateAccountResponse{
		ID:         account.ID,
		Name:       account.Name,
		OfficeID:   account.OfficeID,
		CurrencyID: account.CurrencyID,
	}

	return &response, nil
}
