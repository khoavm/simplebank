package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	
	account1 := testCreateAccount(t)
	account2 := testCreateAccount(t)

	// run a concurrent transfer transactions
	n := int64(3)
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)
	// check results
	existed := make(map[int64]bool)

	for i :=int64(0); i< n; i++{
		txName := fmt.Sprintf("tx %d", i + 1)
		go func(){
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID: account2.ID,
				Amount: amount,
			})
		errs <- err
		results <- result
		}()
	}
	fmt.Println(">>before tx:", account1.Balance.Int64, account2.Balance.Int64)
	// check results
	for i := int64(0); i < n; i++{
		
		err := <- errs
		require.NoError(t, err)

		result := <- results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID.Int64)
		require.Equal(t, account2.ID, transfer.ToAccountID.Int64)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_,err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)
		
		toEntry := result.ToEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_,err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//check account balance
		fromAccount := result.FromAccount
		toAccount := result.ToAccount


		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)
	
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)


		// check accounts balance
		fmt.Println(">> tx:", fromAccount.Balance.Int64, toAccount.Balance.Int64)
		diff1 := account1.Balance.Int64 - fromAccount.Balance.Int64
		diff2 := toAccount.Balance.Int64 - account2.Balance.Int64
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1 % amount == 0) // amount, 2 * amount, 3 * amount

		
		k := int64(diff1 / amount)
		require.True(t, k >= 1, k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true



		//account1.Balance.Int64 -= amount
		//account2.Balance.Int64 += amount
		//require.Equal(t, account1.Balance.Int64 - amount, fromAccount.Balance.Int64 )
		//require.Equal(t, account2.Balance.Int64 + amount, toAccount.Balance.Int64 )
	}
	//check the final updated balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance.Int64 - n * amount, updatedAccount1.Balance.Int64)
	require.Equal(t, account2.Balance.Int64 + n * amount, updatedAccount2.Balance.Int64)
}

func TestTransferTxDeadLock(t *testing.T) {
	store := NewStore(testDB)
	
	account1 := testCreateAccount(t)
	account2 := testCreateAccount(t)

	// run a concurrent transfer transactions
	n := int64(10)
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)
	// check results
	

	for i :=int64(0); i< n; i++{
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i % 2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		txName := fmt.Sprintf("tx %d", i + 1)
		go func(){
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID: toAccountID,
				Amount: amount,
			})
		errs <- err
		results <- result
		}()
	}
	fmt.Println(">>before tx:", account1.Balance.Int64, account2.Balance.Int64)
	// check results
	for i := int64(0); i < n; i++{
		
		err := <- errs
		require.NoError(t, err)





		
	}
	//check the final updated balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance.Int64, updatedAccount1.Balance.Int64)
	require.Equal(t, account2.Balance.Int64, updatedAccount2.Balance.Int64)
}