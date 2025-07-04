package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func createNewTransfer(t *testing.T) Transfer {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        50,
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, account1.ID, transfer.FromAccountID)
	require.Equal(t, account2.ID, transfer.ToAccountID)
	require.Equal(t, int64(50), transfer.Amount)
	return transfer
}
func TestCreateTransfer(t *testing.T) {
	createNewTransfer(t)

}

func TestGetTransfer(t *testing.T) {
	account1 := createNewTransfer(t)

	transfer, err := testQueries.GetTransfer(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, account1.ID, transfer.ID)
	require.Equal(t, account1.Amount, transfer.Amount)
	require.Equal(t, account1.ToAccountID, transfer.ToAccountID)
	require.Equal(t, account1.FromAccountID, transfer.FromAccountID)
}
