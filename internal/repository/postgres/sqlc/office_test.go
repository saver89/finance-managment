package db

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

func createRandomHqOffice(t *testing.T) Office {
	arg := gofakeit.AppName()
	office, err := testQueries.CreateHQOffice(context.Background(), arg)

	require.Equal(t, err, nil)
	require.NotEmpty(t, office)
	require.NotZero(t, office.ID)
	require.Equal(t, arg, office.Name)
	require.Equal(t, OfficeTypeHq, office.Type)
	require.Equal(t, OfficeStateActive, office.State)
	require.NotEmpty(t, office.CreatedAt)
	require.Equal(t, office.ID, office.ParentID)

	return office
}

func TestCreateOffice(t *testing.T) {
	createRandomHqOffice(t)
}

func TestAddOffice(t *testing.T) {
	office1 := createRandomHqOffice(t)

	arg := AddOfficeParams{
		Name:     gofakeit.AppName(),
		ParentID: office1.ID,
		Type:     OfficeTypeStore,
	}

	office2, err := testQueries.AddOffice(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, office2)
	require.NotZero(t, office2.ID)
	require.Equal(t, arg.Name, office2.Name)
	require.Equal(t, arg.Type, office2.Type)
	require.Equal(t, OfficeStateActive, office2.State)
	require.NotEmpty(t, office2.CreatedAt)
	require.Equal(t, office1.ID, office2.ParentID)
}

func TestGetOffice(t *testing.T) {
	office1 := createRandomHqOffice(t)
	office2, err := testQueries.GetOffice(context.Background(), office1.ID)

	require.Equal(t, err, nil)
	require.NotEmpty(t, office2)
	require.Equal(t, office1.ID, office2.ID)
	require.Equal(t, office1.Name, office2.Name)
	require.Equal(t, office1.Type, office2.Type)
	require.Equal(t, office1.State, office2.State)
	require.Equal(t, office1.ParentID, office2.ParentID)
	require.WithinDuration(t, office1.CreatedAt, office2.CreatedAt, time.Second)
}

func TestUpdateOffice(t *testing.T) {
	office1 := createRandomHqOffice(t)

	arg := UpdateOfficeParams{
		ID:   office1.ID,
		Name: gofakeit.AppName(),
	}

	office2, err := testQueries.UpdateOffice(context.Background(), arg)

	require.Equal(t, err, nil)
	require.NotEmpty(t, office2)
	require.Equal(t, office1.ID, office2.ID)
	require.Equal(t, arg.Name, office2.Name)
	require.Equal(t, office1.Type, office2.Type)
	require.Equal(t, office1.State, office2.State)
	require.Equal(t, office1.ParentID, office2.ParentID)
	require.WithinDuration(t, office1.CreatedAt, office2.CreatedAt, time.Second)
}

func TestListOffice(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomHqOffice(t)
	}

	arg := ListOfficeParams{
		Limit:  5,
		Offset: 5,
		Type:   OfficeTypeHq,
	}
	offices, err := testQueries.ListOffice(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, offices, 5)
	for i := 0; i < 5; i++ {
		require.NotEmpty(t, offices[i])
	}
}
