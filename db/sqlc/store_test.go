package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBoh(t *testing.T) {
	fmt.Println("TEST")
}

func TestTransferTxDeadLock(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before ", account1.Balance, account2.Balance)

	// run n transazioni concorrenti
	n := 10
	amount := int64(10)

	// siccome da dentro le routine non posso essere sicuro dei risultati dei test,
	// uso i channels per passare indietro alla main routine il risultato
	errs := make(chan error) // 1 canale per gli errori

	for i := 0; i < n; i++ {
		fromAccountiD := account1.ID
		toAccountiD := account2.ID

		if i%2 == 1 {
			fromAccountiD = account2.ID
			toAccountiD = account1.ID
		}

		go func() {
			_, err := store.TransferTX(context.Background(), TransferTxParams{
				FromAccountiD: fromAccountiD,
				ToAccountiD:   toAccountiD,
				Amount:        amount,
			})
			//mando fuori dalla routine, tramite i canali, errori e risultati
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

	}

	//controllo finale balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
	fmt.Println(">> after ", updatedAccount1.Balance, updatedAccount2.Balance)

}
