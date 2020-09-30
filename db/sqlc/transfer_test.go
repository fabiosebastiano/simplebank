package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

// per testare tutte le operazioni diverse dalla create, occorre poter utilizzare ogni volta
// un oggetto diverso, che va creato randomicamente
func createRandomTransfer(t *testing.T) Transfer {
	arg := CreateTransferParams{
		FromAccountID: 1,
		ToAccountID:   2,
		Amount:        util.RandomMoney(),
	}

	Transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, Transfer)
	require.Equal(t, arg.FromAccountID, Transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, Transfer.ToAccountID)
	require.Equal(t, arg.Amount, Transfer.Amount)
	require.NotZero(t, Transfer.ID)
	require.NotZero(t, Transfer.CreatedAt)
	return Transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	Transfer1 := createRandomTransfer(t)
	Transfer2, err := testQueries.GetTransfer(context.Background(), Transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, Transfer2)
	require.Equal(t, Transfer1.ID, Transfer2.ID)
	require.Equal(t, Transfer1.FromAccountID, Transfer2.FromAccountID)
	require.Equal(t, Transfer1.ToAccountID, Transfer2.ToAccountID)
	require.Equal(t, Transfer1.Amount, Transfer2.Amount)
	require.WithinDuration(t, Transfer1.CreatedAt, Transfer2.CreatedAt, time.Second)

}

func TestListTransfers(t *testing.T) {

	for i := 0; i < 10; i++ {
		createRandomTransfer(t)
	}

	arg := ListTransfersParams{
		FromAccountID: 1,
		ToAccountID:   2,
		Limit:         5,
		Offset:        5,
	}

	Transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, Transfers, 5)

	for _, Transfer := range Transfers {
		require.NotEmpty(t, Transfer)
	}
}
