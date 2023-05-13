package db

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	var err error
	store := NewStore(testDB)

	office := createRandomHqOffice(t)
	currency := createRandomCurrency(t)
	user := createRandomUser(t, office.ID)
	account1 := createRandomAccount(t, createRandomAccountParam{office.ID, currency.ID, user.ID})
	account1, err = testQueries.UpdateAccountBalance(context.Background(), UpdateAccountBalanceParams{
		ID:      account1.ID,
		Balance: fmt.Sprintf("%f", float64(gofakeit.Price(100, 10000))),
	})
	require.NoError(t, err)

	account2 := createRandomAccount(t, createRandomAccountParam{office.ID, currency.ID, user.ID})
	account2, err = testQueries.UpdateAccountBalance(context.Background(), UpdateAccountBalanceParams{
		ID:      account2.ID,
		Balance: fmt.Sprintf("%f", float64(gofakeit.Price(100, 10000))),
	})
	require.NoError(t, err)

	n := 5
	amount := float64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParam{
				OfficeID:      office.ID,
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
				CurrencyID:    currency.ID,
				CreatedBy:     user.ID,
			})
			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)
		require.NotZero(t, result.Transaction.ID)
		require.Equal(t, office.ID, result.Transaction.OfficeID)
		require.Equal(t, account1.ID, result.Transaction.FromAccountID)
		require.Equal(t, sql.NullInt64{Int64: account2.ID, Valid: true}, result.Transaction.ToAccountID)
		require.Equal(t, currency.ID, result.Transaction.CurrencyID)
		require.Equal(t, sql.NullInt64{Int64: user.ID, Valid: true}, result.Transaction.CreatedBy)
		require.NotZero(t, result.Transaction.CreatedAt)

		resultAmount, err := strconv.ParseFloat(result.Transaction.Amount, 64)
		require.NoError(t, err)
		require.Equal(t, amount, resultAmount)

		_, err = store.GetTransaction(context.Background(), result.Transaction.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check account balances
		balance1, err := strconv.ParseFloat(account1.Balance, 64)
		require.NoError(t, err)
		balance2, err := strconv.ParseFloat(account2.Balance, 64)
		require.NoError(t, err)
		fromAccountBalance, err := strconv.ParseFloat(fromAccount.Balance, 64)
		require.NoError(t, err)
		toAccountBalance, err := strconv.ParseFloat(toAccount.Balance, 64)
		require.NoError(t, err)

		diff1 := balance1 - fromAccountBalance
		diff2 := toAccountBalance - balance2

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, math.Mod(diff1, amount) == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	account1Balance, err := strconv.ParseFloat(account1.Balance, 64)
	require.NoError(t, err)
	account2Balance, err := strconv.ParseFloat(account2.Balance, 64)
	require.NoError(t, err)

	updatedBalance1 := fmt.Sprintf("%f", account1Balance-float64(n)*amount)
	updatedBalance2 := fmt.Sprintf("%f", account2Balance+float64(n)*amount)
	require.Equal(t, updatedBalance1, updatedAccount1.Balance)
	require.Equal(t, updatedBalance2, updatedAccount2.Balance)

}
