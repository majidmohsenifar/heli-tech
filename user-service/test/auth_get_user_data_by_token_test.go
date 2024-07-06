package test

import (
	"context"
	"net"
	"testing"

	userpb "github.com/majidmohsenifar/heli-tech/data-contracts/proto/user"
	"github.com/majidmohsenifar/heli-tech/user-service/core"
	"github.com/majidmohsenifar/heli-tech/user-service/handler/usergrpc"
	"github.com/majidmohsenifar/heli-tech/user-service/repository"
	"github.com/majidmohsenifar/heli-tech/user-service/service/auth"
	"github.com/majidmohsenifar/heli-tech/user-service/service/jwt"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestAuth_GetUserDataByToken_InvalidInputs(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	userGrpcServer := usergrpc.NewServer(nil)

	l, err := net.Listen("tcp", "127.0.0.1:0")
	assert.Nil(err)
	defer l.Close()
	googleGrpcServer := grpc.NewServer()

	userpb.RegisterUserServer(googleGrpcServer, userGrpcServer)
	userConn, err := grpc.NewClient(
		l.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.Nil(err)
	go googleGrpcServer.Serve(l)
	client := userpb.NewUserClient(userConn)

	//empty token
	req := userpb.GetUserDataByTokenRequest{
		Token: " ",
		Path:  "",
	}
	res, err := client.GetUserDataByToken(ctx, &req)
	assert.Nil(res)
	e, ok := status.FromError(err)
	assert.True(ok)
	assert.Equal(e.Code(), codes.Code(400))
	assert.Equal(e.Message(), "token is empty")

	//empty path
	req = userpb.GetUserDataByTokenRequest{
		Token: "token",
		Path:  "",
	}
	res, err = client.GetUserDataByToken(ctx, &req)
	assert.Nil(res)
	e, ok = status.FromError(err)
	assert.True(ok)
	assert.Equal(e.Code(), codes.Code(400))
	assert.Equal(e.Message(), "path is empty")
}

func TestAuth_GetUserDataByToken_InvalidToken(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	db := getDB()
	err := truncateDB()
	passwordEncoder := core.NewPasswordEncoder()
	assert.Nil(err)
	repo := repository.New()
	viper := getViperConfig()
	jwtService, err := jwt.NewService(viper)
	assert.Nil(err)
	authService := auth.NewService(db, repo, passwordEncoder, jwtService, getLogger(), nil)
	userGrpcServer := usergrpc.NewServer(authService)

	l, err := net.Listen("tcp", "127.0.0.1:0")
	assert.Nil(err)
	defer l.Close()
	googleGrpcServer := grpc.NewServer()

	userpb.RegisterUserServer(googleGrpcServer, userGrpcServer)
	userConn, err := grpc.NewClient(
		l.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.Nil(err)
	go googleGrpcServer.Serve(l)
	client := userpb.NewUserClient(userConn)

	req := userpb.GetUserDataByTokenRequest{
		Token: "invalidToken",
		Path:  "/api/v1/transactions/withdraw",
	}
	res, err := client.GetUserDataByToken(ctx, &req)
	assert.Nil(res)
	e, ok := status.FromError(err)
	assert.True(ok)
	assert.Equal(e.Code(), codes.Code(403))
	assert.Equal(e.Message(), "invalid token")
}

func TestAuth_GetUserDataByToken_DoesNotHaveAccess(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	db := getDB()
	err := truncateDB()
	passwordEncoder := core.NewPasswordEncoder()
	assert.Nil(err)
	repo := repository.New()
	viper := getViperConfig()
	jwtService, err := jwt.NewService(viper)
	assert.Nil(err)
	roleRouteManager := auth.NewRoleRouteManager(db, repo)
	authService := auth.NewService(db, repo, passwordEncoder, jwtService, getLogger(), roleRouteManager)
	userGrpcServer := usergrpc.NewServer(authService)

	l, err := net.Listen("tcp", "127.0.0.1:0")
	assert.Nil(err)
	defer l.Close()
	googleGrpcServer := grpc.NewServer()

	userpb.RegisterUserServer(googleGrpcServer, userGrpcServer)
	userConn, err := grpc.NewClient(
		l.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.Nil(err)
	go googleGrpcServer.Serve(l)

	u, err := repo.CreateUser(ctx, db, repository.CreateUserParams{
		Email:    "test@test.com",
		Password: "password",
	})
	assert.Nil(err)

	//handling all roles related data
	role1, err := repo.CreateRole(ctx, db, "ROLE1")
	assert.Nil(err)
	role2, err := repo.CreateRole(ctx, db, "ROLE2")
	assert.Nil(err)
	role3, err := repo.CreateRole(ctx, db, "ROLE3")
	assert.Nil(err)

	err = repo.AddRoleToUser(ctx, db, repository.AddRoleToUserParams{
		UserID: u.ID,
		RoleID: role1.ID,
	})
	assert.Nil(err)
	err = repo.AddRoleToUser(ctx, db, repository.AddRoleToUserParams{
		UserID: u.ID,
		RoleID: role2.ID,
	})
	assert.Nil(err)

	route1, err := repo.CreateRoute(ctx, db, "/api/v1/transactions/withdraw")
	assert.Nil(err)

	err = repo.AddRouteToRole(ctx, db, repository.AddRouteToRoleParams{
		RoleID:  role3.ID,
		RouteID: route1.ID,
	})
	assert.Nil(err)

	client := userpb.NewUserClient(userConn)
	token, err := jwtService.GenerateToken("test@test.com")
	assert.Nil(err)

	req := userpb.GetUserDataByTokenRequest{
		Token: token,
		Path:  "/api/v1/transactions/withdraw",
	}
	res, err := client.GetUserDataByToken(ctx, &req)
	//assert.Nil(err)
	//assert.Equal(res.Id, u.ID)
	//assert.Equal(res.Email, "test@test.com")
	assert.Nil(res)
	e, ok := status.FromError(err)
	assert.True(ok)
	assert.Equal(e.Code(), codes.Code(403))
	assert.Equal(e.Message(), "access denied")
}

func TestAuth_GetUserDataByToken_Successful(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	db := getDB()
	err := truncateDB()
	passwordEncoder := core.NewPasswordEncoder()
	assert.Nil(err)
	repo := repository.New()
	viper := getViperConfig()
	jwtService, err := jwt.NewService(viper)
	assert.Nil(err)
	roleRouteManager := auth.NewRoleRouteManager(db, repo)
	authService := auth.NewService(db, repo, passwordEncoder, jwtService, getLogger(), roleRouteManager)
	userGrpcServer := usergrpc.NewServer(authService)

	l, err := net.Listen("tcp", "127.0.0.1:0")
	assert.Nil(err)
	defer l.Close()
	googleGrpcServer := grpc.NewServer()

	userpb.RegisterUserServer(googleGrpcServer, userGrpcServer)
	userConn, err := grpc.NewClient(
		l.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.Nil(err)
	go googleGrpcServer.Serve(l)

	u, err := repo.CreateUser(ctx, db, repository.CreateUserParams{
		Email:    "test@test.com",
		Password: "password",
	})
	assert.Nil(err)

	//handling all roles related data
	role1, err := repo.CreateRole(ctx, db, "ROLE1")
	assert.Nil(err)
	role2, err := repo.CreateRole(ctx, db, "ROLE2")
	assert.Nil(err)
	_, err = repo.CreateRole(ctx, db, "ROLE3")
	assert.Nil(err)

	err = repo.AddRoleToUser(ctx, db, repository.AddRoleToUserParams{
		UserID: u.ID,
		RoleID: role1.ID,
	})
	assert.Nil(err)
	err = repo.AddRoleToUser(ctx, db, repository.AddRoleToUserParams{
		UserID: u.ID,
		RoleID: role2.ID,
	})
	assert.Nil(err)

	route1, err := repo.CreateRoute(ctx, db, "/api/v1/transactions/withdraw")
	assert.Nil(err)

	err = repo.AddRouteToRole(ctx, db, repository.AddRouteToRoleParams{
		RoleID:  role2.ID,
		RouteID: route1.ID,
	})
	assert.Nil(err)

	client := userpb.NewUserClient(userConn)
	token, err := jwtService.GenerateToken("test@test.com")
	assert.Nil(err)

	req := userpb.GetUserDataByTokenRequest{
		Token: token,
		Path:  "/api/v1/transactions/withdraw",
	}
	res, err := client.GetUserDataByToken(ctx, &req)
	assert.Nil(err)
	assert.Equal(res.Id, u.ID)
	assert.Equal(res.Email, "test@test.com")
}
