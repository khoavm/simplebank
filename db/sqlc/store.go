package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

//Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

//NewStore creates a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
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

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

//TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, add account entries, and update accounts balance within a single database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: sql.NullInt64{Int64: arg.FromAccountID, Valid: true},
			ToAccountID:   sql.NullInt64{Int64: arg.ToAccountID, Valid: true},
			Amount:        arg.Amount,
		})
		if err != nil {
			log.Println("create transfer fail")
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			log.Println("create from entry fail")
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			log.Println("create to entry fail")
			return err
		}

		//Update from from account balance
		// result.FromAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		// 	ID: arg.FromAccountID,
		// 	Balance: sql.NullInt64{Int64: -arg.Amount, Valid: true} ,
		// })
		// if err != nil {
		// 	return err
		// }

		// //Update to account balance
		// result.ToAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		// 	ID: arg.ToAccountID,
		// 	Balance: sql.NullInt64{Int64: arg.Amount, Valid: true} ,
		// })
		// if err != nil {
		// 	return err
		// }

		fromAccountID := arg.FromAccountID
		toAccountID := arg.ToAccountID

		if fromAccountID < toAccountID {
			result.FromAccount, result.ToAccount, err = AddMoney(ctx, q, fromAccountID, -arg.Amount, toAccountID, arg.Amount)
			if err != nil {
				return err
			}
		} else {
			result.FromAccount, result.ToAccount, err = AddMoney(ctx, q, toAccountID, arg.Amount, fromAccountID, -arg.Amount)
			if err != nil {
				return err
			}
		}

		return nil
	})
	return result, err
}

func AddMoney(ctx context.Context, q *Queries, acount1ID int64, amount1 int64, account2ID int64, amount2 int64) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     acount1ID,
		Amount: sql.NullInt64{Int64: amount1, Valid: true},
	})

	if err != nil {
		return
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     account2ID,
		Amount: sql.NullInt64{Int64: amount2, Valid: true},
	})

	if err != nil {
		return
	}
	return account1, account2, nil
}
