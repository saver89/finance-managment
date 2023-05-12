// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: account.sql

package db

import (
	"context"
)

const createAccount = `-- name: CreateAccount :one
insert into account (
  office_id, name, currency_id, created_by, state
) values (
  $1, $2, $3, $4, 'active'
) returning id, office_id, name, balance, currency_id, created_by, state, created_at, deleted_at
`

type CreateAccountParams struct {
	OfficeID   int64  `db:"office_id"`
	Name       string `db:"name"`
	CurrencyID int64  `db:"currency_id"`
	CreatedBy  int64  `db:"created_by"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, createAccount,
		arg.OfficeID,
		arg.Name,
		arg.CurrencyID,
		arg.CreatedBy,
	)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.OfficeID,
		&i.Name,
		&i.Balance,
		&i.CurrencyID,
		&i.CreatedBy,
		&i.State,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const deleteAccount = `-- name: DeleteAccount :exec
update account set deleted_at = now() where id = $1
`

func (q *Queries) DeleteAccount(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteAccount, id)
	return err
}

const getAccount = `-- name: GetAccount :one
select id, office_id, name, balance, currency_id, created_by, state, created_at, deleted_at from account where id = $1 and deleted_at is null
`

func (q *Queries) GetAccount(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccount, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.OfficeID,
		&i.Name,
		&i.Balance,
		&i.CurrencyID,
		&i.CreatedBy,
		&i.State,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const listAccount = `-- name: ListAccount :many
select id, office_id, name, balance, currency_id, created_by, state, created_at, deleted_at from account 
where deleted_at is null and office_id = $3
LIMIT $1
OFFSET $2
`

type ListAccountParams struct {
	Limit    int32 `db:"limit"`
	Offset   int32 `db:"offset"`
	OfficeID int64 `db:"office_id"`
}

func (q *Queries) ListAccount(ctx context.Context, arg ListAccountParams) ([]Account, error) {
	rows, err := q.db.QueryContext(ctx, listAccount, arg.Limit, arg.Offset, arg.OfficeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Account
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.OfficeID,
			&i.Name,
			&i.Balance,
			&i.CurrencyID,
			&i.CreatedBy,
			&i.State,
			&i.CreatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAccount = `-- name: UpdateAccount :one
update account set name = $2, currency_id = $3
where id = $1 and deleted_at is null
returning id, office_id, name, balance, currency_id, created_by, state, created_at, deleted_at
`

type UpdateAccountParams struct {
	ID         int64  `db:"id"`
	Name       string `db:"name"`
	CurrencyID int64  `db:"currency_id"`
}

func (q *Queries) UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, updateAccount, arg.ID, arg.Name, arg.CurrencyID)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.OfficeID,
		&i.Name,
		&i.Balance,
		&i.CurrencyID,
		&i.CreatedBy,
		&i.State,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}
