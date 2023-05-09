package db

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func createRandomCurrency(t *testing.T) Currency {
	randCurrency := gofakeit.Currency()
	arg := CreateCurrencyParams{
		Name:      randCurrency.Long,
		ShortName: randCurrency.Short,
	}

	currency, err := testQueries.CreateCurrency(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, currency)

	assert.Equal(t, arg.Name, currency.Name)
	assert.Equal(t, arg.ShortName, currency.ShortName)
	assert.Equal(t, CurrencyStateActive, currency.State)
	assert.NotZero(t, currency.ID)
	assert.NotZero(t, currency.CreatedAt)

	return currency
}

func TestCreateCurrency(t *testing.T) {
	createRandomCurrency(t)
}

func TestGetCurrency(t *testing.T) {
	currency1 := createRandomCurrency(t)
	currency2, err := testQueries.GetCurrency(context.Background(), currency1.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, currency2)

	assert.Equal(t, currency1.ID, currency2.ID)
	assert.Equal(t, currency1.Name, currency2.Name)
	assert.Equal(t, currency1.ShortName, currency2.ShortName)
	assert.Equal(t, currency1.State, currency2.State)
	assert.WithinDuration(t, currency1.CreatedAt, currency2.CreatedAt, time.Second)
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
	assert.NoError(t, err)
	assert.NotEmpty(t, currency2)
	assert.Equal(t, arg.ID, currency2.ID)
	assert.Equal(t, arg.Name, currency2.Name)
	assert.Equal(t, arg.ShortName, currency2.ShortName)
	assert.Equal(t, currency1.State, currency2.State)
	assert.WithinDuration(t, currency1.CreatedAt, currency2.CreatedAt, time.Second)
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
	assert.NoError(t, err)
	assert.Len(t, currencies, 5)
	for _, currency := range currencies {
		assert.NotEmpty(t, currency)
	}
}
