-- name: CreateAccount :one
INSERT INTO accounts (
  owner,
  balance,
  currency
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
Where id = $1 LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts 
SET balance = $2
WHERE id = $1
RETURNING *;


-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;


-- name: CreateTransfer :one
INSERT INTO transfers(
  from_account_id,
  to_account_id,
  amount
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: CreateEntry :one
INSERT INTO entries(
  account_id,
  amount
) VALUES(
  $1, $2
) RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers
Where id = $1 LIMIT 1;

-- name: GetEntry :one
SELECT * FROM entries
Where id = $1 LIMIT 1;

-- name: UpdateAccountBalance :one
UPDATE accounts 
SET balance = balance + $2
WHERE id = $1
RETURNING *;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts 
Where id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;