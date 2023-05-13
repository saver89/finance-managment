package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rollback err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParam struct {
	OfficeID      int64   `json:"office_id"`
	FromAccountID int64   `json:"from_account_id"`
	ToAccountID   int64   `json:"to_account_id"`
	CurrencyID    int64   `json:"currency_id"`
	Amount        float64 `json:"amount"`
	CreatedBy     int64   `json:"created_by"`
}

type TransferTxResult struct {
	Transaction Transaction `json:"transaction"`
	FromAccount Account
	ToAccount   Account
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParam) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transaction, err = q.CreateTransaction(ctx, CreateTransactionParams{
			OfficeID:      arg.OfficeID,
			Type:          TransactionTypeTransfer,
			FromAccountID: arg.FromAccountID,
			ToAccountID:   sql.NullInt64{Int64: arg.ToAccountID, Valid: true},
			Amount:        fmt.Sprintf("%f", arg.Amount),
			CurrencyID:    arg.CurrencyID,
			CreatedBy:     sql.NullInt64{Int64: arg.CreatedBy, Valid: true},
		})
		if err != nil {
			return err
		}

		// update account balance
		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:      arg.FromAccountID,
			Balance: fmt.Sprintf("%f", -arg.Amount),
		})
		if err != nil {
			return err
		}

		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:      arg.ToAccountID,
			Balance: fmt.Sprintf("%f", arg.Amount),
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
