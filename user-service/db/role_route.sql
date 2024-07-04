-- name: GetRouteByPath :one
SELECT * FROM routes
WHERE path = $1;

-- name: AddRoleToUser :exec
INSERT INTO users_roles (
    user_id,
    role_id
) VALUES (
  $1, $2 
);

-- name: GetRoleByCode :one
SELECT * FROM roles 
WHERE code = $1;


-- name: CreateRole :one
INSERT INTO roles (
    code
) VALUES (
  $1 
) RETURNING *;

-- name: GetUserRolesByUserID :many
SELECT * FROM users_roles 
WHERE user_id = $1;
