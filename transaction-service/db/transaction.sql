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
