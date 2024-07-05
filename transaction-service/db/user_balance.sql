-- name: CreateUserBalance :one
INSERT INTO user_balances (
    user_id,
    amount,
    created_at
) VALUES (
  $1, $2, now()
) RETURNING *;

-- name: UpdateUserBalance :exec
UPDATE user_balances 
SET amount = $1 
WHERE user_id = $2;
