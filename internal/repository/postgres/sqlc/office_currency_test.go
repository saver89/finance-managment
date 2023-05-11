package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createRandomOfficeCurrency(t *testing.T, officeID int64) OfficeCurrency {
	currency := createRandomCurrency(t)

	if officeID == 0 {
		office := createRandomHqOffice(t)
		officeID = office.ID
	}
	arg := CreateOfficeCurrencyParams{
		OfficeID:   officeID,
		CurrencyID: currency.ID,
		Type:       OfficeCurrencyTypeRetail,
	}

	officeCurrency, err := testQueries.CreateOfficeCurrency(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, officeCurrency)
	assert.NotZero(t, officeCurrency.ID)
	assert.Equal(t, arg.OfficeID, officeCurrency.OfficeID)
	assert.Equal(t, arg.CurrencyID, officeCurrency.CurrencyID)
	assert.Equal(t, arg.Type, officeCurrency.Type)
	assert.NotEmpty(t, officeCurrency.CreatedAt)

	return officeCurrency
}

func TestOfficeCurrencyCreate(t *testing.T) {
	createRandomOfficeCurrency(t, 0)
}

func TestOfficeCurrencyGet(t *testing.T) {
	officeCurrency1 := createRandomOfficeCurrency(t, 0)

	officeCurrency2, err := testQueries.GetOfficeCurrency(context.Background(), officeCurrency1.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, officeCurrency2)
	assert.Equal(t, officeCurrency1.ID, officeCurrency2.ID)
	assert.Equal(t, officeCurrency1.OfficeID, officeCurrency2.OfficeID)
	assert.Equal(t, officeCurrency1.CurrencyID, officeCurrency2.CurrencyID)
	assert.Equal(t, officeCurrency1.Type, officeCurrency2.Type)
	assert.WithinDuration(t, officeCurrency1.CreatedAt, officeCurrency2.CreatedAt, time.Second)
}

func TestOfficeCurrencyDelete(t *testing.T) {
	officeCurrency := createRandomOfficeCurrency(t, 0)

	err := testQueries.DeleteOfficeCurrency(context.Background(), officeCurrency.ID)
	assert.NoError(t, err)

	officeCurrency2, err := testQueries.GetOfficeCurrency(context.Background(), officeCurrency.ID)
	assert.ErrorIs(t, err, sql.ErrNoRows)
	assert.Empty(t, officeCurrency2)
}

func TestOfficeCurrencyList(t *testing.T) {
	office := createRandomHqOffice(t)
	for i := 0; i < 10; i++ {
		createRandomOfficeCurrency(t, office.ID)
	}

	officeCurrencies, err := testQueries.ListOfficeCurrency(context.Background(), office.ID)
	assert.NoError(t, err)
	assert.Len(t, officeCurrencies, 10)

	for _, officeCurrency := range officeCurrencies {
		assert.NotEmpty(t, officeCurrency)
	}
}
