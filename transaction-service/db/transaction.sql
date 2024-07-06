-- name: CreateTransaction :one
INSERT INTO transactions (
    user_id,
    kind,
    amount,
    created_at
) VALUES (
  $1, $2, $3, now()
) RETURNING *;

-- name: GetTransactionByID :one
SELECT * from transactions WHERE id = $1;

-- name: GetUserTransactionsByPagination :many
SELECT * FROM transactions 
WHERE user_id = $1 
ORDER BY created_at DESC
OFFSET $2
LIMIT $3;

