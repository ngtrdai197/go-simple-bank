package db

import (
	"context"
	"testing"
	"time"

	"github.com/ngtrdai197/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	args := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, args.Username, user.Username)
	require.Equal(t, args.HashedPassword, user.HashedPassword)
	require.Equal(t, args.FullName, user.FullName)
	require.Equal(t, args.Email, user.Email)

	require.NotZero(t, user.Username)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)
	result, err := testQueries.GetUser(context.Background(), user.Username)

	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, user.Username, result.Username)
	require.Equal(t, user.HashedPassword, result.HashedPassword)
	require.Equal(t, user.FullName, result.FullName)
	require.Equal(t, user.Email, result.Email)
	require.WithinDuration(t, user.CreatedAt, result.CreatedAt, time.Second)
}
