package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParam) (TransferTxResult, error)
}

type StoreSQL struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &StoreSQL{
		Queries: New(db),
		db:      db,
	}
}

func (store *StoreSQL) execTx(ctx context.Context, fn func(*Queries) error) error {
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

func (store *StoreSQL) TransferTx(ctx context.Context, arg TransferTxParam) (TransferTxResult, error) {
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
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addBalance(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = addBalance(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}

func addBalance(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 float64,
	accountID2 int64,
	amountID2 float64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: fmt.Sprintf("%f", amount1),
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: fmt.Sprintf("%f", amountID2),
	})
	return
}
