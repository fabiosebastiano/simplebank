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
func createRandomEntry(t *testing.T) Entry {
	arg := CreateEntryParams{
		AccountID: 1,
		Amount:    util.RandomMoney(),
	}

	Entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, Entry)
	require.Equal(t, arg.AccountID, Entry.AccountID)
	require.Equal(t, arg.Amount, Entry.Amount)
	require.NotZero(t, Entry.ID)
	require.NotZero(t, Entry.CreatedAt)
	return Entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	Entry1 := createRandomEntry(t)
	Entry2, err := testQueries.GetEntry(context.Background(), Entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, Entry2)
	require.Equal(t, Entry1.ID, Entry2.ID)
	require.Equal(t, Entry1.Amount, Entry2.Amount)
	require.WithinDuration(t, Entry1.CreatedAt, Entry2.CreatedAt, time.Second)

}

func TestListEntrys(t *testing.T) {

	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}

	arg := ListEntrysParams{
		AccountID: 1,
		Limit:     5,
		Offset:    5,
	}

	Entrys, err := testQueries.ListEntrys(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, Entrys, 5)

	for _, Entry := range Entrys {
		require.NotEmpty(t, Entry)
	}
}
