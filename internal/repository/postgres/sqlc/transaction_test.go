package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

type createRandomTransactionParams struct {
	OfficeID      int64
	Type          TransactionType
	FromAccountID int64
	ToAccountID   int64
	Amount        string
	CurrencyID    int64
	CreatedBy     int64
}

func createRandomTransaction(t *testing.T, param createRandomTransactionParams) Transaction {
	if param.OfficeID == 0 {
		office := createRandomHqOffice(t)
		param.OfficeID = office.ID
	}

	if param.Type == "" {
		param.Type = TransactionTypeTransfer
	}

	if param.CurrencyID == 0 {
		officeCurrency := createRandomOfficeCurrency(t, param.OfficeID)
		param.CurrencyID = officeCurrency.CurrencyID
	}

	if param.CreatedBy == 0 {
		user := createRandomUser(t, param.OfficeID)
		param.CreatedBy = user.ID
	}

	if param.FromAccountID == 0 {
		account := createRandomAccount(t, createRandomAccountParam{param.OfficeID, param.CurrencyID, param.CreatedBy})
		param.FromAccountID = account.ID
	}

	if param.ToAccountID == 0 {
		account := createRandomAccount(t, createRandomAccountParam{param.OfficeID, param.CurrencyID, param.CreatedBy})
		param.ToAccountID = account.ID
	}

	if param.Amount == "" {
		param.Amount = fmt.Sprintf("%f", gofakeit.Float64Range(1, 1000))
	}

	arg := CreateTransactionParams{
		OfficeID:      param.OfficeID,
		Type:          param.Type,
		FromAccountID: param.FromAccountID,
		ToAccountID:   sql.NullInt64{Int64: param.ToAccountID, Valid: true},
		Amount:        param.Amount,
		CurrencyID:    param.CurrencyID,
		CreatedBy:     sql.NullInt64{Int64: param.CreatedBy, Valid: true},
	}
	transaction, err := testQueries.CreateTransaction(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, transaction)
	assert.NotZero(t, transaction.ID)
	assert.Equal(t, arg.OfficeID, transaction.OfficeID)
	assert.Equal(t, arg.Type, transaction.Type)
	assert.Equal(t, arg.FromAccountID, transaction.FromAccountID)
	assert.Equal(t, arg.ToAccountID, transaction.ToAccountID)
	assert.Equal(t, arg.Amount, transaction.Amount)
	assert.Equal(t, arg.CurrencyID, transaction.CurrencyID)
	assert.Equal(t, arg.CreatedBy, transaction.CreatedBy)
	assert.NotZero(t, transaction.CreatedAt)

	return transaction
}

func TestCreateTransaction(t *testing.T) {
	createRandomTransaction(t, createRandomTransactionParams{})
}

func TestGetTransaction(t *testing.T) {
	transaction1 := createRandomTransaction(t, createRandomTransactionParams{})
	transaction2, err := testQueries.GetTransaction(context.Background(), transaction1.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, transaction2)
	assert.Equal(t, transaction1.ID, transaction2.ID)
	assert.Equal(t, transaction1.OfficeID, transaction2.OfficeID)
	assert.Equal(t, transaction1.Type, transaction2.Type)
	assert.Equal(t, transaction1.FromAccountID, transaction2.FromAccountID)
	assert.Equal(t, transaction1.ToAccountID, transaction2.ToAccountID)
	assert.Equal(t, transaction1.Amount, transaction2.Amount)
	assert.Equal(t, transaction1.CurrencyID, transaction2.CurrencyID)
	assert.Equal(t, transaction1.CreatedBy, transaction2.CreatedBy)
	assert.WithinDuration(t, transaction1.CreatedAt, transaction2.CreatedAt, time.Second)
}

func TestDeleteTransaction(t *testing.T) {
	transaction1 := createRandomTransaction(t, createRandomTransactionParams{})
	err := testQueries.DeleteTransaction(context.Background(), transaction1.ID)
	assert.NoError(t, err)
	transaction2, err := testQueries.GetTransaction(context.Background(), transaction1.ID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, sql.ErrNoRows)
	assert.Empty(t, transaction2)
}
