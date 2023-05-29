package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	passwordPkg "github.com/saver89/finance-management/pkg/password"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T, officeID int64) User {
	if officeID == 0 {
		office := createRandomHqOffice(t)
		officeID = office.ID
	}

	person := gofakeit.Person()
	password := gofakeit.Password(true, true, true, false, false, 10)
	hashedPassword, err := passwordPkg.HashPassword(password)
	require.NoError(t, err)

	arg := CreateUserParams{
		OfficeID:     officeID,
		Username:     gofakeit.Username(),
		PasswordHash: hashedPassword,
		FirstName:    sql.NullString{String: person.FirstName, Valid: true},
		LastName:     sql.NullString{String: person.LastName, Valid: true},
		MiddleName:   sql.NullString{String: gofakeit.Name(), Valid: true},
		Birthday:     sql.NullTime{Time: gofakeit.Date(), Valid: true},
		Email:        sql.NullString{String: gofakeit.Email(), Valid: true},
		Phone:        sql.NullString{String: gofakeit.Phone(), Valid: true},
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.Equal(t, err, nil)
	require.NotEmpty(t, user)
	require.NotZero(t, user.ID)
	require.Equal(t, arg.OfficeID, user.OfficeID)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.PasswordHash, user.PasswordHash)
	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.LastName, user.LastName)
	require.Equal(t, arg.MiddleName, user.MiddleName)
	require.WithinDuration(t, arg.Birthday.Time, user.Birthday.Time, time.Second)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Phone, user.Phone)
	require.NotEmpty(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t, 0)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t, 0)

	user2, err := testQueries.GetUser(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.OfficeID, user2.OfficeID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.PasswordHash, user2.PasswordHash)
	require.Equal(t, user1.FirstName, user2.FirstName)
	require.Equal(t, user1.LastName, user2.LastName)
	require.Equal(t, user1.MiddleName, user2.MiddleName)
	require.WithinDuration(t, user1.Birthday.Time, user2.Birthday.Time, time.Second)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Phone, user2.Phone)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
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
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.OfficeID, user2.OfficeID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.PasswordHash, user2.PasswordHash)
	require.Equal(t, arg.FirstName, user2.FirstName)
	require.Equal(t, arg.LastName, user2.LastName)
	require.Equal(t, arg.MiddleName, user2.MiddleName)
	require.WithinDuration(t, arg.Birthday.Time, user2.Birthday.Time, time.Second)
	require.Equal(t, arg.Email, user2.Email)
	require.Equal(t, arg.Phone, user2.Phone)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
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
	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.Len(t, users, 5)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}

func TestDeleteUser(t *testing.T) {
	user1 := createRandomUser(t, 0)

	err := testQueries.DeleteUser(context.Background(), user1.ID)
	require.NoError(t, err)

	user2, err := testQueries.GetUser(context.Background(), user1.ID)
	require.Error(t, err)
	require.ErrorIs(t, err, sql.ErrNoRows)
	require.Empty(t, user2)
}
