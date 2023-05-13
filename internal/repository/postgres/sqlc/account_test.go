package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.NotZero(t, account.ID)
	require.Equal(t, arg.Name, account.Name)
	require.Equal(t, arg.OfficeID, account.OfficeID)
	require.Equal(t, arg.CurrencyID, account.CurrencyID)
	require.Equal(t, arg.CreatedBy, account.CreatedBy)
	require.Equal(t, fmt.Sprintf("%f", float64(0)), account.Balance)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t, createRandomAccountParam{})
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t, createRandomAccountParam{})
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Name, account2.Name)
	require.Equal(t, account1.OfficeID, account2.OfficeID)
	require.Equal(t, account1.CurrencyID, account2.CurrencyID)
	require.Equal(t, account1.CreatedBy, account2.CreatedBy)
	require.Equal(t, account1.Balance, account2.Balance)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
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
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
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
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, arg.Name, account2.Name)
	require.Equal(t, account1.OfficeID, account2.OfficeID)
	require.Equal(t, arg.CurrencyID, account2.CurrencyID)
	require.Equal(t, account1.CreatedBy, account2.CreatedBy)
	require.Equal(t, account1.Balance, account2.Balance)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t, createRandomAccountParam{})
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.ErrorIs(t, err, sql.ErrNoRows)
	require.Empty(t, account2)
}

func TestUpdateAccountBalance(t *testing.T) {
	account1 := createRandomAccount(t, createRandomAccountParam{})

	newBalance := fmt.Sprintf("%f", gofakeit.Price(100, 1000000))
	arg := UpdateAccountBalanceParams{
		ID:      account1.ID,
		Balance: newBalance,
	}

	account2, err := testQueries.UpdateAccountBalance(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Name, account2.Name)
	require.Equal(t, account1.OfficeID, account2.OfficeID)
	require.Equal(t, account1.CurrencyID, account2.CurrencyID)
	require.Equal(t, account1.CreatedBy, account2.CreatedBy)
	require.Equal(t, newBalance, account2.Balance)
}
