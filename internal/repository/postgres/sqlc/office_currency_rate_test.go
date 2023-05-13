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

type createCurrencyRateParam struct {
	OfficeID       int64 `json:"office_id"`
	FromCurrencyID int64 `json:"from_currency_id"`
	ToCurrencyID   int64 `json:"to_currency_id"`
}

func createRandomOfficeCurrencyRate(t *testing.T, param createCurrencyRateParam) OfficeCurrencyRate {
	if param.OfficeID == 0 {
		office := createRandomHqOffice(t)
		param.OfficeID = office.ID
	}
	if param.FromCurrencyID == 0 {
		currency := createRandomCurrency(t)
		param.FromCurrencyID = currency.ID
	}
	if param.ToCurrencyID == 0 {
		currency := createRandomCurrency(t)
		param.ToCurrencyID = currency.ID
	}

	rate := gofakeit.Float64Range(1, 100000)
	arg := CreateOfficeCurrencyRateParams{
		OfficeID:       param.OfficeID,
		FromCurrencyID: param.FromCurrencyID,
		ToCurrencyID:   param.ToCurrencyID,
		Rate:           fmt.Sprintf("%f", rate),
	}

	officeCurrencyRate, err := testQueries.CreateOfficeCurrencyRate(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, officeCurrencyRate)
	require.NotZero(t, officeCurrencyRate.ID)
	require.Equal(t, arg.OfficeID, officeCurrencyRate.OfficeID)
	require.Equal(t, arg.FromCurrencyID, officeCurrencyRate.FromCurrencyID)
	require.Equal(t, arg.ToCurrencyID, officeCurrencyRate.ToCurrencyID)
	require.Equal(t, arg.Rate, officeCurrencyRate.Rate)
	require.NotEmpty(t, officeCurrencyRate.CreatedAt)

	return officeCurrencyRate
}

func TestCreateOfficeCurrencyRate(t *testing.T) {
	createRandomOfficeCurrencyRate(t, createCurrencyRateParam{})
}

func TestGetOfficeCurrencyRate(t *testing.T) {
	officeCurrencyRate1 := createRandomOfficeCurrencyRate(t, createCurrencyRateParam{})

	officeCurrencyRate2, err := testQueries.GetOfficeCurrencyRate(context.Background(), officeCurrencyRate1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, officeCurrencyRate2)
	require.Equal(t, officeCurrencyRate1.ID, officeCurrencyRate2.ID)
	require.Equal(t, officeCurrencyRate1.OfficeID, officeCurrencyRate2.OfficeID)
	require.Equal(t, officeCurrencyRate1.FromCurrencyID, officeCurrencyRate2.FromCurrencyID)
	require.Equal(t, officeCurrencyRate1.ToCurrencyID, officeCurrencyRate2.ToCurrencyID)
	require.Equal(t, officeCurrencyRate1.Rate, officeCurrencyRate2.Rate)
	require.WithinDuration(t, officeCurrencyRate1.CreatedAt, officeCurrencyRate2.CreatedAt, time.Second)
}

func TestListOfficeCurrencyRate(t *testing.T) {
	office := createRandomHqOffice(t)
	fromCurrency := createRandomCurrency(t)
	toCurrency := createRandomCurrency(t)
	for i := 0; i < 10; i++ {
		createRandomOfficeCurrencyRate(t, createCurrencyRateParam{office.ID, fromCurrency.ID, toCurrency.ID})
	}

	arg := ListOfficeCurrencyRateParams{
		OfficeID:       office.ID,
		FromCurrencyID: fromCurrency.ID,
		ToCurrencyID:   toCurrency.ID,
		Limit:          5,
		Offset:         5,
	}

	officeCurrencyRates, err := testQueries.ListOfficeCurrencyRate(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, officeCurrencyRates, 5)

	for _, officeCurrencyRate := range officeCurrencyRates {
		require.NotEmpty(t, officeCurrencyRate)
	}
}

func TestDeleteOfficeCurrencyRate(t *testing.T) {
	officeCurrencyRate1 := createRandomOfficeCurrencyRate(t, createCurrencyRateParam{})

	err := testQueries.DeleteOfficeCurrencyRate(context.Background(), officeCurrencyRate1.ID)
	require.NoError(t, err)

	officeCurrencyRate2, err := testQueries.GetOfficeCurrencyRate(context.Background(), officeCurrencyRate1.ID)
	require.Error(t, err)
	require.ErrorIs(t, err, sql.ErrNoRows)
	require.Empty(t, officeCurrencyRate2)
}

func TestUpdateOfficeCurrencyRate(t *testing.T) {
	officeCurrencyRate1 := createRandomOfficeCurrencyRate(t, createCurrencyRateParam{})

	rate := gofakeit.Float64Range(1, 100000)
	arg := UpdateOfficeCurrencyRateParams{
		ID:   officeCurrencyRate1.ID,
		Rate: fmt.Sprintf("%f", rate),
	}

	officeCurrencyRate2, err := testQueries.UpdateOfficeCurrencyRate(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, officeCurrencyRate2)
	require.Equal(t, officeCurrencyRate1.ID, officeCurrencyRate2.ID)
	require.Equal(t, officeCurrencyRate1.OfficeID, officeCurrencyRate2.OfficeID)
	require.Equal(t, officeCurrencyRate1.FromCurrencyID, officeCurrencyRate2.FromCurrencyID)
	require.Equal(t, officeCurrencyRate1.ToCurrencyID, officeCurrencyRate2.ToCurrencyID)
	require.Equal(t, arg.Rate, officeCurrencyRate2.Rate)
	require.WithinDuration(t, officeCurrencyRate1.CreatedAt, officeCurrencyRate2.CreatedAt, time.Second)
}
