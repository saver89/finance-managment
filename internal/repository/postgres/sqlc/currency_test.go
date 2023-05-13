package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

func createRandomCurrency(t *testing.T) Currency {
	randCurrency := gofakeit.Currency()
	arg := CreateCurrencyParams{
		Name:      randCurrency.Long,
		ShortName: randCurrency.Short,
	}

	currency, err := testQueries.CreateCurrency(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, currency)

	require.Equal(t, arg.Name, currency.Name)
	require.Equal(t, arg.ShortName, currency.ShortName)
	require.Equal(t, CurrencyStateActive, currency.State)
	require.NotZero(t, currency.ID)
	require.NotZero(t, currency.CreatedAt)

	return currency
}

func TestCreateCurrency(t *testing.T) {
	createRandomCurrency(t)
}

func TestGetCurrency(t *testing.T) {
	currency1 := createRandomCurrency(t)
	currency2, err := testQueries.GetCurrency(context.Background(), currency1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, currency2)

	require.Equal(t, currency1.ID, currency2.ID)
	require.Equal(t, currency1.Name, currency2.Name)
	require.Equal(t, currency1.ShortName, currency2.ShortName)
	require.Equal(t, currency1.State, currency2.State)
	require.WithinDuration(t, currency1.CreatedAt, currency2.CreatedAt, time.Second)
}

func TestUpdateCurrency(t *testing.T) {
	currency1 := createRandomCurrency(t)

	randCurrency := gofakeit.Currency()
	arg := UpdateCurrencyParams{
		ID:        currency1.ID,
		Name:      randCurrency.Long,
		ShortName: randCurrency.Short,
	}

	currency2, err := testQueries.UpdateCurrency(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, currency2)
	require.Equal(t, arg.ID, currency2.ID)
	require.Equal(t, arg.Name, currency2.Name)
	require.Equal(t, arg.ShortName, currency2.ShortName)
	require.Equal(t, currency1.State, currency2.State)
	require.WithinDuration(t, currency1.CreatedAt, currency2.CreatedAt, time.Second)
}

func TestListCurrency(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomCurrency(t)
	}

	arg := ListCurrencyParams{
		Limit:  5,
		Offset: 5,
	}

	currencies, err := testQueries.ListCurrency(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, currencies, 5)
	for _, currency := range currencies {
		require.NotEmpty(t, currency)
	}
}

func TestDeleteCurrency(t *testing.T) {
	currency1 := createRandomCurrency(t)

	err := testQueries.DeleteCurrency(context.Background(), currency1.ID)
	require.NoError(t, err)

	currency2, err := testQueries.GetCurrency(context.Background(), currency1.ID)
	require.ErrorIs(t, err, sql.ErrNoRows)
	require.Empty(t, currency2)
}
