package db

import (
	"context"
	"database/sql"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	office := createRandomHqOffice(t)
	currency := createRandomCurrency(t)
	user := createRandomUser(t, office.ID)
	account1 := createRandomAccount(t, createRandomAccountParam{office.ID, currency.ID, user.ID})
	account2 := createRandomAccount(t, createRandomAccountParam{office.ID, currency.ID, user.ID})

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

	for i := 0; i < n; i++ {
		err := <-errs
		assert.NoError(t, err)

		result := <-results
		assert.NotEmpty(t, result)
		assert.NotZero(t, result.Transaction.ID)
		assert.Equal(t, office.ID, result.Transaction.OfficeID)
		assert.Equal(t, account1.ID, result.Transaction.FromAccountID)
		assert.Equal(t, sql.NullInt64{Int64: account2.ID, Valid: true}, result.Transaction.ToAccountID)
		assert.Equal(t, currency.ID, result.Transaction.CurrencyID)
		assert.Equal(t, sql.NullInt64{Int64: user.ID, Valid: true}, result.Transaction.CreatedBy)
		assert.NotZero(t, result.Transaction.CreatedAt)

		resultAmount, err := strconv.ParseFloat(result.Transaction.Amount, 64)
		assert.NoError(t, err)
		assert.Equal(t, amount, resultAmount)

		_, err = store.GetTransaction(context.Background(), result.Transaction.ID)
		assert.NoError(t, err)

		// TODO: add asserts for account balance
	}
}
