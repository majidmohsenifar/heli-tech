package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/majidmohsenifar/heli-tech/user-service/repository"
)

const (
	RoleEndUser = "END_USER"
)

type RoleRouteManager struct {
	db   *pgxpool.Pool
	repo *repository.Queries
}

func (m *RoleRouteManager) GetRoleByCode(ctx context.Context, code string) (repository.Role, error) {
	return m.repo.GetRoleByCode(ctx, m.db, RoleEndUser)
}

func NewRoleRouteManager(
	db *pgxpool.Pool,
	repo *repository.Queries,
) *RoleRouteManager {
	return &RoleRouteManager{
		db:   db,
		repo: repo,
	}

}
