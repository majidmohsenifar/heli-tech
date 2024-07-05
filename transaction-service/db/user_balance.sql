-- name: CreateUserBalanceOrIncreaseAmount :one
INSERT INTO user_balances (
    user_id,
    amount,
    created_at,
    updated_at
) VALUES (
    $1, $2, now(), now()
) ON CONFLICT (user_id) DO UPDATE SET amount = user_balances.amount+EXCLUDED.amount,  updated_at = now() 
RETURNING *;


-- name: CreateUserBalanceOrDecreaseAmount :one
INSERT INTO user_balances (
    user_id,
    amount,
    created_at,
    updated_at
) VALUES (
    $1, $2, now(), now()
) ON CONFLICT (user_id) DO UPDATE SET amount = user_balances.amount-EXCLUDED.amount,  updated_at = now() 
RETURNING *;

