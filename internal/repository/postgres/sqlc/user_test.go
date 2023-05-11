package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func createRandomUser(t *testing.T, officeID int64) User {
	if officeID == 0 {
		office := createRandomHqOffice(t)
		officeID = office.ID
	}

	person := gofakeit.Person()
	arg := CreateUserParams{
		OfficeID:     officeID,
		Username:     gofakeit.Username(),
		PasswordHash: gofakeit.Password(true, true, true, false, false, 10),
		FirstName:    sql.NullString{String: person.FirstName, Valid: true},
		LastName:     sql.NullString{String: person.LastName, Valid: true},
		MiddleName:   sql.NullString{String: gofakeit.Name(), Valid: true},
		Birthday:     sql.NullTime{Time: gofakeit.Date(), Valid: true},
		Email:        sql.NullString{String: gofakeit.Email(), Valid: true},
		Phone:        sql.NullString{String: gofakeit.Phone(), Valid: true},
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	assert.Equal(t, err, nil)
	assert.NotEmpty(t, user)
	assert.NotZero(t, user.ID)
	assert.Equal(t, arg.OfficeID, user.OfficeID)
	assert.Equal(t, arg.Username, user.Username)
	assert.Equal(t, arg.PasswordHash, user.PasswordHash)
	assert.Equal(t, arg.FirstName, user.FirstName)
	assert.Equal(t, arg.LastName, user.LastName)
	assert.Equal(t, arg.MiddleName, user.MiddleName)
	assert.WithinDuration(t, arg.Birthday.Time, user.Birthday.Time, time.Second)
	assert.Equal(t, arg.Email, user.Email)
	assert.Equal(t, arg.Phone, user.Phone)
	assert.NotEmpty(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t, 0)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t, 0)

	user2, err := testQueries.GetUser(context.Background(), user1.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, user2)
	assert.Equal(t, user1.ID, user2.ID)
	assert.Equal(t, user1.OfficeID, user2.OfficeID)
	assert.Equal(t, user1.Username, user2.Username)
	assert.Equal(t, user1.PasswordHash, user2.PasswordHash)
	assert.Equal(t, user1.FirstName, user2.FirstName)
	assert.Equal(t, user1.LastName, user2.LastName)
	assert.Equal(t, user1.MiddleName, user2.MiddleName)
	assert.WithinDuration(t, user1.Birthday.Time, user2.Birthday.Time, time.Second)
	assert.Equal(t, user1.Email, user2.Email)
	assert.Equal(t, user1.Phone, user2.Phone)
	assert.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUser(t *testing.T) {
	user1 := createRandomUser(t, 0)

	arg := UpdateUserParams{
		ID:           user1.ID,
		PasswordHash: user1.PasswordHash,
		FirstName:    sql.NullString{String: gofakeit.FirstName(), Valid: true},
		LastName:     sql.NullString{String: gofakeit.LastName(), Valid: true},
		MiddleName:   sql.NullString{String: gofakeit.Name(), Valid: true},
		Birthday:     sql.NullTime{Time: gofakeit.Date(), Valid: true},
		Email:        sql.NullString{String: gofakeit.Email(), Valid: true},
		Phone:        sql.NullString{String: gofakeit.Phone(), Valid: true},
	}

	user2, err := testQueries.UpdateUser(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, user2)
	assert.Equal(t, user1.ID, user2.ID)
	assert.Equal(t, user1.OfficeID, user2.OfficeID)
	assert.Equal(t, user1.Username, user2.Username)
	assert.Equal(t, user1.PasswordHash, user2.PasswordHash)
	assert.Equal(t, arg.FirstName, user2.FirstName)
	assert.Equal(t, arg.LastName, user2.LastName)
	assert.Equal(t, arg.MiddleName, user2.MiddleName)
	assert.WithinDuration(t, arg.Birthday.Time, user2.Birthday.Time, time.Second)
	assert.Equal(t, arg.Email, user2.Email)
	assert.Equal(t, arg.Phone, user2.Phone)
	assert.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestListUser(t *testing.T) {
	office := createRandomHqOffice(t)
	for i := 0; i < 10; i++ {
		createRandomUser(t, office.ID)
	}

	arg := ListUserParams{
		Limit:    5,
		Offset:   5,
		OfficeID: office.ID,
	}

	users, err := testQueries.ListUser(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, users)
	assert.Len(t, users, 5)

	for _, user := range users {
		assert.NotEmpty(t, user)
	}
}

func TestDeleteUser(t *testing.T) {
	user1 := createRandomUser(t, 0)

	err := testQueries.DeleteUser(context.Background(), user1.ID)
	assert.NoError(t, err)

	user2, err := testQueries.GetUser(context.Background(), user1.ID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, sql.ErrNoRows)
	assert.Empty(t, user2)
}
