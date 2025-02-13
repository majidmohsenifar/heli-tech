// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: user.sql

package repository

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    email,
    password
) VALUES (
  $1, $2 
) RETURNING id, email, password
`

type CreateUserParams struct {
	Email    string
	Password string
}

func (q *Queries) CreateUser(ctx context.Context, db DBTX, arg CreateUserParams) (User, error) {
	row := db.QueryRow(ctx, createUser, arg.Email, arg.Password)
	var i User
	err := row.Scan(&i.ID, &i.Email, &i.Password)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, password FROM users
WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, db DBTX, email string) (User, error) {
	row := db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(&i.ID, &i.Email, &i.Password)
	return i, err
}
