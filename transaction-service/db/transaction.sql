-- name: CreateTransaction :one
INSERT INTO transactions (
    user_id,
    kind,
    amount,
    created_at
) VALUES (
  $1, $2, $3, now()
) RETURNING *;


