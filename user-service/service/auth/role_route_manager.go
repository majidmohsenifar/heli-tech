package auth

import (
	"context"

	"github.com/majidmohsenifar/heli-tech/user-service/core"
	"github.com/majidmohsenifar/heli-tech/user-service/repository"
)

const (
	RoleEndUser = "END_USER"
)

type RoleRouteManager struct {
	db   core.PgxInterface
	repo repository.Querier
}

func (m *RoleRouteManager) GetRoleByCode(ctx context.Context, code string) (repository.Role, error) {
	return m.repo.GetRoleByCode(ctx, m.db, RoleEndUser)
}

func NewRoleRouteManager(
	db core.PgxInterface,
	repo repository.Querier,
) *RoleRouteManager {
	return &RoleRouteManager{
		db:   db,
		repo: repo,
	}

}
