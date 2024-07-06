package auth

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/majidmohsenifar/heli-tech/user-service/core"
	"github.com/majidmohsenifar/heli-tech/user-service/repository"
)

const (
	RoleEndUser = "END_USER"
)

var (
	ErrDefaultRoleDoesNotExist = errors.New("default role does not exist")
)

type RoleRouteManager struct {
	db   core.PgxInterface
	repo repository.Querier

	//for in-memory cache
	mutex       sync.Mutex
	roles       []repository.Role
	routes      []repository.Route
	rolesRoutes []repository.RolesRoute
}

func (m *RoleRouteManager) GetRoleByCode(ctx context.Context, code string) (repository.Role, error) {
	roles, err := m.getAllRoles(ctx)
	if err != nil {
		return repository.Role{}, err
	}
	for _, r := range roles {
		if r.Code == code {
			return r, nil
		}
	}
	return repository.Role{}, ErrDefaultRoleDoesNotExist
}

func (m *RoleRouteManager) HasUserAccessToRoute(ctx context.Context, userID int64, path string) (bool, error) {
	userRoles, err := m.repo.GetUserRolesByUserID(ctx, m.db, userID)
	if err != nil {
		return false, err
	}
	roleRoutes, err := m.getAllRoleRoutes(ctx)
	if err != nil {
		return false, err
	}
	routes, err := m.getAllRoutes(ctx)
	if err != nil {
		return false, err
	}
	routeID := int32(0)
	for _, r := range routes {
		if r.Path == path {
			routeID = r.ID
			break
		}
	}
	fmt.Println("routeID", routeID)
	if routeID == 0 {
		return false, fmt.Errorf("route with path %s not found", path)
	}
	for _, ur := range userRoles {
		for _, rr := range roleRoutes {
			if ur.RoleID == rr.RoleID && rr.RouteID == routeID {
				return true, nil
			}
		}
	}
	return false, nil
}

func (m *RoleRouteManager) getAllRoles(ctx context.Context) ([]repository.Role, error) {
	if len(m.roles) != 0 {
		return m.roles, nil
	}
	roles, err := m.repo.GetAllRoles(ctx, m.db)
	if err != nil {
		return nil, err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.roles = roles
	return roles, nil
}

func (m *RoleRouteManager) getAllRoutes(ctx context.Context) ([]repository.Route, error) {
	if len(m.routes) != 0 {
		return m.routes, nil
	}
	routes, err := m.repo.GetAllRoutes(ctx, m.db)
	if err != nil {
		return nil, err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.routes = routes
	return routes, nil
}

func (m *RoleRouteManager) getAllRoleRoutes(ctx context.Context) ([]repository.RolesRoute, error) {
	if len(m.rolesRoutes) != 0 {
		return m.rolesRoutes, nil
	}
	rolesRoutes, err := m.repo.GetAllRolesRoutes(ctx, m.db)
	if err != nil {
		return nil, err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.rolesRoutes = rolesRoutes
	return rolesRoutes, nil
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
