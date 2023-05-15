package service

import (
	"context"
	"strconv"

	"github.com/saver89/finance-management/internal/delivery/http/v1/handlers/request"
	"github.com/saver89/finance-management/internal/domain"
	db "github.com/saver89/finance-management/internal/repository/postgres/sqlc"
	"github.com/saver89/finance-management/internal/service/response"
)

type AccountService interface {
	CreateAccount(ctx context.Context, req *request.CreateAccountRequest) (*response.CreateAccountResponse, error)
	GetAccount(ctx context.Context, req *request.GetAccountRequest) (*response.GetAccountResponse, error)
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

	resp := response.CreateAccountResponse{
		ID:         account.ID,
		Name:       account.Name,
		OfficeID:   account.OfficeID,
		CurrencyID: account.CurrencyID,
	}

	return &resp, nil
}

func (a *accountService) GetAccount(ctx context.Context, req *request.GetAccountRequest) (*response.GetAccountResponse, error) {

	account, err := a.store.GetAccount(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	resp := response.GetAccountResponse{}

	balanceFloat, err := strconv.ParseFloat(account.Balance, 64)
	if err != nil {
		return nil, err
	}
	resp.Account = domain.Account{
		ID:         account.ID,
		Name:       account.Name,
		Balance:    balanceFloat,
		OfficeID:   account.OfficeID,
		CurrencyID: account.CurrencyID,
		CreatedBy:  account.CreatedBy,
		State:      string(account.State),
		CreatedAt:  account.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return &resp, nil
}
