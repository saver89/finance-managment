package db

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestCurrency(t *testing.T) {
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
}
