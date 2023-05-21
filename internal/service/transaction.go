package service

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/saver89/finance-management/internal/delivery/http/v1/handlers/request"
	db "github.com/saver89/finance-management/internal/repository/postgres/sqlc"
	"github.com/saver89/finance-management/internal/service/response"
	httpErrors "github.com/saver89/finance-management/pkg"
	"github.com/saver89/finance-management/pkg/logger"
)

type TransactionService interface {
	CreateTransfer(ctx context.Context, req *request.CreateTransferRequest) (*response.CreateTransferResponse, error)
	ValidateTransfer(ctx context.Context, req *request.CreateTransferRequest) error
}

type transactionService struct {
	store db.Store
	log   logger.Logger
}

func NewTransactionService(store db.Store, log logger.Logger) TransactionService {
	return &transactionService{
		store: store,
		log:   log,
	}
}

func (ts *transactionService) CreateTransfer(ctx context.Context, req *request.CreateTransferRequest) (*response.CreateTransferResponse, error) {

	arg := db.TransferTxParam{
		OfficeID:      req.OfficeID,
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
		CurrencyID:    req.CurrencyID,
		CreatedBy:     req.CreatedBy,
	}
	transfer, err := ts.store.TransferTx(ctx, arg)
	if err != nil {
		return nil, err
	}

	resp := response.CreateTransferResponse{
		Transaction: transfer.Transaction,
		FromAccount: transfer.FromAccount,
		ToAccount:   transfer.ToAccount,
	}
	return &resp, nil
}

func (ts *transactionService) ValidateTransfer(ctx context.Context, req *request.CreateTransferRequest) error {

	err := ts.validateAccount(ctx, req.FromAccountID, req.CurrencyID, req.OfficeID)
	if err != nil {
		return errors.Wrap(err, "validate from account")
	}

	err = ts.validateAccount(ctx, req.ToAccountID, req.CurrencyID, req.OfficeID)
	if err != nil {
		return errors.Wrap(err, "validate to account")
	}

	return nil
}

func (ts *transactionService) validateAccount(ctx context.Context, accountID, currencyID, officeID int64) error {
	account, err := ts.store.GetAccount(ctx, accountID)
	if err != nil {
		return err
	}

	if account.CurrencyID != currencyID {
		return httpErrors.NewValidationError(fmt.Sprintf("account [%d] currency mismatch", accountID))
	}

	if account.OfficeID != officeID {
		return httpErrors.NewValidationError(fmt.Sprintf("account [%d] office mismatch", accountID))
	}

	return nil
}
