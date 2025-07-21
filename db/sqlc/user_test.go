package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sonzai8/golang-sonzai-bank/utils"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := utils.HashPassword(utils.RandomString(10))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       utils.GenerateVietnameseStyleUsername(),
		HashedPassword: hashedPassword,
		FullName:       utils.RandomOwner(),
		Email:          utils.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.CreatedAt)

	return user

}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.NotZero(t, user1.CreatedAt)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := createRandomUser(t)
	newFullName := utils.RandomOwner()
	arg := UpdateUserParams{
		FullName: pgtype.Text{
			String: newFullName,
			Valid:  true,
		},
		Username: oldUser.Username,
	}

	newUser, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newUser)
	require.NotEqual(t, oldUser.FullName, newUser.FullName)
	require.Equal(t, newFullName, newUser.FullName)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := createRandomUser(t)
	newEmail := utils.RandomEmail()
	arg := UpdateUserParams{
		Email: pgtype.Text{
			String: newEmail,
			Valid:  true,
		},
		Username: oldUser.Username,
	}

	newUser, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newUser)
	require.NotEqual(t, oldUser.Email, newUser.Email)
	require.Equal(t, newEmail, newUser.Email)
	require.Equal(t, oldUser.FullName, newUser.FullName)
	require.Equal(t, oldUser.HashedPassword, newUser.HashedPassword)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	oldUser := createRandomUser(t)
	newPassword, err := utils.HashPassword(utils.RandomString(10))
	require.NoError(t, err)
	arg := UpdateUserParams{
		HashedPassword: pgtype.Text{
			String: newPassword,
			Valid:  true,
		},
		Username: oldUser.Username,
	}

	newUser, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.NotEqual(t, oldUser.HashedPassword, newUser.HashedPassword)
	require.Equal(t, newPassword, newUser.HashedPassword)
	//require.WithinDuration(t, time.Now(), newUser.PasswordChangedAt, time.Second)
	require.Equal(t, oldUser.FullName, newUser.FullName)

}

func TestUpdateUserAllField(t *testing.T) {
	oldUser := createRandomUser(t)

	newPassword, err := utils.HashPassword(utils.RandomString(10))
	newEmail := utils.RandomEmail()
	newFullName := utils.RandomOwner()

	require.NoError(t, err)
	arg := UpdateUserParams{
		HashedPassword: pgtype.Text{
			String: newPassword,
			Valid:  true,
		},
		Email: pgtype.Text{
			String: newEmail,
			Valid:  true,
		},
		FullName: pgtype.Text{
			String: newFullName,
			Valid:  true,
		},

		Username: oldUser.Username,
	}

	newUser, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.NotEqual(t, oldUser.HashedPassword, newUser.HashedPassword)
	require.NotEqual(t, oldUser.Email, newUser.Email)
	require.NotEqual(t, oldUser.FullName, newUser.FullName)

	require.Equal(t, newEmail, newUser.Email)
	require.Equal(t, newPassword, newUser.HashedPassword)
	require.Equal(t, newFullName, newUser.FullName)

}
