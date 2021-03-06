// Code generated by sqlc. DO NOT EDIT.
// source: account.sql

package db

import (
	"context"
	"database/sql"
)

const addAccountBalance = `-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + $1
WHERE id = $2
RETURNING id, owner, balance, currency, created_at, country_code
`

type AddAccountBalanceParams struct {
	Amount sql.NullInt64 `json:"amount"`
	ID     int64         `json:"id"`
}

func (q *Queries) AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, addAccountBalance, arg.Amount, arg.ID)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
		&i.CountryCode,
	)
	return i, err
}

const createAccount = `-- name: CreateAccount :one
INSERT INTO accounts (
  owner,
  balance,
  currency
) VALUES (
  $1, $2, $3
) RETURNING id, owner, balance, currency, created_at, country_code
`

type CreateAccountParams struct {
	Owner    sql.NullString `json:"owner"`
	Balance  sql.NullInt64  `json:"balance"`
	Currency sql.NullString `json:"currency"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, createAccount, arg.Owner, arg.Balance, arg.Currency)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
		&i.CountryCode,
	)
	return i, err
}

const createEntry = `-- name: CreateEntry :one
INSERT INTO entries(
  account_id,
  amount
) VALUES(
  $1, $2
) RETURNING id, account_id, amount, created_at
`

type CreateEntryParams struct {
	AccountID int64 `json:"account_id"`
	Amount    int64 `json:"amount"`
}

func (q *Queries) CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error) {
	row := q.db.QueryRowContext(ctx, createEntry, arg.AccountID, arg.Amount)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const createTransfer = `-- name: CreateTransfer :one
INSERT INTO transfers(
  from_account_id,
  to_account_id,
  amount
) VALUES (
  $1, $2, $3
) RETURNING id, from_account_id, to_account_id, amount, created_at
`

type CreateTransferParams struct {
	FromAccountID sql.NullInt64 `json:"from_account_id"`
	ToAccountID   sql.NullInt64 `json:"to_account_id"`
	Amount        int64         `json:"amount"`
}

func (q *Queries) CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, createTransfer, arg.FromAccountID, arg.ToAccountID, arg.Amount)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const deleteAccount = `-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1
`

func (q *Queries) DeleteAccount(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteAccount, id)
	return err
}

const getAccount = `-- name: GetAccount :one
SELECT id, owner, balance, currency, created_at, country_code FROM accounts
Where id = $1 LIMIT 1
`

func (q *Queries) GetAccount(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccount, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
		&i.CountryCode,
	)
	return i, err
}

const getAccountForUpdate = `-- name: GetAccountForUpdate :one
SELECT id, owner, balance, currency, created_at, country_code FROM accounts 
Where id = $1 LIMIT 1
FOR NO KEY UPDATE
`

func (q *Queries) GetAccountForUpdate(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccountForUpdate, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
		&i.CountryCode,
	)
	return i, err
}

const getEntry = `-- name: GetEntry :one
SELECT id, account_id, amount, created_at FROM entries
Where id = $1 LIMIT 1
`

func (q *Queries) GetEntry(ctx context.Context, id int64) (Entry, error) {
	row := q.db.QueryRowContext(ctx, getEntry, id)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const getTransfer = `-- name: GetTransfer :one
SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers
Where id = $1 LIMIT 1
`

func (q *Queries) GetTransfer(ctx context.Context, id int64) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, getTransfer, id)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const listAccounts = `-- name: ListAccounts :many
SELECT id, owner, balance, currency, created_at, country_code FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListAccountsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error) {
	rows, err := q.db.QueryContext(ctx, listAccounts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Account
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Balance,
			&i.Currency,
			&i.CreatedAt,
			&i.CountryCode,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAccount = `-- name: UpdateAccount :one
UPDATE accounts 
SET balance = $2
WHERE id = $1
RETURNING id, owner, balance, currency, created_at, country_code
`

type UpdateAccountParams struct {
	ID      int64         `json:"id"`
	Balance sql.NullInt64 `json:"balance"`
}

func (q *Queries) UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, updateAccount, arg.ID, arg.Balance)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
		&i.CountryCode,
	)
	return i, err
}

const updateAccountBalance = `-- name: UpdateAccountBalance :one
UPDATE accounts 
SET balance = balance + $2
WHERE id = $1
RETURNING id, owner, balance, currency, created_at, country_code
`

type UpdateAccountBalanceParams struct {
	ID      int64         `json:"id"`
	Balance sql.NullInt64 `json:"balance"`
}

func (q *Queries) UpdateAccountBalance(ctx context.Context, arg UpdateAccountBalanceParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, updateAccountBalance, arg.ID, arg.Balance)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
		&i.CountryCode,
	)
	return i, err
}
