// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package repository

import (
	"context"
)

type Querier interface {
	AddRoleToUser(ctx context.Context, db DBTX, arg AddRoleToUserParams) error
	AddRouteToRole(ctx context.Context, db DBTX, arg AddRouteToRoleParams) error
	CreateRole(ctx context.Context, db DBTX, code string) (Role, error)
	CreateRoute(ctx context.Context, db DBTX, path string) (Route, error)
	CreateUser(ctx context.Context, db DBTX, arg CreateUserParams) (User, error)
	GetAllRoles(ctx context.Context, db DBTX) ([]Role, error)
	GetAllRolesRoutes(ctx context.Context, db DBTX) ([]RolesRoute, error)
	GetAllRoutes(ctx context.Context, db DBTX) ([]Route, error)
	GetRoleByCode(ctx context.Context, db DBTX, code string) (Role, error)
	GetRouteByPath(ctx context.Context, db DBTX, path string) (Route, error)
	GetUserByEmail(ctx context.Context, db DBTX, email string) (User, error)
	GetUserRolesByUserID(ctx context.Context, db DBTX, userID int64) ([]UsersRole, error)
}

var _ Querier = (*Queries)(nil)
