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

type createRandomAccountParam struct {
	OfficeID   int64
	CurrencyID int64
	UserID     int64
}

func createRandomAccount(t *testing.T, param createRandomAccountParam) Account {
	if param.OfficeID == 0 {
		office := createRandomHqOffice(t)
		param.OfficeID = office.ID
	}
	if param.CurrencyID == 0 {
		officeCurrency := createRandomOfficeCurrency(t, param.OfficeID)
		param.CurrencyID = officeCurrency.CurrencyID
	}
	if param.UserID == 0 {
		user := createRandomUser(t, param.OfficeID)
		param.UserID = user.ID
	}

	arg := CreateAccountParams{
		Name:       gofakeit.Name(),
		OfficeID:   param.OfficeID,
		CurrencyID: param.CurrencyID,
		CreatedBy:  param.UserID,
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, account)
	assert.NotZero(t, account.ID)
	assert.Equal(t, arg.Name, account.Name)
	assert.Equal(t, arg.OfficeID, account.OfficeID)
	assert.Equal(t, arg.CurrencyID, account.CurrencyID)
	assert.Equal(t, arg.CreatedBy, account.CreatedBy)
	assert.Equal(t, fmt.Sprintf("%f", float64(0)), account.Balance)
	assert.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t, createRandomAccountParam{})
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t, createRandomAccountParam{})
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, account2)
	assert.Equal(t, account1.ID, account2.ID)
	assert.Equal(t, account1.Name, account2.Name)
	assert.Equal(t, account1.OfficeID, account2.OfficeID)
	assert.Equal(t, account1.CurrencyID, account2.CurrencyID)
	assert.Equal(t, account1.CreatedBy, account2.CreatedBy)
	assert.Equal(t, account1.Balance, account2.Balance)
	assert.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestListAccount(t *testing.T) {
	office := createRandomHqOffice(t)
	for i := 0; i < 10; i++ {
		createRandomAccount(t, createRandomAccountParam{office.ID, 0, 0})
	}

	arg := ListAccountParams{
		Limit:    5,
		Offset:   5,
		OfficeID: office.ID,
	}

	accounts, err := testQueries.ListAccount(context.Background(), arg)
	assert.NoError(t, err)
	assert.Len(t, accounts, 5)

	for _, account := range accounts {
		assert.NotEmpty(t, account)
	}
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t, createRandomAccountParam{})
	currency := createRandomCurrency(t)
	arg := UpdateAccountParams{
		ID:         account1.ID,
		Name:       gofakeit.Name(),
		CurrencyID: currency.ID,
	}
	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, account2)
	assert.Equal(t, account1.ID, account2.ID)
	assert.Equal(t, arg.Name, account2.Name)
	assert.Equal(t, account1.OfficeID, account2.OfficeID)
	assert.Equal(t, arg.CurrencyID, account2.CurrencyID)
	assert.Equal(t, account1.CreatedBy, account2.CreatedBy)
	assert.Equal(t, account1.Balance, account2.Balance)
	assert.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t, createRandomAccountParam{})
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	assert.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	assert.ErrorIs(t, err, sql.ErrNoRows)
	assert.Empty(t, account2)
}
