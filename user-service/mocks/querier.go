// Code generated by mockery v2.33.3. DO NOT EDIT.

package mocks

import (
	context "context"

	repository "github.com/majidmohsenifar/heli-tech/user-service/repository"
	mock "github.com/stretchr/testify/mock"
)

// MockQuerier is an autogenerated mock type for the Querier type
type MockQuerier struct {
	mock.Mock
}

type MockQuerier_Expecter struct {
	mock *mock.Mock
}

func (_m *MockQuerier) EXPECT() *MockQuerier_Expecter {
	return &MockQuerier_Expecter{mock: &_m.Mock}
}

// AddRoleToUser provides a mock function with given fields: ctx, db, arg
func (_m *MockQuerier) AddRoleToUser(ctx context.Context, db repository.DBTX, arg repository.AddRoleToUserParams) error {
	ret := _m.Called(ctx, db, arg)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.DBTX, repository.AddRoleToUserParams) error); ok {
		r0 = rf(ctx, db, arg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockQuerier_AddRoleToUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddRoleToUser'
type MockQuerier_AddRoleToUser_Call struct {
	*mock.Call
}

// AddRoleToUser is a helper method to define mock.On call
//   - ctx context.Context
//   - db repository.DBTX
//   - arg repository.AddRoleToUserParams
func (_e *MockQuerier_Expecter) AddRoleToUser(ctx interface{}, db interface{}, arg interface{}) *MockQuerier_AddRoleToUser_Call {
	return &MockQuerier_AddRoleToUser_Call{Call: _e.mock.On("AddRoleToUser", ctx, db, arg)}
}

func (_c *MockQuerier_AddRoleToUser_Call) Run(run func(ctx context.Context, db repository.DBTX, arg repository.AddRoleToUserParams)) *MockQuerier_AddRoleToUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(repository.DBTX), args[2].(repository.AddRoleToUserParams))
	})
	return _c
}

func (_c *MockQuerier_AddRoleToUser_Call) Return(_a0 error) *MockQuerier_AddRoleToUser_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockQuerier_AddRoleToUser_Call) RunAndReturn(run func(context.Context, repository.DBTX, repository.AddRoleToUserParams) error) *MockQuerier_AddRoleToUser_Call {
	_c.Call.Return(run)
	return _c
}

// AddRouteToRole provides a mock function with given fields: ctx, db, arg
func (_m *MockQuerier) AddRouteToRole(ctx context.Context, db repository.DBTX, arg repository.AddRouteToRoleParams) error {
	ret := _m.Called(ctx, db, arg)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.DBTX, repository.AddRouteToRoleParams) error); ok {
		r0 = rf(ctx, db, arg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockQuerier_AddRouteToRole_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddRouteToRole'
type MockQuerier_AddRouteToRole_Call struct {
	*mock.Call
}

// AddRouteToRole is a helper method to define mock.On call
//   - ctx context.Context
//   - db repository.DBTX
//   - arg repository.AddRouteToRoleParams
func (_e *MockQuerier_Expecter) AddRouteToRole(ctx interface{}, db interface{}, arg interface{}) *MockQuerier_AddRouteToRole_Call {
	return &MockQuerier_AddRouteToRole_Call{Call: _e.mock.On("AddRouteToRole", ctx, db, arg)}
}

func (_c *MockQuerier_AddRouteToRole_Call) Run(run func(ctx context.Context, db repository.DBTX, arg repository.AddRouteToRoleParams)) *MockQuerier_AddRouteToRole_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(repository.DBTX), args[2].(repository.AddRouteToRoleParams))
	})
	return _c
}

func (_c *MockQuerier_AddRouteToRole_Call) Return(_a0 error) *MockQuerier_AddRouteToRole_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockQuerier_AddRouteToRole_Call) RunAndReturn(run func(context.Context, repository.DBTX, repository.AddRouteToRoleParams) error) *MockQuerier_AddRouteToRole_Call {
	_c.Call.Return(run)
	return _c
}

// CreateRole provides a mock function with given fields: ctx, db, code
func (_m *MockQuerier) CreateRole(ctx context.Context, db repository.DBTX, code string) (repository.Role, error) {
	ret := _m.Called(ctx, db, code)

	var r0 repository.Role
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.DBTX, string) (repository.Role, error)); ok {
		return rf(ctx, db, code)
	}
	if rf, ok := ret.Get(0).(func(context.Context, repository.DBTX, string) repository.Role); ok {
		r0 = rf(ctx, db, code)
	} else {
		r0 = ret.Get(0).(repository.Role)
	}

	if rf, ok := ret.Get(1).(func(context.Context, repository.DBTX, string) error); ok {
		r1 = rf(ctx, db, code)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockQuerier_CreateRole_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateRole'
type MockQuerier_CreateRole_Call struct {
	*mock.Call
}

// CreateRole is a helper method to define mock.On call
//   - ctx context.Context
//   - db repository.DBTX
//   - code string
func (_e *MockQuerier_Expecter) CreateRole(ctx interface{}, db interface{}, code interface{}) *MockQuerier_CreateRole_Call {
	return &MockQuerier_CreateRole_Call{Call: _e.mock.On("CreateRole", ctx, db, code)}
}

func (_c *MockQuerier_CreateRole_Call) Run(run func(ctx context.Context, db repository.DBTX, code string)) *MockQuerier_CreateRole_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(repository.DBTX), args[2].(string))
	})
	return _c
}

func (_c *MockQuerier_CreateRole_Call) Return(_a0 repository.Role, _a1 error) *MockQuerier_CreateRole_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockQuerier_CreateRole_Call) RunAndReturn(run func(context.Context, repository.DBTX, string) (repository.Role, error)) *MockQuerier_CreateRole_Call {
	_c.Call.Return(run)
	return _c
}

// CreateRoute provides a mock function with given fields: ctx, db, path
func (_m *MockQuerier) CreateRoute(ctx context.Context, db repository.DBTX, path string) (repository.Route, error) {
	ret := _m.Called(ctx, db, path)

	var r0 repository.Route
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.DBTX, string) (repository.Route, error)); ok {
		return rf(ctx, db, path)
	}
	if rf, ok := ret.Get(0).(func(context.Context, repository.DBTX, string) repository.Route); ok {
		r0 = rf(ctx, db, path)
	} else {
		r0 = ret.Get(0).(repository.Route)
	}

	if rf, ok := ret.Get(1).(func(context.Context, repository.DBTX, string) error); ok {
		r1 = rf(ctx, db, path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockQuerier_CreateRoute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateRoute'
type MockQuerier_CreateRoute_Call struct {
	*mock.Call
}

// CreateRoute is a helper method to define mock.On call
//   - ctx context.Context
//   - db repository.DBTX
//   - path string
func (_e *MockQuerier_Expecter) CreateRoute(ctx interface{}, db interface{}, path interface{}) *MockQuerier_CreateRoute_Call {
	return &MockQuerier_CreateRoute_Call{Call: _e.mock.On("CreateRoute", ctx, db, path)}
}

func (_c *MockQuerier_CreateRoute_Call) Run(run func(ctx context.Context, db repository.DBTX, path string)) *MockQuerier_CreateRoute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(repository.DBTX), args[2].(string))
	})
	return _c
}

func (_c *MockQuerier_CreateRoute_Call) Return(_a0 repository.Route, _a1 error) *MockQuerier_CreateRoute_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockQuerier_CreateRoute_Call) RunAndReturn(run func(context.Context, repository.DBTX, string) (repository.Route, error)) *MockQuerier_CreateRoute_Call {
	_c.Call.Return(run)
	return _c
}

// CreateUser provides a mock function with given fields: ctx, db, arg
func (_m *MockQuerier) CreateUser(ctx context.Context, db repository.DBTX, arg repository.CreateUserParams) (repository.User, error) {
	ret := _m.Called(ctx, db, arg)

	var r0 repository.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.DBTX, repository.CreateUserParams) (repository.User, error)); ok {
		return rf(ctx, db, arg)
	}
	if rf, ok := ret.Get(0).(func(context.Context, repository.DBTX, repository.CreateUserParams) repository.User); ok {
		r0 = rf(ctx, db, arg)
	} else {
		r0 = ret.Get(0).(repository.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, repository.DBTX, repository.CreateUserParams) error); ok {
		r1 = rf(ctx, db, arg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockQuerier_CreateUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateUser'
type MockQuerier_CreateUser_Call struct {
	*mock.Call
}

// CreateUser is a helper method to define mock.On call
//   - ctx context.Context
//   - db repository.DBTX
//   - arg repository.CreateUserParams
func (_e *MockQuerier_Expecter) CreateUser(ctx interface{}, db interface{}, arg interface{}) *MockQuerier_CreateUser_Call {
	return &MockQuerier_CreateUser_Call{Call: _e.mock.On("CreateUser", ctx, db, arg)}
}

func (_c *MockQuerier_CreateUser_Call) Run(run func(ctx context.Context, db repository.DBTX, arg repository.CreateUserParams)) *MockQuerier_CreateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(repository.DBTX), args[2].(repository.CreateUserParams))
	})
	return _c
}

func (_c *MockQuerier_CreateUser_Call) Return(_a0 repository.User, _a1 error) *MockQuerier_CreateUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockQuerier_CreateUser_Call) RunAndReturn(run func(context.Context, repository.DBTX, repository.CreateUserParams) (repository.User, error)) *MockQuerier_CreateUser_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllRoles provides a mock function with given fields: ctx, db
func (_m *MockQuerier) GetAllRoles(ctx context.Context, db repository.DBTX) ([]repository.Role, error) {
	ret := _m.Called(ctx, db)

	var r0 []repository.Role
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.DBTX) ([]repository.Role, error)); ok {
		return rf(ctx, db)
	}
	if rf, ok := ret.Get(0).(func(context.Context, repository.DBTX) []repository.Role); ok {
		r0 = rf(ctx, db)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]repository.Role)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, repository.DBTX) error); ok {
		r1 = rf(ctx, db)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockQuerier_GetAllRoles_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllRoles'
type MockQuerier_GetAllRoles_Call struct {
	*mock.Call
}

// GetAllRoles is a helper method to define mock.On call
//   - ctx context.Context
//   - db repository.DBTX
func (_e *MockQuerier_Expecter) GetAllRoles(ctx interface{}, db interface{}) *MockQuerier_GetAllRoles_Call {
	return &MockQuerier_GetAllRoles_Call{Call: _e.mock.On("GetAllRoles", ctx, db)}
}

func (_c *MockQuerier_GetAllRoles_Call) Run(run func(ctx context.Context, db repository.DBTX)) *MockQuerier_GetAllRoles_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(repository.DBTX))
	})
	return _c
}

func (_c *MockQuerier_GetAllRoles_Call) Return(_a0 []repository.Role, _a1 error) *MockQuerier_GetAllRoles_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockQuerier_GetAllRoles_Call) RunAndReturn(run func(context.Context, repository.DBTX) ([]repository.Role, error)) *MockQuerier_GetAllRoles_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllRolesRoutes provides a mock function with given fields: ctx, db
func (_m *MockQuerier) GetAllRolesRoutes(ctx context.Context, db repository.DBTX) ([]repository.RolesRoute, error) {
	ret := _m.Called(ctx, db)

	var r0 []repository.RolesRoute
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.DBTX) ([]repository.RolesRoute, error)); ok {
		return rf(ctx, db)
	}
	if rf, ok := ret.Get(0).(func(context.Context, repository.DBTX) []repository.RolesRoute); ok {
		r0 = rf(ctx, db)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]repository.RolesRoute)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, repository.DBTX) error); ok {
		r1 = rf(ctx, db)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockQuerier_GetAllRolesRoutes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllRolesRoutes'
type MockQuerier_GetAllRolesRoutes_Call struct {
	*mock.Call
}

// GetAllRolesRoutes is a helper method to define mock.On call
//   - ctx context.Context
//   - db repository.DBTX
func (_e *MockQuerier_Expecter) GetAllRolesRoutes(ctx interface{}, db interface{}) *MockQuerier_GetAllRolesRoutes_Call {
	return &MockQuerier_GetAllRolesRoutes_Call{Call: _e.mock.On("GetAllRolesRoutes", ctx, db)}
}

func (_c *MockQuerier_GetAllRolesRoutes_Call) Run(run func(ctx context.Context, db repository.DBTX)) *MockQuerier_GetAllRolesRoutes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(repository.DBTX))
	})
	return _c
}

func (_c *MockQuerier_GetAllRolesRoutes_Call) Return(_a0 []repository.RolesRoute, _a1 error) *MockQuerier_GetAllRolesRoutes_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockQuerier_GetAllRolesRoutes_Call) RunAndReturn(run func(context.Context, repository.DBTX) ([]repository.RolesRoute, error)) *MockQuerier_GetAllRolesRoutes_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllRoutes provides a mock function with given fields: ctx, db
func (_m *MockQuerier) GetAllRoutes(ctx context.Context, db repository.DBTX) ([]repository.Route, error) {
	ret := _m.Called(ctx, db)

	var r0 []repository.Route
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.DBTX) ([]repository.Route, error)); ok {
		return rf(ctx, db)
	}
	if rf, ok := ret.Get(0).(func(context.Context, repository.DBTX) []repository.Route); ok {
		r0 = rf(ctx, db)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]repository.Route)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, repository.DBTX) error); ok {
		r1 = rf(ctx, db)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockQuerier_GetAllRoutes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllRoutes'
type MockQuerier_GetAllRoutes_Call struct {
	*mock.Call
}

// GetAllRoutes is a helper method to define mock.On call
//   - ctx context.Context
//   - db repository.DBTX
func (_e *MockQuerier_Expecter) GetAllRoutes(ctx interface{}, db interface{}) *MockQuerier_GetAllRoutes_Call {
	return &MockQuerier_GetAllRoutes_Call{Call: _e.mock.On("GetAllRoutes", ctx, db)}
}

func (_c *MockQuerier_GetAllRoutes_Call) Run(run func(ctx context.Context, db repository.DBTX)) *MockQuerier_GetAllRoutes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(repository.DBTX))
	})
	return _c
}

func (_c *MockQuerier_GetAllRoutes_Call) Return(_a0 []repository.Route, _a1 error) *MockQuerier_GetAllRoutes_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockQuerier_GetAllRoutes_Call) RunAndReturn(run func(context.Context, repository.DBTX) ([]repository.Route, error)) *MockQuerier_GetAllRoutes_Call {
	_c.Call.Return(run)
	return _c
}

// GetRoleByCode provides a mock function with given fields: ctx, db, code
func (_m *MockQuerier) GetRoleByCode(ctx context.Context, db repository.DBTX, code string) (repository.Role, error) {
	ret := _m.Called(ctx, db, code)

	var r0 repository.Role
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.DBTX, string) (repository.Role, error)); ok {
		return rf(ctx, db, code)
	}
	if rf, ok := ret.Get(0).(func(context.Context, repository.DBTX, string) repository.Role); ok {
		r0 = rf(ctx, db, code)
	} else {
		r0 = ret.Get(0).(repository.Role)
	}

	if rf, ok := ret.Get(1).(func(context.Context, repository.DBTX, string) error); ok {
		r1 = rf(ctx, db, code)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockQuerier_GetRoleByCode_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRoleByCode'
type MockQuerier_GetRoleByCode_Call struct {
	*mock.Call
}

// GetRoleByCode is a helper method to define mock.On call
//   - ctx context.Context
//   - db repository.DBTX
//   - code string
func (_e *MockQuerier_Expecter) GetRoleByCode(ctx interface{}, db interface{}, code interface{}) *MockQuerier_GetRoleByCode_Call {
	return &MockQuerier_GetRoleByCode_Call{Call: _e.mock.On("GetRoleByCode", ctx, db, code)}
}

func (_c *MockQuerier_GetRoleByCode_Call) Run(run func(ctx context.Context, db repository.DBTX, code string)) *MockQuerier_GetRoleByCode_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(repository.DBTX), args[2].(string))
	})
	return _c
}

func (_c *MockQuerier_GetRoleByCode_Call) Return(_a0 repository.Role, _a1 error) *MockQuerier_GetRoleByCode_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockQuerier_GetRoleByCode_Call) RunAndReturn(run func(context.Context, repository.DBTX, string) (repository.Role, error)) *MockQuerier_GetRoleByCode_Call {
	_c.Call.Return(run)
	return _c
}

// GetRouteByPath provides a mock function with given fields: ctx, db, path
func (_m *MockQuerier) GetRouteByPath(ctx context.Context, db repository.DBTX, path string) (repository.Route, error) {
	ret := _m.Called(ctx, db, path)

	var r0 repository.Route
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.DBTX, string) (repository.Route, error)); ok {
		return rf(ctx, db, path)
	}
	if rf, ok := ret.Get(0).(func(context.Context, repository.DBTX, string) repository.Route); ok {
		r0 = rf(ctx, db, path)
	} else {
		r0 = ret.Get(0).(repository.Route)
	}

	if rf, ok := ret.Get(1).(func(context.Context, repository.DBTX, string) error); ok {
		r1 = rf(ctx, db, path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockQuerier_GetRouteByPath_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRouteByPath'
type MockQuerier_GetRouteByPath_Call struct {
	*mock.Call
}

// GetRouteByPath is a helper method to define mock.On call
//   - ctx context.Context
//   - db repository.DBTX
//   - path string
func (_e *MockQuerier_Expecter) GetRouteByPath(ctx interface{}, db interface{}, path interface{}) *MockQuerier_GetRouteByPath_Call {
	return &MockQuerier_GetRouteByPath_Call{Call: _e.mock.On("GetRouteByPath", ctx, db, path)}
}

func (_c *MockQuerier_GetRouteByPath_Call) Run(run func(ctx context.Context, db repository.DBTX, path string)) *MockQuerier_GetRouteByPath_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(repository.DBTX), args[2].(string))
	})
	return _c
}

func (_c *MockQuerier_GetRouteByPath_Call) Return(_a0 repository.Route, _a1 error) *MockQuerier_GetRouteByPath_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockQuerier_GetRouteByPath_Call) RunAndReturn(run func(context.Context, repository.DBTX, string) (repository.Route, error)) *MockQuerier_GetRouteByPath_Call {
	_c.Call.Return(run)
	return _c
}

// GetUserByEmail provides a mock function with given fields: ctx, db, email
func (_m *MockQuerier) GetUserByEmail(ctx context.Context, db repository.DBTX, email string) (repository.User, error) {
	ret := _m.Called(ctx, db, email)

	var r0 repository.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.DBTX, string) (repository.User, error)); ok {
		return rf(ctx, db, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, repository.DBTX, string) repository.User); ok {
		r0 = rf(ctx, db, email)
	} else {
		r0 = ret.Get(0).(repository.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, repository.DBTX, string) error); ok {
		r1 = rf(ctx, db, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockQuerier_GetUserByEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserByEmail'
type MockQuerier_GetUserByEmail_Call struct {
	*mock.Call
}

// GetUserByEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - db repository.DBTX
//   - email string
func (_e *MockQuerier_Expecter) GetUserByEmail(ctx interface{}, db interface{}, email interface{}) *MockQuerier_GetUserByEmail_Call {
	return &MockQuerier_GetUserByEmail_Call{Call: _e.mock.On("GetUserByEmail", ctx, db, email)}
}

func (_c *MockQuerier_GetUserByEmail_Call) Run(run func(ctx context.Context, db repository.DBTX, email string)) *MockQuerier_GetUserByEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(repository.DBTX), args[2].(string))
	})
	return _c
}

func (_c *MockQuerier_GetUserByEmail_Call) Return(_a0 repository.User, _a1 error) *MockQuerier_GetUserByEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockQuerier_GetUserByEmail_Call) RunAndReturn(run func(context.Context, repository.DBTX, string) (repository.User, error)) *MockQuerier_GetUserByEmail_Call {
	_c.Call.Return(run)
	return _c
}

// GetUserRolesByUserID provides a mock function with given fields: ctx, db, userID
func (_m *MockQuerier) GetUserRolesByUserID(ctx context.Context, db repository.DBTX, userID int64) ([]repository.UsersRole, error) {
	ret := _m.Called(ctx, db, userID)

	var r0 []repository.UsersRole
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.DBTX, int64) ([]repository.UsersRole, error)); ok {
		return rf(ctx, db, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, repository.DBTX, int64) []repository.UsersRole); ok {
		r0 = rf(ctx, db, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]repository.UsersRole)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, repository.DBTX, int64) error); ok {
		r1 = rf(ctx, db, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockQuerier_GetUserRolesByUserID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserRolesByUserID'
type MockQuerier_GetUserRolesByUserID_Call struct {
	*mock.Call
}

// GetUserRolesByUserID is a helper method to define mock.On call
//   - ctx context.Context
//   - db repository.DBTX
//   - userID int64
func (_e *MockQuerier_Expecter) GetUserRolesByUserID(ctx interface{}, db interface{}, userID interface{}) *MockQuerier_GetUserRolesByUserID_Call {
	return &MockQuerier_GetUserRolesByUserID_Call{Call: _e.mock.On("GetUserRolesByUserID", ctx, db, userID)}
}

func (_c *MockQuerier_GetUserRolesByUserID_Call) Run(run func(ctx context.Context, db repository.DBTX, userID int64)) *MockQuerier_GetUserRolesByUserID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(repository.DBTX), args[2].(int64))
	})
	return _c
}

func (_c *MockQuerier_GetUserRolesByUserID_Call) Return(_a0 []repository.UsersRole, _a1 error) *MockQuerier_GetUserRolesByUserID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockQuerier_GetUserRolesByUserID_Call) RunAndReturn(run func(context.Context, repository.DBTX, int64) ([]repository.UsersRole, error)) *MockQuerier_GetUserRolesByUserID_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockQuerier creates a new instance of MockQuerier. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockQuerier(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockQuerier {
	mock := &MockQuerier{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
