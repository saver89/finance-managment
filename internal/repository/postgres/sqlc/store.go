package db

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
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

var txKey = struct{}{}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParam) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey).(string)

		fmt.Println(txName, "create transaction")
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
		fmt.Println(txName, "get fromAccount")
		fromAccount, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}
		fromAccountBalance, err := strconv.ParseFloat(fromAccount.Balance, 64)
		if err != nil {
			return err
		}
		fmt.Println(txName, "update from account")
		result.FromAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
			ID:      fromAccount.ID,
			Balance: fmt.Sprintf("%f", fromAccountBalance-arg.Amount),
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "get to account")
		toAccount, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}
		toAccountBalance, err := strconv.ParseFloat(toAccount.Balance, 64)
		if err != nil {
			return err
		}
		fmt.Println(txName, "update to account")
		result.ToAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
			ID:      toAccount.ID,
			Balance: fmt.Sprintf("%f", toAccountBalance+arg.Amount),
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
