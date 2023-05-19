// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: transactions.sql

package db

import (
	"context"
	"database/sql"
)

const createTransaction = `-- name: CreateTransaction :one
insert into transactions (
  office_id, type, from_account_id, to_account_id, amount, currency_id, created_by
) values (
  $1, $2, $3, $4, $5, $6, $7
)
returning id, office_id, type, from_account_id, to_account_id, amount, currency_id, created_by, created_at, deleted_at
`

type CreateTransactionParams struct {
	OfficeID      int64           `db:"office_id" json:"office_id"`
	Type          TransactionType `db:"type" json:"type"`
	FromAccountID int64           `db:"from_account_id" json:"from_account_id"`
	ToAccountID   sql.NullInt64   `db:"to_account_id" json:"to_account_id"`
	Amount        string          `db:"amount" json:"amount"`
	CurrencyID    int64           `db:"currency_id" json:"currency_id"`
	CreatedBy     sql.NullInt64   `db:"created_by" json:"created_by"`
}

func (q *Queries) CreateTransaction(ctx context.Context, arg CreateTransactionParams) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, createTransaction,
		arg.OfficeID,
		arg.Type,
		arg.FromAccountID,
		arg.ToAccountID,
		arg.Amount,
		arg.CurrencyID,
		arg.CreatedBy,
	)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.OfficeID,
		&i.Type,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CurrencyID,
		&i.CreatedBy,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const deleteTransaction = `-- name: DeleteTransaction :exec
update transactions set deleted_at = now() where id = $1
`

func (q *Queries) DeleteTransaction(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteTransaction, id)
	return err
}

const getTransaction = `-- name: GetTransaction :one
select id, office_id, type, from_account_id, to_account_id, amount, currency_id, created_by, created_at, deleted_at from transactions where id = $1 and deleted_at is null limit 1
`

func (q *Queries) GetTransaction(ctx context.Context, id int64) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, getTransaction, id)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.OfficeID,
		&i.Type,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CurrencyID,
		&i.CreatedBy,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}
