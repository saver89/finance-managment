package db

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func createRandomHqOffice(t *testing.T) Office {
	arg := gofakeit.AppName()
	office, err := testQueries.CreateHQOffice(context.Background(), arg)

	assert.Equal(t, err, nil)
	assert.NotEmpty(t, office)
	assert.NotZero(t, office.ID)
	assert.Equal(t, arg, office.Name)
	assert.Equal(t, OfficeTypeHq, office.Type)
	assert.Equal(t, OfficeStateActive, office.State)
	assert.NotEmpty(t, office.CreatedAt)
	assert.Equal(t, office.ID, office.ParentID)

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
	assert.NoError(t, err)
	assert.NotEmpty(t, office2)
	assert.NotZero(t, office2.ID)
	assert.Equal(t, arg.Name, office2.Name)
	assert.Equal(t, arg.Type, office2.Type)
	assert.Equal(t, OfficeStateActive, office2.State)
	assert.NotEmpty(t, office2.CreatedAt)
	assert.Equal(t, office1.ID, office2.ParentID)
}

func TestGetOffice(t *testing.T) {
	office1 := createRandomHqOffice(t)
	office2, err := testQueries.GetOffice(context.Background(), office1.ID)

	assert.Equal(t, err, nil)
	assert.NotEmpty(t, office2)
	assert.Equal(t, office1.ID, office2.ID)
	assert.Equal(t, office1.Name, office2.Name)
	assert.Equal(t, office1.Type, office2.Type)
	assert.Equal(t, office1.State, office2.State)
	assert.Equal(t, office1.ParentID, office2.ParentID)
	assert.WithinDuration(t, office1.CreatedAt, office2.CreatedAt, time.Second)
}

func TestUpdateOffice(t *testing.T) {
	office1 := createRandomHqOffice(t)

	arg := UpdateOfficeParams{
		ID:   office1.ID,
		Name: gofakeit.AppName(),
	}

	office2, err := testQueries.UpdateOffice(context.Background(), arg)

	assert.Equal(t, err, nil)
	assert.NotEmpty(t, office2)
	assert.Equal(t, office1.ID, office2.ID)
	assert.Equal(t, arg.Name, office2.Name)
	assert.Equal(t, office1.Type, office2.Type)
	assert.Equal(t, office1.State, office2.State)
	assert.Equal(t, office1.ParentID, office2.ParentID)
	assert.WithinDuration(t, office1.CreatedAt, office2.CreatedAt, time.Second)
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
	assert.NoError(t, err)
	assert.Len(t, offices, 5)
	for i := 0; i < 5; i++ {
		assert.NotEmpty(t, offices[i])
	}
}
