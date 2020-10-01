package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store fornisce tutte le funzionalità per query e transactions
type Store struct {
	// sfruttiamo COMPOSIZIONE invece di EREDITARIETA' per estendere le interfacce
	*Queries // tutte le query presenti nel QUERIES saranno disponibili nello store
	db       *sql.DB
}

// TransferTxParams param di input
type TransferTxParams struct {
	FromAccountiD int64 `json:"from_account_id"`
	ToAccountiD   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult param di output della transazione
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// NewStore istanzia un nuovo store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTX crea una nuova transaction partendo dal contesto e passandola alla funzione di callback passata in input
// e facendo commit o rollback alla fine in funzione dell'error ritornato da quella funzione
// NB: non esportata perchè si vuole creare una funzione per ogni specifica transazione
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// TransferTX definisce una transazione legata ad una singola bank transfer
func (store *Store) TransferTX(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// 1) creazione transfer
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountiD,
			ToAccountID:   arg.ToAccountiD,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// 2) creazione entry A
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountiD,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// 3) creazione entry B
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountiD,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		//2a) aggiorno balance
		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.FromAccountiD,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		//2b) aggiorno balance
		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.ToAccountiD,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		//e poi fare aggiornamento
		return nil
	})
	return result, err
}
