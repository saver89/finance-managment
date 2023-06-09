package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)
	require.NotEmpty(t, officeCurrency)
	require.NotZero(t, officeCurrency.ID)
	require.Equal(t, arg.OfficeID, officeCurrency.OfficeID)
	require.Equal(t, arg.CurrencyID, officeCurrency.CurrencyID)
	require.Equal(t, arg.Type, officeCurrency.Type)
	require.NotEmpty(t, officeCurrency.CreatedAt)

	return officeCurrency
}

func TestOfficeCurrencyCreate(t *testing.T) {
	createRandomOfficeCurrency(t, 0)
}

func TestOfficeCurrencyGet(t *testing.T) {
	officeCurrency1 := createRandomOfficeCurrency(t, 0)

	officeCurrency2, err := testQueries.GetOfficeCurrency(context.Background(), officeCurrency1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, officeCurrency2)
	require.Equal(t, officeCurrency1.ID, officeCurrency2.ID)
	require.Equal(t, officeCurrency1.OfficeID, officeCurrency2.OfficeID)
	require.Equal(t, officeCurrency1.CurrencyID, officeCurrency2.CurrencyID)
	require.Equal(t, officeCurrency1.Type, officeCurrency2.Type)
	require.WithinDuration(t, officeCurrency1.CreatedAt, officeCurrency2.CreatedAt, time.Second)
}

func TestOfficeCurrencyDelete(t *testing.T) {
	officeCurrency := createRandomOfficeCurrency(t, 0)

	err := testQueries.DeleteOfficeCurrency(context.Background(), officeCurrency.ID)
	require.NoError(t, err)

	officeCurrency2, err := testQueries.GetOfficeCurrency(context.Background(), officeCurrency.ID)
	require.ErrorIs(t, err, sql.ErrNoRows)
	require.Empty(t, officeCurrency2)
}

func TestOfficeCurrencyList(t *testing.T) {
	office := createRandomHqOffice(t)
	for i := 0; i < 10; i++ {
		createRandomOfficeCurrency(t, office.ID)
	}

	officeCurrencies, err := testQueries.ListOfficeCurrency(context.Background(), office.ID)
	require.NoError(t, err)
	require.Len(t, officeCurrencies, 10)

	for _, officeCurrency := range officeCurrencies {
		require.NotEmpty(t, officeCurrency)
	}
}
